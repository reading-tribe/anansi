package dbmodel

type BookDiversityTags struct {
	BookID                           string `dynamodbav:"book_id"`
	DiversityAndInclusionCatalogueID string `dynamodbav:"dei_id"`
}
