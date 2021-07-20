package halapi

import (
	"fmt"
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

// StartVideoStream
type StartVideoStream struct {
	// how many buffers per second should be sent
	RefreshRate int
}

func (s StartVideoStream) Write(w io.Writer) (err error) {
	b, err := msgpack.Marshal(&s)
	b2, err := msgpack.Marshal(&WrappedMsg{
		Type: StartVideoStreamType,
		Body: b,
	})
	_, err = w.Write(b2)
	return
}

// EndVideoStream
type EndVideoStream struct{}

func (s EndVideoStream) Write(w io.Writer) (err error) {
	b, err := msgpack.Marshal(&WrappedMsg{Type: EndVideoStreamType})
	_, err = w.Write(b)
	return
}

// ShowAudioMsg is sent to HAL to show code msg to user
type ShowAudioMsg struct {
	Audio []byte
}

func (s ShowAudioMsg) Write(w io.Writer) (err error) {
	b, err := msgpack.Marshal(&s)
	b2, err := msgpack.Marshal(&WrappedMsg{Type: ShowAudioMsgType, Body: b})
	_, err = w.Write(b2)
	return
}

// ShowCodeMsg is sent to HAL to play audio msg to user
type ShowCodeMsg struct {
	Code int
}

func (s ShowCodeMsg) Write(w io.Writer) (err error) {
	b, err := msgpack.Marshal(&s)
	b2, err := msgpack.Marshal(&WrappedMsg{Type: ShowCodeMsgType, Body: b})
	_, err = w.Write(b2)
	return
}

func ReadFromUnit(r io.Reader, svs *StartVideoStream, evs *EndVideoStream, sam *ShowAudioMsg, scm *ShowCodeMsg) (recvdType Type, err error) {
	var wrapped WrappedMsg

	dec := msgpack.NewDecoder(r)
	err = dec.Decode(&wrapped)

	recvdType = wrapped.Type

	switch recvdType {
	case StartVideoStreamType:
		err = msgpack.Unmarshal(wrapped.Body, svs)
		break
	case EndVideoStreamType:
		err = msgpack.Unmarshal(wrapped.Body, evs)
		break
	case ShowAudioMsgType:
		err = msgpack.Unmarshal(wrapped.Body, sam)
		break
	case ShowCodeMsgType:
		err = msgpack.Unmarshal(wrapped.Body, scm)
		break
	default:
		err = fmt.Errorf("received unkown msg type = %v", recvdType)
		break
	}

	return
}
