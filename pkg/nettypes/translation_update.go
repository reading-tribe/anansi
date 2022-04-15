package nettypes

import (
	"github.com/reading-tribe/anansi/pkg/dbmodel"
	"github.com/reading-tribe/anansi/pkg/idx"
)

type UpdateTranslationRequest struct {
	BookID         idx.BookID `json:"book_id"`
	LocalisedTitle string     `json:"localised_title"`
	Language       string     `json:"language"`
}

type UpdateTranslationResponse dbmodel.Translation
