package main

import (
	"net"
	"testing"
)

func TestJoinQueryMarshalAndUnmarshal(t *testing.T) {
	var jq JoinQuery
	ip0 := net.IP([]byte{0x01, 0x02, 0x03, 0x04})
	ip1 := net.IP([]byte{0x05, 0x06, 0x07, 0x08})
	ip2 := net.IP([]byte{0x09, 0x0A, 0x0B, 0x0C})
	ip3 := net.IP([]byte{0x0D, 0x0E, 0x0F, 0x10})
	jq.SeqNo = 999
	jq.TTL = 20
	jq.SrcIP = net.IP([]byte{0x09, 0x0A, 0x0B, 0x0C})
	jq.GrpIP = net.IP([]byte{0x0D, 0x0E, 0x0F, 0x10})
	jq.Dests = []net.IP{ip0, ip1, ip2, ip3}

	payload := jq.MarshalBinary()
	newJq, ok := UnmarshalJoinQuery(payload)

	if !ok {
		t.Errorf("Unmarshal should return no erros")
	}
	if len(jq.Dests) != len(newJq.Dests) {
		t.Errorf("Dests ips length are not equal")
	}
	if jq.SeqNo != newJq.SeqNo {
		t.Errorf("SeqNo are not equal")
	}
	if jq.TTL != newJq.TTL {
		t.Errorf("TTL are not equal")
	}
	if !net.IP.Equal(jq.SrcIP, newJq.SrcIP) {
		t.Errorf("src ips not equal: %#v != %#v", jq.SrcIP.String(), jq.GrpIP.String())
	}
	if !net.IP.Equal(jq.GrpIP, newJq.GrpIP) {
		t.Errorf("grp ips not equal: %#v != %#v", jq.SrcIP.String(), jq.GrpIP.String())
	}
	for i := 0; i < len(jq.Dests); i++ {
		if !net.IP.Equal(jq.Dests[i], newJq.Dests[i]) {
			t.Errorf("ips not equal: %#v != %#v", jq.Dests[i].String(), newJq.Dests[i].String())
		}
	}
}

func TestJoinQueryTtlLessThanZero(t *testing.T) {
	var jq JoinQuery
	ip0 := net.IP([]byte{0x01, 0x02, 0x03, 0x04})
	ip1 := net.IP([]byte{0x05, 0x06, 0x07, 0x08})
	ip2 := net.IP([]byte{0x09, 0x0A, 0x0B, 0x0C})
	ip3 := net.IP([]byte{0x0D, 0x0E, 0x0F, 0x10})
	jq.SeqNo = 999
	jq.TTL = -1
	jq.SrcIP = net.IP([]byte{0x09, 0x0A, 0x0B, 0x0C})
	jq.GrpIP = net.IP([]byte{0x0D, 0x0E, 0x0F, 0x10})
	jq.Dests = []net.IP{ip0, ip1, ip2, ip3}

	payload := jq.MarshalBinary()
	newJq, ok := UnmarshalJoinQuery(payload)

	if !(!ok && newJq == nil) {
		t.Errorf("Unmarshal should have erros")
	}
}
