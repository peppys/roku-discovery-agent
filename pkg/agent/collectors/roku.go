package collectors

import (
	"fmt"
	"github.com/peppys/roku-discovery-agent/pkg/roku"
)

type RokuClient interface {
	Discover() (string, error)
	QueryDevice(host string) (roku.Device, error)
	QueryActiveApp(host string) (roku.App, error)
	QueryMediaPlayer(host string) (roku.MediaPlayer, error)
}

type QueryResult struct {
	Data  QueryResultData
	Error error
}

type QueryResultData struct {
	Label string
	Value interface{}
}

func RokuCollector(roku RokuClient) func() (map[string]interface{}, error) {
	return func() (map[string]interface{}, error) {
		return collect(roku)
	}
}

func collect(roku RokuClient) (map[string]interface{}, error) {
	host, err := roku.Discover()
	if err != nil {
		return nil, fmt.Errorf("Roku not found: %s\n", err)
	}

	queryResultChan := make(chan QueryResult)

	go queryDeviceData(func() (interface{}, error) {
		return roku.QueryDevice(host)
	}, "device", queryResultChan)
	go queryDeviceData(func() (interface{}, error) {
		return roku.QueryMediaPlayer(host)
	}, "media_player", queryResultChan)
	go queryDeviceData(func() (interface{}, error) {
		return roku.QueryActiveApp(host)
	}, "active_app", queryResultChan)

	payload := make(map[string]interface{})
	for i := 0; i < 3; i++ {
		result := <-queryResultChan
		if result.Error != nil {
			return nil, fmt.Errorf("error while querying roku device %s", result.Error)
		}

		payload[result.Data.Label] = result.Data.Value
	}

	close(queryResultChan)

	return payload, nil
}

func queryDeviceData(queryFunc func() (interface{}, error), label string, results chan QueryResult) {
	data, err := queryFunc()
	results <- QueryResult{
		QueryResultData{
			label,
			data,
		},
		err,
	}
}
