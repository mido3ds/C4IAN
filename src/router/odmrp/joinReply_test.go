package odmrp

import (
	"net"
	"testing"

	. "github.com/mido3ds/C4IAN/src/router/ip"
)

func TestJoinReplyMarshalAndUnmarshal(t *testing.T) {
	var jr joinReply
	ip0 := net.IP([]byte{0x01, 0x02, 0x03, 0x04})
	ip1 := net.IP([]byte{0x05, 0x06, 0x07, 0x08})
	ip2 := net.IP([]byte{0x09, 0x0A, 0x0B, 0x0C})
	ip3 := net.IP([]byte{0x0D, 0x0E, 0x0F, 0x10})
	mac1 := net.HardwareAddr([]byte{0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12})
	mac2 := net.HardwareAddr([]byte{0x0A, 0x0E, 0x0C, 0x10, 0x11, 0x12})
	mac3 := net.HardwareAddr([]byte{0x0F, 0xFF, 0x0F, 0x10, 0x11, 0x12})
	mac4 := net.HardwareAddr([]byte{0x0F, 0xFF, 0x0F, 0xAA, 0xCC, 0xC2})

	jr.seqNum = 215
	jr.destIP = ip3
	jr.grpIP = ip0
	jr.cost = 5
	jr.prevHop = mac4
	jr.srcIPs = []net.IP{ip1, ip2, ip3}
	jr.nextHops = []net.HardwareAddr{mac1, mac2, mac3}

	payload := jr.marshalBinary()
	newJr, ok := unmarshalJoinReply(payload)

	if !ok {
		t.Errorf("Unmarshal should return no erros")
	}

	if jr.seqNum != newJr.seqNum {
		t.Errorf("SeqNo are not equal")
	}

	if HwAddrToUInt64(jr.prevHop) != HwAddrToUInt64(newJr.prevHop) {
		t.Errorf("ips not equal: %#v != %#v", jr.prevHop.String(), newJr.prevHop.String())
	}

	if !net.IP.Equal(jr.grpIP, newJr.grpIP) {
		t.Errorf("ips not equal: %#v != %#v", jr.grpIP.String(), newJr.grpIP.String())
	}

	if len(jr.srcIPs) != len(newJr.srcIPs) {
		t.Errorf("SrcIPs length are not equal")
	}

	for i := 0; i < len(jr.srcIPs); i++ {
		if !net.IP.Equal(jr.srcIPs[i], newJr.srcIPs[i]) {
			t.Errorf("ips not equal: %#v != %#v", jr.srcIPs[i].String(), newJr.srcIPs[i].String())
		}
	}

	if len(jr.nextHops) != len(newJr.nextHops) {
		t.Errorf("NextHops length are not equal")
	}

	if jr.cost != newJr.cost {
		t.Errorf("Cost are not equal")
	}

	for i := 0; i < len(jr.nextHops); i++ {
		if HwAddrToUInt64(jr.nextHops[i]) != HwAddrToUInt64(newJr.nextHops[i]) {
			t.Errorf("ips not equal: %#v != %#v", jr.nextHops[i].String(), newJr.nextHops[i].String())
		}
	}
}
