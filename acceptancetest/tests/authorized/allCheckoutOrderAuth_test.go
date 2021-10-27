// +build acceptancetest

package authorized

import (
	"code.citik.ru/mic/mobile-api/acceptancetest"
	productv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/catalog/product/v1"
	orderv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/order/order/v1"
	overallv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/overall/v1"
	pages "code.citik.ru/mic/mobile-api/acceptancetest/pages"
	"context"
	"fmt"
	"github.com/dailymotion/allure-go"
	"testing"
)

func init() {
	initializer = acceptancetest.NewInitializer(context.Background())
	initializer.GlobalInit(true)
}

func TestGetOrder(t *testing.T) {
	allure.Test(t, allure.Description("simple get order with self-delivery"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получаю первый товар из рекомендаций с главной
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//выбираю чекбокс "Прислать SMS с номером заказа на телефон получателя"
			pages.SetBoolOption(t, acceptancetest.OrderClient, true, "2")

			//выбираю чекбокс "Отказаться от звонка оператора Call-центра"
			pages.SetBoolOption(t, acceptancetest.OrderClient, true, "1")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//получаю список магазинов
			pup_id := pages.GetStoresWithSelDeliveryFromBasket(t, acceptancetest.SelfDeliveryClient)

			//выбираю магазин для самовывоза
			pages.ChooseShopAtBasket(t, acceptancetest.SelfDeliveryClient, pup_id)

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)
		}),
	)
}

func TestGetOrderStandardDelivery(t *testing.T) {
	allure.Test(t, allure.Description("simple get order with standard delivery"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получаю первый товар из рекомендаций с главной
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//задаю адрес доставки заказа
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_GLOBAL, "", "")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//задаю эмейл для отправки чека при онлайн оплате
			pages.SetContactForCheck(t, acceptancetest.OrderClient, acceptancetest.UserData.Email, "")

			//задаю телефон для отправки чека при онлайн оплате
			pages.SetContactForCheck(t, acceptancetest.OrderClient, "", acceptancetest.UserData.PhoneNumber)

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)
		}),
	)
}

func TestGetOrderSameDayDelivery(t *testing.T) {
	allure.Test(t, allure.Description("simple get order with same day delivery"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получаю первый товар из рекомендаций с главной
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//задаю адрес доставки заказа
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_SAME_DAY, "", "")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)
		}),
	)
}

func TestGetOrderFastDelivery(t *testing.T) {
	allure.Test(t, allure.Description("simple get order with fast delivery"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price, _ := pages.FastSearch(t, acceptancetest.CatalogClient, "Недорогие смартфоны")

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//получаю доступный для доставки город и улицу
			_, cityKladrId, _, streetKladrId := pages.GetAvailableCityAndStreetForAddress(t, acceptancetest.CourierDeliveryClient)

			//получаю доступные виды доставк
			deliveryTypes := pages.GetAvailableDeliveryTypes(t, acceptancetest.CourierDeliveryClient, cityKladrId, streetKladrId)

			//проверяю, есть ли нужный тип доставки
			pages.CheckDeliveryTypeAvailability(t, deliveryTypes, orderv1.Order_Delivery_CitilinkCourier_ID_FAST)

			//задаю адрес доставки заказа
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_FAST, "", "")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)
		}),
	)
}

func TestGetOrderSubcontractService(t *testing.T) {
	allure.Test(t, allure.Description("get order with subcontract service, self-delivery"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

			//добавляю товар в корзину
			uniqId := pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//выбор услуги установки у товара
			ids := []string{id}
			serviceId, _ := pages.SelectService(t, acceptancetest.ProductClient, ids, productv1.ProductServiceInfo_GROUP_ID_PRODUCT_SUBCONTRACT_SERVICE)

			//добавляю услугу установки в корзину
			serviceUniqId := pages.AddService(t, acceptancetest.OrderClient, serviceId, orderv1.Basket_Item_TYPE_PRODUCT_SUBCONTRACT_SERVICE, overallv1.ItemType_ITEM_TYPE_PRODUCT_SUBCONTRACT_SERVICE, uniqId)

			//получаю доступный для доставки город
			city, cityKladrId, _, _ := pages.GetAvailableCityAndStreetForAddress(t, acceptancetest.CourierDeliveryClient)

			//добавляю данные для услуги установки в корзине
			pages.SetSubcontractAdditionalInfo(t, acceptancetest.OrderClient, serviceUniqId, city, cityKladrId)

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//выбираю чекбокс "Прислать SMS с номером заказа на телефон получателя"
			pages.SetBoolOption(t, acceptancetest.OrderClient, true, fmt.Sprint(int(orderv1.Order_Modules_BoolOption_TYPE_SMS_WITH_ORDER_ID)))

			//выбираю чекбокс "Нужен звонок оператора Call-центра"
			pages.SetBoolOption(t, acceptancetest.OrderClient, true, fmt.Sprint(int(orderv1.Order_Modules_BoolOption_TYPE_NEED_CALL_CENTER_CALL)))

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//получаю список магазинов
			pup_id := pages.GetStoresWithSelDeliveryFromBasket(t, acceptancetest.SelfDeliveryClient)

			//выбираю магазин для самовывоза
			pages.ChooseShopAtBasket(t, acceptancetest.SelfDeliveryClient, pup_id)

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)
		}),
	)
}

func TestGetOrderInsuranceService(t *testing.T) {
	allure.Test(t, allure.Description("get order with insurance service, self-delivery"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price, _ := pages.FastSearch(t, acceptancetest.CatalogClient, "Холодильники")

			//добавляю товар в корзину
			uniqId := pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//выбор услуги страховки у товара
			ids := []string{id}
			serviceId, _ := pages.SelectService(t, acceptancetest.ProductClient, ids, productv1.ProductServiceInfo_GROUP_ID_PRODUCT_INSURANCE_SERVICE)

			//добавляю цифровую услугу в корзину
			pages.AddService(t, acceptancetest.OrderClient, serviceId, orderv1.Basket_Item_TYPE_PRODUCT_INSURANCE_SERVICE, overallv1.ItemType_ITEM_TYPE_PRODUCT_INSURANCE_SERVICE, uniqId)

			//добавляю данные для услуги страховки в корзине
			pages.SetInsuranceAdditionalInfo(t, acceptancetest.OrderClient)

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//получаю список магазинов
			pup_id := pages.GetStoresWithSelDeliveryFromBasket(t, acceptancetest.SelfDeliveryClient)

			//выбираю магазин для самовывоза
			pages.ChooseShopAtBasket(t, acceptancetest.SelfDeliveryClient, pup_id)

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)
		}),
	)
}

func TestGetOrderDigitalService(t *testing.T) {
	allure.Test(t, allure.Description("get order with digital service, self-delivery"),
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

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//получаю список магазинов
			pup_id := pages.GetStoresWithSelDeliveryFromBasket(t, acceptancetest.SelfDeliveryClient)

			//выбираю магазин для самовывоза
			pages.ChooseShopAtBasket(t, acceptancetest.SelfDeliveryClient, pup_id)

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)
		}),
	)
}

func TestGetOrderWithProductsInBox(t *testing.T) {
	allure.Test(t, allure.Description("get order with products in box"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price, _ := pages.FastSearch(t, acceptancetest.CatalogClient, "Папки-уголки")

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//получаю список магазинов
			pup_id := pages.GetStoresWithSelDeliveryFromBasket(t, acceptancetest.SelfDeliveryClient)

			//выбираю магазин для самовывоза
			pages.ChooseShopAtBasket(t, acceptancetest.SelfDeliveryClient, pup_id)

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)
		}),
	)
}

func TestAddDeliveryAddressToLkAndSelectDeliveryAddressFromSaved(t *testing.T) {
	allure.Test(t, allure.Description("add delivery address to lk and select delivery address from saved"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//задаю адрес доставки в личном кабинете
			pages.AddNewAddress(t, acceptancetest.AddressClient, "Москва", "Арбатская площадь")

			//получаю сохраненные адреса доставки
			address := pages.GetSavedAddress(t, acceptancetest.CourierDeliveryClient)

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//задаю адрес доставки заказа
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_GLOBAL, address[0].GetAddress().GetCity().GetName(), address[0].GetAddress().GetStreet().GetName())

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//проверяю возможность оформления заказа
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}

func TestAddPromoSetProductsToCart(t *testing.T) {
	allure.Test(t, allure.Description("add promo-set products to cart"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получаю промо-комплект
			promoSetId := pages.GetPromoKits(t, acceptancetest.ProductClient, "1365670")

			//добавляю товары промо-комплекта в корзину
			pages.AddItems(t, acceptancetest.OrderClient, acceptancetest.ProductClient, []string{promoSetId[0].GetMainProduct().GetId(), promoSetId[0].GetAdditionalProducts()[0].GetId()},
				[]int64{promoSetId[0].GetMainProduct().GetPrice(), promoSetId[0].GetAdditionalProducts()[0].GetPrice()})

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//задаю адрес доставки заказа
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_SAME_DAY, "", "")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)
		}),
	)
}

func TestAddSeveralTypeServices(t *testing.T) {
	allure.Test(t, allure.Description("add several type services"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price, _ := pages.FastSearch(t, acceptancetest.CatalogClient, "Холодильники")

			//добавляю товар в корзину
			uniqId := pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//выбор цифровой услуги у товара
			ids := []string{id}
			digitalServiceId, _ := pages.SelectService(t, acceptancetest.ProductClient, ids, productv1.ProductServiceInfo_GROUP_ID_DIGITAL_SERVICE)

			//добавляю цифровую услугу в корзину
			pages.AddService(t, acceptancetest.OrderClient, digitalServiceId, orderv1.Basket_Item_TYPE_DIGITAL_SERVICE, overallv1.ItemType_ITEM_TYPE_DIGITAL_SERVICE, uniqId)

			//выбор услуги страховки у товара
			insuranceServiceId, _ := pages.SelectService(t, acceptancetest.ProductClient, ids, productv1.ProductServiceInfo_GROUP_ID_PRODUCT_INSURANCE_SERVICE)

			//добавляю услугу страховки в корзину
			pages.AddService(t, acceptancetest.OrderClient, insuranceServiceId, orderv1.Basket_Item_TYPE_PRODUCT_INSURANCE_SERVICE, overallv1.ItemType_ITEM_TYPE_PRODUCT_INSURANCE_SERVICE, uniqId)

			//выбор услуги установки у товара
			subcontractServiceId, _ := pages.SelectService(t, acceptancetest.ProductClient, ids, productv1.ProductServiceInfo_GROUP_ID_PRODUCT_SUBCONTRACT_SERVICE)

			//добавляю услугу установки в корзину
			serviceUniqId := pages.AddService(t, acceptancetest.OrderClient, subcontractServiceId, orderv1.Basket_Item_TYPE_PRODUCT_SUBCONTRACT_SERVICE, overallv1.ItemType_ITEM_TYPE_PRODUCT_SUBCONTRACT_SERVICE, uniqId)

			//добавляю данные для услуги страховки в корзине
			pages.SetInsuranceAdditionalInfo(t, acceptancetest.OrderClient)

			//получаю доступный для доставки город
			city, cityKladrId, _, _ := pages.GetAvailableCityAndStreetForAddress(t, acceptancetest.CourierDeliveryClient)

			//добавляю данные для услуги установки в корзине
			pages.SetSubcontractAdditionalInfo(t, acceptancetest.OrderClient, serviceUniqId, city, cityKladrId)

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//задаю адрес доставки заказа
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_SAME_DAY, "", "")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)
		}),
	)
}

func TestCheckOnlinePaymentAvailability(t *testing.T) {
	allure.Test(t, allure.Description("check online payment availability"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получаю первый товар из рекомендаций с главной
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//задаю адрес доставки заказа
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_GLOBAL, "", "")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CARD_ONLINE)

			//выбираю чекбокс "Подтверждение контактных данных"
			pages.SetBoolOption(t, acceptancetest.OrderClient, true, fmt.Sprint(int(orderv1.Order_Modules_BoolOption_TYPE_CONTACT_CONFIRMATION)))

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю возможность оплаты заказа
			pages.CheckPaymentAvailability(t, acceptancetest.OrderHistoryClient, []string{orderId})

			//получаю  ссылку для онлайн оплаты
			pages.GetPaymentUrl(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)
		}),
	)
}

func TestRepeatOrder(t *testing.T) {
	allure.Test(t, allure.Description("repeat order"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получение id товара
			id, price, _ := pages.FastSearch(t, acceptancetest.CatalogClient, "Ноутбуки")

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, id, price)

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//задаю адрес доставки заказа
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_GLOBAL, "", "")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)

			//получаю из списка заказов id заказа
			orderIdFromList := pages.Get(t, acceptancetest.OrderHistoryClient, orderId)

			//повторяю заказ
			pages.RepeatOrder(t, acceptancetest.OrderHistoryClient, orderIdFromList.GetId())

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//задаю адрес доставки заказа
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_GLOBAL, "", "")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//оформляю повторный заказ
			reorderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, reorderId)

			//отменяю повторный заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, reorderId)
		}),
	)
}

func TestUsePromoCode(t *testing.T) {
	allure.Test(t, allure.Description("use promo code"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)
			//сделать получение промокода и id товара со страницы акции, как реализуют возвращение списка товаров текущей акции и получение промокода
			//на данный момент промокод есть на странице акции, но он не приходит отдельным полем,
			//и нет возможности применить его к товару. т.к нет возможности получить точно тот товар, к которому применяется промокод из акции

			//добавляю товар в корзину
			pages.AddItem(t, acceptancetest.OrderClient, "1159231", 2290)

			//применяю промокод
			pages.UsePromoCode(t, acceptancetest.OrderClient, "TOOLS")

			//удаляю промокод
			pages.RemovePromoCode(t, acceptancetest.OrderClient)

			//применяю промокод
			pages.UsePromoCode(t, acceptancetest.OrderClient, "TOOLS")

			//задаю контактные данные для заказа
			pages.SetContactsForOrder(t, acceptancetest.OrderClient)

			//задаю адрес доставки заказа
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_SAME_DAY, "", "")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//оформляю заказ
			orderId := pages.CheckoutOrder(t, acceptancetest.CheckoutClient)

			//проверяю список заказов
			pages.CheckOrder(t, acceptancetest.OrderHistoryClient, orderId)

			//отменяю заказ
			pages.OrderCancel(t, acceptancetest.OrderHistoryClient, orderId)
		}),
	)
}
