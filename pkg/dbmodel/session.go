package dbmodel

type Session struct {
	Key           string `dynamodbav:"key"`
	EmailAddress  string `dynamodbav:"email_address"`
	ExpiresAtUnix int64  `dynamodbav:"expires_at"`
}
