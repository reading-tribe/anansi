package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reading-tribe/anansi/pkg/dbmodel"
	"github.com/reading-tribe/anansi/pkg/headers"
	"github.com/reading-tribe/anansi/pkg/idx"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/reading-tribe/anansi/pkg/repository"
)

const FuncName = "Anansi.API-Book.FuncUpdateBook"

func main() {
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
		localLogger.Error("Error occurred while updating book", notOkErr)
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

	var parsedRequest nettypes.UpdateBookRequest
	parseError := json.Unmarshal([]byte(request.Body), &parsedRequest)
	if parseError != nil {
		localLogger.Error("Error occurred while trying to parse body as nettypes.UpdateBookRequest", parseError)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, parseError
	}

	updatedBook := dbmodel.Book{
		ID:            idx,
		InternalTitle: parsedRequest.InternalTitle,
		Authors:       parsedRequest.Authors,
	}

	bookRepository := repository.NewBookRepository()

	bookUpdateErr := bookRepository.UpdateBook(ctx, updatedBook)
	if bookUpdateErr != nil {
		localLogger.Error("Error occurred while updating book", bookUpdateErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, bookUpdateErr
	}

	responseBody := nettypes.UpdateBookResponse{
		ID:            updatedBook.ID,
		InternalTitle: updatedBook.InternalTitle,
		Authors:       updatedBook.Authors,
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
