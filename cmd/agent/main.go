package main

import (
	"flag"
	"github.com/peppys/roku-discovery-agent/pkg/agent"
	"github.com/peppys/roku-discovery-agent/pkg/agent/collectors"
	"github.com/peppys/roku-discovery-agent/pkg/agent/transports"
	"github.com/peppys/roku-discovery-agent/pkg/roku"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var topic string
	flag.StringVar(&topic, "t", "", "google pubsub topic destination to publish roku stats")
	flag.Parse()

	if topic == "" {
		flag.Usage()
		os.Exit(1)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	a := agent.New(collectors.RokuCollector(roku.NewClient(http.DefaultClient)),
		agent.WithInterval(5*time.Second),
		agent.WithTransport(transports.NewStandardOutput()),
	)

	go func() {
		<-sigs
		log.Println("Stopping agent...")
		a.Stop()
		os.Exit(1)
	}()

	log.Println("Starting agent...")
	a.Start()
}
