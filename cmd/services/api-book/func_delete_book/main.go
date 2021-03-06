package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reading-tribe/anansi/pkg/headers"
	"github.com/reading-tribe/anansi/pkg/idx"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/repository"
)

const FuncName = "Anansi.API-Book.FuncDeleteBook"

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

	id, ok := request.PathParameters["id"]
	if !ok {
		notOkErr := fmt.Errorf("unable to read id from path")
		localLogger.Error("Error occurred while deleting book", notOkErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, notOkErr
	}

	idx := idx.BookID(id)

	if validationErr := idx.Validate(); validationErr != nil {
		localLogger.Error("Invalid book ID", idx.String())
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, validationErr.GetError()
	}

	bookRepository := repository.NewBookRepository()

	translations, listTranslationsErr := translationRepository.ListTranslationsByBookID(ctx, idx)
	if listTranslationsErr != nil {
		localLogger.Error("Error occurred while listing translations", listTranslationsErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, listTranslationsErr
	}

	for _, translation := range translations {
		pages, listPagesErr := pageRepository.ListPagesByTranslationID(ctx, translation.ID)
		if listPagesErr != nil {
			localLogger.Error("Error occurred while listing pages", listPagesErr)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
			}, listPagesErr
		}

		for _, page := range pages {
			pageDeletionErr := pageRepository.DeletePage(ctx, page.ID)
			if pageDeletionErr != nil {
				localLogger.Error("Error occurred while deleting page", pageDeletionErr)
				return events.APIGatewayV2HTTPResponse{
					StatusCode: http.StatusInternalServerError,
				}, pageDeletionErr
			}
		}

		translationDeletionErr := translationRepository.DeleteTranslation(ctx, translation.ID)
		if translationDeletionErr != nil {
			localLogger.Error("Error occurred while deleting translation", translationDeletionErr)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
			}, translationDeletionErr
		}
	}

	bookDeletionErr := bookRepository.DeleteBook(ctx, idx)
	if bookDeletionErr != nil {
		localLogger.Error("Error occurred while deleting book", bookDeletionErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, bookDeletionErr
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers:    headers.NewMapHeader().ContentTypeJSON().GetMap(),
	}, nil
}
