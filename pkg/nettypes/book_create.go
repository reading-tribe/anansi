package nettypes

import "github.com/reading-tribe/anansi/pkg/dbmodel"

type CreateBookRequest struct {
	InternalTitle string `json:"internal_title"`
	Authors       string `json:"authors"`
	Translations  []struct {
		LocalisedTitle string           `json:"localised_title"`
		Language       dbmodel.Language `json:"lang"`
		Pages          []struct {
			ImageURL string `json:"image_url"`
			Number   int    `json:"page_number"`
		} `json:"pages"`
	} `json:"translations"`
}

type CreateBookResponse GetBookResponse
