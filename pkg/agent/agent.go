package agent

import (
	"fmt"
	"github.com/peppys/roku-discovery-agent/pkg/roku"
	"log"
	"sync"
	"time"
)

type Agent struct {
	roku        RokuClient
	pubSubTopic string
	transports  []Transport
}

type RokuClient interface {
	Discover() (*roku.Client, error)
}

type Transport interface {
	Send(data interface{}) error
	ID() string
}

func New(topic string, roku RokuClient, transports []Transport) *Agent {
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

		log.Printf("Discovered roku with IP %s...\n", client.Host)
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

func (a *Agent) queryDevice(client *roku.Client, results chan QueryResult) {
	device, err := client.QueryDevice()
	results <- QueryResult{
		QueryResultData{
			"device",
			device,
		},
		err,
	}
}

func (a *Agent) queryActiveApp(client *roku.Client, results chan QueryResult) {
	activeApp, err := client.QueryActiveApp()
	results <- QueryResult{
		QueryResultData{
			"active_app",
			activeApp,
		},
		err,
	}
}

func (a *Agent) queryMediaPlayer(client *roku.Client, results chan QueryResult) {
	mediaPlayer, err := client.QueryMediaPlayer()
	results <- QueryResult{
		QueryResultData{
			"media_player",
			mediaPlayer,
		},
		err,
	}
}

func (a *Agent) collect(client *roku.Client) (map[string]interface{}, error) {
	queryResultChan := make(chan QueryResult)

	go a.queryDevice(client, queryResultChan)
	go a.queryActiveApp(client, queryResultChan)
	go a.queryMediaPlayer(client, queryResultChan)

	payload := make(map[string]interface{})
	for i := 0; i < 3; i++ {
		result := <-queryResultChan
		if result.Error != nil {
			return nil, fmt.Errorf("error while querying roku device %s", result.Error)
		}

		payload[result.Data.Type] = result.Data.Value
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
