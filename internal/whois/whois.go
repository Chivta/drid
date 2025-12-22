package whois

import (
	"bufio"
	"log"
	"net"
	"drid/pkg/db"
)

func NewWhoisServer(db *db.DB) *WhoisServer {
	return &WhoisServer{db: db}
}

type WhoisServer struct {
	db *db.DB
}

func (s *WhoisServer) ListenForWhoisConnections(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Could not start WHOIS server: %v", err)
	}
	defer ln.Close()

	log.Printf("WHOIS server started successfully on %s", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Could not accept connection: %v", err)
			continue
		}

		remoteAddr := conn.RemoteAddr().String()
		log.Printf("New WHOIS connection accepted from %s", remoteAddr)

		go func(c net.Conn, addr string) {
			defer func() {
				c.Close()
				log.Printf("WHOIS connection closed from %s", addr)
			}()
			s.handleWhoisRequest(c)
		}(conn, remoteAddr)
	}
}

func (s *WhoisServer) handleWhoisRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	request,err:= reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading request: %v", err)
		return
	}

	log.Printf("Request: %s\n",request)

	response,err := s.db.Whois(request)
	if err != nil {
		log.Printf("Error processing WHOIS request: %v", err)
		response = "Error processing request\n"
	}
	writer.WriteString(response)
	writer.Flush()
}