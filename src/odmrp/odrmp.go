package odmrp

type Node struct {
	down                     bool
	is_ready                 bool
	ip_address               IP
	multicast_source_address IP
	multicast_group          []IP // addresses of the multicast group the node is part of
}
