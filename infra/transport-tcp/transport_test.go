package transporttcp

import (
	"encoding/gob"
	"log"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/analytics77/services/sanalytics"
	"github.com/tudorhulban/analytics77/services/sgeo"
	"github.com/tudorhulban/analytics77/services/sstorage"
	"github.com/tudorhulban/analytics77/shared"
)

func TestTransport_TCP(t *testing.T) {
	dummyURL, _ := url.Parse("https://example.com/analytics")

	offsets := helpers.TimestampOffsets{
		OffsetUTC: -3,
	}

	apiKey := ""

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

				serviceStorage, errCrServiceStorage := sstorage.NewServiceStorage(".")
				require.NoError(t, errCrServiceStorage)
				require.NotNil(t, serviceStorage)

				serviceGeo, errCrServiceGeo := sgeo.NewServiceGeo(
					&sgeo.ParamsNewServiceGeo{
						APIKeyGeolocation: apiKey,
					},
					serviceStorage,
				)
				require.NoError(t, errCrServiceGeo)
				require.NotNil(t, serviceGeo)

				serviceAnalytics := sanalytics.NewServiceAnalytics(
					&sanalytics.PiersNewServiceAnalytics{
						ServiceGeo: serviceGeo,
					},
					&offsets,
				)

				server := NewTransportTCP(
					listener,
					&PiersNewTransportTCP{
						ServiceAnalytics: serviceAnalytics,
					},
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
