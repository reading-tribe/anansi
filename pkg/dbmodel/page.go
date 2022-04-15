package dbmodel

import "github.com/reading-tribe/anansi/pkg/idx"

type Page struct {
	ID            idx.PageID        `dynamodbav:"id"`
	ImageURL      string            `dynamodbav:"image_url"`
	Number        int               `dynamodbav:"number"`
	TranslationID idx.TranslationID `dynamodbav:"translation_id"`
}
