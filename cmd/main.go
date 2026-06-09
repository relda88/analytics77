package main

import (
	"log"
	"net"

	"github.com/TudorHulban/analytics77/domain"
	transporttcp "github.com/TudorHulban/analytics77/infra/transport-tcp"
	"github.com/TudorHulban/analytics77/services/sanalytics"
	// Ensure you import your domain or storage implementation package here
)

func main() {
	log.Println("Initializing Analytics Application...")

	dc := domain.NewDataCenter()

	serviceAnalytics := sanalytics.NewServiceAnalytics(dc)

	listener, errListener := net.Listen("tcp", "127.0.0.1:8000")
	if errListener != nil {
		log.Fatalf("failed to create listener: %v", errListener)
	}

	transportTCP := transporttcp.NewServer(
		listener,
		serviceAnalytics,
	)

	if err := transportTCP.Start(); err != nil {
		log.Fatalf("Fatal error running TCP server: %v", err)
	}
}
