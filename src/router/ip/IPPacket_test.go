package ip

import (
	"net"
	"testing"
)

func Benchmark_ipv4Checksum(b *testing.B) {
	a := []byte{0x45, 0x00, 0x00, 0x73, 0x00, 0x00, 0x40, 0x00, 0x40, 0x11, 0x00, 0x00, 0xc0, 0xa8, 0x00, 0x01, 0xc0, 0xa8, 0x00, 0xc7}
	for n := 0; n < b.N; n++ {
		ipv4Checksum(a)
	}
}

func TestIsBroadcastIP(t *testing.T) {
	tests := []struct {
		name string
		args net.IP
		want bool
	}{
		{"", net.ParseIP("255.255.0.14"), true},
		{"", net.ParseIP("255.255.0.1"), true},
		{"", net.ParseIP("255.255.0.0"), true},
		{"", net.ParseIP("255.255.255.255"), true},
		{"", net.ParseIP("255.0.255.255"), false},
		{"", net.ParseIP("0.0.255.255"), false},
		{"", net.ParseIP("2.0.0.255"), false},
		{"", net.ParseIP("0.0.0.0"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBroadcast(tt.args); got != tt.want {
				t.Errorf("IsBroadcastIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBroadcastRadius(t *testing.T) {
	tests := []struct {
		name string
		args net.IP
		want uint16
	}{
		{"", net.ParseIP("255.255.0.0"), 0},
		{"", net.ParseIP("255.255.0.1"), 1},
		{"", net.ParseIP("255.255.1.0"), 256},
		{"", net.ParseIP("255.255.1.1"), 257},
		{"", net.ParseIP("255.255.255.255"), 65535},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BroadcastRadius(tt.args); got != tt.want {
				t.Errorf("BroadcastRadius() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecrementBroadcastRadius(t *testing.T) {
	tests := []struct {
		name string
		args []byte
		want bool
	}{
		{"", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false},
		{"", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecrementBroadcastRadius(tt.args); got != tt.want {
				t.Errorf("DecrementBroadcastRadius() = %v, want %v", got, tt.want)
			}
		})
	}

	b := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	DecrementBroadcastRadius(b)
	if b[19] != 0 {
		t.Error("failed to decrement radius")
	}

	b = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1}
	DecrementBroadcastRadius(b)
	if b[19] != 0 || b[18] != 1 {
		t.Error("failed to decrement radius")
	}
}
