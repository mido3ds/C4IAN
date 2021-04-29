package zid

var myZlen uint8
var myZoneID ZoneID

func MyZone() *Zone {
	return &Zone{Len: myZlen, ID: myZoneID}
}

func MyZIDHeader(dstZID ZoneID) *ZIDHeader {
	return &ZIDHeader{ZLen: myZlen, SrcZID: myZoneID, DstZID: dstZID}
}

func OnZoneIDChanged(newZoneID ZoneID) {
	myZoneID = newZoneID
}
