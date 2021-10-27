// +build acceptancetest

package unauthorized

import (
	"code.citik.ru/mic/mobile-api/acceptancetest"
	"code.citik.ru/mic/mobile-api/acceptancetest/pages"
	"context"
	"github.com/dailymotion/allure-go"
	"testing"
)

func init() {
	initializer = acceptancetest.NewInitializer(context.Background())
	initializer.GlobalInit(false)
}

func TestGetLastSeenProducts(t *testing.T) {
	allure.Test(t, allure.Description("get last seen products"),
		allure.Action(func() {
			//получаю первый товар из рекомендаций с главной
			id, _ := pages.GetItem(t, acceptancetest.HomepageClient)

			//получить данные о товаре, чтобы товар добавился в блок "Последние просмотренные товары"
			pages.GetItemData(t, acceptancetest.ProductClient, id)

			//проверяю добавление товара в блоке "Последние просмотренные товары"
			pages.CheckLastSeenProducts(t, acceptancetest.HomepageClient, id)
		}),
	)
}
