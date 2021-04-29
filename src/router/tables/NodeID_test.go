package tables

import (
	"fmt"
	"net"
	"testing"

	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

func BenchmarkNewNodeID(t *testing.B) {
	var a net.IP = net.IPv4(0, 0, 0, 55).To4()
	var b ZoneID = 55
	S := ToNodeID(a)
	S2 := ToNodeID(b)
	fmt.Println(S)
	fmt.Println(S2)
}
