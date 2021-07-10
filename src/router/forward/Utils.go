package forward

import (
	"net"

	. "github.com/mido3ds/C4IAN/src/router/tables"
)

func imDestination(ip, destIP net.IP) bool {
	return destIP.Equal(ip) || destIP.IsLoopback()
}

func (forwarder *Forwarder) GetUnicastNextHop(dst NodeID) (net.HardwareAddr, bool) {
	// Destination is a direct neighbor, send directly to it
	ne, ok := forwarder.neighborsTable.Get(dst)
	if ok {
		return ne.MAC, true
	}
	//log.Println("Get nextHop of", dst)
	// Otherwise look for the n, ext hop in the forwarding table
	forwarder.updateUnicastForwardingTable(forwarder.UniForwTable)
	//log.Println(forwarder.UniForwTable)
	nextHopMAC, ok := forwarder.UniForwTable.Get(dst)
	if ok {
		return nextHopMAC, true
	}
	return nil, false
}
