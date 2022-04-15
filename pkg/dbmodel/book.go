package dbmodel

import "github.com/reading-tribe/anansi/pkg/idx"

type Book struct {
	ID            idx.BookID `dynamodbav:"id" json:"id"`
	InternalTitle string     `dynamodbav:"internal_title" json:"internal_title"`
	Authors       string     `dynamodbav:"authors" json:"authors"`
}
