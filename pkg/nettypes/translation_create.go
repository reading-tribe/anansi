package nettypes

import "github.com/reading-tribe/anansi/pkg/dbmodel"

type CreateTranslationRequest struct {
	BookID         string `json:"book_id"`
	LocalisedTitle string `json:"localised_title"`
	Language       string `json:"language"`
}

type CreateTranslationResponse dbmodel.Translation
