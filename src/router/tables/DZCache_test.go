package tables

import (
	"testing"
	"time"

	. "github.com/mido3ds/C4IAN/src/router/ip"
	. "github.com/mido3ds/C4IAN/src/router/zhls/zid"
)

func BenchmarkDZCache(t *testing.B) {
	zoneCache := NewDZCache()

	dst1 := UInt32ToIPv4(1)

	zoneCache.Set(dst1, MyZone().ID)

	print(zoneCache.String(), "\n")

	time.Sleep(6 * time.Second)

	print(zoneCache.String(), "\n")
}
