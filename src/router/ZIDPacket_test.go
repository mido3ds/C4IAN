package main

import (
	"testing"
)

func Test_zidChecksum(t *testing.T) {
	tests := []struct {
		name string
		args []byte
		want int16
	}{
		{"basic", []byte{1, 1, 1, 1, 1}, 5},
		{"basic2", []byte{1, 1, 1, 2, 1}, 6},
		{"zero", []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0},
		{"big", []byte{255, 255, 255, 255, 255}, 1275},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := basicChecksum(tt.args); got != tt.want {
				t.Errorf("basicChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_zidChecksum(b *testing.B) {
	arr := []byte{255, 255, 255, 255, 255}
	for n := 0; n < b.N; n++ {
		basicChecksum(arr)
	}
}

func BenchmarkMarshalBinary(b *testing.B) {
	pm, err := NewZIDPacketMarshaler(1500)
	if err != nil {
		panic(err)
	}
	header := &ZIDHeader{ZLen: 1, DestZID: 5, SrcZID: 12}
	for n := 0; n < b.N; n++ {
		_, err := pm.MarshalBinary(header, pm.buffer[:1000])
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkUnpackZIDHeader(b *testing.B) {
	packet := []byte{0, 218, 136, 65, 0, 0, 0, 5, 0, 0, 0, 12}
	for n := 0; n < b.N; n++ {
		_, _, err := UnpackZIDHeader(packet)
		if err != nil {
			panic(err)
		}
	}
}
