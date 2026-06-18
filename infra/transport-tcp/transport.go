package transporttcp

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/tudorhulban/analytics77/services/sanalytics"
	"github.com/tudorhulban/analytics77/shared"
)

type TransportTCP struct {
	listener         net.Listener
	serviceAnalytics *sanalytics.ServiceAnalytics
}

func NewTransportTCP(l net.Listener, service *sanalytics.ServiceAnalytics) *TransportTCP {
	return &TransportTCP{
		listener:         l,
		serviceAnalytics: service,
	}
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
			// Expected EOF or connection reset when client disconnects
			break
		}

		errsValidationEvents, errsProcessEvents := s.serviceAnalytics.RecordEvents(batch)
		if errsValidationEvents != nil {
			log.Printf(
				"validation error(s) from %s: %v",
				conn.RemoteAddr(),
				errsValidationEvents,
			)
		}
		if errsProcessEvents != nil {
			log.Printf(
				"processing error(s) from %s: %v",
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
