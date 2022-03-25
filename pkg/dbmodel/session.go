package dbmodel

type Session struct {
	Key              string `dynamodbav:"key"`
	CurrentProfileID string `dynamodbav:"current_child_profile_id"`
	UserID           string `dynamodbav:"user_id"`
	EmailAddress     string `dynamodbav:"email_address"`
	ExpiresAtUnix    int64  `dynamodbav:"expires_at"`
}
