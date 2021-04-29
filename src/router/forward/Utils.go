package forward

import (
	"net"

	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

func imDestination(ip, destIP net.IP, destZoneID ZoneID) bool {
	// TODO: use destZID with the ip
	return destIP.Equal(ip) || destIP.IsLoopback()
}

func imInMulticastGrp(destGrpIP net.IP) bool {
	// TODO
	return false
}

func getUnicastNextHop(destIP net.IP, forwarder *Forwarder) (*UniForwardingEntry, bool) {
	// TODO: Get destination zone to check if we should look for its zone or ip in the forwarding table
	// Destination is a direct neighbor
	if ne, ok := forwarder.neighborsTable.Get(ToNodeID(destIP)); ok {
		return &UniForwardingEntry{NextHopMAC: ne.MAC, DestZoneID: uint32(MyZone().ID)}, true
	}
	forwarder.updateUnicastForwardingTable(forwarder.UniForwTable)
	if fe, ok := forwarder.UniForwTable.Get(ToNodeID(destIP)); ok {
		return fe, true
	}
	return nil, false
}
