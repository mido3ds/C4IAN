package odmrp

import (
	"time"

	"github.com/mido3ds/C4IAN/src/utils"
)

type Node struct {
	down                     bool
	is_ready                 bool
	ip_address               IP
	multicast_source_address IP
	multicast_group          []IP // addresses of the multicast group the node is part of
	multicast_receivers      utils.Interface
}

type ForwardingTableEntry struct {
	groupID         string
	lastRefreshTime int64
}

type MessageCacheEntry struct {
	packet_id      int64
	source_address IP
}

type ODMRPAgent struct {
	// Trace *logtarget;
	// join_query_timers   JoinQueryTimer
	// join_reply_timers   JoinReplyTimer
	//TODO our link layer output (ll)
	//TODO output interface priority queue (ifq)
	//TODO mac layer (Mac_ *mac)
	mcast_target_          int
	ll                     int
	ifq                    int
	mac_                   int
	off_mac_               int
	off_ll_                int
	off_ip_                int
	off_odmrp_             int
	mcast_base_address_    int
	jq_seqno               int
	rrep_pending_ack       int
	cache                  Cache
	packet_buffer          PacketBuffer
	jr_packets             *JoinRequestPacket
	net_id                 int // ip address
	MAC_id                 int // mac address
	mcast_membership_table MemTable
	jq_table               JoinQueryTable
	jr_table               JoinReplyTable
}

type JoinReplyPacket struct {
	pkt      *Packet
	out_time Time
	grp      int
}

func newODMRPAgent() *ODMRPAgent {
	var new_agent ODMRPAgent
	new_agent.mcast_target_ = 0
	new_agent.ll = 0
	new_agent.ifq = 0
	new_agent.mac_ = 0
	new_agent.jq_seqno = 0
	new_agent.rrep_pending_ack = 0
	new_agent.jr_packets = make([]JoinReplyPacket, MAX_NUM_GROUPS)

	for i := 0; i < len(new_agent.jr_packets); i++ {
		new_agent.jr_packets[i].pkt = nil
		new_agent.jr_packets[i].out_time = 0
		new_agent.jr_packets[i].grp = 0
	}

	return &new_agent
}

func (ag *ODMRPAgent) GetJoinReplyPacket(mcast_grp int) *JoinReplyPacket {
	current_time := time.Now()

	for i := 0; i < len(ag.jr_packets); i++ {
		if ag.jr_packets[i].grp == mcast_grp && ag.jr_packets[i].out_time > Time(current_time.Second()) {
			return &ag.jr_packets[i]
		} else if ag.jr_packets[i].out_time < current_time {
			ag.jr_packets[i].pkt = 0
			if ag.jr_packets[i].grp == mcast_grp {
				return nil
			}
		}
	}

	return nil
}

func (ag *ODMRPAgent) AddJoinReplyPacket(pkt *Packet, mcast_grp int, out_time Time) {
	new_index := 0

	for i := 0; i < len(ag.jr_packets); i++ {
		if ag.jr_packets[i].grp == mcast_grp {
			ag.jr_packets[i].pkt = pkt
			ag.jr_packets[i].out_time = out_time
			return
		} else if mcast_grp == 0 {
			new_index = i
		}
	}
	ag.jr_packets[new_index].grp = mcast_grp
	ag.jr_packets[new_index].pkt = pkt
	ag.jr_packets[new_index].out_time = out_time
}

func (ag *ODMRPAgent) McastAddress(dest int) bool {
	return dest > ag.mcast_base_address_
}
