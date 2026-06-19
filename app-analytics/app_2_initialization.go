package appanalytics

import (
	"fmt"
	"net"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/analytics77/infra/initialization"
	transporttcp "github.com/tudorhulban/analytics77/infra/transport-tcp"
	"github.com/tudorhulban/analytics77/services/slogging"
	"github.com/tudorhulban/hxerrors"
)

type ParamsInitializeApp struct {
	ConfigPort        string
	KeyGeolocationAPI string
	PathLogFile       string
}

func InitializeApp(params *ParamsInitializeApp) *App {
	listener, errListener := net.Listen(
		"tcp",
		fmt.Sprintf(
			"127.0.0.1:%s",
			params.ConfigPort,
		),
	)
	if errListener != nil {
		fmt.Printf(
			"error create listener: %s\n",
			errListener.Error(),
		)

		os.Exit(hxerrors.OSExitForConnectivityIssues)
	}

	serviceAnalytics := initialization.Services(
		&initialization.ParamsServices{
			Offsets: helpers.TimestampOffsets{
				OffsetUTC: -3,
			},
			APIKeyGeolocation: params.KeyGeolocationAPI,
		},
	)

	serviceLogging, fnCloseLogging, erCrServiceLogging := slogging.NewServiceLogging(params.PathLogFile)
	if erCrServiceLogging != nil {
		fmt.Printf(
			"error create listener: %s\n",
			erCrServiceLogging.Error(),
		)

		os.Exit(hxerrors.OSExitForLoggingIssues)
	}

	return &App{
		transportHTTP: fiber.New(
			fiber.Config{
				BodyLimit: 1 * 1024 * 1024, // in mb
			},
		),

		transportTCP: transporttcp.NewTransportTCP(
			listener,
			serviceAnalytics,
		),

		serviceAnalytics: serviceAnalytics,
		serviceLogging:   serviceLogging,

		fnFreeResources: fnCloseLogging,
	}
}
