package nettypes

import "github.com/reading-tribe/anansi/pkg/idx"

type ListUserResponseItem struct {
	ID           idx.UserID `json:"id"`
	EmailAddress string     `json:"email_address"`
	Confirmed    bool       `json:"confirmed"`
}

type ListUsersResponse []ListUserResponseItem
