package pages

import (
	productv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/catalog/product/v1"
	homepagev1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/homepage/v1"
	"context"
	"github.com/dailymotion/allure-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"testing"
)

type HomepageClient struct {
	client homepagev1.HomepageAPIClient
	ctx    context.Context
}

func NewHomepageClient(ctx context.Context, connection *grpc.ClientConn) *HomepageClient {
	return &HomepageClient{client: homepagev1.NewHomepageAPIClient(connection), ctx: ctx}
}

func GetItem(t *testing.T, homepageClient *HomepageClient) (string, int64) {
	var id string
	var price int64

	allure.Step(allure.Description("get item"),
		allure.Action(func() {
			homepageResponse, err := homepageClient.client.GetHomepagePage(homepageClient.ctx, &homepagev1.GetHomepagePageRequest{
				PageNumber: 1,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			for _, item := range homepageResponse.GetPage().GetItems() {
				if item != nil {
					for _, product := range item.GetProducts().GetProducts() {
						if product != nil {
							id = product.GetId()
							price = product.GetPrice()
						}
					}
				}
			}

			assert.NotEmpty(t, id)
			assert.NotEmpty(t, price)
		}),
	)

	return id, price
}

func GetLastSeenProducts(t *testing.T, homepageClient *HomepageClient) []*productv1.ShortProduct {
	var products []*productv1.ShortProduct

	allure.Step(allure.Description("get last seen product"),
		allure.Action(func() {
			lastSeenProductsResponse, err := homepageClient.client.GetLastSeenProducts(homepageClient.ctx, &homepagev1.GetLastSeenProductsRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			products = lastSeenProductsResponse.GetLastSeen()
			assert.NotEmpty(t, products)
		}),
	)

	return products
}

func CheckLastSeenProducts(t *testing.T, homepageClient *HomepageClient, id string) {
	allure.Step(allure.Description("check last seen product"),
		allure.Action(func() {
			products := GetLastSeenProducts(t, homepageClient)
			assert.NotEmpty(t, products[0].GetId())
			assert.Equal(t, products[0].GetId(), id)
			assert.NotEmpty(t, products[0].GetPrice())
		}),
	)
}
