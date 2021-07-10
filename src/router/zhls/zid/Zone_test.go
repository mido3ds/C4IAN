package zid

import (
	"testing"
)

func Test_zlenMask(t *testing.T) {
	if zlenMask(0) != 0b0000000000000000 {
		t.Errorf("failed to mask %d, output: %#b", 0, zlenMask(0))
	}
	if zlenMask(1) != 0b1000000000000000 {
		t.Errorf("failed to mask %d, output: %#b", 1, zlenMask(1))
	}
	if zlenMask(2) != 0b1100000000000000 {
		t.Errorf("failed to mask %d, output: %#b", 2, zlenMask(2))
	}
	if zlenMask(3) != 0b1110000000000000 {
		t.Errorf("failed to mask %d, output: %#b", 3, zlenMask(3))
	}
	if zlenMask(4) != 0b1111000000000000 {
		t.Errorf("failed to mask %d, output: %#b", 4, zlenMask(4))
	}
	if zlenMask(5) != 0b1111100000000000 {
		t.Errorf("failed to mask %d, output: %#b", 5, zlenMask(5))
	}
	if zlenMask(6) != 0b1111110000000000 {
		t.Errorf("failed to mask %d, output: %#b", 6, zlenMask(6))
	}
	if zlenMask(7) != 0b1111111000000000 {
		t.Errorf("failed to mask %d, output: %#b", 7, zlenMask(7))
	}
	if zlenMask(8) != 0b1111111100000000 {
		t.Errorf("failed to mask %d, output: %#b", 8, zlenMask(8))
	}
	if zlenMask(9) != 0b1111111110000000 {
		t.Errorf("failed to mask %d, output: %#b", 9, zlenMask(9))
	}
	if zlenMask(10) != 0b1111111111000000 {
		t.Errorf("failed to mask %d, output: %#b", 10, zlenMask(10))
	}
	if zlenMask(11) != 0b1111111111100000 {
		t.Errorf("failed to mask %d, output: %#b", 11, zlenMask(11))
	}
	if zlenMask(12) != 0b1111111111110000 {
		t.Errorf("failed to mask %d, output: %#b", 12, zlenMask(12))
	}
	if zlenMask(13) != 0b1111111111111000 {
		t.Errorf("failed to mask %d, output: %#b", 13, zlenMask(13))
	}
	if zlenMask(14) != 0b1111111111111100 {
		t.Errorf("failed to mask %d, output: %#b", 14, zlenMask(14))
	}
	if zlenMask(15) != 0b1111111111111110 {
		t.Errorf("failed to mask %d, output: %#b", 15, zlenMask(15))
	}
	if zlenMask(16) != 0b1111111111111111 {
		t.Errorf("failed to mask %d, output: %#b", 16, zlenMask(16))
	}
}

func TestNewZoneID(t *testing.T) {
	l := gpsLocation{Lon: 31.20713052751232, Lat: 30.027653972850516}
	if uint32(newZoneID(l, 0)) != 0 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 0)
	}
	if uint32(newZoneID(l, 1)) != 2147516416 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 1)
	}
	if uint32(newZoneID(l, 2)) != 3221274624 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 2)
	}
	if uint32(newZoneID(l, 3)) != 3758153728 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 3)
	}
	if uint32(newZoneID(l, 4)) != 4026593280 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 4)
	}
	if uint32(newZoneID(l, 5)) != 4026593280 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 5)
	}
	if uint32(newZoneID(l, 6)) != 4026593280 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 6)
	}
	if uint32(newZoneID(l, 7)) != 4060148224 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 7)
	}
	if uint32(newZoneID(l, 8)) != 4060148224 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 8)
	}
	if uint32(newZoneID(l, 9)) != 4068536960 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 9)
	}
	if uint32(newZoneID(l, 10)) != 4072731264 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 10)
	}
	if uint32(newZoneID(l, 11)) != 4074828448 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 11)
	}
	if uint32(newZoneID(l, 12)) != 4075877024 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 12)
	}
	if uint32(newZoneID(l, 13)) != 4076401320 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 13)
	}
	if uint32(newZoneID(l, 14)) != 4076401320 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 14)
	}
	if uint32(newZoneID(l, 15)) != 4076532394 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 15)
	}
	if uint32(newZoneID(l, 16)) != 4076532394 {
		t.Errorf("failed to create zone id for l={%#v}, zlen={%d}", l, 16)
	}
}

func TestZoneIDToLen(t *testing.T) {
	var len byte = 16
	l := gpsLocation{Lon: 31.20713052751232, Lat: 30.027653972850516}
	zid := newZoneID(l, len)
	z := Zone{ID: zid, Len: len}
	for i := 0; i < 16; i++ {
		z1 := z.ID.ToLen(byte(i))
		z2 := newZoneID(l, byte(i))
		if z1 != z2 {
			t.Error("mismatch", z1, z2)
		}
		t.Log(z1, z2)
	}
}

func TestZoneID_DistTo(t *testing.T) {
	tests := []struct {
		name string
		z1   ZoneID
		z2   ZoneID
		want uint16
	}{
		{"1", gridLocation{1, 1}.toZoneID(), gridLocation{1, 1}.toZoneID(), 0},
		{"2", gridLocation{1, 1}.toZoneID(), gridLocation{1, 2}.toZoneID(), 1},
		{"3", gridLocation{2, 1}.toZoneID(), gridLocation{1, 1}.toZoneID(), 1},
		{"4", gridLocation{0, 0}.toZoneID(), gridLocation{50, 100}.toZoneID(), 150},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.z1.DistTo(tt.z2); got != tt.want {
				t.Errorf("ZoneID.DistTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
