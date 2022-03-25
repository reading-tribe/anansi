package nettypes

import "github.com/reading-tribe/anansi/pkg/dbmodel"

type UpdateBookRequest struct {
	InternalTitle string `json:"internal_title"`
	Authors       string `json:"authors"`
}

type UpdateBookResponse dbmodel.Book
