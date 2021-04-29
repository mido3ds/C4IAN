package zid

import (
	"log"
	"math"
)

const earthRadiusKM = 6371 * 1000

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

	x := uint16(z.ID >> 16)
	y := uint16(z.ID)

	if len < z.Len {
		// z is smaller
		// enlarge z, will lose details

		// ID >> (z.Len - len)
		x >>= z.Len - len
		y >>= z.Len - len
	} else {
		// z is bigger
		// reduce z, will get arbitrary smaller zone,
		// but its part of the original nevertheless

		// ID << (len - z.Len)
		x <<= len - z.Len
		y <<= len - z.Len
	}

	return ZoneID(uint32(x)<<16 | uint32(y))
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
