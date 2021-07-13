package halapi

// StartVideoStream
type StartVideoStream struct {
	// how many buffers per second should be sent
	RefreshRate int
}

// EndVideoStream
type EndVideoStream struct{}

// ShowAudioMsg
type ShowAudioMsg struct {
	Audio []byte
}

// ShowCodeMsg
type ShowCodeMsg struct {
	Code int
}
