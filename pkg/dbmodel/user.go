package dbmodel

import "github.com/reading-tribe/anansi/pkg/idx"

type User struct {
	ID           idx.UserID `dynamodb:"id" json:"id"`
	EmailAddress string     `dynamodbav:"email_address"`
	PasswordHash string     `dynamodbav:"password_hash"`
	InviteCode   string     `dynamodbav:"invite_code"`
	Confirmed    bool       `dynamodb:"confirmed"`
}

func (u User) AsMap() map[string]any {
	return map[string]any{
		"id":            u.ID,
		"email_address": u.EmailAddress,
		"password_hash": u.PasswordHash,
		"invite_code":   u.InviteCode,
		"confirmed":     u.Confirmed,
	}
}
