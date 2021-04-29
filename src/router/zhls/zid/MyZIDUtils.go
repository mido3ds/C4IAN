package zid

// TODO: Is it okay to use global variables (package-level) here?
var myZlen uint8
var myZoneID ZoneID

func MyZone() *Zone {
	return &Zone{Len: myZlen, ID: myZoneID}
}

func MyZIDHeader(dstZID ZoneID) *ZIDHeader {
	return &ZIDHeader{ZLen: myZlen, SrcZID: myZoneID, DstZID: dstZID}
}
