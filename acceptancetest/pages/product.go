package pages

import (
	productv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/catalog/product/v1"
	opinionv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/opinion/v1"
	questionv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/question/v1"
	"context"
	"github.com/dailymotion/allure-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"testing"
)

type ProductClient struct {
	client productv1.ProductAPIClient
	ctx    context.Context
}

func NewProductClient(ctx context.Context, connection *grpc.ClientConn) *ProductClient {
	return &ProductClient{client: productv1.NewProductAPIClient(connection), ctx: ctx}
}

func SelectService(t *testing.T, productClient *ProductClient, productId []string, serviceType productv1.ProductServiceInfo_GroupId) (string, int64) {
	var serviceId string
	var servicePrice int64

	allure.Step(allure.Description("select service"),
		allure.Action(func() {
			selectServiceResponse, err := productClient.client.GetServices(productClient.ctx, &productv1.GetServicesRequest{
				ProductIds: productId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			services := selectServiceResponse.GetProductServices()

			for i, service := range services {
				if service.GetService().GetGroup().GetGroupId() == serviceType {
					serviceId = services[i].GetService().GetId()
					servicePrice = services[i].GetService().GetPrice()
					break
				}
			}

			assert.NotEmpty(t, serviceId)
			assert.NotEmpty(t, servicePrice)
		}),
	)

	return serviceId, servicePrice
}

func GetItemData(t *testing.T, productClient *ProductClient, id string) *productv1.Product {
	var itemData *productv1.Product

	allure.Step(allure.Description("get item data"),
		allure.Action(func() {
			productResponse, err := productClient.client.Get(productClient.ctx, &productv1.GetRequest{
				Id: id,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			itemData = productResponse.GetProduct()
			assert.NotEmpty(t, itemData)
		}),
	)

	return itemData
}

func NotifyAboutProductArrival(t *testing.T, productClient *ProductClient, id string, email string) {
	allure.Step(allure.Description("notify about product arrival"),
		allure.Action(func() {
			_, err := productClient.client.NotifyAboutProductArrival(productClient.ctx, &productv1.NotifyAboutProductArrivalRequest{
				ProductId: id,
				Email:     email,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
		}),
	)
}

func GetOpinions(t *testing.T, productClient *ProductClient, id string, sortId string) []*opinionv1.Opinion {
	var opinions []*opinionv1.Opinion

	allure.Step(allure.Description("get opinions"),
		allure.Action(func() {
			productResponse, err := productClient.client.GetOpinions(productClient.ctx, &productv1.GetOpinionsRequest{
				ProductId: id,
				SortId:    sortId,
				Offset:    10,
				Limit:     10,
			})

			if err != nil {
				log.Println(err)
			}
			assert.Nil(t, err)
			opinions = productResponse.GetOpinions()
			assert.NotEmpty(t, opinions)
		}),
	)

	return opinions
}

func GetQuestions(t *testing.T, productClient *ProductClient, id string) []*questionv1.Question {
	var questions []*questionv1.Question

	allure.Step(allure.Description("get questions"),
		allure.Action(func() {
			productResponse, err := productClient.client.GetQuestions(productClient.ctx, &productv1.GetQuestionsRequest{
				ProductId: id,
				Offset:    10,
				Limit:     10,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			questions = productResponse.GetQuestions()
			assert.NotEmpty(t, questions)
		}),
	)

	return questions
}

func GetVideoReviews(t *testing.T, productClient *ProductClient, id string) []*productv1.VideoReview {
	var videoReviews []*productv1.VideoReview

	allure.Step(allure.Description("get video reviews"),
		allure.Action(func() {
			productResponse, err := productClient.client.GetVideoReviews(productClient.ctx, &productv1.GetVideoReviewsRequest{
				ProductId: id,
				Offset:    10,
				Limit:     10,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			videoReviews = productResponse.GetVideos()
			assert.NotEmpty(t, videoReviews)
		}),
	)

	return videoReviews
}

func GetRecommendedProducts(t *testing.T, productClient *ProductClient, id string) []*productv1.ShortProduct {
	var products []*productv1.ShortProduct

	allure.Step(allure.Description("get recommended products"),
		allure.Action(func() {
			productResponse, err := productClient.client.GetRecommendedProducts(productClient.ctx, &productv1.GetRecommendedProductsRequest{
				ProductId: id,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			products = productResponse.GetBlock()[1].GetProducts()
			assert.NotEmpty(t, products)
		}),
	)

	return products
}

func GetAccessories(t *testing.T, productClient *ProductClient, id string) []*productv1.ShortProduct {
	var products []*productv1.ShortProduct

	allure.Step(allure.Description("get accessories"),
		allure.Action(func() {
			productResponse, err := productClient.client.GetAccessories(productClient.ctx, &productv1.GetAccessoriesRequest{
				ProductId: id,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			products = productResponse.GetAccessories()
			assert.NotEmpty(t, products)
		}),
	)

	return products
}

func GetPromoKits(t *testing.T, productClient *ProductClient, id string) []*productv1.PromoKit {
	var promoKit []*productv1.PromoKit

	allure.Step(allure.Description("get promo kits"),
		allure.Action(func() {
			productResponse, err := productClient.client.GetPromoKits(productClient.ctx, &productv1.GetPromoKitsRequest{
				ProductId: id,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			promoKit = productResponse.GetPromoKits()
			assert.NotEmpty(t, promoKit)
		}),
	)

	return promoKit
}

func CheckProductInStock(t *testing.T, productClient *ProductClient, id string) {
	allure.Step(allure.Description("check product in stock"),
		allure.Action(func() {
			itemData := GetItemData(t, productClient, id)
			quantityStock := itemData.GetAvailableAmount()
			assert.NotEqual(t, int32(0), quantityStock)
		}),
	)
}

func CheckGetOpinionsProduct(t *testing.T, productClient *ProductClient, products []*productv1.ShortProduct) {
	allure.Step(allure.Description("check get opinions product"),
		allure.Action(func() {
			for _, product := range products {
				if product.GetCountOpinions() > 0 {
					GetOpinions(t, productClient, product.GetId(), "")
					break
				}
			}
		}),
	)
}

func CheckGetQuestionsProduct(t *testing.T, productClient *ProductClient, products []*productv1.ShortProduct) {
	allure.Step(allure.Description("check get questions product"),
		allure.Action(func() {
			for _, product := range products {
				itemData := GetItemData(t, productClient, product.GetId())
				if itemData.GetCountQuestionsAnswers() > 0 {
					GetQuestions(t, productClient, itemData.GetId())
					break
				}
			}
		}),
	)
}

func CheckGetVideoReviews(t *testing.T, productClient *ProductClient, products []*productv1.ShortProduct) {
	allure.Step(allure.Description("check get video reviews"),
		allure.Action(func() {
			for _, product := range products {
				itemData := GetItemData(t, productClient, product.GetId())
				if len(itemData.GetVideos()) > 0 {
					GetVideoReviews(t, productClient, itemData.GetId())
					break
				}
			}
		}),
	)
}
