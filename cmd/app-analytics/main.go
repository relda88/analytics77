package main

import (
	"log"
	"net"

	"github.com/TudorHulban/analytics77/helpers"
	"github.com/TudorHulban/analytics77/infra/initialization"
	transporttcp "github.com/TudorHulban/analytics77/infra/transport-tcp"
)

func main() {
	listener, errListener := net.Listen("tcp", "127.0.0.1:8000")
	if errListener != nil {
		log.Fatalf("failed to create listener: %v", errListener)
	}

	serviceAnalytics := initialization.Services(
		&initialization.ParamsServices{
			Offsets: helpers.TimestampOffsets{
				OffsetUTC: -3,
			},
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
