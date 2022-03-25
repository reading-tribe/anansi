package authorizer

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
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

	session, err := sessionRepo.GetSession(ctx, string(token))
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token")
	}

	if session.ExpiresAtUnix < timex.GetCurrentUTCUnixNano() {
		return generatePolicy("user", "Allow", event.MethodArn), nil
	}

	_ = sessionRepo.DeleteSession(ctx, token)

	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
}
