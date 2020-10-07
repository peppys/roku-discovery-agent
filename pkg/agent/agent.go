package agent

import (
	"fmt"
	"github.com/peppys/roku-discovery-agent/pkg/roku"
	"log"
	"sync"
	"time"
)

type Agent struct {
	roku        RokuDiscoveryClient
	pubSubTopic string
	transports  []Transport
}

type RokuDiscoveryClient interface {
	Discover() (RokuClient, error)
}

type RokuClient interface {
	GetHost() string
	QueryDevice() (roku.Device, error)
	QueryActiveApp() (roku.ActiveApp, error)
	QueryMediaPlayer() (roku.MediaPlayer, error)
}

type Transport interface {
	Send(data interface{}) error
	ID() string
}

func New(topic string, roku RokuDiscoveryClient, transports []Transport) *Agent {
	return &Agent{
		roku,
		topic,
		transports,
	}
}

func (a *Agent) Start() {
	log.Println("Starting agent")
	for {
		log.Println("Searching for roku...")
		client, err := a.roku.Discover()
		if err != nil {
			log.Printf("Roku not found: %s\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("Discovered roku with IP %s...\n", client.GetHost())
		log.Println("Collecting stats...")
		payload, err := a.collect(client)
		if err != nil {
			log.Printf("Error while collecting stats: %s\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		a.transport(payload)

		log.Println("Finished collecting stats")
		time.Sleep(5 * time.Second)
	}
}

func (a *Agent) queryDeviceData(queryFunc func() (interface{}, error), label string, results chan QueryResult) {
	data, err := queryFunc()
	results <- QueryResult{
		QueryResultData{
			label,
			data,
		},
		err,
	}
}

func (a *Agent) collect(client RokuClient) (map[string]interface{}, error) {
	queryResultChan := make(chan QueryResult)

	go a.queryDeviceData(func() (interface{}, error) {
		return client.QueryDevice()
	}, "device", queryResultChan)
	go a.queryDeviceData(func() (interface{}, error) {
		return client.QueryMediaPlayer()
	}, "media_player", queryResultChan)
	go a.queryDeviceData(func() (interface{}, error) {
		return client.QueryActiveApp()
	}, "active_app", queryResultChan)

	payload := make(map[string]interface{})
	for i := 0; i < 3; i++ {
		result := <-queryResultChan
		if result.Error != nil {
			return nil, fmt.Errorf("error while querying roku device %s", result.Error)
		}

		payload[result.Data.Label] = result.Data.Value
	}

	return payload, nil
}

func (a *Agent) transport(payload map[string]interface{}) {
	var wg sync.WaitGroup
	for _, transport := range a.transports {
		wg.Add(1)

		go func(transport Transport) {
			defer wg.Done()
			err := transport.Send(payload)
			if err != nil {
				log.Printf("Error while sending transport %s\n", err)
			}
		}(transport)
	}
	wg.Wait()
}
