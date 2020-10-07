package transports

import (
	"log"
)

type StandardOutputTransport struct {
}

func NewStandardOutput() *StandardOutputTransport {
	return &StandardOutputTransport{}
}

func (t *StandardOutputTransport) Send(data interface{}) error {
	log.Printf("%s", data)
	return nil
}

func (t *StandardOutputTransport) ID() string {
	return "standard-output"
}
