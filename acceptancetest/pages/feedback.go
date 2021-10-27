package pages

import (
	captchav1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/captcha/v1"
	feedbackv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/feedback/v1"
	"context"
	"github.com/dailymotion/allure-go"
	"google.golang.org/grpc"
	"log"
	"testing"
)

type FeedbackClient struct {
	client feedbackv1.FeedbackAPIClient
	ctx    context.Context
}

func NewFeedbackClient(ctx context.Context, connection *grpc.ClientConn) *FeedbackClient {
	return &FeedbackClient{client: feedbackv1.NewFeedbackAPIClient(connection), ctx: ctx}
}

func SendFeedback(t *testing.T, feedbackClient *FeedbackClient) {
	allure.Step(allure.Description("send feedback"),
		allure.Action(func() {
			_, err := feedbackClient.client.Send(feedbackClient.ctx, &feedbackv1.SendRequest{
				TopicId: "Test",
				Contact: &feedbackv1.SendRequest_Email{Email: "go-autotester@citilink.ru"},
				Message: "test",
				Name:    "tester",
				Captcha: &captchav1.CaptchaVerify{Token: ""},
			})

			if err != nil {
				log.Println(err)
			}

			//assert.Nil(t, err)
		}),
	)
}
