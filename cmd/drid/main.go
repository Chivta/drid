package main

import (
	"drid/internal/epp"
	"drid/internal/epp/handlers"
	"drid/internal/epp/services"
	"log"
	"net"
)

func main() {
	authService := services.NewAuthService()
	commandHandler := handlers.NewCommandHandler(authService)
	router := handlers.NewRouter(commandHandler)
	processor := epp.NewRequestProcessor(router)

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
			processor.Handle(c)
		}(conn, remoteAddr)
	}
}
