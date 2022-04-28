package nettypes

import (
	"github.com/reading-tribe/anansi/pkg/dbmodel"
	"github.com/reading-tribe/anansi/pkg/idx"
)

type GetBookResponse_Translation_Page struct {
	ID       idx.PageID `json:"page_id"`
	ImageURL string     `json:"image_url"`
	Number   int        `json:"number"`
}

type GetBookResponse_Translation struct {
	ID             idx.TranslationID                  `json:"id"`
	LocalisedTitle string                             `json:"localised_title"`
	Language       dbmodel.Language                   `json:"lang"`
	Pages          []GetBookResponse_Translation_Page `json:"pages"`
}

type GetBookResponse struct {
	ID            idx.BookID                    `json:"id"`
	InternalTitle string                        `json:"internal_title"`
	Authors       string                        `json:"authors"`
	Translations  []GetBookResponse_Translation `json:"translations"`
}
