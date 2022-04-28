package dbmodel

import (
	"fmt"
	"strconv"

	"github.com/reading-tribe/anansi/pkg/idx"
)

type DynamoPage struct {
	ID            idx.PageID        `dynamodbav:"id"`
	ImageURL      string            `dynamodbav:"image_url"`
	Number        string            `dynamodbav:"page_number"`
	TranslationID idx.TranslationID `dynamodbav:"translation_id"`
}

func (p DynamoPage) ToPage() Page {
	number, conversionErr := strconv.Atoi(p.Number)
	if conversionErr != nil {
		fmt.Println(conversionErr)
		number = -1
	}
	return Page{
		ID:            p.ID,
		ImageURL:      p.ImageURL,
		Number:        number,
		TranslationID: p.TranslationID}
}

type Page struct {
	ID            idx.PageID        `dynamodbav:"id"`
	ImageURL      string            `dynamodbav:"image_url"`
	Number        int               `dynamodbav:"page_number" json:"page_number,string"`
	TranslationID idx.TranslationID `dynamodbav:"translation_id"`
}
