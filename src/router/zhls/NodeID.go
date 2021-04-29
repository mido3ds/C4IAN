package zhls

import (
	"fmt"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type NodeID uint64

func ToNodeID(ID interface{}) (nodeID NodeID) {
	nodeID = 0
	switch id := ID.(type) {
	case net.IP:
		nodeID |= NodeID(IPv4ToUInt32(id))
		nodeID |= 1 << 63 // Set MSB to 1 to avoid collisions with node ids
	case ZoneID:
		nodeID |= NodeID(id)
	default:
		log.Panic("Invalid type to NodeId")
	}
	return
}

func (nodeID NodeID) String() string {
	s := "NodeID"
	if nodeID>>63 == 1 {
		s += (" (IP): ")
		s += UInt32ToIPv4(uint32(nodeID)).String()
	} else {
		s += (" (Zone ID): ")
		s += fmt.Sprint(uint32(nodeID))
	}
	return s
}
