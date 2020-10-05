package main

import (
	"flag"
	"github.com/peppys/roku-discovery-agent/pkg/agent"
	"github.com/peppys/roku-discovery-agent/pkg/roku"
	"net/http"
	"os"
)

func main() {
	var topic string
	flag.StringVar(&topic, "t", "", "google pubsub topic destination to publish roku stats")
	flag.Parse()

	if topic == "" {
		flag.Usage()
		os.Exit(1)
	}

	a := agent.New(topic, roku.NewClient(http.DefaultClient))
	a.Start()
}
