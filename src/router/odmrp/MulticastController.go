package odmrp

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	. "github.com/mido3ds/C4IAN/src/router/flood"
	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
)

const JQ_REFRESH_TIME = 400 * time.Millisecond

type MulticastController struct {
	gmTable        *GroupMembersTable
	queryFlooder   *GlobalFlooder
	jrConn         *MACLayerConn
	ip             net.IP
	mac            net.HardwareAddr
	routingTable   *RoutingTable
	cacheTable     *CacheTable
	memberTable    *MemberTable
	msec           *MSecLayer
	jQRefreshTimer *time.Timer
	packetSeqNo    uint64
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
		cacheTable:   newCacheTable(),
		memberTable:  newMemberTable(),
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
	members, ok := c.gmTable.Get(grpIP)
	if !ok {
		log.Panic("must have the members!")
	}

	c.sendJoinQuery(grpIP, members)

	return nil, false
}

func jqRefreshfireTimerHelper(srcIP net.IP, members []net.IP, c *MulticastController) {
	c.sendJoinQuery(srcIP, members)
}

func jqRefreshfireTimer(srcIP net.IP, members []net.IP, c *MulticastController) func() {
	return func() {
		jqRefreshfireTimerHelper(srcIP, members, c)
	}
}

func (c *MulticastController) sendJoinQuery(grpIP net.IP, members []net.IP) {
	c.packetSeqNo++
	jq := JoinQuery{
		// TODO encode time to seqNo (Not Sure!!)
		SeqNo:   c.packetSeqNo,
		TTL:     ODMRPDefaultTTL,
		SrcIP:   c.ip,
		PrevHop: c.mac,
		GrpIP:   grpIP,
		Dests:   members,
	}

	// insert in cache in case it use broadcast
	cached := &cacheEntry{SeqNo: jq.SeqNo, GrpIP: jq.GrpIP, PrevHop: jq.PrevHop}
	c.cacheTable.Set(jq.SrcIP, cached)

	c.queryFlooder.Flood(jq.MarshalBinary())
	log.Println("sent join query to", grpIP) // TODO remove

	// fireFunc := jqRefreshfireTimer(grpIP, members, c)
	// c.jQRefreshTimer = time.AfterFunc(JQ_REFRESH_TIME, fireFunc)
	// TODO important to stop the timer once the sender stop sending packets to group address
}

func (c *MulticastController) Start(ft *MultiForwardTable) {
	log.Println("~~ MulticastController started ~~")
	go c.queryFlooder.ReceiveFloodedMsgs(c.onRecvJoinQuery)
	go c.recvJoinReplyMsgs(ft)
}

func (c *MulticastController) onRecvJoinQuery(fldHdr *FloodHeader, payload []byte) ([]byte, bool) {
	jq, valid := UnmarshalJoinQuery(payload)
	log.Println(jq) // TODO: remove this

	if !valid {
		log.Panicln("Corrupted JoinQuery msg received") // TODO: no panicing!
	}

	// if the join query allready sent
	// Check if it is a duplicate by comparing the (Source IP Address, Sequence Number) in the cache. DONE
	cache, ok := c.cacheTable.Get(jq.SrcIP)
	if ok && cache.SeqNo >= jq.SeqNo {
		return nil, false
	}
	// else insert in cache==
	cached := &cacheEntry{SeqNo: jq.SeqNo, GrpIP: jq.GrpIP, PrevHop: jq.PrevHop}
	c.cacheTable.Set(jq.SrcIP, cached)

	// jr's nextHop is this jq's prevHop
	c.routingTable.Set(jq.SrcIP, &routingEntry{nextHop: jq.PrevHop})

	// im the prev hop for the next one
	jq.PrevHop = c.mac

	if c.imInDests(jq) {
		log.Println("im in dests! :'D") // TODO: remove this

		// fill member table
		entry := &memberEntry{grpIP: jq.GrpIP}
		c.memberTable.Set(jq.SrcIP, entry)

		// send back join reply to prevHop
		jr := c.destGenerateJoinReply(jq)
		encJR := c.msec.Encrypt(jr.MarshalBinary())
		c.jrConn.Write(encJR, jq.PrevHop)
	}

	// If the TTL field value is less than  0, then discard. DONE
	jq.TTL--
	if jq.TTL < 0 {
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

		if c.imInSrcs(jr) {
			log.Println("Source Recieved Join Reply!!")
		} else {
			// jr.Forwarders = append(jr.Forwarders, c.ip)
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

func (c *MulticastController) destGenerateJoinReply(jq *JoinQuery) *JoinReply {
	jr := &JoinReply{
		SeqNo:    jq.SeqNo,
		DestIP:   c.ip,
		GrpIP:    jq.GrpIP,
		PrevHop:  c.mac,
		SrcIPs:   []net.IP{},
		NextHops: []net.HardwareAddr{},
	}

	// TODO think for better/faster way
	for item := range c.memberTable.m.Iter() {
		src := item.Key.(uint32)
		member, ok := c.routingTable.m.Get(src)
		if ok {
			nextHop := member.(*routingEntry).nextHop
			jr.SrcIPs = append(jr.SrcIPs, UInt32ToIPv4(src))
			jr.NextHops = append(jr.NextHops, nextHop)
		}
	}

	return jr
}
