// +build acceptancetest

package unauthorized

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

var initializer *acceptancetest.Initializer

func init() {
	initializer = acceptancetest.NewInitializer(context.Background())
	initializer.GlobalInit(false)
}

func TestAddItemSelfDelivery(t *testing.T) {
	allure.Test(t, allure.Description("add item with self-delivery"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получаю первый товар из рекомендаций с главной
			id, price := pages.GetItem(t, acceptancetest.HomepageClient)

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

			//выбираю чекбокс "Отправить SMS при поступлении товара в точку выдачи"
			pages.SetBoolOption(t, acceptancetest.OrderClient, true, fmt.Sprint(int(orderv1.Order_Modules_BoolOption_TYPE_SMS_ORDER_IN_STORE)))

			//проверяю возможность оформления
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}

func TestAddItemStandardDelivery(t *testing.T) {
	allure.Test(t, allure.Description("add item with standard delivery"),
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

			//проверяю возможность оформления
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}

func TestAddItemSameDayDelivery(t *testing.T) {
	allure.Test(t, allure.Description("add item with same day delivery"),
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

			//проверяю возможность оформления
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}

func TestAddItemFastDelivery(t *testing.T) {
	allure.Test(t, allure.Description("add item with fast delivery"),
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
			pages.SetCourierDelivery(t, acceptancetest.CourierDeliveryClient, orderv1.Order_Delivery_CitilinkCourier_ID_FAST, "", "")

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//проверяю возможность оформления
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}

func TestAddItemSubcontractService(t *testing.T) {
	allure.Test(t, allure.Description("add item with subcontract service, courier delivery"),
		allure.Action(func() {
			//очищаю корзину
			pages.ClearBasket(t, acceptancetest.OrderClient)

			//получаю первый товар из рекомендаций с главной
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

			//выбираю тип оплаты
			pages.ChoosePaymentType(t, acceptancetest.OrderClient, orderv1.PaymentId_PAYMENT_ID_CASH_WITH_CARD)

			//получаю список магазинов
			pup_id := pages.GetStoresWithSelDeliveryFromBasket(t, acceptancetest.SelfDeliveryClient)

			//выбираю магазин для самовывоза
			pages.ChooseShopAtBasket(t, acceptancetest.SelfDeliveryClient, pup_id)

			//проверяю возможность оформления
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}

func TestAddItemInsuranceService(t *testing.T) {
	allure.Test(t, allure.Description("add item with insurance service, self-delivery"),
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

			//проверяю возможность оформления
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}

func TestAddItemDigitalService(t *testing.T) {
	allure.Test(t, allure.Description("add item with digital service, self-delivery"),
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

			//проверяю возможность оформления
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}

func TestAddItemWithProductsInBox(t *testing.T) {
	allure.Test(t, allure.Description("add item with products in box"),
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

			//проверяю возможность оформления
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}

func TestAddPromoSetItemsToCart(t *testing.T) {
	allure.Test(t, allure.Description("add promo-set items to cart"),
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

			//проверяю возможность оформления
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}

func TestAddItemSeveralTypeServices(t *testing.T) {
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

			//проверяю возможность оформления
			pages.AvailabilityOfCheckout(t, acceptancetest.OrderClient)
		}),
	)
}
