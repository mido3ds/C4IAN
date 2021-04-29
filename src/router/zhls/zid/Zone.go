package zid

import (
	"log"
	"math"
)

const earthRadiusKM = 6371 * 1000

type gridLocation struct {
	norths, easts uint16
}

// GPSLocation is gps position
// where Lat and Lon are in degrees
type GPSLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (p *GPSLocation) toGridPosition() gridLocation {
	return gridLocation{norths: degreesToCartesian(p.Lat), easts: degreesToCartesian(p.Lon)}
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
	grid.easts >>= shifts
	grid.norths >>= shifts

	// pack
	return ZoneID(uint32(grid.easts)<<16 | uint32(grid.norths))
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
	} else if len < z.Len {
		// z is smaller
		// enlarge z, will lose details
		return z.ID >> (z.Len - len)
	}

	// z is bigger
	// reduce z, will get arbitrary smaller zone,
	// but its part of the original nevertheless
	return z.ID << (len - z.Len)
}

// Intersects returns true if z2 & z1 are the same zone
// or one part of the other
func (z1 *Zone) Intersects(z2 *Zone) bool {
	return z1.ToLen(z2.Len) == z2.ID
}

// Area returns pseudo areas for comparisons
func (z *Zone) Area() byte {
	return -z.Len
}
