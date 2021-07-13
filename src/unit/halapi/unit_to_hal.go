package halapi

import (
	"encoding/gob"
	"fmt"
)

// StartVideoStream
type StartVideoStream struct {
	// how many buffers per second should be sent
	RefreshRate int
}

func (s *StartVideoStream) Send(enc *gob.Encoder) error {
	err := enc.Encode(byte(StartVideoStreamType))
	if err != nil {
		return err
	}
	return enc.Encode(s.RefreshRate)
}

// EndVideoStream
type EndVideoStream struct{}

func (s *EndVideoStream) Send(enc *gob.Encoder) error {
	return enc.Encode(byte(EndVideoStreamType))
}

// ShowAudioMsg
type ShowAudioMsg struct {
	Audio []byte
}

func (s *ShowAudioMsg) Send(enc *gob.Encoder) error {
	err := enc.Encode(byte(ShowAudioMsgType))
	if err != nil {
		return err
	}
	return enc.Encode(s.Audio)
}

// ShowCodeMsg
type ShowCodeMsg struct {
	Code int
}

func (s *ShowCodeMsg) Send(enc *gob.Encoder) error {
	err := enc.Encode(byte(ShowCodeMsgType))
	if err != nil {
		return err
	}
	return enc.Encode(s.Code)
}

func RecvFromUnit(dec *gob.Decoder, svs *StartVideoStream, evs *EndVideoStream, sam *ShowAudioMsg, scm *ShowCodeMsg) (Type, error) {
	var sentType Type
	err := dec.Decode(&sentType)
	if err != nil {
		return 0, err
	}

	switch sentType {
	case StartVideoStreamType:
		err = dec.Decode(&svs.RefreshRate)
		break
	case EndVideoStreamType:
		break
	case ShowAudioMsgType:
		err = dec.Decode(&sam.Audio)
		break
	case ShowCodeMsgType:
		err = dec.Decode(&scm.Code)
		break
	default:
		err = fmt.Errorf("received unkown msg type = %v", sentType)
		break
	}

	return sentType, err
}
