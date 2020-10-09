package transports

import (
	"encoding/json"
	"fmt"
	"log"
)

func NewStandardOutputPrinter() func(interface{}) error {
	return func(data interface{}) error {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("error while json marshalling data %s", err)
		}

		log.Println(string(jsonBytes))
		return nil
	}
}
