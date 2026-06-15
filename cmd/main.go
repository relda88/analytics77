package main

import (
	"log"
	"net"

	"github.com/TudorHulban/analytics77/domain"
	"github.com/TudorHulban/analytics77/helpers"
	transporttcp "github.com/TudorHulban/analytics77/infra/transport-tcp"
	"github.com/TudorHulban/analytics77/services/sanalytics"
	"github.com/TudorHulban/analytics77/services/sgeo"
	"github.com/TudorHulban/analytics77/services/sstorage"
)

func main() {
	log.Println("Initializing Analytics Application...")

	serviceGeo, errSvcGeo := sgeo.NewServiceGeo(sstorage.NewServiceStorage())
	if errSvcGeo != nil {
		log.Printf(
			"service geo creation: %s",
			errSvcGeo.Error(),
		)
	}

	offsets := helpers.TimestampOffsets{
		OffsetUTC: -3,
	}

	serviceAnalytics := sanalytics.NewServiceAnalytics(
		&sanalytics.PiersNewServiceAnalytics{
			DC:         domain.NewDataCenter(),
			ServiceGeo: serviceGeo,
		},
		&offsets,
	)

	listener, errListener := net.Listen("tcp", "127.0.0.1:8000")
	if errListener != nil {
		log.Fatalf("failed to create listener: %v", errListener)
	}

	transportTCP := transporttcp.NewTransportTCP(
		listener,
		serviceAnalytics,
	)

	if err := transportTCP.Start(); err != nil {
		log.Fatalf("Fatal error running TCP server: %v", err)
	}
}
