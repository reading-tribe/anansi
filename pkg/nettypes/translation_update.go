package nettypes

import "github.com/reading-tribe/anansi/pkg/dbmodel"

type UpdateTranslationRequest struct {
	BookID         string `json:"book_id"`
	LocalisedTitle string `json:"localised_title"`
	Language       string `json:"language"`
}

type UpdateTranslationResponse dbmodel.Translation
