package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reading-tribe/anansi/pkg/headers"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/reading-tribe/anansi/pkg/repository"
)

const FuncName = "Anansi.API-Book.FuncListBooks"

var (
	translationRepository repository.TranslationRepository
	pageRepository        repository.PageRepository
	bookRepository        repository.BookRepository
)

func main() {
	translationRepository = repository.NewTranslationRepository()
	pageRepository = repository.NewPageRepository()
	bookRepository = repository.NewBookRepository()

	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	localLogger := logging.NewLogger(map[string]interface{}{
		logging.FunctionName:      FuncName,
		logging.RequestIdentifier: logging.UniqueRequestName(),
		logging.FieldEvent:        request,
	})
	localLogger.Info("Request received!")

	books, getBookErr := bookRepository.ListBooks(ctx)
	if getBookErr != nil {
		localLogger.Error("Error occurred while listing books", getBookErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, getBookErr
	}

	responseBody := nettypes.ListBooksResponse{}

	for _, book := range books {
		tempBook := nettypes.GetBookResponse{
			ID:            book.ID,
			InternalTitle: book.InternalTitle,
			Authors:       book.Authors,
			Translations:  []nettypes.GetBookResponse_Translation{},
		}

		translations, listTranslationsErr := translationRepository.ListTranslationsByBookID(ctx, book.ID)
		if listTranslationsErr != nil {
			localLogger.Error("Error occurred while listing translations", listTranslationsErr)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
			}, listTranslationsErr
		}

		for _, translation := range translations {
			tempTranslation := nettypes.GetBookResponse_Translation{
				ID:             translation.ID,
				LocalisedTitle: translation.LocalisedTitle,
				Language:       translation.Language,
				Pages:          []nettypes.GetBookResponse_Translation_Page{},
			}

			pages, listPagesErr := pageRepository.ListPagesByTranslationID(ctx, translation.ID)
			if listPagesErr != nil {
				localLogger.Error("Error occurred while listing pages", listPagesErr)
				return events.APIGatewayV2HTTPResponse{
					StatusCode: http.StatusInternalServerError,
				}, listPagesErr
			}

			for _, page := range pages {
				tempTranslation.Pages = append(tempTranslation.Pages, nettypes.GetBookResponse_Translation_Page{
					ID:       page.ID,
					ImageURL: page.ImageURL,
					Number:   page.Number,
				})
			}

			tempBook.Translations = append(tempBook.Translations, tempTranslation)
		}

		responseBody = append(responseBody, tempBook)
	}

	responseJSON, marshalErr := json.Marshal(responseBody)
	if marshalErr != nil {
		localLogger.Error("Error occurred while trying to marshal response as json", marshalErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, marshalErr
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers:    headers.NewMapHeader().ContentTypeJSON().GetMap(),
		Body:       string(responseJSON),
	}, nil
}
