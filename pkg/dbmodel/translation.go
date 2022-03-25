package dbmodel

type Translation struct {
	ID             string   `dynamodbav:"id" json:"id"`
	BookID         string   `dynamodbav:"book_id" json:"book_id"`
	LocalisedTitle string   `dynamodbav:"localised_title" json:"localised_title"`
	Language       Language `dynamodbav:"lang" json:"language"`
}
