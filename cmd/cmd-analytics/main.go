package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/tudorhulban/analytics77/cmd"
	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/analytics77/infra/initialization"
	transporttcp "github.com/tudorhulban/analytics77/infra/transport-tcp"
	"github.com/tudorhulban/hxerrors"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(
			"Error: Please provide the geolocation API key as the first argument.",
		)
		fmt.Println(
			"Usage: go run main.go <API_KEY>",
		)

		os.Exit(
			hxerrors.OSExitForApplicationIssues,
		)
	}

	keyGeolocationAPI := os.Args[1]

	configRaw := initialization.Configuration(cmd.PathConfig)

	configPort, errParse := extractConfiguration(configRaw)
	if errParse != nil {
		fmt.Printf(
			"error extract configuration: %s\n",
			errParse.Error(),
		)

		os.Exit(
			hxerrors.OSExitForConfigurationIssues,
		)
	}

	listener, errListener := net.Listen(
		"tcp",
		fmt.Sprintf(
			"127.0.0.1:%s",
			configPort,
		),
	)
	if errListener != nil {
		log.Fatalf("failed to create listener: %v", errListener)
	}

	serviceAnalytics := initialization.Services(
		&initialization.ParamsServices{
			Offsets: helpers.TimestampOffsets{
				OffsetUTC: -3,
			},
			APIKeyGeolocation: keyGeolocationAPI,
		},
	)

	transportTCP := transporttcp.NewTransportTCP(
		listener,
		serviceAnalytics,
	)

	if errTransportStart := transportTCP.Start(); errTransportStart != nil {
		log.Fatalf(
			"Fatal error running TCP server: %v",
			errTransportStart,
		)
	}
}
