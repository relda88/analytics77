package main

import (
	"fmt"
	"strconv"
)

func extractConfiguration(raw map[string]any) (string, error) {
	server, exists := raw["server"].(map[string]any)
	if !exists {
		return "",
			fmt.Errorf(
				"invalid or missing server configuration",
			)
	}

	// Go defaults JSON numbers to float64.
	// Type assertions return the zero-value (empty string / 0) if they fail.
	port, couldCastPort := server["port"].(float64)
	if !couldCastPort {
		return "",
			fmt.Errorf(
				"invalid port as %v",
				server["port"],
			)
	}

	return strconv.Itoa(int(port)),
		nil
}
