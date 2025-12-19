package handlers

import (
	"drid/internal/epp/types"
	"drid/internal/epp/services"
	"log"
)

func NewCommandHandler(authService *services.AuthService) *CommandHandler {
	return &CommandHandler{
		authService: authService,
	}
}

type CommandHandler struct {
	authService *services.AuthService
}

func (h *CommandHandler) HandleLogin(client *types.Client, loginCmd *types.LoginCommand) (*types.EPPResponse, error) {
	log.Printf("Authenticating client ID: %s", loginCmd.ClientID)

	err := h.authService.Authenticate(loginCmd.ClientID, loginCmd.Password)

	var responseCode types.ResultCode
	switch err {
	case services.ErrEppInvalidAuthInfo:
		responseCode = types.EppInvalidAuthInfo
	case services.ErrAuthErrorServerClosesConnection:
		responseCode = types.EppAuthFailedBye
	case nil:
		// Success - update client state
		client.Authenticated = true
		client.ClientID = loginCmd.ClientID
		responseCode = types.EppOk
	default:
		log.Printf("Unexpected authentication error for client %s: %v", loginCmd.ClientID, err)
		responseCode = types.EppCommandFailedBye
	}

	return &types.EPPResponse{
		Response: &types.Response{
			Result: []types.Result{
				{
					Code: responseCode.Code(),
					Msg: types.Msg{
						Lang: "en",
						Text: responseCode.Message(),
					},
				},
			},
		},
	}, nil
}
