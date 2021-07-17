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
