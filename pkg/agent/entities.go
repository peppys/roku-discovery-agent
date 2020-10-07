package agent

import "github.com/peppys/roku-discovery-agent/pkg/roku"

type QueryDeviceResult struct {
	Device roku.Device
	Error  error
}
type QueryActiveAppResult struct {
	ActiveApp roku.App
	Error     error
}
type QueryMediaPlayerResult struct {
	MediaPlayer roku.MediaPlayer
	Error       error
}

type QueryResult struct {
	Data QueryResultData
	Error error
}

type QueryResultData struct {
	Type  string
	Value interface{}
}
