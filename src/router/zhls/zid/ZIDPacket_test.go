package zid

import (
	"testing"

	. "github.com/mido3ds/C4IAN/src/router/msec"
)

func TestBasicChecksum(t *testing.T) {
	tests := []struct {
		name string
		args []byte
		want uint16
	}{
		{"basic", []byte{1, 1, 1, 1, 1}, 5},
		{"basic2", []byte{1, 1, 1, 2, 1}, 6},
		{"zero", []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0},
		{"big", []byte{255, 255, 255, 255, 255}, 1275},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BasicChecksum(tt.args); got != tt.want {
				t.Errorf("BasicChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_zidChecksum(b *testing.B) {
	arr := []byte{255, 255, 255, 255, 255}
	for n := 0; n < b.N; n++ {
		BasicChecksum(arr)
	}
}

func BenchmarkMarshalBinary(b *testing.B) {
	header := &ZIDHeader{ZLen: 1, DstZID: 5, SrcZID: 12}
	for n := 0; n < b.N; n++ {
		header.MarshalBinary()
	}
}

func BenchmarkUnmarshalZIDHeader(b *testing.B) {
	packet := []byte{0, 218, 136, 65, 0, 0, 0, 5, 0, 0, 0, 12}
	for n := 0; n < b.N; n++ {
		UnmarshalZIDHeader(packet)
	}
}
