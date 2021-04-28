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
	routingTable *RoutingTable
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
		routingTable: newRoutingTable(),
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

	if c.imInDests(jq) {
		log.Println("im in dests! :'D") // TODO: remove this

		// send back join reply to prevHop
		jr := &JoinReply{
			SeqNo:  jq.SeqNo,
			SrcIPs: []net.IP{jq.SrcIP},
			GrpIPs: []net.IP{jq.GrpIP},
		}
		encJR := c.msec.Encrypt(jr.MarshalBinary())
		c.jrConn.Write(encJR, jq.PrevHop)

		return nil, false
	}

	// jr's nextHop is this jq's prevHop
	c.routingTable.Set(jq.SrcIP, &routingEntry{nextHop: jq.PrevHop})

	// im the prev hop for the next one
	jq.PrevHop = c.mac

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

		if c.imInSrcs(jr) {
			log.Println("Source Recieved Join Reply!!")
		} else {
			jr.Forwarders = append(jr.Forwarders, c.ip)
			msg = c.msec.Encrypt(jr.MarshalBinary())
			for _, srcIP := range jr.SrcIPs {
				entryroutingTable, ok := c.routingTable.Get(srcIP)
				if ok {
					c.jrConn.Write(msg, entryroutingTable.nextHop)
				}
			}
		}
		log.Println(jr)
	}
}

func (c *MulticastController) imInSrcs(jr *JoinReply) bool {
	for i := 0; i < len(jr.SrcIPs); i++ {
		if c.ip.Equal(jr.SrcIPs[i]) {
			return true
		}
	}
	return false
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
