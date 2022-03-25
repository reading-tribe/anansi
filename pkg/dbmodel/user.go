package dbmodel

type User struct {
	ID           int64  `dynamodb:"id"`
	EmailAddress string `dynamodbav:"email_address"`
	PasswordHash string `dynamodbav:"password_hash"`
}
