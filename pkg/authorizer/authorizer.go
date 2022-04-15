package authorizer

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/reading-tribe/anansi/pkg/idx"
	"github.com/reading-tribe/anansi/pkg/repository"
	"github.com/reading-tribe/anansi/pkg/timex"
)

func generatePolicy(principalId, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalId}
	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	return authResponse
}

func GetAuthorizer() func(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	return authorize
}

func authorize(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	sessionRepo := repository.NewSessionRepository()
	token := event.AuthorizationToken

	sessionID := idx.SessionID(token)

	if validationErr := sessionID.Validate(); validationErr != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token")
	}

	session, err := sessionRepo.GetSession(ctx, sessionID)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token")
	}

	userID := idx.UserID(session.UserID)
	if userIDValidationErr := userID.Validate(); userIDValidationErr != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid user ID associated with session")
	}

	if session.ExpiresAtUnix < timex.GetCurrentUTCUnixNano() {
		return generatePolicy("user", "Allow", event.MethodArn), nil
	}

	_ = sessionRepo.DeleteSession(ctx, sessionID)

	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
}
