package dbmodel

type Page struct {
	ID            string `dynamodbav:"id"`
	ImageURL      string `dynamodbav:"image_url"`
	Number        int    `dynamodbav:"number"`
	TranslationID string `dynamodbav:"translation_id"`
}
