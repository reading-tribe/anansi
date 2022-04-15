package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reading-tribe/anansi/pkg/idx"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/reading-tribe/anansi/pkg/repository"
	"github.com/reading-tribe/anansi/pkg/timex"
)

const FuncName = "Anansi.API-Auth.FuncRefresh"

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

	var parsedRequest nettypes.RefreshRequest
	parseError := json.Unmarshal([]byte(request.Body), &parsedRequest)
	if parseError != nil {
		localLogger.Error("Error occurred while trying to parse body as nettypes.RefreshRequest", parseError)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, parseError
	}

	key := idx.SessionID(parsedRequest.Key)
	if validateKeyErr := key.Validate(); validateKeyErr != nil {
		localLogger.Error("Invalid session key", validateKeyErr.GetError())
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, validateKeyErr.GetError()
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

	session, getSessionErr := sessionRepository.GetSession(ctx, key)
	if getSessionErr != nil {
		localLogger.Error("Error occurred while trying to get session", getSessionErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, getSessionErr
	}

	session.ExpiresAtUnix = timex.GetFutureUTCUnixNano(timex.ThirtyMinutes())

	updateSessionErr := sessionRepository.UpdateSession(ctx, session)
	if updateSessionErr != nil {
		localLogger.Error("Error occurred while trying to update session", updateSessionErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, updateSessionErr
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
	}, nil
}
