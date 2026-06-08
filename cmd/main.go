package main

import (
	"log"

	"github.com/TudorHulban/analytics77/domain"
	transporttcp "github.com/TudorHulban/analytics77/infra/transport-tcp"
	"github.com/TudorHulban/analytics77/services/sanalytics"
	// Ensure you import your domain or storage implementation package here
)

func main() {
	log.Println("Initializing Analytics Application...")

	dc := domain.NewDataCenter()

	serviceAnalytics := sanalytics.NewServiceAnalytics(dc)

	transportTCP := transporttcp.NewServer(
		":8080",
		serviceAnalytics,
	)

	// 4. Start listening
	if err := transportTCP.Start(); err != nil {
		log.Fatalf("Fatal error running TCP server: %v", err)
	}
}
