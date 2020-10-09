package transports

import (
	"encoding/json"
	"fmt"
)

func NewStandardOutputPrinter() func(interface{}) error {
	return func(data interface{}) error {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("error while json marshalling data %s", err)
		}

		fmt.Println(string(jsonBytes))
		return nil
	}
}
