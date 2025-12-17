package main

import (
	"drid/internal/epp"
	"log"
	"net"
)

func main() {
	eppHandler := epp.NewEPPHandler()

	ln, err := net.Listen("tcp", ":7000")
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
	defer ln.Close()

	log.Println("EPP server started successfully on :7000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Could not accept connection: %v", err)
			continue
		}

		remoteAddr := conn.RemoteAddr().String()
		log.Printf("New connection accepted from %s", remoteAddr)

		go func(c net.Conn, addr string) {
			defer func() {
				c.Close()
				log.Printf("Connection closed from %s", addr)
			}()
			eppHandler.Handle(c)
		}(conn, remoteAddr)
	}
}
