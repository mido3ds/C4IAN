package odmrp

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	. "github.com/mido3ds/C4IAN/src/router/flood"
	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
)

type MulticastController struct {
	gmTable      *GroupMembersTable
	queryFlooder *GlobalFlooder
	jrConn       *MACLayerConn
	ip           net.IP
	mac          net.HardwareAddr
	jrFTable     *jrForwardTable
	msec         *MSecLayer
}

func NewMulticastController(iface *net.Interface, ip net.IP, mac net.HardwareAddr, msec *MSecLayer, mgrpFilePath string) (*MulticastController, error) {
	queryFlooder, err := NewGlobalFlooder(ip, iface, JoinQueryEtherType, msec)
	if err != nil {
		log.Panic("failed to initiate query flooder, err: ", err)
	}

	jrConn, err := NewMACLayerConn(iface, JoinReplyEtherType)
	if err != nil {
		log.Panic("failed to initiate mac conn, err: ", err)
	}

	// read mgroup
	var mgrpContent string
	if os.Getenv("MTEST") == "1" {
		address := "224.0.2.1"
		log.Printf("multicast test mode, ping address={%s} from any node to start mcasting\n", address)
		mgrpContent = "{\"" + address + "\": [\"10.0.0.14\", \"10.0.0.15\", \"10.0.0.16\"]}"
	} else {
		mgrpContent = readOptionalJsonFile(mgrpFilePath)
	}

	log.Println("initalized multicast controller")

	return &MulticastController{
		gmTable:      NewGroupMembersTable(mgrpContent),
		queryFlooder: queryFlooder,
		jrConn:       jrConn,
		ip:           ip,
		mac:          mac,
		jrFTable:     newJRForwardTable(),
		msec:         msec,
	}, nil
}

// GetMissingEntries called by forwarder when it doesn't find and entry
// for given grpIP in the forwarding table
//
// forwarder should put the returned entries in the forwarding table
//
// it may return false in case it can't find any path to the grpIP
// or can't find the grpIP itself
func (c *MulticastController) GetMissingEntries(grpIP net.IP) (*MultiForwardingEntry, bool) {
	// TODO

	// for now i will just send join queries
	members, ok := c.gmTable.Get(grpIP)
	if !ok {
		log.Panic("must have the members!")
	}
	jq := JoinQuery{
		SeqNo:   1,
		TTL:     ODMRPDefaultTTL,
		SrcIP:   c.ip,
		PrevHop: c.mac,
		GrpIP:   grpIP,
		Dests:   members,
	}
	c.queryFlooder.Flood(jq.MarshalBinary())
	log.Println("sent join query to", grpIP) // TODO remove

	return nil, false
}

func (c *MulticastController) Start(ft *MultiForwardTable) {
	log.Println("MulticastController started listening for control packets from the forwarder")
	go c.queryFlooder.ReceiveFloodedMsgs(c.onRecvJoinQuery)
	go c.recvJoinReplyMsgs(ft)
}

func (c *MulticastController) onRecvJoinQuery(fldHdr *FloodHeader, payload []byte) ([]byte, bool) {
	jq, valid := UnmarshalJoinQuery(payload)
	log.Println(jq) // TODO: remove this
	if !valid {
		log.Panicln("Corrupted JoinQuery msg received") // TODO: no panicing!
	}

	jq.TTL--
	if jq.TTL < 0 {
		return nil, false
	}

	// register jq for its jr to get back to src
	// jr's nextHop is this jq's prevHop
	c.jrFTable.Set(jq.SrcIP, &jrForwardEntry{seqNum: jq.SeqNo, nextHop: jq.PrevHop})

	jq.PrevHop = c.mac

	if c.imInDests(jq) {
		log.Println("im in dests! :'D") // TODO: remove this
		// TODO: reply with join reply
		jrEntry, ok := c.jrFTable.Get(jq.SrcIP)
		if ok && jrEntry.seqNum <= jq.SeqNo {
			jr := &JoinReply{
				SeqNo:  jq.SeqNo,
				SrcIPs: []net.IP{jq.SrcIP},
				GrpIPs: []net.IP{jq.GrpIP},
			}
			msg := jr.MarshalBinary()
			entryJrFTable, ok := c.jrFTable.Get(jq.SrcIP)
			if ok {
				c.jrConn.Write(c.msec.Encrypt(msg), entryJrFTable.nextHop)
			}
		}
		return nil, false
	}

	return jq.MarshalBinary(), true
}

func (c *MulticastController) imInDests(jq *JoinQuery) bool {
	for i := 0; i < len(jq.Dests); i++ {
		if c.ip.Equal(jq.Dests[i]) {
			return true
		}
	}
	return false
}

func (c *MulticastController) recvJoinReplyMsgs(ft *MultiForwardTable) {
	for {
		msg := c.jrConn.Read()
		log.Println("Recieved Join Reply!!")
		pd := c.msec.NewPacketDecrypter(msg)
		decryptedJR := pd.DecryptAll()

		jr, valid := UnmarshalJoinReply(decryptedJR)
		if !valid {
			log.Panicln("Corrupted JoinReply msg received")
		}
		log.Println(jr)
		// TODO: store msg
		// TODO: resend to next hop, unless im source
	}
}

func (c *MulticastController) Close() {
	c.jrConn.Close()
	c.queryFlooder.Close()
}

func readOptionalJsonFile(path string) string {
	if path != "" {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Panic(err)
		}
		return string(content)
	}
	return "{}"
}
