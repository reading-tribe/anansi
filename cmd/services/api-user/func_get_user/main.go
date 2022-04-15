package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reading-tribe/anansi/pkg/headers"
	"github.com/reading-tribe/anansi/pkg/idx"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/reading-tribe/anansi/pkg/repository"
)

const FuncName = "Anansi.API-Translation.FuncGetTranslation"

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
		localLogger.Error("Error occurred while getting user", notOkErr)
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

	userRepository := repository.NewUserRepository()

	user, getUserErr := userRepository.GetUser(ctx, idx)
	if getUserErr != nil {
		localLogger.Error("Error occurred while getting user", getUserErr)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, getUserErr
	}

	responseBody := nettypes.GetUserResponse{
		ID:           user.ID,
		EmailAddress: user.EmailAddress,
		Confirmed:    user.Confirmed,
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
