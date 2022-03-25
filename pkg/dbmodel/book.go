package dbmodel

type Book struct {
	ID            string `dynamodbav:"id" json:"id"`
	InternalTitle string `dynamodbav:"internal_title" json:"internal_title"`
	Authors       string `dynamodbav:"authors" json:"authors"`
}
