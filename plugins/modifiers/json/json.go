// plugin.go
package main

import (
	"encoding/json"
	"fmt"
)

var ModifierJson ModifierPlugin

type ModifierPlugin struct{}

func (m ModifierPlugin) Transform(input interface{}) (interface{}, error) {
	// Call the Transform method of the plugin
	// Convert input to a JSON-formatted string
	jsonString, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("input is not a string")
	}

	// Unmarshal the JSON string into an interface{}
	var result interface{}
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return result, nil
}
