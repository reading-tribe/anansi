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

const FuncName = "Anansi.API-Translation.FuncUpdateTranslation"

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
		localLogger.Error("Error occurred while updating translation", notOkErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, notOkErr
	}

	idx := idx.TranslationID(id)

	if validationErr := idx.Validate(); validationErr != nil {
		localLogger.Error("Invalid translation ID", idx.String())
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, validationErr.GetError()
	}

	var parsedRequest nettypes.UpdateTranslationRequest
	parseError := json.Unmarshal([]byte(request.Body), &parsedRequest)
	if parseError != nil {
		localLogger.Error("Error occurred while trying to parse body as nettypes.UpdateTranslationRequest", parseError)
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

	updatedTranslation := dbmodel.Translation{
		ID:             idx,
		BookID:         parsedRequest.BookID,
		LocalisedTitle: parsedRequest.LocalisedTitle,
		Language:       lang,
	}

	TranslationRepository := repository.NewTranslationRepository()

	TranslationUpdateErr := TranslationRepository.UpdateTranslation(ctx, updatedTranslation)
	if TranslationUpdateErr != nil {
		localLogger.Error("Error occurred while updating Translation", TranslationUpdateErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, TranslationUpdateErr
	}

	responseBody := nettypes.UpdateTranslationResponse{
		ID:             updatedTranslation.ID,
		BookID:         updatedTranslation.BookID,
		LocalisedTitle: updatedTranslation.LocalisedTitle,
		Language:       updatedTranslation.Language,
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
