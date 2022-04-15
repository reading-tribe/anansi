package nettypes

import "github.com/reading-tribe/anansi/pkg/dbmodel"

type CreateUserRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
}

type CreateUserResponse dbmodel.User
