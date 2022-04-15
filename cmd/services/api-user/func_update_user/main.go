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
	"github.com/reading-tribe/anansi/pkg/headers"
	"github.com/reading-tribe/anansi/pkg/idx"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/reading-tribe/anansi/pkg/repository"
)

const FuncName = "Anansi.API-Translation.FuncUpdateUser"

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

	idx := idx.UserID(id)

	if validationErr := idx.Validate(); validationErr != nil {
		localLogger.Error("Invalid user ID", idx.String())
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, validationErr.GetError()
	}

	var parsedRequest nettypes.UpdateUserRequest
	parseError := json.Unmarshal([]byte(request.Body), &parsedRequest)
	if parseError != nil {
		localLogger.Error("Error occurred while trying to parse body as nettypes.UpdateUserRequest", parseError)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, parseError
	}

	passwordBytes := []byte(parsedRequest.Password)
	hashedPassword := cryptography.HashAndSalt(passwordBytes)

	userRepository := repository.NewUserRepository()

	foundUser, getUserErr := userRepository.GetUser(ctx, idx)
	if getUserErr != nil {
		localLogger.Error("Error occurred while trying to get user", getUserErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, getUserErr
	}

	updatedUser := dbmodel.User{
		ID:           foundUser.ID,
		EmailAddress: parsedRequest.EmailAddress,
		PasswordHash: hashedPassword,
		InviteCode:   foundUser.InviteCode,
		Confirmed:    parsedRequest.Confirmed,
	}

	updateUserErr := userRepository.UpdateUser(ctx, updatedUser)
	if updateUserErr != nil {
		localLogger.Error("Error occurred while updating user", updateUserErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, updateUserErr
	}

	responseBody := nettypes.UpdateUserResponse{
		ID:           updatedUser.ID,
		EmailAddress: updatedUser.EmailAddress,
		Confirmed:    updatedUser.Confirmed,
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
