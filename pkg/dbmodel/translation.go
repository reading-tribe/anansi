package dbmodel

import "github.com/reading-tribe/anansi/pkg/idx"

type Translation struct {
	ID             idx.TranslationID `dynamodbav:"id" json:"id"`
	BookID         idx.BookID        `dynamodbav:"book_id" json:"book_id"`
	LocalisedTitle string            `dynamodbav:"localised_title" json:"localised_title"`
	Language       Language          `dynamodbav:"lang" json:"language"`
}
