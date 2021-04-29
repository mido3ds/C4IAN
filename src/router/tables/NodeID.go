package tables

import (
	"fmt"
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

type NodeID uint64

func ToNodeID(Id interface{}) (nodeId NodeID) {
	switch Id.(type) {
	case net.IP:
		nodeId = NodeID("0" + fmt.Sprint(IPv4ToUInt32(Id.(net.IP))))
	case ZoneID:
		nodeId = NodeID("1" + fmt.Sprint(Id))
	default:
		log.Panic("Invalid type to NodeId")
	}
	return
}
