package main

import (
	"context"
	"fmt"
	"net/http"

	//go get -u github.com/aws/aws-sdk-go
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/repository"
)

const FuncName = "Anansi.API-Auth.FuncAcceptInvite"

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

	code, ok := request.PathParameters["code"]
	if !ok {
		notOkErr := fmt.Errorf("unable to read code from path")
		localLogger.Error("Error occurred while accepting invite", notOkErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, notOkErr
	}

	userRepository := repository.NewUserRepository()

	matchingUsers, listUsersErr := userRepository.ListUsersByCode(ctx, code)
	if listUsersErr != nil {
		localLogger.Error("Error occurred while trying to list users", listUsersErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, listUsersErr
	}

	if len(matchingUsers) != 1 {
		err := fmt.Errorf("Unexpected number of matching users %d", len(matchingUsers))
		localLogger.Error("Unexpected number of matching users", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user := matchingUsers[0]

	user.Confirmed = true

	updateErr := userRepository.UpdateUser(ctx, user)
	if updateErr != nil {
		localLogger.Error("Failed to mark user as confirmed", updateErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, updateErr
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
	}, nil
}
