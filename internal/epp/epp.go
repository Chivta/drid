package epp

import (
	"encoding/xml"
	"io"
	"log"
	"net"
	"drid/internal/epp/handlers"
	"drid/internal/epp/services"
	"drid/internal/epp/types"
)

type RequestProcessor struct {
	router *handlers.Router
}

func NewRequestProcessor(router *handlers.Router) *RequestProcessor {
	return &RequestProcessor{
		router: router,
	}
}

func ListenForEPPConnections(address string) error {
	authService := services.NewAuthService()
	commandHandler := handlers.NewCommandHandler(authService)
	router := handlers.NewRouter(commandHandler)
	processor := NewRequestProcessor(router)

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

func (p *RequestProcessor) Handle(conn net.Conn) {
	connection := NewConnection(conn)
	client := &types.Client{}
	
	greeting := &types.EPPResponse{
		Greeting: &types.Greeting{
			ServerID:   "EPP Server v1.0",
			ServerDate: "2025-12-19T00:00:00Z",
		},
	}
	greetingXML, err := xml.Marshal(greeting)
	if err != nil {
		log.Printf("Failed to marshal greeting: %v", err)
		return
	}
	
	err = connection.WriteMessage(greetingXML)
	if err != nil {
		log.Printf("Failed to send greeting: %v", err)
		return
	}
	
	for {
		rawXML, err := connection.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read message: %v", err)
			}
			return
		}
		
		responseXML, err := p.ProcessRequest(client, string(rawXML))
		if err != nil {
			log.Printf("Failed to process request: %v", err)
			return
		}
		
		if responseXML != nil {
			err = connection.WriteMessage(responseXML)
			if err != nil {
				log.Printf("Failed to send response: %v", err)
				return
			}
		}
	}
}

func (p *RequestProcessor) ProcessRequest(client *types.Client, rawXML string) ([]byte, error) {
	request := &types.EPPRequest{}
	err := xml.Unmarshal([]byte(rawXML), request)
	if err != nil {
		return nil, err
	}

	response, err := p.router.HandleRequest(client, request)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	responseXML, err := xml.Marshal(response)
	if err != nil {
		return nil, err
	}

	return responseXML, nil
}