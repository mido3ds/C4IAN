package dzd

import (
	"fmt"
	"testing"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

func BenchmarkDZDHeader(t *testing.B) {
	src := UInt32ToIPv4(1)
	dst := UInt32ToIPv4(2)

	visitedZones := []ZoneID{ZoneID(1), ZoneID(2), ZoneID(3)}
	dzdHeader := &DZDHeader{srcIP: src, dstIP: dst, visitedZones: visitedZones}

	packet := dzdHeader.MarshalBinary()
	newHeader, valid := UnmarshalDZDPHeader(packet)
	visitedZones = newHeader.visitedZones

	fmt.Println("Valid: ", valid)
	fmt.Println("SrcIP: ", newHeader.srcIP)
	fmt.Println("DstIP: ", newHeader.dstIP)
	fmt.Println("Num of Visited Zones: ", len(visitedZones))
	fmt.Println("1st Zone: ", visitedZones[0])
	fmt.Println("2nd Zone: ", visitedZones[1])
	fmt.Println("3rd Zone: ", visitedZones[2])

}
