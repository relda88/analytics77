package integrationtests

import (
	"encoding/gob"
	"log"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/TudorHulban/analytics77/domain"
	"github.com/TudorHulban/analytics77/helpers"
	transporttcp "github.com/TudorHulban/analytics77/infra/transport-tcp"
	"github.com/TudorHulban/analytics77/services/sanalytics"
	"github.com/TudorHulban/analytics77/services/sgeo"
	"github.com/TudorHulban/analytics77/services/sstorage"
	"github.com/TudorHulban/analytics77/shared"
	"github.com/stretchr/testify/require"
)

func TestAnalytics_E2E(t *testing.T) {
	dc := domain.NewDataCenter()
	require.NotNil(t, dc)

	offsetsROU := helpers.TimestampOffsets{
		OffsetUTC: -3,
	}

	apiKey := ""

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
			DC:         dc,
			ServiceGeo: serviceGeo,
		},
		&offsetsROU,
	)
	require.NotNil(t, serviceAnalytics)

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
					RemoteAddr: "82.77.237.37",
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
					RemoteAddr: "82.77.237.38",
					Host:       "api.com",
					Method:     "GET",
					URL:        dummyURL,
				},
				{
					RemoteAddr: "82.77.237.39",
					Host:       "metrics.com",
					Method:     "PUT",
					URL:        dummyURL,
				},
			},
			expectedCount: 3, // adding test case 1 requests
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				listener, errListener := net.Listen("tcp", "127.0.0.1:0")
				require.NoError(t, errListener)

				transportTCP := transporttcp.NewTransportTCP(
					listener,
					serviceAnalytics,
				)

				go func() {
					if errServerStart := transportTCP.Start(); errServerStart != nil {
						log.Printf("server stopped: %v", errServerStart)
					}
				}()

				// Give the OS a tiny moment to bind the socket
				time.Sleep(10 * time.Millisecond)

				connClient, errListener := net.Dial(
					"tcp",
					transportTCP.GetListeningAddress(),
				)
				require.NoError(t, errListener)
				require.NotZero(t, connClient)

				require.NoError(t,
					gob.NewEncoder(connClient).Encode(&tc.inputRequests),
				)

				connClient.Close()

				require.Eventually(t,
					func() bool {
						return len(serviceAnalytics.DC.GetLastHourRecordsPerSite(&offsetsROU)) == tc.expectedCount
					},

					1*time.Second,
					10*time.Millisecond,
					serviceAnalytics.
						DC.
						GetLastHourRecordsPerSite(&offsetsROU).String(),

					serviceAnalytics.DC.String(),
				)
			},
		)
	}
}
