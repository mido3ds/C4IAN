package odmrp

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
	mac1 := net.HardwareAddr([]byte{0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12})
	mac2 := net.HardwareAddr([]byte{0x0A, 0x0E, 0x0C, 0x10, 0x11, 0x12})
	mac3 := net.HardwareAddr([]byte{0x0F, 0xFF, 0x0F, 0x10, 0x11, 0x12})

	jr.SeqNo = 215
	jr.GrpIP = ip0
	jr.SrcIPs = []net.IP{ip1, ip2, ip3}
	jr.NextHops = []net.HardwareAddr{mac1, mac2, mac3}

	payload := jr.MarshalBinary()
	newJr, ok := UnmarshalJoinReply(payload)

	if !ok {
		t.Errorf("Unmarshal should return no erros")
	}

	if jr.SeqNo != newJr.SeqNo {
		t.Errorf("SeqNo are not equal")
	}

	if !net.IP.Equal(jr.GrpIP, newJr.GrpIP) {
		t.Errorf("ips not equal: %#v != %#v", jr.GrpIP.String(), newJr.GrpIP.String())
	}

	if len(jr.SrcIPs) != len(newJr.SrcIPs) {
		t.Errorf("SrcIPs length are not equal")
	}

	for i := 0; i < len(jr.SrcIPs); i++ {
		if !net.IP.Equal(jr.SrcIPs[i], newJr.SrcIPs[i]) {
			t.Errorf("ips not equal: %#v != %#v", jr.SrcIPs[i].String(), newJr.SrcIPs[i].String())
		}
	}

	if len(jr.NextHops) != len(newJr.NextHops) {
		t.Errorf("NextHops length are not equal")
	}

	for i := 0; i < len(jr.NextHops); i++ {
		if hwAddrToUInt64(jr.NextHops[i]) != hwAddrToUInt64(newJr.NextHops[i]) {
			t.Errorf("ips not equal: %#v != %#v", jr.NextHops[i].String(), newJr.NextHops[i].String())
		}
	}
}
