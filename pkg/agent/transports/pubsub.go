package transports

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
)

type Topic interface {
	Publish(ctx context.Context, message *pubsub.Message) *pubsub.PublishResult
}

func NewPubsub(ctx context.Context, t Topic) func(data interface{}) error {
	return func(data interface{}) error {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("error while json marshalling data %s", err)
		}

		_, err = t.Publish(ctx, &pubsub.Message{Data: jsonBytes}).Get(ctx)
		if err != nil {
			return fmt.Errorf("failed publishing %s to topic %s: %s", data, t, err)
		}

		return nil
	}
}
