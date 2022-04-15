package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reading-tribe/anansi/pkg/cryptography"
	"github.com/reading-tribe/anansi/pkg/dbmodel"
	"github.com/reading-tribe/anansi/pkg/headers"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/reading-tribe/anansi/pkg/repository"
)

const FuncName = "Anansi.API-User.FuncCreateUser"

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

	var parsedRequest nettypes.CreateUserRequest
	parseError := json.Unmarshal([]byte(request.Body), &parsedRequest)
	if parseError != nil {
		localLogger.Error("Error occurred while trying to parse body as nettypes.CreateUserResponse", parseError)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, parseError
	}

	passwordBytes := []byte(parsedRequest.Password)
	hashedPassword := cryptography.HashAndSalt(passwordBytes)

	userRepository := repository.NewUserRepository()

	createErr := userRepository.CreateUser(ctx, dbmodel.User{
		ID:           "",
		EmailAddress: parsedRequest.EmailAddress,
		PasswordHash: hashedPassword,
	})
	if createErr != nil {
		localLogger.Error("Error occurred while trying to create user", createErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, createErr
	}

	responseBody := nettypes.RegisterResponse{
		Message: "Successfully created user",
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
