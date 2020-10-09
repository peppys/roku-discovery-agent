package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"flag"
	"github.com/peppys/roku-discovery-agent/pkg/agent"
	"github.com/peppys/roku-discovery-agent/pkg/agent/collectors"
	"github.com/peppys/roku-discovery-agent/pkg/agent/transports"
	"github.com/peppys/roku-discovery-agent/pkg/roku"
	"github.com/peppys/roku-discovery-agent/pkg/ssdp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var projectId string
	var topicId string
	flag.StringVar(&projectId, "p", "", "google pubsub project destination to publish roku stats")
	flag.StringVar(&topicId, "t", "", "google pubsub topic destination to publish roku stats")
	flag.Parse()

	t := []agent.Transport{
		transports.NewStandardOutputPrinter(),
	}

	if (projectId != "") != (topicId != "") {
		log.Fatalf("both project and topic must be provided to publish data to pubsub")
	}

	if projectId != "" && topicId != "" {
		t = append(t, buildPubsubPublisher(projectId, topicId))
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	a := agent.New(collectors.RokuCollector(roku.NewClient(http.DefaultClient, ssdp.DefaultClient)),
		agent.WithInterval(5*time.Second),
		agent.WithTransport(transports.NewBulkTransporter(t)),
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

func buildPubsubPublisher(projectId string, topicId string) agent.Transport {
	client, err := pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		log.Fatalf("failed to instantiate pubsub service: %s", err)
	}

	topic := client.Topic(topicId)
	return transports.NewPubsubPublisher(context.Background(), topic)
}
