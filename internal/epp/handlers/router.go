package handlers

import (
	"drid/internal/epp/types"
	"log"
)

type Router struct {
	handler *CommandHandler
}

func NewRouter(handler *CommandHandler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) HandleRequest(client *types.Client, request *types.EPPRequest) (*types.EPPResponse, error) {
	if greeting := request.Greeting; greeting != nil {
		log.Printf("Greeting from server: %s at %s", greeting.ServerID, greeting.ServerDate)
		return nil, nil
	}

	if command := request.Command; command != nil {
		switch {
		case command.Login != nil:
			return r.handler.HandleLogin(client, command.Login)
		case command.Logout != nil:
			log.Println("Logout command")
		case command.Check != nil:
			log.Println("Check command")
		case command.Info != nil:
			log.Println("Info command")
		case command.Create != nil:
			log.Println("Create command")
		case command.Delete != nil:
			log.Println("Delete command")
		case command.Renew != nil:
			log.Println("Renew command")
		case command.Transfer != nil:
			log.Println("Transfer command")
		case command.Update != nil:
			log.Println("Update command")
		default:
			log.Println("Unknown command type")
		}
		return nil, nil
	}

	log.Println("No valid EPP message component found")
	return nil, nil
}
