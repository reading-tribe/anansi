package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/reading-tribe/anansi/pkg/emailx"
	"github.com/reading-tribe/anansi/pkg/logging"
	"github.com/reading-tribe/anansi/pkg/nettypes"
	"github.com/reading-tribe/anansi/pkg/repository"
	"github.com/reading-tribe/anansi/pkg/sesx"
)

const FuncName = "Anansi.API-Auth.FuncInviteUser"

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.SQSEvent) error {
	localLogger := logging.NewLogger(map[string]interface{}{
		logging.FunctionName:      FuncName,
		logging.RequestIdentifier: logging.UniqueRequestName(),
		logging.FieldEvent:        request,
	})
	localLogger.Info("SQS event received!")

	for _, event := range request.Records {
		var parsedRequest nettypes.InviteUserRequest
		parseError := json.Unmarshal([]byte(event.Body), &parsedRequest)
		if parseError != nil {
			localLogger.Error("Error occurred while trying to parse body as nettypes.InviteUserRequest", parseError)
			continue
		}

		client, getClientErr := sesx.GetClient()
		if getClientErr != nil {
			localLogger.Error("Error occurred while trying to instantiate SES client", getClientErr)
			continue
		}

		emailAddress := parsedRequest.EmailAddress
		if emailAddress == "" {
			localLogger.Error("Invalid email address")
			continue
		}

		userRepository := repository.NewUserRepository()

		user, getUserErr := userRepository.GetUserByEmailAddress(ctx, parsedRequest.EmailAddress)
		if getUserErr != nil {
			localLogger.Error("Error occurred while trying to get user", getUserErr)
			continue
		}

		subject, htmlBody, textBody := emailx.InviteUserEmail(fmt.Sprintf("https://81qrzgok36.execute-api.eu-central-1.amazonaws.com/auth/accept/%s", user.InviteCode))

		// Assemble the email.
		input := &ses.SendEmailInput{
			Destination: &ses.Destination{
				CcAddresses: []*string{},
				ToAddresses: []*string{
					aws.String(emailx.Sender), //aws.String(emailAddress),
				},
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Html: &ses.Content{
						Charset: aws.String(emailx.CharSet),
						Data:    aws.String(htmlBody),
					},
					Text: &ses.Content{
						Charset: aws.String(emailx.CharSet),
						Data:    aws.String(textBody),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String(emailx.CharSet),
					Data:    aws.String(subject),
				},
			},
			Source: aws.String(emailx.Sender),
		}

		// Attempt to send the email.
		result, err := client.SendEmail(input)

		// Display error messages if they occur.
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case ses.ErrCodeMessageRejected:
					localLogger.Error(ses.ErrCodeMessageRejected, aerr.Error())
				case ses.ErrCodeMailFromDomainNotVerifiedException:
					localLogger.Error(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
				case ses.ErrCodeConfigurationSetDoesNotExistException:
					localLogger.Error(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
				default:
					localLogger.Error(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				localLogger.Error(err.Error())
			}
			continue
		}

		localLogger.Info("Email Sent to address: " + emailAddress)
		localLogger.Debug(result)
	}

	return nil
}
