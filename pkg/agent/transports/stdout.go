package transports

import (
	"encoding/json"
	"fmt"
)

func NewStandardOutput() func(data interface{}) error {
	return transport
}

func transport(data interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error while json marshalling data %s", err)
	}

	fmt.Println(string(jsonBytes))
	return nil
}
