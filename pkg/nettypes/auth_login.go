package nettypes

import "github.com/reading-tribe/anansi/pkg/idx"

type LoginRequest struct {
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
}

type LoginResponse struct {
	Token   idx.SessionID `json:"token"`
	Message string        `json:"message"`
}
