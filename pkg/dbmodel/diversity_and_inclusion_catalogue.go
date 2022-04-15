package dbmodel

import "github.com/reading-tribe/anansi/pkg/idx"

type DiversityAndInclusionCatalogue struct {
	ID          idx.DiversityAndInclusionID `dynamodbav:"id"`
	Key         string                      `dynamodbav:"key"`
	Value       string                      `dynamodbav:"value"`
	Description string                      `dynamodbav:"description"`
}
