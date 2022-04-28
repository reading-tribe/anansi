package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reading-tribe/anansi/pkg/dbmodel"
	"github.com/reading-tribe/anansi/pkg/headers"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/reading-tribe/anansi/pkg/repository"
)

const FuncName = "Anansi.API-Book.FuncCreateBook"

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

	var parsedRequest nettypes.CreateBookRequest
	parseError := json.Unmarshal([]byte(request.Body), &parsedRequest)
	if parseError != nil {
		localLogger.Error("Error occurred while trying to parse body as nettypes.CreateBookRequest", parseError)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, parseError
	}

	// Create stubbed response
	response := &nettypes.CreateBookResponse{
		ID:            "",
		InternalTitle: "",
		Authors:       "",
		Translations:  []nettypes.GetBookResponse_Translation{},
	}

	// Create Book
	newBook := dbmodel.Book{
		ID:            "",
		InternalTitle: parsedRequest.InternalTitle,
		Authors:       parsedRequest.Authors,
	}

	createdBook, bookCreationErr := bookRepository.CreateBook(ctx, newBook)
	if bookCreationErr != nil {
		localLogger.Error("Error occurred while creating book", bookCreationErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, bookCreationErr
	}

	response.ID = createdBook.ID
	response.InternalTitle = createdBook.InternalTitle
	response.Authors = createdBook.Authors

	// Create Translations
	for _, translation := range parsedRequest.Translations {
		newTranslation := dbmodel.Translation{
			ID:             "",
			BookID:         createdBook.ID,
			LocalisedTitle: translation.LocalisedTitle,
			Language:       translation.Language,
		}

		createdTranslation, translationCreationErr := translationRepository.CreateTranslation(ctx, newTranslation)
		if translationCreationErr != nil {
			localLogger.Error("Error occurred while creating translation", translationCreationErr)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
			}, translationCreationErr
		}

		responseTranslation := nettypes.GetBookResponse_Translation{
			ID:             createdTranslation.ID,
			LocalisedTitle: createdTranslation.LocalisedTitle,
			Language:       createdTranslation.Language,
			Pages:          []nettypes.GetBookResponse_Translation_Page{},
		}

		// Create Pages
		for _, page := range translation.Pages {
			newPage := dbmodel.Page{
				ID:            "",
				ImageURL:      page.ImageURL,
				Number:        page.Number,
				TranslationID: createdTranslation.ID,
			}

			createdPage, pageCreationErr := pageRepository.CreatePage(ctx, newPage)
			if pageCreationErr != nil {
				localLogger.Error("Error occurred while creating page", pageCreationErr)
				return events.APIGatewayV2HTTPResponse{
					StatusCode: http.StatusInternalServerError,
				}, pageCreationErr
			}

			responsePage := nettypes.GetBookResponse_Translation_Page{
				ID:       createdPage.ID,
				ImageURL: createdPage.ImageURL,
				Number:   createdPage.Number,
			}

			responseTranslation.Pages = append(responseTranslation.Pages, responsePage)
		}

		response.Translations = append(response.Translations, responseTranslation)
	}

	responseJSON, marshalErr := json.Marshal(response)
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
