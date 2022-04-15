package dbmodel

import "github.com/reading-tribe/anansi/pkg/idx"

type BookDiversityTags struct {
	BookID                           idx.BookID                  `dynamodbav:"book_id"`
	DiversityAndInclusionCatalogueID idx.DiversityAndInclusionID `dynamodbav:"dei_id"`
}
