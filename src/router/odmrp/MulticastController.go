package odmrp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"

	. "github.com/mido3ds/C4IAN/src/router/flood"
	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
)

const (
	jqRefreshTime       = 400 * time.Millisecond
	fillFwdTableTimeout = 5 * time.Second
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
	refJoinQuery    *Timer
	timers          *TimersQueue
	// TODO make dest doesn't send join reply if it doesn't want to this grpIP
}

func NewMulticastController(iface *net.Interface, ip net.IP, mac net.HardwareAddr, msec *MSecLayer, mgrpFilePath string, timers *TimersQueue) (*MulticastController, error) {
	queryFlooder := NewGlobalFlooder(ip, iface, JoinQueryEtherType, msec)

	jrConn, err := NewMACLayerConn(iface, JoinReplyEtherType)
	if err != nil {
		log.Panic("failed to initiate mac conn, err: ", err)
	}

	// read mgroup
	var mgrpContent string
	if os.Getenv("MTEST") == "1" {
		address := "224.0.2.1"

		// pass members ids in MEMS env var
		// like MEMS=5,14,20
		var membersIPs []string
		for _, i := range strings.Split(os.Getenv("MEMS"), ",") {
			membersIPs = append(membersIPs, "\"10.0.0."+i+"\"")
		}

		// pass src sender in MSRC
		// like MSRC=3
		src := "10.0.0.1"
		if msrc := os.Getenv("MSRC"); msrc != "" {
			src = "10.0.0." + msrc
		}

		// start sender & receivers
		if ip.String() == src {
			go sendUDPPackets(address)
		} else {
			for _, ip2 := range membersIPs {
				if "\""+ip.String()+"\"" == ip2 {
					go receiveUDPPackets(address)
				}
			}
		}

		log.Printf("multicast test mode, adr={%v}, members={%v}\n", address, membersIPs)
		mgrpContent = "{\"" + address + "\": [" + strings.Join(membersIPs, ", ") + "]}"
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

func sendUDPPackets(address string) {
	adr, err := net.ResolveUDPAddr("udp", address+":1234")
	if err != nil {
		log.Panic(err)
	}

	conn, err := net.DialUDP("udp", nil, adr)
	if err != nil {
		log.Panic(err)
	}

	i := 0
	log.Println("******sending udp packets")
	for {
		time.Sleep(5 * time.Second)
		msg := fmt.Sprintf("message #%v", i)
		i++
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Panic(err)
		}
		log.Println("******sent:", msg)
	}
}

func receiveUDPPackets(address string) {
	adr, err := net.ResolveUDPAddr("udp", address+":1234")
	if err != nil {
		log.Panic(err)
	}

	conn, err := net.ListenUDP("udp", adr)
	if err != nil {
		log.Panic(err)
	}

	log.Println("+++++++receiving udp packets")
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			log.Panic(err)
		}
		log.Println("+++++++received", string(buf))
	}
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
	t1.Stop()
	log.Println("finished")
	return flag
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

	encryptedJQ := c.msec.Encrypt(jq.marshalBinary())
	c.queryFlooder.Flood(encryptedJQ)
	log.Println("sent join query to", grpIP) // TODO remove

	// TODO important stop timer when you want to stop sending to this grpIP
	c.refJoinQuery = c.timers.Add(jqRefreshTime, func() {
		// c.sendJoinQuery(grpIP, members)
	})
}

func (c *MulticastController) onRecvJoinQuery(encryptedPayload []byte) []byte {
	payload := c.msec.Decrypt(encryptedPayload)

	jq, valid := unmarshalJoinQuery(payload)
	if !valid {
		log.Panicln("Corrupted JoinQuery msg received") // TODO: no panicing!
	}

	log.Printf("(ip:%#v, mac:%#v), Recieved JoinQuery form %#v\n", c.ip.String(), c.mac.String(), jq.PrevHop.String())

	log.Println(jq) // TODO: remove this

	// // if the join query allready sent
	// // Check if it is a duplicate by comparing the (Source IP Address, Sequence Number) in the cache. DONE
	// cache, ok := c.cacheTable.Get(jq.SrcIP)
	// if ok && cache.SeqNo >= jq.SeqNo {
	// 	return nil, false
	// }

	// If the TTL field value is less than  0, then discard. DONE
	jq.TTL--
	if jq.TTL < 0 {
		return nil
	}

	// else insert join query in cache
	cached := &cacheEntry{seqNo: jq.SeqNo, grpIP: jq.GrpIP, cost: odmrpDefaultTTL - jq.TTL}
	cached.prevHop = make(net.HardwareAddr, len(jq.PrevHop))
	copy(cached.prevHop, jq.PrevHop)
	isCached := c.cacheTable.set(jq.SrcIP, cached)
	if !isCached {
		return nil
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

	return c.msec.Encrypt(jq.marshalBinary())
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
		Cost:     1, // intialize cost, here it is hop count
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
	jr.Cost = calcNewJRCost(jr)
	jr.PrevHop = c.mac
	newSrcIPs := []net.IP{}
	newNextHops := []net.HardwareAddr{}

	// Fill srcIPs and nextHops of the JoinReply
	log.Println("Debug before")
	log.Println(c.cacheTable.String())
	for i := 0; i < len(jr.SrcIPs); i++ {
		if !jr.SrcIPs[i].Equal(c.ip) { // TODO check if I can remove this if
			cacheEntry, ok := c.cacheTable.get(jr.SrcIPs[i])
			if ok && cacheEntry.grpIP.Equal(jr.GrpIP) {
				newSrcIPs = append(newSrcIPs, jr.SrcIPs[i])
				newNextHops = append(newNextHops, cacheEntry.prevHop)
			}
		}
	}

	if len(newSrcIPs) > 0 {
		jr.SrcIPs = newSrcIPs
		jr.NextHops = newNextHops
		return jr
	}
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
	msg := fmt.Sprintf("ip: %#v sends JoinReply to", c.ip.String())
	logMacIPs(msg, jr.NextHops)
}

func (c *MulticastController) onRecvJoinReply(ft *MultiForwardTable) {
	for {
		msg := c.jrConn.Read()
		go c.handleJoinReply(msg, ft)
	}
}

func (c *MulticastController) handleJoinReply(msg []byte, ft *MultiForwardTable) {
	decryptedJR := c.msec.Decrypt(msg)

	jr, valid := unmarshalJoinReply(decryptedJR)
	if !valid {
		log.Panicln("Corrupted JoinReply msg received")
	}

	// TODO remove log
	log.Printf("Recieved JoinReply %#v\n", jr.PrevHop.String())

	// update forwarding table
	forwardingEntry := &forwardingEntry{nextHop: jr.PrevHop, cost: jr.Cost}
	refreshForwarder := c.forwardingTable.set(jr.DestIP, forwardingEntry)
	if refreshForwarder {
		ft.Set(jr.GrpIP, jr.PrevHop)
	}

	if c.imInSrcs(jr) {
		log.Println("Source Recieved JoinReply !!")
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
	log.Println("Cache After Recieve JoinReply")
	log.Println(c.cacheTable.String())

	log.Println("Forwarding Tables After Recieve JoinReply")
	log.Println(c.forwardingTable.String())
	log.Println(ft.String())
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
