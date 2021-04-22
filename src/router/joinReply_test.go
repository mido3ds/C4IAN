package main

import (
	"net"
	"testing"
)

func TestJoinReplyMarshalAndUnmarshal(t *testing.T) {
	var jr JoinReply
	ip0 := net.IP([]byte{0x01, 0x02, 0x03, 0x04})
	ip1 := net.IP([]byte{0x05, 0x06, 0x07, 0x08})
	ip2 := net.IP([]byte{0x09, 0x0A, 0x0B, 0x0C})
	ip3 := net.IP([]byte{0x0D, 0x0E, 0x0F, 0x10})
	jr.SeqNo = 999
	jr.SrcIPs = []net.IP{ip0, ip1}
	jr.GrpIPs = []net.IP{ip2}
	jr.Forwarders = []net.IP{ip1, ip3, ip0}

	payload := jr.MarshalBinary()
	newJr, ok := UnmarshalJoinReply(payload)

	if !ok {
		t.Errorf("Unmarshal should return no erros")
	}
	if jr.SeqNo != newJr.SeqNo {
		t.Errorf("SeqNo are not equal")
	}
	if len(jr.SrcIPs) != len(newJr.SrcIPs) {
		t.Errorf("SrcIPs length are not equal")
	}
	for i := 0; i < len(jr.SrcIPs); i++ {
		if !net.IP.Equal(jr.SrcIPs[i], newJr.SrcIPs[i]) {
			t.Errorf("ips not equal: %#v != %#v", jr.SrcIPs[i].String(), newJr.SrcIPs[i].String())
		}
	}
	if len(jr.GrpIPs) != len(newJr.GrpIPs) {
		t.Errorf("GrpIPs length are not equal")
	}
	for i := 0; i < len(jr.GrpIPs); i++ {
		if !net.IP.Equal(jr.GrpIPs[i], newJr.GrpIPs[i]) {
			t.Errorf("ips not equal: %#v != %#v", jr.GrpIPs[i].String(), newJr.GrpIPs[i].String())
		}
	}
	if len(jr.Forwarders) != len(newJr.Forwarders) {
		t.Errorf("Forwarders length are not equal")
	}
	for i := 0; i < len(jr.Forwarders); i++ {
		if !net.IP.Equal(jr.Forwarders[i], newJr.Forwarders[i]) {
			t.Errorf("ips not equal: %#v != %#v", jr.Forwarders[i].String(), newJr.Forwarders[i].String())
		}
	}
}
