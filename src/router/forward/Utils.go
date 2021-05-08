package forward

import (
	"net"

	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

func imDestination(ip, destIP net.IP) bool {
	return destIP.Equal(ip) || destIP.IsLoopback()
}

func imInMulticastGrp(destGrpIP net.IP) bool {
	// TODO
	return false
}

func (forwarder *Forwarder) getUnicastNextHop(destIP net.IP, destZID ZoneID) (*UniForwardingEntry, bool) {
	// Destination is a direct neighbor, send directly to it
	if ne, ok := forwarder.neighborsTable.Get(ToNodeID(destIP)); ok {
		return &UniForwardingEntry{NextHopMAC: ne.MAC, DestZoneID: uint32(MyZone().ID)}, true
	}

	// Otherwise look for the next hop in the forwarding table
	forwarder.updateUnicastForwardingTable(forwarder.UniForwTable)

	var fe *UniForwardingEntry
	var ok bool
	if destZID == MyZone().ID {
		// The destination is in my zone, search in the forwarding table by its ip
		// TODO: If the IP is not found in the forwarding table then the dest may have moved out of this zone, discover its new zone
		fe, ok = forwarder.UniForwTable.Get(ToNodeID(destIP))
	} else {
		// The destination is in a different zone, search in the forwarding table by its zone
		fe, ok = forwarder.UniForwTable.Get(ToNodeID(destZID))
	}
	if ok {
		return fe, true
	}
	return nil, false
}
