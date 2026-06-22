package transporttcp

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/tudorhulban/analytics77/services/sanalytics"
	"github.com/tudorhulban/analytics77/services/slogging"
	"github.com/tudorhulban/analytics77/shared"
	"github.com/tudorhulban/arenalog"
	"github.com/tudorhulban/hxhelpers/piers"
)

type TransportTCP struct {
	listener net.Listener

	serviceAnalytics *sanalytics.ServiceAnalytics
	serviceLogging   *slogging.ServiceLogging
	logContext       *arenalog.LogContext
}

type PiersNewTransportTCP struct {
	ServiceLogging   *slogging.ServiceLogging
	ServiceAnalytics *sanalytics.ServiceAnalytics
}

func NewTransportTCP(l net.Listener, dependencies *PiersNewTransportTCP) (*TransportTCP, error) {
	if errValidate := piers.ValidateDependencies(dependencies); errValidate != nil {
		return nil,
			errValidate
	}

	logContext := arenalog.
		NewLogContext(dependencies.ServiceLogging.Logger).
		WithRoot("transport", "TCP")

	return &TransportTCP{
			listener:         l,
			serviceAnalytics: dependencies.ServiceAnalytics,
			serviceLogging:   dependencies.ServiceLogging,
			logContext:       logContext,
		},
		nil
}

func (s *TransportTCP) GetListeningAddress() string {
	return s.listener.Addr().String()
}

func (s *TransportTCP) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := gob.NewDecoder(conn)

	for {
		var batch shared.Requests

		if errDecode := decoder.Decode(&batch); errDecode != nil {
			s.serviceLogging.
				Logger.
				Printf(
					"gob decoder error: %s\n",
					errDecode.Error(),
				)

			// Expected EOF or connection reset when client disconnects
			break
		}

		// TODO: take out
		log.Printf(
			"received %d requests\n",
			len(batch),
		)

		errsValidationEvents, errsProcessEvents := s.serviceAnalytics.RecordEvents(batch)
		if errsValidationEvents != nil {
			log.Printf(
				"handleConnection - validation error(s) from %s: %v",
				conn.RemoteAddr(),
				errsValidationEvents,
			)
		}

		if errsProcessEvents != nil {
			log.Printf(
				"handleConnection - processing error(s) from %s: %v",
				conn.RemoteAddr(),
				errsProcessEvents,
			)
		}
	}
}

func (s *TransportTCP) Start() error {
	log.Printf(
		"TCP transport server listening on %s",
		s.listener.Addr().String(),
	)

	for {
		conn, errListenAccept := s.listener.Accept()
		if errListenAccept != nil {
			return errListenAccept
		}

		go s.handleConnection(conn)
	}
}
