package token

import "time"

// Maker is an interface responsible for managing tokens
type Maker interface {
	// CreateToken creates and signs a token based on specific username and duration 
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}