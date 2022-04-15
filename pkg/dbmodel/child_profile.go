package dbmodel

import "github.com/reading-tribe/anansi/pkg/idx"

type ChildProfile struct {
	ID       idx.ChildProfileID `dynamodbav:"id" json:"id"`
	UserID   idx.UserID         `dynamodbav:"user_id" json:"user_id"`
	Nickname string             `dynamodbav:"nickname" json:"nickname"`
}
