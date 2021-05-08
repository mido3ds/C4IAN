package tables

import (
	"testing"
	"time"

	. "github.com/mido3ds/C4IAN/src/router/ip"
)

func BenchmarkBuffer(t *testing.B) {
	buffer := NewPacketsBuffer()

	dst1 := UInt32ToIPv4(1)
	dst2 := UInt32ToIPv4(2)

	packet1 := []byte{0, 0, 0, 0, 1}
	packet2 := []byte{0, 0, 0, 0, 0, 0, 0, 2}
	packet3 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 3}
	packet4 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4}

	buffer.AppendPacket(dst1, packet1)
	buffer.AppendPacket(dst1, packet2)
	buffer.AppendPacket(dst2, packet3)
	buffer.AppendPacket(dst2, packet4)

	print(buffer.String(), "\n")

	time.Sleep(6 * time.Second)

	print(buffer.String(), "\n")
}
