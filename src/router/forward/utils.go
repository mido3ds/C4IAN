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

func getNextHop(destIP net.IP, ft *UniForwardTable, nt *NeighborsTable, zoneID ZoneID) (*UniForwardingEntry, bool) {
	fe, ok := ft.Get(destIP)
	if !ok {
		ne, ok := nt.Get(destIP)
		if !ok {
			return nil, false
		}
		return &UniForwardingEntry{NextHopMAC: ne.MAC, DestZoneID: uint32(zoneID)}, true
	}
	return fe, true
}
