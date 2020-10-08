package agent

import (
	"log"
	"sync"
	"time"
)

type Agent struct {
	collect    Collector
	transports []Transport
	interval   time.Duration
	stop       chan string
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

func WithTransports(transports []Transport) Option {
	return func(agent *Agent) {
		agent.transports = transports
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
		default:
			// TODO - check against time instead of pausing goroutine
			time.Sleep(a.interval)
			break
		}

		payload, err := a.collect()
		if err != nil {
			log.Printf("Error while collecting stats: %s\n", err)
			continue
		}

		a.transport(payload)
	}
}

func (a *Agent) Stop() {
	a.stop <- "stop"
	<-a.stop
	close(a.stop)
}

func (a *Agent) transport(payload map[string]interface{}) {
	var wg sync.WaitGroup
	for _, transport := range a.transports {
		wg.Add(1)

		go func(transport Transport) {
			defer wg.Done()
			err := transport(payload)
			if err != nil {
				log.Printf("Error while sending transport %s\n", err)
			}
		}(transport)
	}
	wg.Wait()
}
