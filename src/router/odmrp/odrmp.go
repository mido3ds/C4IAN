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
	multicastReceivers       utils.Interface
}
