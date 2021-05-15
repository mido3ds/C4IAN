package dzd

import (
	"fmt"
	"testing"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/tables"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

func BenchmarkDZDController(t *testing.B) {
	myIP := UInt32ToIPv4(15)
	requiredDstIP := UInt32ToIPv4(2)

	visitedZones := []ZoneID{ZoneID(4), ZoneID(5), ZoneID(6)}
	srcZone := Zone{ID: ZoneID(10), Len: 16}
	nextZone := ToNodeID(ZoneID(1))

	dzdController, _ := NewDZDController(myIP, nil, nil)
	dzRequestPacket := dzdController.createDZRequestPacket(myIP, srcZone, nextZone, requiredDstIP, visitedZones)
	zidHeader, dzRequestHeader := dzdController.unpackDZRequestPacket(dzRequestPacket)

	fmt.Println("DstZID: ", zidHeader.DstZID)
	fmt.Println("srcIP: ", dzRequestHeader.srcIP)
	fmt.Println("requiredDstIP: ", dzRequestHeader.requiredDstIP)
	fmt.Println("VisitedZones: ", dzRequestHeader.visitedZones)

	requiredDstZone := ZoneID(1)

	dstIP := myIP
	dstZone := ZoneID(10)

	dzResponsePacket := dzdController.createDZResponsePacket(requiredDstIP, requiredDstZone, dstIP, dstZone)
	ZIDHeader, dzResponseHeader := dzdController.unpackDZResponsePacket(dzResponsePacket)

	fmt.Println("----------------------------")
	fmt.Println("requiredDstIP: ", dzResponseHeader.requiredDstIP)
	fmt.Println("requiredDstZone: ", dzResponseHeader.requiredDstZone)
	fmt.Println("DSTIP: ", dzResponseHeader.dstIP)
	fmt.Println("DSTZone: ", ZIDHeader.DstZID)

}
