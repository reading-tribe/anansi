package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reading-tribe/anansi/pkg/cryptography"
	"github.com/reading-tribe/anansi/pkg/dbmodel"
	"github.com/reading-tribe/anansi/pkg/idx"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/reading-tribe/anansi/pkg/repository"
	"github.com/reading-tribe/anansi/pkg/timex"
	"github.com/sirupsen/logrus"
)

const FuncName = "Anansi.API-Auth.FuncLogin"

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	logging.SetupLogger(map[string]interface{}{
		logging.FunctionName:      FuncName,
		logging.RequestIdentifier: logging.UniqueRequestName(),
		logging.FieldEvent:        request,
	})
	logrus.Info("Request received!")

	var parsedRequest nettypes.LoginRequest
	parseError := json.Unmarshal([]byte(request.Body), &parsedRequest)
	if parseError != nil {
		logrus.Error("Error occurred while trying to parse body as nettypes.LoginRequest", parseError)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, parseError
	}

	passwordBytes := []byte(parsedRequest.Password)

	userRepository := repository.NewUserRepository()

	user, getUserErr := userRepository.GetUser(ctx, parsedRequest.EmailAddress)
	if getUserErr != nil {
		logrus.Error("Error occurred while trying to get user", getUserErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, getUserErr
	}

	if comparisonErr := cryptography.ComparePasswords(
		user.PasswordHash, passwordBytes,
	); comparisonErr != nil {
		logrus.Error("Password did not match", comparisonErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("Bad login attempt")
	}

	sessionRepository := repository.NewSessionRepository()

	sessionKey, sessionKeyErr := idx.NewSessionID()
	if sessionKeyErr != nil {
		logrus.Error("Error occurred while trying to generate session key", sessionKeyErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, sessionKeyErr
	}

	createSessionErr := sessionRepository.CreateSession(ctx, dbmodel.Session{
		Key:           sessionKey,
		EmailAddress:  parsedRequest.EmailAddress,
		ExpiresAtUnix: timex.GetFutureUTCUnixNano(timex.ThirtyMinutes()),
	})
	if createSessionErr != nil {
		logrus.Error("Error occurred while trying to create session", createSessionErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, createSessionErr
	}

	responseBody := nettypes.LoginResponse{
		Token:   sessionKey,
		Message: "Successfully logged user in",
	}

	responseJSON, marshalErr := json.Marshal(responseBody)
	if marshalErr != nil {
		logrus.Error("Error occurred while trying to marshal response as json", marshalErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, marshalErr
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJSON),
	}, nil
}
