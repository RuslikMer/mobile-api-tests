package pages

import (
	categoryv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/catalog/category/v1"
	productv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/catalog/product/v1"
	productFilterv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/catalog/productfilter/v1"
	"context"
	"github.com/dailymotion/allure-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"
)

type ProductFilterClient struct {
	client productFilterv1.ProductfilterAPIClient
	ctx    context.Context
}

func NewProductFilterClient(ctx context.Context, connection *grpc.ClientConn) *ProductFilterClient {
	return &ProductFilterClient{client: productFilterv1.NewProductfilterAPIClient(connection), ctx: ctx}
}

type CategoryClient struct {
	client categoryv1.CategoryAPIClient
	ctx    context.Context
}

func NewCategoryClient(ctx context.Context, connection *grpc.ClientConn) *CategoryClient {
	return &CategoryClient{client: categoryv1.NewCategoryAPIClient(connection), ctx: ctx}
}

func FastSearch(t *testing.T, productFilterClient *ProductFilterClient, searchText string) (string, int64, []*categoryv1.CategoryWithParents) {
	var id string
	var price int64
	var categories []*categoryv1.CategoryWithParents

	allure.Step(allure.Description("fast search"),
		allure.Action(func() {
			productFilterResponse, err := productFilterClient.client.FastSearch(productFilterClient.ctx, &productFilterv1.FastSearchRequest{
				Text: searchText,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			id = productFilterResponse.GetProducts()[0].GetId()
			price = productFilterResponse.GetProducts()[0].GetPrice()
			categories = productFilterResponse.GetCategories()
			assert.NotEmpty(t, id)
			assert.NotEmpty(t, price)
		}),
	)

	return id, price, categories
}

func Filter(t *testing.T, productFilterClient *ProductFilterClient, searchText string) []*productv1.ShortProduct {
	var products []*productv1.ShortProduct

	allure.Step(allure.Description("get a filtered list of products"),
		allure.Action(func() {
			productFilterResponse, err := productFilterClient.client.Filter(productFilterClient.ctx, &productFilterv1.FilterRequest{
				CompilationId: "",
				FilterIds:     nil,
				SortId:        "",
				PriceFrom:     0,
				PriceTo:       100000,
				Offset:        1,
				Limit:         100,
				SearchText:    searchText,
				IsDiscount:    false,
				Flags:         nil,
				IsPromo:       false,
				PromoIds:      nil,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			products = productFilterResponse.GetProducts()
			assert.NotEmpty(t, products)
		}),
	)

	return products
}

func GetAllCategories(t *testing.T, categoryClient *CategoryClient) []*categoryv1.Category {
	var categories []*categoryv1.Category

	allure.Step(allure.Description("get all categories"),
		allure.Action(func() {
			categoriesResponse, err := categoryClient.client.GetAll(categoryClient.ctx, &categoryv1.GetAllRequest{})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			categories = categoriesResponse.GetCategories()
			assert.NotEmpty(t, categories)
		}),
	)

	return categories
}

func GetItemThatNotInStock(t *testing.T, productFilterClient *ProductFilterClient, searchText string) (string, int64) {
	var id string
	var price int64

	allure.Step(allure.Description("get item that not stock"),
		allure.Action(func() {
			products := Filter(t, productFilterClient, searchText)
			for _, product := range products {
				if product.GetAvailableAmount() == 0 {
					id = product.GetId()
					price = product.GetPrice()
					break
				}
			}
		}),
	)

	return id, price
}

func GetRandomCategory(t *testing.T, categoryClient *CategoryClient) string {
	var categoriesName string
	var randCategoryName string

	allure.Step(allure.Description("get random category"),
		allure.Action(func() {
			categories := GetAllCategories(t, categoryClient)
			for i := range categories {
				categoriesName += categories[i].GetName()
				categoriesName += ","
			}

			assert.NotEmpty(t, categoriesName, "Отсутствуют категории")
			sliceCategoriesName := strings.Split(categoriesName, ",")
			rand.Seed(time.Now().UnixNano())
			randCategoryName = sliceCategoriesName[rand.Intn(len(sliceCategoriesName)-1)]
			assert.NotEmpty(t, randCategoryName, "Отсутствует наименование категории")
		}),
	)

	return randCategoryName
}
