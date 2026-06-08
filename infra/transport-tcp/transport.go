package transporttcp

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/TudorHulban/analytics77/services/sanalytics"
	"github.com/TudorHulban/analytics77/shared"
)

type Server struct {
	listenAddr string
	svc        *sanalytics.ServiceAnalytics
}

func NewServer(listenAddr string, service *sanalytics.ServiceAnalytics) *Server {
	return &Server{
		listenAddr: listenAddr,
		svc:        service,
	}
}

// Start kicks off the TCP listener and blocks while accepting connections
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("TCP transport server listening on %s", s.listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)

	for {
		var req shared.FlatRequest

		// Decode straight from the stream into memory
		if errDecode := decoder.Decode(&req); errDecode != nil {
			// Expected EOF or connection reset when client disconnects
			break
		}

		// Hand off data to the service layer safely
		if errProcessEvent := s.svc.RecordEvent(&req); errProcessEvent != nil {
			log.Printf("error recording event from %s: %v", conn.RemoteAddr(), errProcessEvent)
		}
	}
}
