package transporttcp

import (
	"encoding/gob"
	"log"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/TudorHulban/analytics77/domain"
	"github.com/TudorHulban/analytics77/services/sanalytics"
	"github.com/TudorHulban/analytics77/shared"
	"github.com/stretchr/testify/require"
)

func TestTransport_TCP(t *testing.T) {
	dummyURL, _ := url.Parse("https://example.com/analytics")

	tests := []struct {
		description   string
		inputRequests shared.Requests
		expectedCount int
	}{
		{
			description: "1. Send single request",
			inputRequests: shared.Requests{
				{
					RemoteAddr: "192.168.1.1:5000",
					Host:       "example.com",
					Method:     "POST",
					URL:        dummyURL,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
				},
			},
			expectedCount: 1,
		},
		{
			description: "2. Send multiple requests in one batch",
			inputRequests: shared.Requests{
				{
					RemoteAddr: "192.168.1.2:5001",
					Host:       "api.com",
					Method:     "GET",
					URL:        dummyURL,
				},
				{
					RemoteAddr: "192.168.1.3:5002",
					Host:       "metrics.com",
					Method:     "PUT",
					URL:        dummyURL,
				},
			},
			expectedCount: 2,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				listener, errListener := net.Listen("tcp", "127.0.0.1:0")
				if errListener != nil {
					t.Fatalf("failed to create listener: %v", errListener)
				}

				serviceAnalytics := sanalytics.NewServiceAnalytics(domain.NewDataCenter())

				server := NewServer(
					listener,
					serviceAnalytics,
				)

				go func() {
					if errServerStart := server.Start(); errServerStart != nil {
						log.Printf("server stopped: %v", errServerStart)
					}
				}()

				// Give the OS a tiny moment to bind the socket
				time.Sleep(10 * time.Millisecond)

				connClient, errListener := net.Dial("tcp", server.listener.Addr().String())
				if errListener != nil {
					t.Fatalf("failed to dial server: %v", errListener)
				}

				if err := gob.
					NewEncoder(connClient).
					Encode(&tc.inputRequests); err != nil {
					t.Fatalf("transport encoding failed: %v", err)
				}

				// Close to flush and trigger EOF on the server side
				connClient.Close()

				require.Eventually(t,
					func() bool {
						return len(serviceAnalytics.DC.GetLastHourRecordsPerSite()) == tc.expectedCount
					},
					1*time.Second,
					10*time.Millisecond,
					serviceAnalytics.DC.GetLastHourRecordsPerSite().String(),
				)
			},
		)
	}
}
