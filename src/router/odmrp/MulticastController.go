package odmrp

import (
	"fmt"
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

const (
	jqRefreshTime       = 400 * time.Millisecond
	fillFwdTableTimeout = 2 * time.Second
)

type MulticastController struct {
	gmTable         *GroupMembersTable // have grpIP: [destIP1, destIP2, destIP3, ...] where to send?
	queryFlooder    *GlobalFlooder     // control message flooder
	jrConn          *MACLayerConn
	ip              net.IP
	mac             net.HardwareAddr
	forwardingTable *forwardingTable
	cacheTable      *cache        // cache for duplicate checks and for building forwarding table
	memberTable     *membersTable // member table group ips, Am I a destination?
	msec            *MSecLayer    // decreption & encryption
	packetSeqNo     uint64
	ch              chan bool
	refJoinQuery    Timer
	timers          *TimersQueue
}

func NewMulticastController(iface *net.Interface, ip net.IP, mac net.HardwareAddr, msec *MSecLayer, mgrpFilePath string, timers *TimersQueue) (*MulticastController, error) {
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
		mgrpContent = "{\"" + address + "\": [\"10.0.0.15\", \"10.0.0.26\", \"10.0.0.9\"]}"
	} else {
		mgrpContent = readOptionalJsonFile(mgrpFilePath)
	}

	log.Println("initalized multicast controller")

	return &MulticastController{
		gmTable:         NewGroupMembersTable(mgrpContent),
		queryFlooder:    queryFlooder,
		jrConn:          jrConn,
		ip:              ip,
		mac:             mac,
		forwardingTable: newForwardingTable(timers),
		cacheTable:      newCache(timers),
		memberTable:     newMembersTable(timers),
		msec:            msec,
		ch:              make(chan bool),
		packetSeqNo:     0,
		timers:          timers,
	}, nil
}

// GetMissingEntries called by forwarder when it doesn't find and entry
// for given grpIP in the forwarding table
//
// forwarder should put the returned entries in the forwarding table
//
// it may return false in case it can't find any path to the grpIP
// or can't find the grpIP itself
func (c *MulticastController) GetMissingEntries(grpIP net.IP) bool {
	// TODO
	destsIPs, ok := c.gmTable.Get(grpIP)
	if !ok {
		log.Panic("must have the destsIPs!")
	}

	c.sendJoinQuery(grpIP, destsIPs)

	t1 := c.timers.Add(fillFwdTableTimeout, func() {
		for i := 0; i < len(destsIPs); i++ {
			c.ch <- false
		}
	})
	flag := false
	for i := 0; i < len(destsIPs); i++ {
		flag = flag || <-c.ch
	}
	if flag {
		// TODO important stop timer when you want to stop sending to this grpIP
		c.refJoinQuery = c.timers.Add(jqRefreshTime, func() {
			c.sendJoinQuery(grpIP, destsIPs)
		})
	}
	t1.Stop()
	return flag
}

func jqRefreshfireTimer(srcIP net.IP, members []net.IP, c *MulticastController) func() {
	return func() {
		c.sendJoinQuery(srcIP, members)
	}
}

func (c *MulticastController) Start(ft *MultiForwardTable) {
	log.Println("~~ MulticastController started ~~")
	go c.queryFlooder.ListenForFloodedMsgs(c.onRecvJoinQuery)
	go c.onRecvJoinReply(ft)
}

func (c *MulticastController) sendJoinQuery(grpIP net.IP, members []net.IP) {
	c.packetSeqNo++
	jq := joinQuery{
		// TODO encode time to seqNo (Not Sure!!)
		SeqNo:   c.packetSeqNo,
		TTL:     odmrpDefaultTTL,
		SrcIP:   c.ip,
		PrevHop: c.mac,
		GrpIP:   grpIP,
		Dests:   members,
	}

	// insert in cache in case it use broadcast
	cached := &cacheEntry{seqNo: jq.SeqNo, grpIP: jq.GrpIP, prevHop: jq.PrevHop, cost: odmrpDefaultTTL - jq.TTL}
	c.cacheTable.set(jq.SrcIP, cached)

	c.queryFlooder.Flood(jq.marshalBinary())
	log.Println("sent join query to", grpIP) // TODO remove

	c.timers.Add(jqRefreshTime, func() {
		c.sendJoinQuery(grpIP, members)
	})
	// TODO important to stop the timer once the sender stop sending packets to group address
}

func (c *MulticastController) onRecvJoinQuery(fldHdr *FloodHeader, payload []byte) ([]byte, bool) {
	jq, valid := unmarshalJoinQuery(payload)
	log.Printf("(ip:%#v, mac:%#v), Recieved JoinQuery form %#v\n", c.ip.String(), c.mac.String(), jq.PrevHop.String())

	if "10.0.0.13" == c.ip.String() {
		log.Println("JQ important debug")
	}
	log.Println(jq) // TODO: remove this

	if !valid {
		log.Panicln("Corrupted JoinQuery msg received") // TODO: no panicing!
	}

	// // if the join query allready sent
	// // Check if it is a duplicate by comparing the (Source IP Address, Sequence Number) in the cache. DONE
	// cache, ok := c.cacheTable.Get(jq.SrcIP)
	// if ok && cache.SeqNo >= jq.SeqNo {
	// 	return nil, false
	// }

	// else insert join query in cache
	cached := &cacheEntry{seqNo: jq.SeqNo, grpIP: jq.GrpIP, cost: odmrpDefaultTTL - jq.TTL}
	cached.prevHop = make(net.HardwareAddr, len(jq.PrevHop))
	copy(cached.prevHop, jq.PrevHop)
	isCached := c.cacheTable.set(jq.SrcIP, cached)
	if !isCached {
		return nil, false
	}
	// im the prev hop for the next one
	jq.PrevHop = c.mac
	log.Println("Cache after change prev hop")
	log.Println(c.cacheTable.String())

	// grpIPExists := c.memberTable.Get(jq.GrpIP)
	// memberTable for faster recieving
	// if grpIPExists || c.imInDests(jq) {
	if c.imInDests(jq) {
		log.Printf("im in dests, (ip:%#v, mac: %#v)\n", c.ip.String(), c.mac.String()) // TODO: remove this
		// fill member table
		c.memberTable.set(jq.GrpIP)

		// send back join reply to prevHop
		jr := c.generateJoinReply(jq)
		if jr != nil { // impossible equals nil
			c.sendJoinReply(jr)
		}
	}

	// If the TTL field value is less than  0, then discard. DONE
	jq.TTL--
	if jq.TTL < 0 {
		return nil, false
	}

	return jq.marshalBinary(), true
}

func (c *MulticastController) generateJoinReply(jq *joinQuery) *joinReply {
	log.Println("Generate JoinReply")
	jr := &joinReply{
		SeqNo:    jq.SeqNo,
		DestIP:   c.ip,
		GrpIP:    jq.GrpIP,
		PrevHop:  c.mac,
		SrcIPs:   []net.IP{},
		NextHops: []net.HardwareAddr{},
		Cost:     0, // intialize cost, here it is hop count
	}

	// Fill srcIPs and nextHops of the JoinReply
	for item := range c.cacheTable.m.Iter() {
		val := item.Value.(*cacheEntry)
		if val.grpIP.Equal(jr.GrpIP) {
			jr.SrcIPs = append(jr.SrcIPs, UInt32ToIPv4(item.Key.(uint32)))
			jr.NextHops = append(jr.NextHops, val.prevHop)
		}
	}
	if len(jr.SrcIPs) > 0 {
		return jr
	}
	return nil
}

func (c *MulticastController) updateJoinReply(jr *joinReply, ft *MultiForwardTable) *joinReply {
	newJR := &joinReply{
		SeqNo:    jr.SeqNo,
		DestIP:   jr.DestIP,
		GrpIP:    jr.GrpIP,
		PrevHop:  c.mac,
		SrcIPs:   []net.IP{},
		NextHops: []net.HardwareAddr{},
		Cost:     calcNewJRCost(jr),
	}

	// Fill srcIPs and nextHops of the JoinReply
	log.Println("Debug before")
	log.Println(c.cacheTable.String())
	for i := 0; i < len(jr.SrcIPs); i++ {
		if !jr.SrcIPs[i].Equal(c.ip) { // TODO check if I can remove this if
			cacheEntry, ok := c.cacheTable.get(jr.SrcIPs[i])
			if ok && cacheEntry.grpIP.Equal(jr.GrpIP) {
				newJR.SrcIPs = append(newJR.SrcIPs, jr.SrcIPs[i])
				newJR.NextHops = append(newJR.NextHops, cacheEntry.prevHop)
			}
		}
	}
	logIPs("Old join reply source ips", jr.SrcIPs)
	logMacIPs("Old join reply next hops ips", jr.NextHops)
	logIPs("New join reply source ips", newJR.SrcIPs)
	logMacIPs("New join reply next hops ips", newJR.NextHops)
	log.Println("print forwarding table")
	log.Println(c.forwardingTable.String())
	log.Println("print multiforwarding table")
	log.Println(ft.String())

	if len(newJR.SrcIPs) > 0 {
		log.Println("Debug cache after")
		log.Println(c.cacheTable.String())
		return newJR
	}
	log.Println("Debug cache after")
	log.Println(c.cacheTable.String())
	return nil
}

func calcNewJRCost(jr *joinReply) uint16 {
	return jr.Cost + 1
}

func (c *MulticastController) sendJoinReply(jr *joinReply) {
	encJR := c.msec.Encrypt(jr.marshalBinary())
	for i := 0; i < len(jr.NextHops); i++ {
		c.jrConn.Write(encJR, jr.NextHops[i])
	}
	msg := fmt.Sprintf("ip: %#v sends Join Reply to", c.ip.String())
	logMacIPs(msg, jr.NextHops)
}

func (c *MulticastController) onRecvJoinReply(ft *MultiForwardTable) {
	for {
		msg := c.jrConn.Read()
		log.Println("Recieved Join Reply!!")
		decryptedJR := c.msec.Decrypt(msg)

		jr, valid := unmarshalJoinReply(decryptedJR)
		if !valid {
			log.Panicln("Corrupted JoinReply msg received")
		}

		// TODO remove log
		log.Printf("(ip:%#v, mac:%#v), Recieved JoinReply form %#v\n", c.ip.String(), c.mac.String(), jr.PrevHop.String())
		log.Println("Before")
		log.Println(jr.String())
		log.Println(c.cacheTable.String())

		// update forwarding table
		forwardingEntry := &forwardingEntry{nextHop: jr.PrevHop, cost: jr.Cost}
		refreshForwarder := c.forwardingTable.set(jr.DestIP, forwardingEntry)
		if refreshForwarder {
			ft.Set(jr.GrpIP, jr.PrevHop)
		}

		if c.imInSrcs(jr) {
			log.Println("Source Recieved Join Reply!!")
			newJR := c.updateJoinReply(jr, ft)
			if newJR != nil {
				c.sendJoinReply(newJR)
			}
			c.ch <- true
		} else {
			newJR := c.updateJoinReply(jr, ft)
			if newJR != nil {
				c.sendJoinReply(newJR)
			}
		}
		log.Println("Final Cache Debug After Recieve JoinReply")
		log.Println(c.cacheTable.String())
	}
}

func (c *MulticastController) imInDests(jq *joinQuery) bool {
	for i := 0; i < len(jq.Dests); i++ {
		if c.ip.Equal(jq.Dests[i]) {
			return true
		}
	}
	return false
}

func (c *MulticastController) imInSrcs(jr *joinReply) bool {
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

func logMacIPs(msg string, macIPs []net.HardwareAddr) {
	msg += ": {"
	for i := 0; i < len(macIPs); i++ {
		msg += fmt.Sprintf("%#v, ", macIPs[i].String())
	}
	msg += "}"
	log.Println(msg)
}

func logIPs(msg string, ips []net.IP) {
	msg += ": {"
	for i := 0; i < len(ips); i++ {
		msg += fmt.Sprintf("%#v, ", ips[i].String())
	}
	msg += "}"
	log.Println(msg)
}

func (c *MulticastController) IsDest(grpIP net.IP) bool {
	return c.memberTable.get(grpIP)
}
