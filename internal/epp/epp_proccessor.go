package epp

import (
	"encoding/xml"
	"io"
	"log"
	"net"
	"drid/internal/epp/handlers"
	"drid/internal/epp/types"
)

type EPPProccessor struct {
	router *handlers.Router
}

func NewRequestProcessor(router *handlers.Router) *EPPProccessor {
	return &EPPProccessor{
		router: router,
	}
}

func (p *EPPProccessor) Handle(conn net.Conn) {
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

func (p *EPPProccessor) ProcessRequest(client *types.Client, rawXML string) ([]byte, error) {
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