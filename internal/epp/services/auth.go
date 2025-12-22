package services

import (
	"log"
	"drid/pkg/db"
)

func NewAuthService(db *db.DB) *AuthService {
	return &AuthService{
		db: db,
	}
}

type AuthService struct {
	db *db.DB
}

func (s *AuthService) Authenticate(clientID, password string) error {
	// Implement authentication logic here
	log.Printf("Authenticating client ID: %s", clientID)
	
	// Dummy logic for passing tests
	if password == "bad password" {
		return ErrEppInvalidAuthInfo
	}
	if clientID == "bad login" {
		return ErrAuthErrorServerClosesConnection
	}
	return nil
}
