package main

import "testing"

func Test_getNumTSToSend(t *testing.T) {
	m3um8 := []byte(`#EXTM3U
	#EXT-X-VERSION:3
	#EXT-X-TARGETDURATION:2
	#EXT-X-MEDIA-SEQUENCE:0
	#EXT-X-PLAYLIST-TYPE:EVENT
	#EXTINF:2.000000,
	index0.ts
	#EXTINF:2.000000,
	index1.ts
	#EXTINF:2.000000,
	index2.ts
	#EXTINF:2.000000,
	index3.ts
	#EXTINF:2.100000,
	index4.ts
	#EXT-X-ENDLIST`)
	if getNumTS(m3um8)-0 != 4+1 {
		t.Error("failed")
	}
	if getNumTS(m3um8)-1 != 4+0 {
		t.Error("failed")
	}
	if getNumTS(m3um8)-2 != 4+-1 {
		t.Error("failed")
	}
	if getNumTS(m3um8)-3 != 4+-2 {
		t.Error("failed")
	}
	if getNumTS(m3um8)-4 != 4+-3 {
		t.Error("failed")
	}
	if getNumTS(m3um8)-5 != 4+-4 {
		t.Error("failed")
	}
}
