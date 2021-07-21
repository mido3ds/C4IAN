package zid

import "sync"

var myZoneMutex sync.RWMutex
var myGpsLocationMutex sync.RWMutex

var myZone Zone
var myGpsLocation GpsLocation

func MyZone() Zone {
	myZoneMutex.RLock()
	defer myZoneMutex.RUnlock()
	return myZone
}

func MyLocation() GpsLocation {
	myGpsLocationMutex.RLock()
	defer myGpsLocationMutex.RUnlock()
	return myGpsLocation
}

func MyZIDHeader(dstZID ZoneID) *ZIDHeader {
	myZoneMutex.RLock()
	defer myZoneMutex.RUnlock()
	return &ZIDHeader{ZLen: myZone.Len, SrcZID: myZone.ID, DstZID: dstZID}
}
