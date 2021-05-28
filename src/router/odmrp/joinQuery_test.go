package odmrp

import (
	"bytes"
	"net"
	"testing"
)

func TestJoinQueryMarshalAndUnmarshal(t *testing.T) {
	var jq joinQuery
	ip0 := net.IP([]byte{0x01, 0x02, 0x03, 0x04})
	ip1 := net.IP([]byte{0x05, 0x06, 0x07, 0x08})
	ip2 := net.IP([]byte{0x09, 0x0A, 0x0B, 0x0C})
	ip3 := net.IP([]byte{0x0D, 0x0E, 0x0F, 0x10})
	jq.seqNum = 999
	jq.ttl = 20
	jq.srcIP = net.IP([]byte{0x09, 0x0A, 0x0B, 0x0C})
	jq.grpIP = net.IP([]byte{0x0D, 0x0E, 0x0F, 0x10})
	jq.dests = []net.IP{ip0, ip1, ip2, ip3}
	prevhop, err := net.ParseMAC("00:26:bb:15:31:dd")
	if err != nil {
		t.Error(err)
	}
	jq.prevHop = prevhop

	payload := jq.marshalBinary()
	newJq, ok := unmarshalJoinQuery(payload)

	if !ok {
		t.Errorf("Unmarshal should return no erros")
	}
	if len(jq.dests) != len(newJq.dests) {
		t.Errorf("Dests ips length are not equal")
	}
	if jq.seqNum != newJq.seqNum {
		t.Errorf("SeqNo are not equal")
	}
	if jq.ttl != newJq.ttl {
		t.Errorf("TTL are not equal")
	}
	if !net.IP.Equal(jq.srcIP, newJq.srcIP) {
		t.Errorf("src ips not equal: %#v != %#v", jq.srcIP.String(), jq.grpIP.String())
	}
	if !net.IP.Equal(jq.grpIP, newJq.grpIP) {
		t.Errorf("grp ips not equal: %#v != %#v", jq.srcIP.String(), jq.grpIP.String())
	}
	if !bytes.Equal(jq.prevHop, newJq.prevHop) {
		t.Errorf("prev hops not equal: %#v != %#v", jq.prevHop, newJq.prevHop)
	}
	for i := 0; i < len(jq.dests); i++ {
		if !net.IP.Equal(jq.dests[i], newJq.dests[i]) {
			t.Errorf("ips not equal: %#v != %#v", jq.dests[i].String(), newJq.dests[i].String())
		}
	}
}
