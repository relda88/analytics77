package main

import (
	"errors"
	"fmt"
	"strconv"
)

type configuration struct {
	port        string
	nameLogfile string
}

func extractConfiguration(raw map[string]any) (*configuration, error) {
	server, exists := raw["server"].(map[string]any)
	if !exists {
		return nil,
			errors.New(
				"invalid or missing server configuration",
			)
	}

	// Go defaults JSON numbers to float64.
	// Type assertions return the zero-value (empty string / 0) if they fail.
	port, couldCastPort := server["port"].(float64)
	if !couldCastPort {
		return nil,
			fmt.Errorf(
				"invalid port as %v",
				server["port"],
			)
	}

	nameLogfile, couldCastNameLogfile := server["logfile"].(string)
	if !couldCastNameLogfile {
		return nil,
			fmt.Errorf(
				"invalid name log file as %v",
				server["logfile"],
			)
	}

	return &configuration{
			port:        strconv.Itoa(int(port)),
			nameLogfile: nameLogfile,
		},
		nil
}
