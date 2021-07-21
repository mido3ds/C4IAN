package zid

import (
	"fmt"
	"log"
	"math"
)

const (
	earthTotalSurfaceAreaKM = 510.1e6
	meridianLengthKM        = 40_007.86
	equatorLengthKM         = 40_075.0
)

type gridLocation struct {
	x, y uint16
}

func (g gridLocation) String() string {
	return fmt.Sprintf("gridLocation{x:%v,y:%v}", g.x, g.y)
}

func (g gridLocation) toGPSLocation() GpsLocation {
	return GpsLocation{Lat: indexToDegrees(g.y), Lon: indexToDegrees(g.x)}
}

func (g gridLocation) toZoneID() ZoneID {
	return ZoneID(uint32(g.x)<<16 | uint32(g.y))
}

// gpsLocation is gps position
// where Lat and Lon are in degrees
type GpsLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (p GpsLocation) toGridLocation() gridLocation {
	return gridLocation{y: degreesToIndex(p.Lat), x: degreesToIndex(p.Lon)}
}

func (p GpsLocation) String() string {
	return fmt.Sprintf("GPSLocation{Lon:%v,Lat:%v}", p.Lon, p.Lat)
}

// [0, 2] -> [-1, 1]
func unwrap(f float64) float64 {
	if f > 1 {
		return f - 2
	}
	return f
}

func indexToDegrees(i uint16) float64 {
	/*i*/                    // [0, 0xFFFF]
	d := float64(i) / 0xFFFF // [0, 1]
	d *= 2                   // [0, 2]
	d = unwrap(d)            // [-1, 1]
	d *= 180                 // [-180, 180]
	return d
}

// [-1, 1] -> [0, 2]
func wrap(f float64) float64 {
	if f < 0 {
		return f + 2
	}
	return f
}

func degreesToIndex(d float64) uint16 {
	/*d*/       // [-180, 180]
	d /= 180    // [-1, 1]
	d = wrap(d) // [0, 2]
	d /= 2      // [0, 1]
	d *= 0xFFFF // [0, 0xFFFF]
	return uint16(d)
}

func zlenMask(zlen byte) uint16 {
	return ^(0xFFFF >> zlen)
}

type ZoneID uint32

func newZoneID(l GpsLocation, zlen byte) ZoneID {
	if zlen < 0 || zlen > 16 {
		log.Panic("zlen must be between 0 and 16")
	}

	// transform
	grid := l.toGridLocation()

	// mask
	mask := zlenMask(zlen)
	grid.x &= mask
	grid.y &= mask

	// pack
	return grid.toZoneID()
}

func (z ZoneID) toGridLocation() gridLocation {
	return gridLocation{x: uint16(z >> 16), y: uint16(z)}
}

func abs(i int32) int32 {
	if i >= 0 {
		return i
	}
	return -i
}

func absub(i, j uint16) uint16 {
	return uint16(abs(int32(i) - int32(j)))
}

// DistTo manhatten distance between 2 zones
func (z1 ZoneID) DistTo(z2 ZoneID) uint16 {
	g1 := z1.toGridLocation()
	g2 := z2.toGridLocation()
	return absub(g1.x, g2.x) + absub(g1.y, g2.y)
}

func (z ZoneID) String() string {
	g := z.toGridLocation()
	return fmt.Sprintf("%04X.%04X", g.x, g.y)
}

type Zone struct {
	ID  ZoneID
	Len byte // from 0 to 16
}

// ToLen returns new ZoneID as this zone but with given Len
func (z ZoneID) ToLen(len byte) ZoneID {
	mask16 := zlenMask(len)
	mask := ZoneID(mask16)<<16 | ZoneID(mask16)
	return z & mask
}

// Intersects casts z2 to z1 area
// and returns true if z2 & z1 are the same zone
// or one is part of the other
func (z1 Zone) Intersection(z2 Zone) (ZoneID, bool) {
	z3 := z2.ID.ToLen(z1.Len)
	return z3, z3 == z2.ID
}

// Equal zones have same id and len
func (z1 Zone) Equal(z2 *Zone) bool {
	return z1.ID == z2.ID && z1.Len == z2.Len
}

func (z Zone) String() string {
	return fmt.Sprintf("%v/%d", z.ID, z.Len)
}

func round3(x float64) float64 {
	return math.Round(x*1000) / 1000
}

// ZLenToAreaKMs returns area of the zone in kms^2
// of the given zlen
func ZLenToAreaKMs(zlen byte) float64 {
	precision := 1 / float64(int(1<<zlen))
	xl := equatorLengthKM * precision
	yl := meridianLengthKM * precision
	return round3(xl * yl)
}
