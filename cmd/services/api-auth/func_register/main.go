package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/reading-tribe/anansi/pkg/cryptography"
	"github.com/reading-tribe/anansi/pkg/dbmodel"
	"github.com/reading-tribe/anansi/pkg/headers"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/reading-tribe/anansi/pkg/repository"
	"github.com/reading-tribe/anansi/pkg/sqsx"
	"github.com/reading-tribe/anansi/pkg/timex"
)

const FuncName = "Anansi.API-Auth.FuncRegister"

var InviteUsersQueueUrl = os.Getenv("INVITE_USER_QUEUE_URL")

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

	var parsedRequest nettypes.RegisterRequest
	parseError := json.Unmarshal([]byte(request.Body), &parsedRequest)
	if parseError != nil {
		localLogger.Error("Error occurred while trying to parse body as nettypes.RegisterRequest", parseError)
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

	sqsClient, _ := sqsx.GetClient()

	inviteUsersRequest := nettypes.InviteUserRequest{
		EmailAddress: parsedRequest.EmailAddress,
	}

	sqsMessageBody, marshalJSONErr := json.Marshal(inviteUsersRequest)
	if marshalJSONErr != nil {
		localLogger.Error("Error occurred while trying to marshal JSON for SQS", marshalJSONErr)
	}

	sqsOutput, sqsErr := sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody:            aws.String(string(sqsMessageBody)),
		MessageDeduplicationId: aws.String(strconv.FormatInt(timex.GetCurrentUTCUnixNano(), 10)),
		MessageGroupId:         aws.String(parsedRequest.EmailAddress),
		QueueUrl:               aws.String(InviteUsersQueueUrl),
	})
	if sqsErr != nil {
		localLogger.Error("Error occurred while trying to send SQS message", sqsErr)
	}

	localLogger.Debug("Got SQS Output", sqsOutput)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers:    headers.NewMapHeader().ContentTypeJSON().GetMap(),
		Body:       string(responseJSON),
	}, nil
}
