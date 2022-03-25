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

const FuncName = "Anansi.API-Translation.FuncCreateTranslation"

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

	var parsedRequest nettypes.CreateTranslationRequest
	parseError := json.Unmarshal([]byte(request.Body), &parsedRequest)
	if parseError != nil {
		localLogger.Error("Error occurred while trying to parse body as nettypes.CreateTranslationRequest", parseError)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, parseError
	}

	var lang = dbmodel.Language(parsedRequest.Language)
	if validationErr := lang.Validate(); validationErr != nil {
		localLogger.Error("Invalid language specified", validationErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, validationErr
	}

	newTranslation := dbmodel.Translation{
		BookID:         parsedRequest.BookID,
		LocalisedTitle: parsedRequest.LocalisedTitle,
		Language:       lang,
	}

	translationRepository := repository.NewTranslationRepository()

	createdTranslation, TranslationCreationErr := translationRepository.CreateTranslation(ctx, newTranslation)
	if TranslationCreationErr != nil {
		localLogger.Error("Error occurred while creating translation", TranslationCreationErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, TranslationCreationErr
	}

	responseBody := nettypes.CreateTranslationResponse{
		ID:             createdTranslation.ID,
		BookID:         createdTranslation.BookID,
		LocalisedTitle: createdTranslation.LocalisedTitle,
		Language:       createdTranslation.Language,
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
