package dbmodel

import "github.com/reading-tribe/anansi/pkg/idx"

type Session struct {
	Key              idx.SessionID `dynamodbav:"key"`
	CurrentProfileID string        `dynamodbav:"current_child_profile_id"`
	UserID           string        `dynamodbav:"user_id"`
	EmailAddress     string        `dynamodbav:"email_address"`
	ExpiresAtUnix    int64         `dynamodbav:"expires_at"`
}
