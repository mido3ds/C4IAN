package zid

import "sync"

var myZoneMutex sync.RWMutex
var myZone Zone

func MyZone() Zone {
	myZoneMutex.RLock()
	defer myZoneMutex.RUnlock()
	return myZone
}

func MyLocation() GpsLocation {
	myZoneMutex.RLock()
	defer myZoneMutex.RUnlock()
	return MyZone().ID.toGridLocation().toGPSLocation()
}

func MyZIDHeader(dstZID ZoneID) *ZIDHeader {
	myZoneMutex.RLock()
	defer myZoneMutex.RUnlock()
	return &ZIDHeader{ZLen: myZone.Len, SrcZID: myZone.ID, DstZID: dstZID}
}
