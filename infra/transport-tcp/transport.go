package transporttcp

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/TudorHulban/analytics77/services/sanalytics"
	"github.com/TudorHulban/analytics77/shared"
)

type Server struct {
	listener net.Listener
	service  *sanalytics.ServiceAnalytics
}

func NewServer(l net.Listener, service *sanalytics.ServiceAnalytics) *Server {
	return &Server{
		listener: l,
		service:  service,
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)

	for {
		var batch shared.Requests

		if errDecode := decoder.Decode(&batch); errDecode != nil {
			// Expected EOF or connection reset when client disconnects
			break
		}

		errsValidationEvents, errsProcessEvents := s.service.RecordEvents(batch)
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

func (s *Server) Start() error {
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
