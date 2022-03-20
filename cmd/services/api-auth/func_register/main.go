package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reading-tribe/anansi/pkg/cryptography"
	"github.com/reading-tribe/anansi/pkg/dbmodel"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/reading-tribe/anansi/pkg/repository"
	"github.com/sirupsen/logrus"
)

const FuncName = "Anansi.API-Auth.FuncRegister"

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

	var parsedRequest nettypes.RegisterRequest
	parseError := json.Unmarshal([]byte(request.Body), &parsedRequest)
	if parseError != nil {
		logrus.Error("Error occurred while trying to parse body as nettypes.RegisterRequest", parseError)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, parseError
	}

	passwordBytes := []byte(parsedRequest.Password)
	hashedPassword := cryptography.HashAndSalt(passwordBytes)

	userRepository := repository.NewUserRepository()

	createErr := userRepository.CreateUser(ctx, dbmodel.User{
		EmailAddress: parsedRequest.EmailAddress,
		PasswordHash: hashedPassword,
	})
	if createErr != nil {
		logrus.Error("Error occurred while trying to create user", createErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, createErr
	}

	responseBody := nettypes.RegisterResponse{
		Message: "Successfully created user",
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
