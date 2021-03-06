package transports

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type Topic interface {
	Publish(ctx context.Context, message *pubsub.Message) *pubsub.PublishResult
	String() string
}

func NewPubsubPublisher(ctx context.Context, t Topic) func(interface{}) error {
	return func(data interface{}) error {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("error while json marshalling data %s", err)
		}

		_, err = t.Publish(ctx, &pubsub.Message{Data: jsonBytes}).Get(ctx)
		if err != nil {
			return fmt.Errorf("failed publishing %s to topic %s: %s", data, t, err)
		}

		log.Printf("Successfully published data to topic %s\n", t.String())
		return nil
	}
}
