package halapi

type Type byte

const (
	VideoFragmentType    Type = 0
	SensorDataType       Type = 1
	StartVideoStreamType Type = 2
	EndVideoStreamType   Type = 3
	ShowAudioMsgType     Type = 4
	ShowCodeMsgType      Type = 5
	AudioMsgType         Type = 6
	CodeMsgType          Type = 7
)

// WrappedMsg is wrapping the sent/received msg
// with its type so you can know the body
// when you receive a WrappedMsg you should decode the body as indicated
// by the provided type
type WrappedMsg struct {
	Type Type
	Body []byte
}
