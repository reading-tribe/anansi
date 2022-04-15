package nettypes

import "github.com/reading-tribe/anansi/pkg/idx"

type UpdateUserRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
	Confirmed    bool   `json:"confirmed"`
}

type UpdateUserResponse struct {
	ID           idx.UserID `json:"id"`
	EmailAddress string     `json:"email_address"`
	Confirmed    bool       `json:"confirmed"`
}
