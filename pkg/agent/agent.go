package agent

import (
	"fmt"
	"github.com/peppys/roku-discovery-agent/pkg/roku"
	"time"
)

type Agent struct {
	roku        RokuClient
	pubSubTopic string
}

type RokuClient interface {
	Discover() (*roku.Client, error)
}

func New(topic string, roku RokuClient) *Agent {
	return &Agent{
		roku,
		topic,
	}
}

func (a *Agent) Start() {
	fmt.Printf("Starting agent - publishing stats to topic %s\n", a.pubSubTopic)
	for {
		fmt.Println("Searching for roku...")
		client, err := a.roku.Discover()
		if err != nil {
			fmt.Printf("Roku not found: %s", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf("Discovered roku with IP %s...\n", client.Host)

		deviceChan := make(chan *QueryDeviceResult, 1)
		activeAppChan := make(chan *QueryActiveAppResult, 1)
		mediaPlayerChan := make(chan *QueryMediaPlayerResult, 1)

		fmt.Println("Collecting stats...")
		go a.queryDevice(client, deviceChan)
		go a.queryActiveApp(client, activeAppChan)
		go a.queryMediaPlayer(client, mediaPlayerChan)

		queryDeviceResult := <-deviceChan
		if queryDeviceResult.Error != nil {
			fmt.Printf("Error while querying for device: %s", queryDeviceResult.Error)
			time.Sleep(5 * time.Second)
			continue
		}
		queryActiveAppResult := <-activeAppChan
		if queryActiveAppResult.Error != nil {
			fmt.Printf("Error while querying for active app: %s", queryActiveAppResult.Error)
			time.Sleep(5 * time.Second)
			continue
		}
		queryMediaPlayerResult := <-mediaPlayerChan
		if queryMediaPlayerResult.Error != nil {
			fmt.Printf("Error while querying for media player: %s", queryMediaPlayerResult.Error)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf("Device: %s\n", queryDeviceResult.Device)
		fmt.Printf("Active app: %s\n", queryActiveAppResult.ActiveApp)
		fmt.Printf("Media player: %s\n", queryMediaPlayerResult.MediaPlayer)
		time.Sleep(5 * time.Second)
	}
}

func (a *Agent) queryDevice(client *roku.Client, results chan *QueryDeviceResult) {
	device, err := client.QueryDevice()
	results <- &QueryDeviceResult{
		device,
		err,
	}
}

func (a *Agent) queryActiveApp(client *roku.Client, results chan *QueryActiveAppResult) {
	activeApp, err := client.QueryActiveApp()
	results <- &QueryActiveAppResult{
		activeApp,
		err,
	}
}

func (a *Agent) queryMediaPlayer(client *roku.Client, results chan *QueryMediaPlayerResult) {
	mediaPlayer, err := client.QueryMediaPlayer()
	results <- &QueryMediaPlayerResult{
		mediaPlayer,
		err,
	}
}
