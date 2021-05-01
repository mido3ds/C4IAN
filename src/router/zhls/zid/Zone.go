package zid

import (
	"fmt"
	"log"
)

const (
	earthTotalSurfaceAreaKM = 510.1e6
)

type gridLocation struct {
	x, y uint16
}

func (g gridLocation) String() string {
	return fmt.Sprintf("gridLocation{x:%v,y:%v}", g.x, g.y)
}

func (g *gridLocation) toGPSLocation() GPSLocation {
	return GPSLocation{Lat: indexToDegrees(g.y), Lon: indexToDegrees(g.x)}
}

// GPSLocation is gps position
// where Lat and Lon are in degrees
type GPSLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (p *GPSLocation) toGridPosition() gridLocation {
	return gridLocation{y: degreesToIndex(p.Lat), x: degreesToIndex(p.Lon)}
}

func (p GPSLocation) String() string {
	return fmt.Sprintf("GPSLocation{Lon:%v,Lat:%v}", p.Lon, p.Lat)
}

func indexToDegrees(i uint16) float64 {
	/*i*/                    // [0, 0xFFFF]
	d := float64(i) / 0xFFFF // [0, 1]
	d *= 2                   // [0, 2]
	d -= 1                   // [-1, 1]
	d *= 180                 // [-180, 180]
	return d
}

func degreesToIndex(d float64) uint16 {
	/*d*/       // [-180, 180]
	d /= 180    // [-1, 1]
	d += 1      // [0, 2]
	d /= 2      // [0, 1]
	d *= 0xFFFF // [0, 0xFFFF]
	return uint16(d)
}

func zlenMask(zlen byte) uint16 {
	return ^(0xFFFF >> zlen)
}

type ZoneID uint32

func NewZoneID(l GPSLocation, zlen byte) ZoneID {
	if zlen < 0 || zlen > 16 {
		log.Panic("zlen must be between 0 and 16")
	}

	// transform
	grid := l.toGridPosition()

	// mask
	mask := zlenMask(zlen)
	grid.x &= mask
	grid.y &= mask

	// pack
	return ZoneID(uint32(grid.x)<<16 | uint32(grid.y))
}

func (z ZoneID) toGridPosition() gridLocation {
	return gridLocation{x: uint16(z >> 16), y: uint16(z)}
}

func (z ZoneID) String() string {
	g := z.toGridPosition()
	return fmt.Sprintf("%04X.%04X", g.x, g.y)
}

type Zone struct {
	ID  ZoneID
	Len byte // from 0 to 16
}

// ToLen returns new ZoneID as this zone but with given len
func (z *Zone) ToLen(len byte) ZoneID {
	if len == z.Len {
		return z.ID
	}

	mask16 := zlenMask(len) & zlenMask(z.Len)
	mask := ZoneID(mask16)<<16 | ZoneID(mask16)

	return z.ID & mask
}

// Intersects returns true if z2 & z1 are the same zone
// or one part of the other
func (z1 *Zone) Intersects(z2 *Zone) bool {
	return z1.ToLen(z2.Len) == z2.ID
}

// Equal zones have same id and len
func (z1 *Zone) Equal(z2 *Zone) bool {
	return z1.ID == z2.ID && z1.Len == z2.Len
}

// Area returns pseudo areas for comparisons
func (z *Zone) Area() byte {
	return -z.Len
}

func (z Zone) String() string {
	g := z.ID.toGridPosition()
	return fmt.Sprintf("%04X.%04X/%d", g.x, g.y, z.Len)
}

// ZLenToAreaKMs returns area of the zone in kms^2
// of the given zlen
func ZLenToAreaKMs(zlen byte) float64 {
	var numZones uint32 = 0xFFFF_FFFF >> (32 - 2*zlen)
	var area float64 = earthTotalSurfaceAreaKM / float64(numZones)
	return area
}
