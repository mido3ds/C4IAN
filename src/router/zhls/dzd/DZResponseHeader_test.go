package dzd

import (
	"fmt"
	"testing"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

func BenchmarkDZResponseHeader(t *testing.B) {
	dst := UInt32ToIPv4(5)
	requiredDst := UInt32ToIPv4(2)
	requiredDstZone := ZoneID(1)
	dzResponseHeader := &DZResponseHeader{dstIP: dst, requiredDstIP: requiredDst, requiredDstZone: requiredDstZone}

	packet := dzResponseHeader.MarshalBinary()

	newHeader, valid := UnmarshalDZResponseHeader(packet)

	fmt.Println("Valid: ", valid)
	fmt.Println("DstIP: ", newHeader.dstIP)
	fmt.Println("RequiredDstIP: ", newHeader.requiredDstIP)
	fmt.Println("RequiredDstZone: ", newHeader.requiredDstZone)
}
