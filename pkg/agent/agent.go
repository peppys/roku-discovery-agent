package agent

import (
	"log"
	"time"
)

type Agent struct {
	collect   Collector
	transport Transport
	interval  time.Duration
	stop      chan string
}

type Option func(agent *Agent)

type Collector func() (map[string]interface{}, error)
type Transport func(interface{}) error

func New(c Collector, opts ...Option) *Agent {
	a := &Agent{
		collect:  c,
		interval: time.Second * 5,
		stop:     make(chan string),
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func WithTransport(t Transport) Option {
	return func(agent *Agent) {
		agent.transport = t
	}
}

func WithInterval(i time.Duration) Option {
	return func(agent *Agent) {
		agent.interval = i
	}
}

func (a *Agent) Start() {
	for {
		select {
		case <-a.stop:
			a.stop <- "stopped"
			return
		case <-time.After(a.interval):
			break
		}

		payload, err := a.collect()
		if err != nil {
			log.Printf("Error while collecting stats: %s\n", err)
			continue
		}

		err = a.transport(payload)
		if err != nil {
			log.Printf("Error while collecting transporting metrics: %s\n", err)
			continue
		}
	}
}

func (a *Agent) Stop() {
	a.stop <- "stop"
	<-a.stop
	close(a.stop)
}
