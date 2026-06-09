package transporttcp

import (
	"encoding/gob"
	"log"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/TudorHulban/analytics77/domain"
	"github.com/TudorHulban/analytics77/helpers"
	"github.com/TudorHulban/analytics77/services/sanalytics"
	"github.com/TudorHulban/analytics77/shared"
	"github.com/stretchr/testify/require"
)

func TestTransport_TCP(t *testing.T) {
	dummyURL, _ := url.Parse("https://example.com/analytics")

	offsets := helpers.TimestampOffsets{
		OffsetUTC: -3,
	}

	tests := []struct {
		description   string
		inputRequests shared.Requests
		expectedCount int
	}{
		{
			description: "1. Send single request",
			inputRequests: shared.Requests{
				{
					RemoteAddr: "1.1.1.1",
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
					RemoteAddr: "2.2.2.2",
					Host:       "api.com",
					Method:     "GET",
					URL:        dummyURL,
				},
				{
					RemoteAddr: "8.8.8.8",
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
				require.NoError(t, errListener)

				serviceAnalytics := sanalytics.NewServiceAnalytics(
					domain.NewDataCenter(),
					&offsets,
				)

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
				require.NoError(t, errListener)

				require.NoError(t,
					gob.NewEncoder(connClient).Encode(&tc.inputRequests),
				)

				// Close to flush and trigger EOF on the server side
				connClient.Close()

				require.Eventually(t,
					func() bool {
						return len(serviceAnalytics.DC.GetLastHourRecordsPerSite(&offsets)) == tc.expectedCount
					},

					1*time.Second,
					10*time.Millisecond,
					serviceAnalytics.
						DC.
						GetLastHourRecordsPerSite(&offsets).String(),
				)
			},
		)
	}
}
