package nettypes

import "github.com/reading-tribe/anansi/pkg/dbmodel"

type CreateBookRequest_Translation_Page struct {
	ImageURL string `json:"image_url"`
	Number   int    `json:"page_number"`
}

type CreateBookRequest_Translation struct {
	LocalisedTitle string                               `json:"localised_title"`
	Language       dbmodel.Language                     `json:"lang"`
	Pages          []CreateBookRequest_Translation_Page `json:"pages"`
}

type CreateBookRequest struct {
	InternalTitle string                          `json:"internal_title"`
	Authors       string                          `json:"authors"`
	Translations  []CreateBookRequest_Translation `json:"translations"`
}

type CreateBookResponse GetBookResponse
