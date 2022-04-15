package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/reading-tribe/anansi/pkg/repository"
)

const FuncName = "Anansi.API-Auth.FuncLogout"

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

	var parsedRequest nettypes.LogoutRequest
	parseError := json.Unmarshal([]byte(request.Body), &parsedRequest)
	if parseError != nil {
		localLogger.Error("Error occurred while trying to parse body as nettypes.LogoutRequest", parseError)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, parseError
	}

	userRepository := repository.NewUserRepository()

	_, getUserErr := userRepository.GetUserByEmailAddress(ctx, parsedRequest.EmailAddress)
	if getUserErr != nil {
		localLogger.Error("Error occurred while trying to get user", getUserErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, getUserErr
	}

	sessionRepository := repository.NewSessionRepository()

	deleteSessionErr := sessionRepository.DeleteSession(ctx, "")
	if deleteSessionErr != nil {
		localLogger.Error("Error occurred while trying to delete session", deleteSessionErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, deleteSessionErr
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
	}, nil
}
