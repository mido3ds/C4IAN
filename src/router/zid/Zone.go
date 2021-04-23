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

// largerZone returns ZoneID with smaller ZLen by a difference of `diff`
// this means you will get a zone with a bigger area
func (z ZoneID) largerZone(diff byte) ZoneID {
	if diff < 0 || diff > 16 {
		log.Panic("diff must be between 0 and 16")
	}

	easts := uint16(z>>16) >> diff
	norths := uint16(z) >> diff

	// pack
	return ZoneID(uint32(easts)<<16 | uint32(norths))
}

type Zone struct {
	ID  ZoneID
	Len byte // from 0 to 16
}

// InZone returns true if either z or z2 is inside the other
// the other doesn't matter
func (z1 *Zone) InZone(z2 *Zone) bool {
	diff := z2.Len - z1.Len

	if diff > 0 { // z2 is smaller in area
		return z1.ID == z2.ID.largerZone(diff)
	} else if diff < 0 { // z1 is smaller in area
		return z1.ID.largerZone(-diff) == z2.ID
	}

	return z1.ID == z2.ID
}
