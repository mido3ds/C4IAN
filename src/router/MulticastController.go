package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/mdlayher/ethernet"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
)

const (
	// Make use of an unassigned EtherType to differentiate between odmrp traffic and other traffic
	// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml
	joinQueryEtherType = ethernet.EtherType(0x0901)
	joinReplyEtherType = ethernet.EtherType(0x0902)
)

type MulticastController struct {
	gmTable      *GroupMembersTable
	queryFlooder *GlobalFlooder
	jrConn       *MACLayerConn
}

func NewMulticastController(iface *net.Interface, ip net.IP, msec *MSecLayer, mgrpFilePath string) (*MulticastController, error) {
	queryFlooder, err := NewGlobalFlooder(ip, iface, joinQueryEtherType, msec)
	if err != nil {
		log.Panic("failed to initiate query flooder, err: ", err)
	}

	jrConn, err := NewMACLayerConn(iface, joinReplyEtherType)
	if err != nil {
		log.Panic("failed to initiate mac conn, err: ", err)
	}

	// read mgroup
	var mgrpContent string
	if os.Getenv("MTEST") == "1" {
		mgrpContent = startMCastTestMode()
	} else {
		mgrpContent = readOptionalJsonFile(mgrpFilePath)
	}

	log.Println("initalized multicast controller")

	return &MulticastController{
		gmTable:      NewGroupMembersTable(mgrpContent),
		queryFlooder: queryFlooder,
		jrConn:       jrConn,
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
	return nil, false
}

func (c *MulticastController) Start(ft *MultiForwardTable) {
	go c.queryFlooder.ReceiveFloodedMsgs(c.onRecvJoinQuery)
	go c.recvJoinReplyMsgs(ft)
}

func (c *MulticastController) onRecvJoinQuery(fldHdr *FloodHeader, payload []byte) bool {
	// TODO: reply with join reply
	// TODO: store msg in cache
	jq, valid := UnmarshalJoinQuery(payload)
	if !valid {
		log.Panicln("Corrupted JoinQuery msg received")
	}
	log.Println(jq)
	// TODO: continue or stop flooding?
	return true
}

func (c *MulticastController) recvJoinReplyMsgs(ft *MultiForwardTable) {
	for {
		msg := c.jrConn.Read()

		jr, valid := UnmarshalJoinReply(msg)
		if !valid {
			log.Panicln("Corrupted JoinReply msg received")
		}
		log.Println(jr)
		// TODO: store msg
		// TODO: resend to next hop, unless im source
	}
}

func startMCastTestMode() string {
	log.Print("start in multicast test mode, assuming im working in rings.topo")
	address := "224.0.2.1"
	go startSendingMCastMsgs(7, address)

	// return groups members table json
	return "{\"" + address + "\": [\"10.0.0.20\", \"10.0.0.21\", \"10.0.0.22\"]}"
}

func startSendingMCastMsgs(secs int, address string) {
	log.Println("started sending multicast dgrams to ", address)

	raddr, err := net.ResolveUDPAddr("udp", address+":8080")
	if err != nil {
		log.Panic(err)
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	for {
		_, err := conn.Write([]byte("hello world message"))
		if err != nil {
			log.Panic(err)
		}

		time.Sleep(time.Duration(secs) * time.Second)
	}
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
