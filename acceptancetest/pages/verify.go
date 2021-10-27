package pages

import (
	profilev1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/profile/v1"
	verifyv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/verify/v1"
	"context"
	"github.com/dailymotion/allure-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"testing"
)

type VerifyClient struct {
	client verifyv1.VerifyAPIClient
	ctx    context.Context
}

func NewVerifyClient(ctx context.Context, connection *grpc.ClientConn) *VerifyClient {
	return &VerifyClient{client: verifyv1.NewVerifyAPIClient(connection), ctx: ctx}
}

func SendSms(t *testing.T, verifyClient *VerifyClient, operation verifyv1.SendSmsRequest_Operation, phone string) string {
	var requestId string

	allure.Step(allure.Description("send sms"),
		allure.Action(func() {
			sendSmsResponse, err := verifyClient.client.SendSms(verifyClient.ctx, &verifyv1.SendSmsRequest{
				Operation: operation,
				Phone:     phone,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			requestId = sendSmsResponse.GetRequestId()
			assert.NotEmpty(t, requestId)
		}),
	)

	return requestId
}

func SmsCheck(t *testing.T, verifyClient *VerifyClient, requestId string, code string) string {
	var token string

	allure.Step(allure.Description("sms check"),
		allure.Action(func() {
			smsCheckResponse, err := verifyClient.client.Check(verifyClient.ctx, &verifyv1.CheckRequest{
				RequestId: requestId,
				Code:      code,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			token = smsCheckResponse.GetVerification().GetToken()
			assert.True(t, true, token)
		}),
	)

	return token
}

func SendVerificationSms(t *testing.T, verifyClient *VerifyClient, operation verifyv1.SendVerificationSmsRequest_Operation) string {
	var verificationRequestId string

	allure.Step(allure.Description("check verification sms"),
		allure.Action(func() {
			sendVerificationSmsResponse, err := verifyClient.client.SendVerificationSms(verifyClient.ctx, &verifyv1.SendVerificationSmsRequest{
				Operation: operation,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			verificationRequestId = sendVerificationSmsResponse.GetVerificationRequestId()
			assert.NotEmpty(t, verificationRequestId)
		}),
	)

	return verificationRequestId
}

func CheckSmsVerification(t *testing.T, verifyClient *VerifyClient, verificationRequestId string, code string) string {
	var verificationToken string

	allure.Step(allure.Description("send verification sms"),
		allure.Action(func() {
			checkSmsVerificationResponse, err := verifyClient.client.CheckSmsVerification(verifyClient.ctx, &verifyv1.CheckSmsVerificationRequest{
				VerificationRequestId: verificationRequestId,
				Code:                  code,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			verificationToken = checkSmsVerificationResponse.GetVerification().GetToken()
			assert.True(t, true, verificationToken)
		}),
	)

	return verificationToken
}

func SetPhone(t *testing.T, profileClient *ProfileClient, token string) {
	allure.Step(allure.Description("set phone"),
		allure.Action(func() {
			setResponse, err := profileClient.client.SetPhone(profileClient.ctx, &profilev1.SetPhoneRequest{
				VerifyToken: &verifyv1.Verification{Token: token},
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			set := setResponse.GetProfile().GetIsPhoneConfirmed()
			assert.True(t, set)
		}),
	)
}
