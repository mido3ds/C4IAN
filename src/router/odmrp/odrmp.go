package odmrp

import (
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

type ODMRPPacket struct {
	packet_type        byte
	source_ddr         string
	multicast_group_ip IP
	prev_hop_ip        IP
	sequence_number    int
}
