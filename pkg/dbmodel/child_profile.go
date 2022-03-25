package dbmodel

type ChildProfile struct {
	ID       string `dynamodbav:"id" json:"id"`
	UserID   string `dynamodbav:"user_id" json:"user_id"`
	Nickname string `dynamodbav:"nickname" json:"nickname"`
}
