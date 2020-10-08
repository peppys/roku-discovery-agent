package main

import (
	"cloud.google.com/go/pubsub"
	"context"
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
	var project string
	var topic string
	flag.StringVar(&project, "p", "", "google pubsub project destination to publish roku stats")
	flag.StringVar(&topic, "t", "", "google pubsub topic destination to publish roku stats")
	flag.Parse()

	t := []agent.Transport{
		transports.NewStandardOutput(),
	}

	if (project != "") != (topic != "") {
		log.Fatalf("both project and topic must be provided to publish data to pubsub")
	}

	if project != "" && topic != "" {
		client, err := pubsub.NewClient(context.Background(), project)
		if err != nil {
			log.Fatalf("failed to instantiate pubsub service: %s", err)
		}

		topic := client.Topic(topic)
		t = append(t, transports.NewPubsub(context.Background(), topic))
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	a := agent.New(collectors.RokuCollector(roku.NewClient(http.DefaultClient)),
		agent.WithInterval(5*time.Second),
		agent.WithTransports(t),
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
