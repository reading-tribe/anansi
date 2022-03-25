package dbmodel

type DiversityAndInclusionCatalogue struct {
	ID          string `dynamodbav:"id"`
	Key         string `dynamodbav:"key"`
	Value       string `dynamodbav:"value"`
	Description string `dynamodbav:"description"`
}
