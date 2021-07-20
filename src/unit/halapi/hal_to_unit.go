package halapi

import (
	"fmt"
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

// VideoFragment is sent to unit periodically
// as configured in StartVideoStream sent by unit
type VideoFragment struct {
	Video    []byte
	Metadata []byte
	Filename string
}

func (s VideoFragment) Write(w io.Writer) (err error) {
	b, err := msgpack.Marshal(&s)
	b2, err := msgpack.Marshal(&WrappedMsg{
		Type: VideoFragmentType,
		Body: b,
	})
	_, err = w.Write(b2)
	return
}

// AudioMsg is sent to unit when the owner presses record
// to record their voice and send immediately to cmd
type AudioMsg struct {
	Audio []byte
}

func (a AudioMsg) Write(w io.Writer) (err error) {
	b, err := msgpack.Marshal(&a)
	b2, err := msgpack.Marshal(&WrappedMsg{
		Type: AudioMsgType,
		Body: b,
	})
	_, err = w.Write(b2)
	return
}

// CodeMsg is sent to unit when the owner preses some
// number to send as a code to cmd
type CodeMsg struct {
	Code int
}

func (c CodeMsg) Write(w io.Writer) (err error) {
	b, err := msgpack.Marshal(&c)
	b2, err := msgpack.Marshal(&WrappedMsg{
		Type: CodeMsgType,
		Body: b,
	})
	_, err = w.Write(b2)
	return
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

func (s SensorData) Write(w io.Writer) (err error) {
	b, err := msgpack.Marshal(&s)
	b2, err := msgpack.Marshal(&WrappedMsg{
		Type: SensorDataType,
		Body: b,
	})
	_, err = w.Write(b2)
	return
}

func ReadFromHAL(r io.Reader, vp *VideoFragment, s *SensorData, a *AudioMsg, c *CodeMsg) (recvdType Type, err error) {
	var wrapped WrappedMsg

	dec := msgpack.NewDecoder(r)
	err = dec.Decode(&wrapped)

	recvdType = wrapped.Type

	switch recvdType {
	case VideoFragmentType:
		err = msgpack.Unmarshal(wrapped.Body, vp)
		break
	case SensorDataType:
		err = msgpack.Unmarshal(wrapped.Body, s)
		break
	case AudioMsgType:
		err = msgpack.Unmarshal(wrapped.Body, a)
		break
	case CodeMsgType:
		err = msgpack.Unmarshal(wrapped.Body, c)
		break
	default:
		err = fmt.Errorf("received unkown msg type = %v", recvdType)
		break
	}

	return
}
