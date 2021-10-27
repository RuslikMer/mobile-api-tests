package pages

import (
	promov1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/promo/v1"
	"context"
	"github.com/dailymotion/allure-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"testing"
)

type PromoClient struct {
	client promov1.PromoAPIClient
	ctx    context.Context
}

func NewPromoClient(ctx context.Context, connection *grpc.ClientConn) *PromoClient {
	return &PromoClient{client: promov1.NewPromoAPIClient(connection), ctx: ctx}
}

func GetInfo(t *testing.T, promoClient *PromoClient, id string) *promov1.Promo {
	var promo *promov1.Promo

	allure.Step(allure.Description("get promo info by id"),
		allure.Action(func() {
			promoResponse, err := promoClient.client.Get(promoClient.ctx, &promov1.GetRequest{
				Id: id,
			})

			if err != nil {
				log.Println(err)
				assert.Nil(t, err)
			}

			promo = promoResponse.GetPromo()
			assert.NotEmpty(t, promo)
		}),
	)

	return promo
}

func GetActual(t *testing.T, promoClient *PromoClient) *promov1.Promo {
	var promo *promov1.Promo

	allure.Step(allure.Description("get promo info by type"),
		allure.Action(func() {
			actualResponse, err := promoClient.client.GetActual(promoClient.ctx, &promov1.GetActualRequest{})
			if err != nil {
				log.Println(err)
				assert.Nil(t, err)
			}

			n := rand.Intn(len(actualResponse.GetPromo()))
			promo = actualResponse.GetPromo()[n]
			assert.NotEmpty(t, promo)
		}),
	)

	return promo
}
