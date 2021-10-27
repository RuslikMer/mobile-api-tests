package pages

import (
	authv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/auth/v1"
	"context"
	"github.com/dailymotion/allure-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"testing"
)

type AuthClient struct {
	client authv1.AuthAPIClient
	ctx    context.Context
}

func NewAuthClient(ctx context.Context, connection *grpc.ClientConn) *AuthClient {
	return &AuthClient{client: authv1.NewAuthAPIClient(connection), ctx: ctx}
}

func AuthByEmail(t *testing.T, authClient *AuthClient, email string, password string) *authv1.AuthResponse {
	var authResponse *authv1.AuthResponse

	allure.Step(allure.Description("auth by email"),
		allure.Action(func() {
			authByEmailResponse, err := authClient.client.AuthByPassword(authClient.ctx, &authv1.AuthByPasswordRequest{
				Credentials: &authv1.AuthByPasswordRequest_Email{
					Email: email,
				},
				Password: password,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			authResponse = authByEmailResponse.GetResponse()
			assert.NotEmpty(t, authResponse)
		}),
	)

	return authResponse
}

func Logout(t *testing.T, authClient *AuthClient) *authv1.AuthResponse {
	var response *authv1.AuthResponse

	allure.Step(allure.Description("logout"),
		allure.Action(func() {
			logoutResponse, err := authClient.client.Logout(authClient.ctx, &authv1.LogoutRequest{})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			response = logoutResponse.GetResponse()
			assert.NotEmpty(t, response)
		}),
	)

	return response
}
