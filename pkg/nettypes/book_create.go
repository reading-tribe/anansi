package nettypes

import "github.com/reading-tribe/anansi/pkg/dbmodel"

type CreateBookRequest struct {
	InternalTitle string `json:"internal_title"`
	Authors       string `json:"authors"`
}

type CreateBookResponse dbmodel.Book
