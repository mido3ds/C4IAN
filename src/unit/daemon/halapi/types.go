package halapi

type Type byte

const (
	VideoFragmentType        Type = 0
	HeartBeatType        Type = 1
	LocationType         Type = 2
	StartVideoStreamType Type = 3
	EndVideoStreamType   Type = 4
	ShowAudioMsgType     Type = 5
	ShowCodeMsgType      Type = 6
)
