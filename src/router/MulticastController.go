package main

import (
	"log"
	"net"

	"github.com/mdlayher/ethernet"
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

func NewMulticastController(router *Router, mgroupContent string) (*MulticastController, error) {
	queryFlooder, err := NewGlobalFlooder(router.ip, router.iface, joinQueryEtherType, router.msec)
	if err != nil {
		log.Panic("failed to initiate query flooder, err: ", err)
	}

	jrConn, err := NewMACLayerConn(router.iface, joinReplyEtherType)
	if err != nil {
		log.Panic("failed to initiate mac conn, err: ", err)
	}

	log.Println("initalized multicast controller")

	return &MulticastController{
		gmTable:      NewGroupMembersTable(mgroupContent),
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
