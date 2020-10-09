package transports

import (
	"fmt"
	"github.com/peppys/roku-discovery-agent/pkg/agent"
	"strings"
	"sync"
)

func NewBulkTransporter(transports []agent.Transport) func(interface{}) error {
	return func(payload interface{}) error {
		var wg sync.WaitGroup
		var errors []string

		for _, transport := range transports {
			wg.Add(1)

			go func(transport agent.Transport) {
				defer wg.Done()
				err := transport(payload)
				if err != nil {
					errors = append(errors, err.Error())
				}
			}(transport)
		}
		wg.Wait()

		if len(errors) > 0 {
			return fmt.Errorf("error while transporting: %s", strings.Join(errors, ", "))
		}

		return nil
	}
}
