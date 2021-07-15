package halapi

import (
	"encoding/gob"
	"fmt"
)

// VideoFragment is sent to unit periodically
// as configured in StartVideoStream sent by unit
type VideoFragment struct {
	Video []byte
}

func (s VideoFragment) Send(enc *gob.Encoder) error {
	err := enc.Encode(byte(VideoFragmentType))
	if err != nil {
		return err
	}
	return enc.Encode(s.Video)
}

// HeartBeat is sent to unit periodically
// after connection establishment between HAL and unit
// sent each 1 second
type HeartBeat struct {
	// from 0 (dead or no sensor) to 100
	BeatsPerMinut int
}

func (s HeartBeat) Send(enc *gob.Encoder) error {
	err := enc.Encode(byte(HeartBeatType))
	if err != nil {
		return err
	}
	return enc.Encode(s.BeatsPerMinut)
}

// Location is a gps coordinates
// which are sent to unit periodically
// after connection establishment between HAL and unit
// sent each 1 second
type Location struct {
	Lon, Lat float64
}

func (s Location) Send(enc *gob.Encoder) error {
	err := enc.Encode(byte(LocationType))
	if err != nil {
		return err
	}
	return enc.Encode(s)
}

func RecvFromHAL(dec *gob.Decoder, vp *VideoFragment, hb *HeartBeat, loc *Location) (Type, error) {
	var sentType Type
	err := dec.Decode(&sentType)
	if err != nil {
		return 0, err
	}

	switch sentType {
	case VideoFragmentType:
		err = dec.Decode(&vp.Video)
		break
	case HeartBeatType:
		err = dec.Decode(&hb.BeatsPerMinut)
		break
	case LocationType:
		err = dec.Decode(loc)
		break
	default:
		err = fmt.Errorf("received unkown msg type = %v", sentType)
		break
	}

	return sentType, err
}
