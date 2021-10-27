// +build acceptancetest

package authorized

import (
	"code.citik.ru/mic/mobile-api/acceptancetest"
	"code.citik.ru/mic/mobile-api/acceptancetest/pages"
	"context"
	"github.com/dailymotion/allure-go"
	"testing"
)

func init() {
	initializer = acceptancetest.NewInitializer(context.Background())
	initializer.GlobalInit(true)
}

func TestGetOpinionsAndQuestions(t *testing.T) {
	allure.Test(t, allure.Description("get opinions and questions"),
		allure.Action(func() {
			//получаю список товаров
			products := pages.Filter(t, acceptancetest.CatalogClient, "Apple")

			//проверяю, что у товара есть отзыв
			pages.CheckGetOpinionsProduct(t, acceptancetest.ProductClient, products)

			//проверяю, что у товара есть вопрос
			pages.CheckGetQuestionsProduct(t, acceptancetest.ProductClient, products)
		}),
	)
}

func TestGetVideoReviews(t *testing.T) {
	allure.Test(t, allure.Description("get opinions, questions"),
		allure.Action(func() {
			//получаю список товаров
			products := pages.Filter(t, acceptancetest.CatalogClient, "Ноутбуки")

			//проверяю, что у товара есть видео обзор . Тест будет падать, пока не исправят метод
			pages.CheckGetVideoReviews(t, acceptancetest.ProductClient, products)
		}),
	)
}

func TestGetRecommendedProducts(t *testing.T) {
	allure.Test(t, allure.Description("get recommended products"),
		allure.Action(func() {
			//получаю первый товар c главной страницы
			id, _ := pages.GetItem(t, acceptancetest.HomepageClient)

			//проверяю товары из блока "Рекомендации"
			pages.GetRecommendedProducts(t, acceptancetest.ProductClient, id)
		}),
	)
}
