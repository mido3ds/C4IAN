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

// AudioMsg is sent to unit when the owner presses record
// to record their voice and send immediately to cmd
type AudioMsg struct {
	Audio []byte
}

func (a AudioMsg) Send(enc *gob.Encoder) error {
	err := enc.Encode(byte(AudioMsgType))
	if err != nil {
		return err
	}
	return enc.Encode(a.Audio)
}

// CodeMsg is sent to unit when the owner preses some
// number to send as a code to cmd
type CodeMsg struct {
	Code int
}

func (c CodeMsg) Send(enc *gob.Encoder) error {
	err := enc.Encode(byte(CodeMsgType))
	if err != nil {
		return err
	}
	return enc.Encode(c.Code)
}

// HeartBeat is sent to unit periodically
// after connection establishment between HAL and unit
// sent each 1 second
type HeartBeat struct {
	// from 0 (dead or no sensor) to 100
	BeatsPerMinut int
}

// Location is a gps coordinates
// which are sent to unit periodically
// after connection establishment between HAL and unit
// sent each 1 second
type Location struct {
	Lon, Lat float64
}

type SensorData struct {
	HeartBeat
	Location
}

func (s SensorData) Send(enc *gob.Encoder) error {
	err := enc.Encode(byte(SensorDataType))
	if err != nil {
		return err
	}
	return enc.Encode(s)
}

func RecvFromHAL(dec *gob.Decoder, vp *VideoFragment, s *SensorData, a *AudioMsg, c *CodeMsg) (Type, error) {
	var sentType Type
	err := dec.Decode(&sentType)
	if err != nil {
		return 0, err
	}

	switch sentType {
	case VideoFragmentType:
		err = dec.Decode(&vp.Video)
		break
	case SensorDataType:
		err = dec.Decode(s)
		break
	case AudioMsgType:
		err = dec.Decode(&a.Audio)
		break
	case CodeMsgType:
		err = dec.Decode(&c.Code)
		break
	default:
		err = fmt.Errorf("received unkown msg type = %v", sentType)
		break
	}

	return sentType, err
}
