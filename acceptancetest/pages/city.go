package pages

import (
	cityv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/city/v1"
	"context"
	"github.com/dailymotion/allure-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"testing"
)

type CityClient struct {
	client cityv1.CityAPIClient
	ctx    context.Context
}

func NewCityClient(ctx context.Context, connection *grpc.ClientConn) *CityClient {
	return &CityClient{client: cityv1.NewCityAPIClient(connection), ctx: ctx}
}

func GetAllCities(t *testing.T, cityClient *CityClient, searchCity string) ([]*cityv1.City, string) {
	var allCities []*cityv1.City
	var cityId string

	allure.Step(allure.Description("get all cities"),
		allure.Action(func() {
			allCitiesResponse, err := cityClient.client.GetAll(cityClient.ctx, &cityv1.GetAllRequest{})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			allCities = allCitiesResponse.GetCities()
			assert.NotEmpty(t, allCities)
			for _, city := range allCities {
				if city.GetName() == searchCity {
					cityId = city.GetId()
					break
				}
			}
		}),
	)

	return allCities, cityId
}
