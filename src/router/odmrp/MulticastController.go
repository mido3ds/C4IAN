package odmrp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"

	. "github.com/mido3ds/C4IAN/src/router/constants"
	. "github.com/mido3ds/C4IAN/src/router/flood"
	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
)

type MulticastController struct {
	gmTable         *GroupMembersTable // have grpIP: [destIP1, destIP2, destIP3, ...] where to send?
	queryFlooder    *GlobalFlooder     // control message flooder
	jrConn          *MACLayerConn      // nac layer connection for join reply
	ip              net.IP             // ip of the router
	mac             net.HardwareAddr   // mac address of the router
	forwardingTable *forwardingTable   // forwarding table using in routing
	cacheTable      *cache             // cache for duplicate checks and for building forwarding table
	memberTable     *membersTable      // member table group ips, Am I a destination?
	msec            *MSecLayer         // decreption & encryption
	packetSeqNo     uint64             // packet unique identifier
	ch              chan byte          // channel used by threads to exit the get missing entries of the multi forward table function
	refJoinQuery    *Timer             // timer to send join query periodically until router calls stop
	timers          *TimersQueue       // timers queue used by tables to delete the expired entries
}

// NewMulticastController creates a new multicast controller to do multicast routing and filling the multi forward table
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
		// pass members ids of the multicast group in MEMS env var used like this MEMS=5,14,20
		var membersIPs []string
		for _, i := range strings.Split(os.Getenv("MEMS"), ",") {
			membersIPs = append(membersIPs, "\"10.0.0."+i+"\"")
		}

		// pass src sender in MSRC used like this MSRC=3
		// sudo MTEST=1 MSRC=1, MEMS=5,6,10 ./start routers
		src := "10.0.0.1"
		if msrc := os.Getenv("MSRC"); msrc != "" {
			src = "10.0.0." + msrc
		}

		// start sender & receivers
		if ip.String() == src {
			// start sending udp packets
			go sendUDPPackets(address)
		} else {
			for _, ip2 := range membersIPs {
				if "\""+ip.String()+"\"" == ip2 {
					go receiveUDPPackets(address)
				}
			}
		}

		// log.Printf("multicast test mode, adr={%v}, members={%v}\n", address, membersIPs)
		mgrpContent = "{\"" + address + "\": [" + strings.Join(membersIPs, ", ") + "]}"
	} else {
		// if not in the test mode read group members table from json file
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
		ch:              make(chan byte),
		packetSeqNo:     0,
		timers:          timers,
	}, nil
}

// start sending udp packets periodically
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
		// every 5 seconds send a packet
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

// when recieving udp packets this function is called to do some logs
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
// forwarder should put the returned entries in the multi forwarding table
//
// it may return false in case it can't find any path to the grpIP
// or can't find the grpIP itself
func (c *MulticastController) GetMissingEntries(grpIP net.IP) bool {
	destsIPs, ok := c.gmTable.Get(grpIP)
	if !ok {
		log.Panic("must have the destsIPs!")
	}

	// To get missing entries start sending join query from the source
	c.sendJoinQuery(grpIP, destsIPs)

	// Add max timeout to fill Forward Table
	t1 := c.timers.Add(FillForwardTableTimeout, func() {
		for i := 0; i < len(destsIPs); i++ {
			c.ch <- 0
		}
	})

	// Wait until timeout or recieve join reply from a destination
	var dstsCount byte = 0
	for i := 0; i < len(destsIPs); i++ {
		dstsCount += <-c.ch
	}
	// stop timer
	t1.Stop()
	// true if destination(s) is/are found, false if didn't recieve a join reply from a destination
	return dstsCount > 0
}

// Start called by forwarder when to start the multicast controller
// and filling the required tables
//
// forwarder should put the returned entries in the forwarding table
//
// it starts two threads for listening for the flooded join query messages
// and listening for join reply messages
func (c *MulticastController) Start(ft *MultiForwardTable) {
	log.Println("~~ MulticastController started ~~")
	go c.queryFlooder.ListenForFloodedMsgs(c.onRecvJoinQuery)
	go c.onRecvJoinReply(ft)
}

// sendJoinQuery start generating and sending join query messages and it called periodically
// on-demand when the source tries to send packets
//
// it starts two threads for listening for the flooded join query messages
// and listening for join reply messages
func (c *MulticastController) sendJoinQuery(grpIP net.IP, members []net.IP) {
	// change the unique identifier for the join query packet
	c.packetSeqNo++
	jq := joinQuery{
		seqNum:  c.packetSeqNo,
		ttl:     ODMRPDefaultTTL,
		srcIP:   c.ip,
		prevHop: c.mac,
		grpIP:   grpIP,
		dests:   members,
	}

	// insert in cache in case it use broadcast
	cached := &cacheEntry{seqNo: jq.seqNum, grpIP: jq.grpIP, prevHop: jq.prevHop, cost: ODMRPDefaultTTL - jq.ttl}
	c.cacheTable.set(jq.srcIP, cached)

	encryptedJQ := c.msec.Encrypt(jq.marshalBinary())
	c.queryFlooder.Flood(encryptedJQ)
	// log.Println("sent join query to", grpIP)

	// To keep table up to date consistantly send join query and recieve join replies to fill tables
	// When you wants to stop call stopSending() func
	c.refJoinQuery = c.timers.Add(JQRefreshTime, func() {
		c.sendJoinQuery(grpIP, members)
	})
}

// stop sending join query packet it is called whenever the router wants to stop sending packets
func (c *MulticastController) stopSending() {
	c.refJoinQuery.Stop()
}

// onRecvJoinQuery when the router receives join query control packet it encrypts it
// then updates it and fills the required tables then continues the flooding process
// on-demand when the source tries to send packets
//
// and it checks if the router is a destination and if true starts generating the join reply
// control message and pass it back to the source router
func (c *MulticastController) onRecvJoinQuery(encryptedPayload []byte) []byte {
	payload := c.msec.Decrypt(encryptedPayload)

	jq, valid := unmarshalJoinQuery(payload)
	if !valid {
		log.Panicln("Corrupted JoinQuery msg received") // TODO: no panicing!
	}
	// log.Printf("(ip:%#v, mac:%#v), Recieved JoinQuery form %#v\n", c.ip.String(), c.mac.String(), jq.prevHop.String())
	// log.Println(jq)
	// If the TTL field value is less than  0, then discard. DONE
	jq.ttl--
	if jq.ttl < 0 {
		return nil
	}

	// else insert join query in cache
	cached := &cacheEntry{seqNo: jq.seqNum, grpIP: jq.grpIP, cost: ODMRPDefaultTTL - jq.ttl}
	cached.prevHop = make(net.HardwareAddr, len(jq.prevHop))
	copy(cached.prevHop, jq.prevHop)
	isCached := c.cacheTable.set(jq.srcIP, cached)
	if !isCached {
		return nil
	}
	// im the prev hop for the next one
	jq.prevHop = c.mac
	// log.Println("Cache after change prev hop")
	// log.Println(c.cacheTable.String())

	grpIPExists := c.memberTable.get(jq.grpIP)
	// if c.imInDests(jq) {
	// memberTable for faster recieving
	if grpIPExists || c.imInDests(jq) {
		// log.Printf("im in dests, (ip:%#v, mac: %#v)\n", c.ip.String(), c.mac.String()) // TODO: remove this
		// fill member table
		c.memberTable.set(jq.grpIP)

		// send back join reply to prevHop
		jr := c.generateJoinReply(jq)
		if jr != nil { // impossible equals nil
			c.sendJoinReply(jr)
		}
	}

	return c.msec.Encrypt(jq.marshalBinary())
}

// onRecvJoinQuery when the router receives join query control packet it decrypts it
// then updates it and fills the required tables then continues the flooding process
// on-demand when the source tries to send packets
//
// and it checks if the router is a destination and if true starts generating the join reply
// control message and pass it back to the source router
func (c *MulticastController) generateJoinReply(jq *joinQuery) *joinReply {
	// log.Println("Generate JoinReply")
	jr := &joinReply{
		seqNum:   jq.seqNum,
		destIP:   c.ip,
		grpIP:    jq.grpIP,
		prevHop:  c.mac,
		srcIPs:   []net.IP{},
		nextHops: []net.HardwareAddr{},
		cost:     1, // intialize cost, here it is hop count
	}

	// Fill srcIPs and nextHops of the JoinReply
	for item := range c.cacheTable.m.Iter() {
		val := item.Value.(*cacheEntry)
		if val.grpIP.Equal(jr.grpIP) {
			jr.srcIPs = append(jr.srcIPs, UInt32ToIPv4(item.Key.(uint32)))
			jr.nextHops = append(jr.nextHops, val.prevHop)
		}
	}
	if len(jr.srcIPs) > 0 {
		return jr
	}
	return nil
}

// updateJoinReply when the router receives join reply control packet it decrypts it
// then updates it and fills the required tables then returns nil if no next hops found
func (c *MulticastController) updateJoinReply(jr *joinReply, ft *MultiForwardTable) *joinReply {
	jr.cost = calcNewJRCost(jr)
	jr.prevHop = c.mac
	newSrcIPs := []net.IP{}
	newNextHops := []net.HardwareAddr{}

	// Fill srcIPs and nextHops of the JoinReply
	// log.Println(c.cacheTable.String())
	for i := 0; i < len(jr.srcIPs); i++ {
		cacheEntry, ok := c.cacheTable.get(jr.srcIPs[i])
		if ok && cacheEntry.grpIP.Equal(jr.grpIP) {
			newSrcIPs = append(newSrcIPs, jr.srcIPs[i])
			newNextHops = append(newNextHops, cacheEntry.prevHop)
		}
	}

	if len(newSrcIPs) > 0 {
		jr.srcIPs = newSrcIPs
		jr.nextHops = newNextHops
		return jr
	}
	return nil
}

func calcNewJRCost(jr *joinReply) uint16 {
	return jr.cost + 1
}

// sendJoinReply start sending join replies packet using a set of next hops
func (c *MulticastController) sendJoinReply(jr *joinReply) {
	encJR := c.msec.Encrypt(jr.marshalBinary())
	for i := 0; i < len(jr.nextHops); i++ {
		c.jrConn.Write(encJR, jr.nextHops[i])
	}
	// msg := fmt.Sprintf("ip: %#v sends JoinReply to", c.ip.String())
	// logMacAddresses(msg, jr.nextHops)
}

// onRecvJoinReply reads from the mac layer and start handling the recevied join reply control packet
func (c *MulticastController) onRecvJoinReply(ft *MultiForwardTable) {
	for {
		msg := c.jrConn.Read()
		go c.handleJoinReply(msg, ft)
	}
}

// handleJoinReply when the router receives join reply it calls this function
// to updates it the join reply packet and fills the required tables then continues
// sending the join reply to the sources process
//
// and it checks if the router is a souce and if true this means it found route between source
// and this one destination and  sends 1 to the get missing entries function
// through the multicast controller channel
func (c *MulticastController) handleJoinReply(msg []byte, ft *MultiForwardTable) {
	decryptedJR := c.msec.Decrypt(msg)

	jr, valid := unmarshalJoinReply(decryptedJR)
	if !valid {
		log.Panicln("Corrupted JoinReply msg received")
	}

	// log.Printf("Recieved JoinReply %#v\n", jr.prevHop.String())
	// fills forward table
	forwardingEntry := &forwardingEntry{nextHop: jr.prevHop, cost: jr.cost}
	refreshForwarder := c.forwardingTable.set(jr.destIP, forwardingEntry)
	// fills multi forward table if forward table is updated
	if refreshForwarder {
		ft.Set(jr.grpIP, jr.prevHop)
	}

	// check if a i am a source node
	if c.imInSrcs(jr) {
		// log.Println("Source Recieved JoinReply !!")
		newJR := c.updateJoinReply(jr, ft)
		if newJR != nil {
			c.sendJoinReply(newJR)
		}
		// pass 1 as I recevied a join reply from a destination and I am a source node
		c.ch <- 1
	} else {
		newJR := c.updateJoinReply(jr, ft)
		if newJR != nil {
			c.sendJoinReply(newJR)
		}
	}

	// log.Println("Cache After Recieve JoinReply")
	// log.Println(c.cacheTable)
	// log.Println("Forwarding Tables After Recieve JoinReply")
	// log.Println(c.forwardingTable)
	// log.Println(ft)
}

// checks if the router is a destination node
func (c *MulticastController) imInDests(jq *joinQuery) bool {
	for i := 0; i < len(jq.dests); i++ {
		if c.ip.Equal(jq.dests[i]) {
			return true
		}
	}
	return false
}

// checks if the router is a source node
func (c *MulticastController) imInSrcs(jr *joinReply) bool {
	for i := 0; i < len(jr.srcIPs); i++ {
		if c.ip.Equal(jr.srcIPs[i]) {
			return true
		}
	}
	return false
}

// Close closes multicast controller with its connections
func (c *MulticastController) Close() {
	c.jrConn.Close()
	c.queryFlooder.Close()
}

// readOptionalJsonFile reads json file from a file path
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

// logMacAddresses logs mac addresses ipv4
func logMacAddresses(msg string, macIPs []net.HardwareAddr) {
	msg += ": {"
	for i := 0; i < len(macIPs); i++ {
		msg += fmt.Sprintf("%#v, ", macIPs[i].String())
	}
	msg += "}"
	log.Println(msg)
}

// logIPs logs ip addresses ipv4
func logIPs(msg string, ips []net.IP) {
	msg += ": {"
	for i := 0; i < len(ips); i++ {
		msg += fmt.Sprintf("%#v, ", ips[i].String())
	}
	msg += "}"
	log.Println(msg)
}

// IsDest checks if the router is a destination from the members table
func (c *MulticastController) IsDest(grpIP net.IP) bool {
	return c.memberTable.get(grpIP)
}
