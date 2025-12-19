package services

import "errors"

var (
	ErrEppInvalidAuthInfo              = errors.New("") // 2202 Invalid authorization information
	ErrAuthErrorServerClosesConnection = errors.New("") // 2501 Authentication error; server closing connection
)
