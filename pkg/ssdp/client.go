package ssdp

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const multicastAddress = "239.255.255.250:1900"

type Client struct {
}

var DefaultClient = &Client{}

func (c *Client) Search(searchType string) (*http.Response, error) {
	pc, err := net.ListenPacket("udp4", "")
	if err != nil {
		return &http.Response{}, fmt.Errorf("unable to listen to packets: %s", err)
	}
	defer pc.Close()
	pc.SetDeadline(time.Now().Add(time.Second * 2))

	addr, err := net.ResolveUDPAddr("udp4", multicastAddress)
	if err != nil {
		return &http.Response{}, fmt.Errorf("unable to resolve UDP addr: %s", err)
	}

	_, err = pc.WriteTo([]byte(buildSearchQuery(searchType)), addr)
	if err != nil {
		return &http.Response{}, fmt.Errorf("unable to write packets: %s", err)
	}

	for {
		buf := make([]byte, 8192)
		_, _, err = pc.ReadFrom(buf)
		if err != nil {
			return &http.Response{}, fmt.Errorf("unable to read packets: %s", err)
		}

		reader := bufio.NewReader(strings.NewReader(string(buf)))
		res, err := http.ReadResponse(reader, nil)
		if err != nil {
			return &http.Response{}, fmt.Errorf("unable to read http response: %s", err)
		}

		return res, nil
	}
}

func buildSearchQuery(searchType string) string {
	return fmt.Sprintf(`M-SEARCH * HTTP/1.1
Host: %s
Man: "ssdp:discover"
ST: %s
`, multicastAddress, searchType)
}
