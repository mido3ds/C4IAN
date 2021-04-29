package zid

import (
	"fmt"
	"log"
	"math"
)

const (
	earthRadiusKM           = 6371e3
	earthTotalSurfaceAreaKM = 510.1e6
)

type gridLocation struct {
	x, y uint16
}

// GPSLocation is gps position
// where Lat and Lon are in degrees
type GPSLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (p *GPSLocation) toGridPosition() gridLocation {
	return gridLocation{y: degreesToCartesian(p.Lat), x: degreesToCartesian(p.Lon)}
}

func degreesToCartesian(d float64) uint16 {
	if d < 0 {
		d += 360
	}
	return uint16((math.Pi / 180) * d * earthRadiusKM)
}

type ZoneID uint32

func NewZoneID(l GPSLocation, zlen byte) ZoneID {
	if zlen < 0 || zlen > 16 {
		log.Panic("zlen must be between 0 and 16")
	}

	// transform
	grid := l.toGridPosition()

	// shift
	shifts := uint16(16 - zlen)
	grid.x >>= shifts
	grid.y >>= shifts

	// pack
	return ZoneID(uint32(grid.x)<<16 | uint32(grid.y))
}

func (z ZoneID) toGridPosition() gridLocation {
	return gridLocation{x: uint16(z >> 16), y: uint16(z)}
}

type Zone struct {
	ID  ZoneID
	Len byte // from 0 to 16
}

// ToLen returns new ZoneID as this zone but with given len
func (z *Zone) ToLen(len byte) ZoneID {
	if len < 0 || len > 16 {
		log.Panic("zlen must be between 0 and 16")
	}

	if len == z.Len {
		return z.ID
	}

	grid := z.ID.toGridPosition()

	if len < z.Len {
		// z is smaller
		// enlarge z, will lose details

		// ID >> (z.Len - len)
		grid.x >>= z.Len - len
		grid.y >>= z.Len - len
	} else {
		// z is bigger
		// reduce z, will get arbitrary smaller zone,
		// but its part of the original nevertheless

		// ID << (len - z.Len)
		grid.x <<= len - z.Len
		grid.y <<= len - z.Len
	}

	return ZoneID(uint32(grid.x)<<16 | uint32(grid.y))
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
	return fmt.Sprintf("%X.%X/%d", g.x, g.y, z.Len)
}

// ZLenToAreaKMs returns area of the zone in kms^2
// of the given zlen
func ZLenToAreaKMs(zlen byte) float64 {
	var numZones uint32 = 0xFFFF_FFFF >> (32 - 2*zlen)
	var area float64 = earthTotalSurfaceAreaKM / float64(numZones)
	return area
}
