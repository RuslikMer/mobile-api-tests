// +build acceptancetest

package authorized

import (
	"code.citik.ru/mic/mobile-api/acceptancetest"
	productv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/catalog/product/v1"
	orderv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/order/order/v1"
	overallv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/overall/v1"
	verifyv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/verify/v1"
	"code.citik.ru/mic/mobile-api/acceptancetest/pages"
	"context"
	"github.com/dailymotion/allure-go"
	"testing"
)

var initializer *acceptancetest.Initializer

func init() {
	initializer = acceptancetest.NewInitializer(context.Background())
	initializer.GlobalInit(true)
}

func TestChangeItemQuantityAtBasket(t *testing.T) {
	allure.Test(t, allure.Description("change item quantity at basket"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

			//добавляю товар в корзину
			uniqId := pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//обновляю позицию товара
			pages.SetQuantity(t, acceptancetest.OrderClient, uniqId, 2)

			//убавляю позицию товара
			pages.SetQuantity(t, acceptancetest.OrderClient, uniqId, 1)
		}),
	)
}

func TestDeleteItem(t *testing.T) {
	allure.Test(t, allure.Description("delete item"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//удаляю позицию товара
			pages.DeleteItem(t, acceptancetest.OrderClient, []string{id})
		}),
	)
}

func TestSetSavedContactInCheckOut(t *testing.T) {
	allure.Test(t, allure.Description("set contacts for order by id"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получаю первый товар из рекомендаций с главной
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//задаю контактные данные для заказа по id
			pages.SetContactsForOrderById(t, acceptancetest.OrderClient, acceptancetest.ContactClient)

			//задаю адрес доставки заказа
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_SAME_DAY, "", "")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//проверяю возможность оформления заказа
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}

func TestSetComment(t *testing.T) {
	allure.Test(t, allure.Description("set comment for order"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получаю первый товар из рекомендаций с главной
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//оставляю комментарий
			pages.SetComment(t, acceptancetest.OrderClient, "Тестовый комментарий")

			//задаю адрес доставки заказа
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_SAME_DAY, "", "")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//проверяю возможность оформления заказа
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}

func TestAddMultipleProductsAndDeleteOne(t *testing.T) {
	allure.Test(t, allure.Description("add multiple products and delete one"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получаю наименование рандомной категории
			randCategoryName := pages.GetRandomCategory(t, acceptancetest.CategoryClient)

			//получение id товара
			id, price, _ := pages.FastSearch(t, acceptancetest.CatalogClient, randCategoryName)

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//получаю наименование рандомной категории
			randCategoryName2 := pages.GetRandomCategory(t, acceptancetest.CategoryClient)

			//получение id товара
			id2, price2, _ := pages.FastSearch(t, acceptancetest.CatalogClient, randCategoryName2)

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id2, price2)

			//удаляю позицию товара
			pages.DeleteItem(t, acceptancetest.OrderClient, []string{id})
		}),
	)
}

func TestAddAccessoryToCart(t *testing.T) {
	allure.Test(t, allure.Description("add accessory to cart"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price, _ := pages.FastSearch(t, acceptancetest.CatalogClient, "Холодильники")

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//получаю аксессуары
			accessories := pages.GetAccessories(t, acceptancetest.ProductClient, id)

			//добавляю аксессуар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, accessories[0].GetId(), accessories[0].GetPrice())
		}),
	)
}

func TestAddProductAndServiceDeleteItem(t *testing.T) {
	allure.Test(t, allure.Description("add product and service delete item"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

			//добавляю товар в корзину
			uniqId := pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//выбор цифровой услуги у товара
			ids := []string{id}
			serviceId, _ := pages.SelectService(t, acceptancetest.ProductClient, ids, productv1.ProductServiceInfo_GROUP_ID_DIGITAL_SERVICE)

			//добавляю цифровую услугу в корзину
			pages.AddService(t, acceptancetest.OrderClient, serviceId, orderv1.Basket_Item_TYPE_DIGITAL_SERVICE, overallv1.ItemType_ITEM_TYPE_DIGITAL_SERVICE, uniqId)

			//удаляю позицию товара
			pages.DeleteItem(t, acceptancetest.OrderClient, []string{id})
		}),
	)
}

func TestChangeCityAndCheckProductStock(t *testing.T) {
	allure.Test(t, allure.Description("change city and check product stock"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получаю id города
			_, cityId := pages.GetAllCities(t, acceptancetest.CityClient, "Куса")

			//меняю город где мало товаров в наличии (например город Куса)
			pages.ChangeCity(t, acceptancetest.ProfileClient, cityId)

			//проверяю что товар не в наличии
			id, price := pages.GetItemThatNotInStock(t, acceptancetest.CatalogClient, "Наушники")

			//добавляю товар в корзину, который не в наличии
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//меняю город где есть товары в наличие(например город Москва)
			pages.ChangeCity(t, acceptancetest.ProfileClient, "")

			//проверить что товар стал в наличии
			pages.CheckProductInStock(t, acceptancetest.ProductClient, id)
		}),
	)
}

func TestSelectProductOneMoreThanAvailable(t *testing.T) {
	allure.Test(t, allure.Description("select product one more than available"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

			//получить данные о количество товара на складе
			itemData := pages.GetItemData(t, acceptancetest.ProductClient, id)

			//добавляю товар в корзину
			uniqId := pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//выбираю количество товара на 1 больше чем доступно
			//ошибку должны запилить, пока не сделали!!!
			pages.SetQuantity(t, acceptancetest.OrderClient, uniqId, itemData.GetAvailableAmount()+1)
		}),
	)
}

func TestUseBonus(t *testing.T) {
	allure.Test(t, allure.Description("use bonus"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price, _ := pages.FastSearch(t, acceptancetest.CatalogClient, "Холодильники")

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//получаю токен
			verificationRequestId := pages.SendVerificationSms(t, acceptancetest.VerifyClient, verifyv1.SendVerificationSmsRequest_OPERATION_BONUS_USE)

			//верифицирую введенный код в ЛК
			verificationToken := pages.CheckSmsVerification(t, acceptancetest.VerifyClient, verificationRequestId, acceptancetest.UserRegistrationData.Code)

			//подтверждаю возможность применения бонусов
			pages.SetBonusApprove(t, acceptancetest.OrderClient, verificationToken)

			//подтверждаю возможность применения бонусов
			pages.SetBonus(t, acceptancetest.OrderClient, 200)
		}),
	)
}
