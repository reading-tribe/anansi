package dbmodel

type User struct {
	EmailAddress string `dynamodbav:"email_address"`
	PasswordHash string `dynamodbav:"password_hash"`
}
