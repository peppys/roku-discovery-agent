package roku

import (
	"encoding/xml"
	"fmt"
	"github.com/peppys/roku-discovery-agent/pkg/ssdp"
	"io/ioutil"
	"net/http"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient}
}

func (c *Client) Discover() (string, error) {
	resp, err := ssdp.Search("roku:ecp")
	if err != nil {
		return "", fmt.Errorf("error via ssdp: %s", err)
	}

	host, err := resp.Location()
	if err != nil {
		return "", fmt.Errorf("error parsing location: %s", err)
	}

	return host.String(), nil
}

func (c *Client) QueryDevice(host string) (Device, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/query/device-info", host))
	if err != nil {
		return Device{}, fmt.Errorf("error while querying for device info: %s", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Device{}, fmt.Errorf("error while reading body: %v", err)
	}

	var device Device
	err = xml.Unmarshal(data, &device)
	if err != nil {
		return Device{}, fmt.Errorf("error while unmarshalling: %v", err)
	}

	return device, nil
}

func (c *Client) QueryActiveApp(host string) (App, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/query/active-app", host))
	if err != nil {
		return App{}, fmt.Errorf("error while querying for active app: %s", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return App{}, fmt.Errorf("error while reading body: %v", err)
	}

	var app App
	err = xml.Unmarshal(data, &app)
	if err != nil {
		return App{}, fmt.Errorf("error while unmarshalling: %v", err)
	}

	return app, nil
}

func (c *Client) QueryMediaPlayer(host string) (MediaPlayer, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/query/media-player", host))
	if err != nil {
		return MediaPlayer{}, fmt.Errorf("error while querying for media player: %s", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return MediaPlayer{}, fmt.Errorf("error while reading body: %v", err)
	}

	var mediaPlayer MediaPlayer
	err = xml.Unmarshal(data, &mediaPlayer)
	if err != nil {
		return MediaPlayer{}, fmt.Errorf("error while unmarshalling: %v", err)
	}

	return mediaPlayer, nil
}
