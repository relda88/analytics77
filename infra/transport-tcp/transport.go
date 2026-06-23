package transporttcp

import (
	"encoding/gob"
	"net"

	"github.com/tudorhulban/analytics77/services/sanalytics"
	"github.com/tudorhulban/analytics77/services/slogging"
	"github.com/tudorhulban/analytics77/shared"
	"github.com/tudorhulban/arenalog"
	"github.com/tudorhulban/hxhelpers/piers"
)

type TransportTCP struct {
	listener net.Listener

	serviceLogging   *slogging.ServiceLogging
	serviceAnalytics *sanalytics.ServiceAnalytics

	logContext *arenalog.LogContext
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

	// 1. Create decoder for this specific payload
	decoder := gob.NewDecoder(conn)

	var batch shared.Requests

	// 2. Expect exactly ONE decode operation
	if errDecode := decoder.Decode(&batch); errDecode != nil {
		s.serviceLogging.
			Logger.
			Printf(
				"failed to decode payload from %s: %s\n",
				conn.RemoteAddr(),
				errDecode.Error(),
			)

		return // Exit immediately, defer handles the connection close
	}

	s.serviceLogging.
		Logger.
		Printf(
			"received %d request(s) from %s\n",
			len(batch),
			conn.RemoteAddr(),
		)

	// 3. Process the data
	errsValidationEvents, errsProcessEvents := s.serviceAnalytics.RecordEvents(batch)
	if len(errsValidationEvents) > 0 {
		s.serviceLogging.
			Logger.
			Printf(
				"handleConnection - validation error(s)(%d) from %s: %v",
				len(errsValidationEvents),
				conn.RemoteAddr(),
				errsValidationEvents,
			)

		return
	}

	if len(errsProcessEvents) > 0 {
		s.serviceLogging.
			Logger.
			Printf(
				"handleConnection - processing error(s)(%d) from %s: %v",
				len(errsProcessEvents),
				conn.RemoteAddr(),
				errsProcessEvents,
			)
	}
}

func (s *TransportTCP) Start() error {
	s.serviceLogging.
		Logger.
		Printf(
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
