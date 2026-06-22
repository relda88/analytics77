package slogging

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/tudorhulban/arenalog"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
)

type ServiceLogging struct {
	Logger *arenalog.Logger
}

// NewServiceLogging provides a closing function.
func NewServiceLogging(pathLogFile string) (*ServiceLogging, func(), error) {
	fileHTTPServer, errFileHTTP := os.OpenFile(
		pathLogFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if errFileHTTP != nil {
		return nil, nil,
			errFileHTTP
	}

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		fileHTTPServer,

		helpers.TernaryWithValueIn(
			[]int{1},
			runtime.NumCPU(),
			nil,
			bytearena.WithCounterCoreCPU(),
		),
	)
	if errCrIngestor != nil {
		return nil, nil,
			errCrIngestor
	}

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	fnClosing := func() {
		fileHTTPServer.Close()

		cancel()
		<-chIngestionEnd
	}

	l, errCrLogger := arenalog.NewLogger(
		&arenalog.ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: arenalog.LevelInfo,

			WithFatalWriter: os.Stdout,
			WithJSON:        true,
		},

		arenalog.WithTimestampRFC3339UTC(ctx),
	)
	if errCrLogger != nil {
		fnClosing()

		return nil, nil,
			fmt.Errorf(
				"failed to create service logging: %w",
				errCrLogger,
			)
	}

	return &ServiceLogging{
			Logger: l,
		},
		fnClosing,
		nil
}
