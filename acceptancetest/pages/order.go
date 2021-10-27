package pages

import (
	checkoutv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/order/checkout/v1"
	courierv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/order/delivery/citilink_courier/v1"
	selfv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/order/delivery/self/v1"
	orderv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/order/order/v1"
	orderhistoryv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/orderhistory/v1"
	overallv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/overall/v1"
	verifyv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/verify/v1"
	"context"
	"fmt"
	"github.com/dailymotion/allure-go"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

func NewOrderClient(ctx context.Context, conn *grpc.ClientConn) *OrderClient {
	return &OrderClient{client: orderv1.NewOrderAPIClient(conn), ctx: ctx}
}

type OrderClient struct {
	client orderv1.OrderAPIClient
	ctx    context.Context
}

func NewOrderHistoryClient(ctx context.Context, conn *grpc.ClientConn) *OrderHistoryClient {
	return &OrderHistoryClient{client: orderhistoryv1.NewOrderhistoryAPIClient(conn), ctx: ctx}
}

type OrderHistoryClient struct {
	client orderhistoryv1.OrderhistoryAPIClient
	ctx    context.Context
}

func NewSelfDeliveryClient(ctx context.Context, conn *grpc.ClientConn) *SelfDeliveryClient {
	return &SelfDeliveryClient{client: selfv1.NewSelfAPIClient(conn), ctx: ctx}
}

type SelfDeliveryClient struct {
	client selfv1.SelfAPIClient
	ctx    context.Context
}

func NewCourierDeliveryClient(ctx context.Context, conn *grpc.ClientConn) *CourierDeliveryClient {
	return &CourierDeliveryClient{client: courierv1.NewCitilinkCourierAPIClient(conn), ctx: ctx}
}

type CourierDeliveryClient struct {
	client courierv1.CitilinkCourierAPIClient
	ctx    context.Context
}

func NewCheckoutClient(ctx context.Context, conn *grpc.ClientConn) *CheckoutClient {
	return &CheckoutClient{client: checkoutv1.NewCheckoutAPIClient(conn), ctx: ctx}
}

type CheckoutClient struct {
	client checkoutv1.CheckoutAPIClient
	ctx    context.Context
}

func ClearBasket(t *testing.T, orderClient *OrderClient) {
	allure.Step(allure.Description("clear basket"),
		allure.Action(func() {
			clearBasketResponse, err := orderClient.client.ClearBasket(orderClient.ctx, &orderv1.ClearBasketRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			emptyCartPrice := clearBasketResponse.GetOrder().GetBasket().GetCost()
			itemCount := clearBasketResponse.GetOrder().GetBasket().GetCount().GetAll()
			assert.Empty(t, emptyCartPrice)
			assert.Equal(t, int32(0), itemCount)
		}),
	)
}

func AddItem(t *testing.T, orderClient *OrderClient, id string, price int64) string {
	var uniqId string

	allure.Step(allure.Description("add product to basket"),
		allure.Action(func() {
			basketResponse, err := orderClient.client.AddItem(orderClient.ctx, &orderv1.AddItemRequest{
				ItemId:   id,
				Type:     orderv1.Basket_Item_TYPE_PRODUCT,
				Count:    1,
				ItemType: overallv1.ItemType_ITEM_TYPE_PRODUCT,
			})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			basket := basketResponse.GetOrder().GetBasket().GetItems()
			for i := range basket {
				ids := basket[i].GetItemId()
				if id == ids {
					cartId := basketResponse.GetOrder().GetBasket().GetItems()[i].GetItemId()
					cartPrice := basketResponse.GetOrder().GetBasket().GetItems()[i].GetCost()
					uniqId = basketResponse.GetOrder().GetBasket().GetItems()[i].GetUniqId()
					assert.Equal(t, id, cartId)
					assert.Equal(t, price, int64(cartPrice))
				}
				break
			}

			assert.NotEmpty(t, uniqId)
		}),
	)

	return uniqId
}

func AddItems(t *testing.T, orderClient *OrderClient, productClient *ProductClient, itemId []string, price []int64) string {
	var uniqId string
	var itemList []string
	var priceList []int

	allure.Step(allure.Description("add products to basket"),
		allure.Action(func() {
			for i := range itemId {
				itemData := GetItemData(t, productClient, itemId[i])
				assert.NotEqual(t, 1, itemData.GetMultiplicity())
			}

			basketResponse, err := orderClient.client.AddItems(orderClient.ctx, &orderv1.AddItemsRequest{
				Items: []*orderv1.AddItemsRequest_Item{{
					ItemId:   itemId[0],
					Count:    1,
					ItemType: overallv1.ItemType_ITEM_TYPE_PRODUCT,
				},
					{
						ItemId:   itemId[1],
						Count:    1,
						ItemType: overallv1.ItemType_ITEM_TYPE_PRODUCT,
					},
				},
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			for i := range itemId {
				itemReverse := []string{itemId[len(itemId)-1-i]}
				itemList = append(itemList, itemReverse...)
			}

			for i := range price {
				priceReverse := []int{int(price[len(price)-1-i])}
				priceList = append(priceList, priceReverse...)
			}

			basket := basketResponse.GetOrder().GetBasket().GetItems()
			for i := range basket {
				ids := basket[i].GetItemId()
				if itemList[i] == ids {
					cartId := basketResponse.GetOrder().GetBasket().GetItems()[i].GetItemId()
					cartPrice := basketResponse.GetOrder().GetBasket().GetItems()[i].GetCost()
					uniqId = basketResponse.GetOrder().GetBasket().GetItems()[i].GetUniqId()
					assert.Equal(t, itemList[i], cartId)
					assert.Equal(t, priceList[i], int(cartPrice))
				}
				break
			}

			assert.NotEmpty(t, uniqId)
		}),
	)

	return uniqId
}

func AddService(t *testing.T, orderClient *OrderClient, id string, Type orderv1.Basket_Item_Type, itemType overallv1.ItemType, parentUniqId string) string {
	uniqId := ""

	allure.Step(allure.Description("add service to basket"),
		allure.Action(func() {
			basketResponse, err := orderClient.client.AddItem(orderClient.ctx, &orderv1.AddItemRequest{
				ItemId:       id,
				Type:         Type,
				Count:        1,
				ParentUniqId: parentUniqId,
				ItemType:     itemType,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			fmt.Println(err)
			uniqId = basketResponse.GetOrder().GetBasket().GetItems()[0].GetUniqId()
			assert.NotNil(t, uniqId)
		}),
	)

	return uniqId
}

func SetContactsForOrder(t *testing.T, orderClient *OrderClient) {
	allure.Step(allure.Description("set contacts for order"),
		allure.Action(func() {
			contactResponse, err := orderClient.client.SetContact(orderClient.ctx, &orderv1.SetContactRequest{
				Contact: &orderv1.Order_Modules_Contact{
					FirstName: "Tester",
					LastName:  "Tester",
					Phone:     "79351001036",
				},
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			needToShow := contactResponse.GetOrder().GetModules().GetContact().GetNeedToShow()
			needToFill := contactResponse.GetOrder().GetModules().GetContact().GetNeedToFill()
			assert.True(t, needToShow)
			assert.False(t, needToFill)
			needToShow = contactResponse.GetOrder().GetModules().GetDelivery().GetNeedToShow()
			needToFill = contactResponse.GetOrder().GetModules().GetDelivery().GetNeedToFill()
			assert.True(t, needToShow)
			assert.True(t, needToFill)
		}),
	)
}

func ChoosePaymentType(t *testing.T, orderClient *OrderClient, paymentId orderv1.PaymentId) {
	allure.Step(allure.Description("choose payment type"),
		allure.Action(func() {
			paymentResponse, err := orderClient.client.SetPaymentType(orderClient.ctx, &orderv1.SetPaymentTypeRequest{
				Id: paymentId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			needToShow := paymentResponse.GetOrder().GetModules().GetPayment().GetNeedToShow()
			needToFill := paymentResponse.GetOrder().GetModules().GetPayment().GetNeedToFill()
			assert.True(t, needToShow)
			assert.False(t, needToFill)
		}),
	)
}

func GetStoresWithSelDeliveryFromBasket(t *testing.T, selfDeliveryClient *SelfDeliveryClient) string {
	var pup_id string

	allure.Step(allure.Description("get stores with self delivery"),
		allure.Action(func() {
			storesResponse, err := selfDeliveryClient.client.GetStores(selfDeliveryClient.ctx, &selfv1.GetStoresRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			stores := storesResponse.GetStores()

			//ищем магазин с возможностью самовывоза
			for i, store := range stores {
				if store.GetStatus() == orderv1.AllowStatus_ALLOW_STATUS_ALLOW {
					pup_id = stores[i].GetPupId()
					break
				}
			}

			assert.NotEmpty(t, pup_id)
		}),
	)

	return pup_id
}

func GetAvailableCityAndStreetForAddress(t *testing.T, courierDeliveryClient *CourierDeliveryClient) (string, string, string, string) {
	var city string
	var cityKladrId string
	var street string
	var streetKladrId string

	allure.Step(allure.Description("get available city and street"),
		allure.Action(func() {
			citiesResponse, err := courierDeliveryClient.client.GetCitiesForAddress(courierDeliveryClient.ctx, &courierv1.GetCitiesForAddressRequest{
				CityPart: "Москва",
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			city = citiesResponse.GetCities()[0].GetName()
			cityKladrId = citiesResponse.GetCities()[0].GetKladrId()
			assert.NotEmpty(t, city)
			assert.NotEmpty(t, cityKladrId)
			streetsResponse, err := courierDeliveryClient.client.GetStreetsForAddress(courierDeliveryClient.ctx, &courierv1.GetStreetsForAddressRequest{
				StreetPart:  "Красногвардейская",
				CityKladrId: cityKladrId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			street = streetsResponse.GetStreets()[0].GetName()
			streetKladrId = streetsResponse.GetStreets()[0].GetKladrId()
			assert.NotEmpty(t, street)
			assert.NotEmpty(t, streetKladrId)
		}),
	)

	return city, cityKladrId, street, streetKladrId
}

func GetAvailableDateAndTimeForDelivery(t *testing.T, courierDeliveryClient *CourierDeliveryClient, deliveryType orderv1.Order_Delivery_CitilinkCourier_Id) (string, string) {
	var date string
	var Time string

	allure.Step(allure.Description("get available date and time"),
		allure.Action(func() {
			dateResponse, err := courierDeliveryClient.client.GetDateAndTimes(courierDeliveryClient.ctx, &courierv1.GetDateAndTimesRequest{
				Id: deliveryType,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			date = dateResponse.GetDateAndTimes()[0].GetId()
			Time = dateResponse.GetDateAndTimes()[0].GetTimes()[0].GetId()
			assert.NotEmpty(t, date)
			assert.NotEmpty(t, Time)
		}),
	)

	return date, Time
}

func SetCourierDelivery(t *testing.T, courierDeliveryClient *CourierDeliveryClient, deliveryType orderv1.Order_Delivery_CitilinkCourier_Id, city string, street string) {
	var cityKladrId string
	var streetKladrId string

	allure.Step(allure.Description("set courier delivery"),
		allure.Action(func() {
			if city == "" {
				city, cityKladrId, street, streetKladrId = GetAvailableCityAndStreetForAddress(t, courierDeliveryClient)
			} else {
				_, cityKladrId, _, streetKladrId = GetAvailableCityAndStreetForAddress(t, courierDeliveryClient)
			}

			date, Time := GetAvailableDateAndTimeForDelivery(t, courierDeliveryClient, deliveryType)
			deliveryResponse, err := courierDeliveryClient.client.Set(courierDeliveryClient.ctx, &courierv1.SetRequest{
				Id: deliveryType,
				DeliveryAddress: &courierv1.SetRequest_AddressDetails{
					Address: &courierv1.SetRequest_AddressDetails_New_{
						&courierv1.SetRequest_AddressDetails_New{
							Address: &orderv1.Address{
								City:                &orderv1.Address_City{Name: city, KladrId: cityKladrId},
								Street:              &orderv1.Address_Street{Name: street, KladrId: streetKladrId},
								House:               "4",
								Corpus:              "2",
								Building:            "2",
								Flat:                "1",
								Porch:               "2",
								Floor:               "3",
								FeSavedId:           "",
								IsPassRequired:      false,
								IsAddressNormalized: true,
							},
						},
					},
				},
				DateId: date,
				TimeId: Time,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			errors := deliveryResponse.GetFieldErrors()
			assert.Empty(t, errors)
		}),
	)
}

func ChooseShopAtBasket(t *testing.T, selfDeliveryClient *SelfDeliveryClient, pup_id string) {
	allure.Step(allure.Description("choose shop"),
		allure.Action(func() {
			storeResponse, err := selfDeliveryClient.client.Set(selfDeliveryClient.ctx, &selfv1.SetRequest{
				PupId: pup_id,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			pupId := storeResponse.GetOrder().GetModules().GetDelivery().GetSelf().GetPupId()
			assert.Equal(t, pup_id, pupId)
		}),
	)
}

func CheckoutOrder(t *testing.T, checkoutClient *CheckoutClient) string {
	var orderId string

	allure.Step(allure.Description("get order"),
		allure.Action(func() {
			checkoutResponse, err := checkoutClient.client.Checkout(checkoutClient.ctx, &checkoutv1.CheckoutRequest{
				GaClientId: "GA1.1.904941809.1556093647",
				UserAgent:  "Mozilla/5.0 (Android 4.3; Mobile; rv:54.0) Gecko/54.0 Firefox/54.0",
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			orderId = checkoutResponse.GetOrderIds()[0]
			assert.NotEmpty(t, orderId)
		}),
	)

	return orderId
}

func AvailabilityOfCheckout(t *testing.T, orderClient *OrderClient) {
	allure.Step(allure.Description("check availability of checkout"),
		allure.Action(func() {
			orderResponse, err := orderClient.client.Get(orderClient.ctx, &orderv1.GetRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			availability := orderResponse.GetOrder().GetCheckoutAvailable()
			assert.True(t, availability)
		}),
	)
}

func CheckOrder(t *testing.T, orderHistoryClient *OrderHistoryClient, orderId string) {
	var status string

	allure.Step(allure.Description("check order"),
		allure.Action(func() {
			orderHistoryResponse, err := orderHistoryClient.client.Filter(orderHistoryClient.ctx, &orderhistoryv1.FilterRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			orders := orderHistoryResponse.GetOrders()

			//ищу новый заказ
			for _, order := range orders {
				if order.GetId() == orderId {
					status = order.GetStatus().GetName()
					break
				}
			}

			assert.Equal(t, "Заказ создан", status)
		}),
	)
}

func OrderCancel(t *testing.T, orderHistoryClient *OrderHistoryClient, orderId string) {
	allure.Step(allure.Description("order cancel"),
		allure.Action(func() {
			orderResponse, err := orderHistoryClient.client.Cancel(orderHistoryClient.ctx, &orderhistoryv1.CancelRequest{
				Id: orderId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			fmt.Println(err)
			status := orderResponse.GetOrder().GetStatus().GetName()
			assert.Equal(t, "Заказ удален", status)
		}),
	)
}

func SetSubcontractAdditionalInfo(t *testing.T, orderClient *OrderClient, uniqId string, city string, cityKladrId string) {
	allure.Step(allure.Description("set subcontract additional info"),
		allure.Action(func() {
			_, err := orderClient.client.SetSubcontractAdditionalInfo(orderClient.ctx, &orderv1.SetSubcontractAdditionalInfoRequest{
				Additions: []*orderv1.SetSubcontractAdditionalInfoRequest_Addition{{
					ItemUniqId:  uniqId,
					CityName:    city,
					CityKladrId: cityKladrId,
					Address:     "ул. Красногвардейская 4",
					DateTime:    &timestamp.Timestamp{Seconds: time.Now().UnixNano()},
				}},
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
		}),
	)
}

func SetInsuranceAdditionalInfo(t *testing.T, orderClient *OrderClient) {
	allure.Step(allure.Description("set insurance additional info"),
		allure.Action(func() {
			_, err := orderClient.client.SetInsuranceAdditionalInfo(orderClient.ctx, &orderv1.SetInsuranceAdditionalInfoRequest{
				FirstName:  "Тестер",
				LastName:   "Тестер",
				Patronymic: "Тестер",
				Email:      "go-autotester@citilink.ru",
				Phone:      "79351001036",
				Passport: &orderv1.SetInsuranceAdditionalInfoRequest_Passport{
					Series:             "2812",
					Number:             "090090",
					DateOfIssueTime:    &timestamp.Timestamp{Seconds: time.Date(2015, time.February, 4, 0, 0, 0, 0, time.UTC).Unix()},
					PlaceOfIssue:       "УФМС по г. Москва",
					AddressOfResidence: "г. Москва, ул. Ленина 1",
					BirthDateTime:      &timestamp.Timestamp{Seconds: time.Date(1995, time.February, 4, 0, 0, 0, 0, time.UTC).Unix()},
				},
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
		}),
	)
}

func GetAvailableDeliveryTypes(t *testing.T, courierDeliveryClient *CourierDeliveryClient, cityKladrId string, streetKladrId string) []*courierv1.CitilinkCourierDeliverySubtype {
	var types []*courierv1.CitilinkCourierDeliverySubtype

	allure.Step(allure.Description("get available delivery types"),
		allure.Action(func() {
			typesResponse, err := courierDeliveryClient.client.GetSubtypes(courierDeliveryClient.ctx, &courierv1.GetSubtypesRequest{
				CityKladrId:   cityKladrId,
				StreetKladrId: streetKladrId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			types = typesResponse.GetSubtypes()
			assert.NotEmpty(t, types)
		}),
	)

	return types
}

func CheckDeliveryTypeAvailability(t *testing.T, deliveryTypes []*courierv1.CitilinkCourierDeliverySubtype, deliveryType orderv1.Order_Delivery_CitilinkCourier_Id) orderv1.Order_Delivery_CitilinkCourier_Id {
	var deliveryId orderv1.Order_Delivery_CitilinkCourier_Id

	allure.Step(allure.Description("Check delivery type availability"),
		allure.Action(func() {
			for i := range deliveryTypes {
				id := deliveryTypes[i].GetId()
				if id == deliveryType {
					deliveryId = id
					break
				}
			}

			assert.NotEmpty(t, deliveryId, "Данный тип доставки недоступен")
		}),
	)

	return deliveryId
}

func UsePromoCode(t *testing.T, orderClient *OrderClient, promoCode string) {
	allure.Step(allure.Description("use promo code"),
		allure.Action(func() {
			promoCodeResponse, err := orderClient.client.UsePromoCode(orderClient.ctx, &orderv1.UsePromoCodeRequest{
				PromoCode: promoCode,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			usedPromoCode := promoCodeResponse.GetOrder().GetModules().GetPromoCode().GetPromoCode()
			assert.NotEmpty(t, usedPromoCode)
		}),
	)
}

func RemovePromoCode(t *testing.T, orderClient *OrderClient) {
	allure.Step(allure.Description("remove promo code"),
		allure.Action(func() {
			promoCodeResponse, err := orderClient.client.RemovePromoCode(orderClient.ctx, &orderv1.RemovePromoCodeRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			usedPromoCode := promoCodeResponse.GetOrder().GetModules().GetPromoCode().GetPromoCode()
			assert.Empty(t, usedPromoCode)
		}),
	)
}

func GetSavedAddress(t *testing.T, courierDeliveryClient *CourierDeliveryClient) []*orderv1.SavedAddress {
	var address []*orderv1.SavedAddress

	allure.Step(allure.Description("get saved address"),
		allure.Action(func() {
			addressesResponse, err := courierDeliveryClient.client.GetSavedAddresses(courierDeliveryClient.ctx, &courierv1.GetSavedAddressesRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			address = addressesResponse.GetSavedAddresses()
			assert.NotEmpty(t, address)
		}),
	)

	return address
}

func Get(t *testing.T, orderHistoryClient *OrderHistoryClient, orderId string) *orderhistoryv1.Order {
	var order *orderhistoryv1.Order

	allure.Step(allure.Description("get order"),
		allure.Action(func() {
			orderResponse, err := orderHistoryClient.client.Get(orderHistoryClient.ctx, &orderhistoryv1.GetRequest{
				Id: orderId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			order = orderResponse.GetOrder()
			assert.NotEmpty(t, order)
		}),
	)

	return order
}

func RepeatOrder(t *testing.T, orderHistoryClient *OrderHistoryClient, orderId string) *orderv1.Cost {
	var cost *orderv1.Cost

	allure.Step(allure.Description("repeat order"),
		allure.Action(func() {
			orderResponse, err := orderHistoryClient.client.Repeat(orderHistoryClient.ctx, &orderhistoryv1.RepeatRequest{
				Id: orderId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			cost = orderResponse.GetOrder().GetBasket().GetCost()
			assert.NotEmpty(t, cost)
		}),
	)

	return cost
}

func CheckPaymentAvailability(t *testing.T, orderHistoryClient *OrderHistoryClient, orderIds []string) map[string]bool {
	var paymentAvailability map[string]bool

	allure.Step(allure.Description("check payment availability"),
		allure.Action(func() {
			orderResponse, err := orderHistoryClient.client.CheckPaymentAvailability(orderHistoryClient.ctx, &orderhistoryv1.CheckPaymentAvailabilityRequest{
				OrderIds: orderIds,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			paymentAvailability = orderResponse.GetCanPayOnline()
			assert.NotEmpty(t, paymentAvailability)
		}),
	)

	return paymentAvailability
}

func SetQuantity(t *testing.T, orderClient *OrderClient, uniqId string, count int32) {
	allure.Step(allure.Description("set quantity"),
		allure.Action(func() {
			updateResponse, err := orderClient.client.UpdateItem(orderClient.ctx, &orderv1.UpdateItemRequest{
				Items: []*orderv1.UpdateItemRequest_Item{{
					UniqId: uniqId,
					Count:  count,
				}},
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			quantity := updateResponse.GetOrder().GetBasket().GetItems()[0].GetCount()
			assert.Equal(t, count, quantity)
		}),
	)
}

func DeleteItem(t *testing.T, orderClient *OrderClient, ItemUniqIds []string) {
	allure.Step(allure.Description("delete item"),
		allure.Action(func() {
			deleteResponse, err := orderClient.client.DeleteItem(orderClient.ctx, &orderv1.DeleteItemRequest{
				ItemUniqIds: ItemUniqIds,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			items := deleteResponse.GetOrder().GetBasket().GetItems()
			assert.NotEmpty(t, items)
			for _, item := range items {
				assert.NotEqual(t, item.GetUniqId(), ItemUniqIds[0])
			}
		}),
	)
}

func SetContactsForOrderById(t *testing.T, orderClient *OrderClient, contactClient *ContactClient) {
	allure.Step(allure.Description("set contacts for order by id"),
		allure.Action(func() {
			contactId := GetAllContacts(t, contactClient)[0].GetId()
			contactResponse, err := orderClient.client.SetContactById(orderClient.ctx, &orderv1.SetContactByIdRequest{
				ContactId: contactId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			getContactId := contactResponse.GetOrder().GetModules().GetContact().GetId()
			assert.Equal(t, contactId, getContactId)
		}),
	)
}

func SetComment(t *testing.T, orderClient *OrderClient, comment string) {
	allure.Step(allure.Description("set comment"),
		allure.Action(func() {
			commentResponse, err := orderClient.client.SetComment(orderClient.ctx, &orderv1.SetCommentRequest{
				Comment: comment,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			getComment := commentResponse.GetOrder().GetModules().GetComment().GetText()
			assert.Equal(t, comment, getComment)
		}),
	)
}

func UseCoupon(t *testing.T, orderClient *OrderClient, coupon string, pin int32) {
	allure.Step(allure.Description("use coupon"),
		allure.Action(func() {
			couponResponse, err := orderClient.client.UseCoupon(orderClient.ctx, &orderv1.UseCouponRequest{
				Coupon: coupon,
				Pin:    pin,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			checkCoupon := couponResponse.GetOrder().GetModules().GetCoupon().GetCoupon()
			assert.NotEmpty(t, checkCoupon)
		}),
	)
}

func RemoveCoupon(t *testing.T, orderClient *OrderClient) {
	allure.Step(allure.Description("remove coupon"),
		allure.Action(func() {
			removeCouponResponse, err := orderClient.client.RemoveCoupon(orderClient.ctx, &orderv1.RemoveCouponRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			checkRemoveCoupon := removeCouponResponse.GetOrder().GetModules().GetCoupon().GetCoupon()
			assert.Empty(t, checkRemoveCoupon)
		}),
	)
}

func SetConsultantCode(t *testing.T, orderClient *OrderClient, consultantCode string) {
	allure.Step(allure.Description("set consultant code"),
		allure.Action(func() {
			_, err := orderClient.client.SetConsultantCode(orderClient.ctx, &orderv1.SetConsultantCodeRequest{
				Code: consultantCode,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
		}),
	)
}

func SetContactForCheck(t *testing.T, orderClient *OrderClient, email string, phone string) {
	allure.Step(allure.Description("set contact for check"),
		allure.Action(func() {
			contactForCheckResponse, err := orderClient.client.SetContactForCheck(orderClient.ctx, &orderv1.SetContactForCheckRequest{
				Email: email,
				Phone: phone,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			if email == "" {
				checkPhoneContactForCheck := contactForCheckResponse.GetOrder().GetModules().GetContactForCheck().GetPhone()
				assert.Equal(t, phone, checkPhoneContactForCheck)
			} else {
				checkEmailContactForCheck := contactForCheckResponse.GetOrder().GetModules().GetContactForCheck().GetEmail()
				assert.Equal(t, email, checkEmailContactForCheck)
			}
		}),
	)
}

func GetConsigneeList(t *testing.T, orderClient *OrderClient) []*orderv1.Consignee {
	var consigneeList []*orderv1.Consignee

	allure.Step(allure.Description("get consignee list"),
		allure.Action(func() {
			consigneeListResponse, err := orderClient.client.GetConsigneeList(orderClient.ctx, &orderv1.GetConsigneeListRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			consigneeList = consigneeListResponse.GetConsignees()
			assert.NotEmpty(t, consigneeList)
		}),
	)

	return consigneeList
}

func SetConsigneeId(t *testing.T, orderClient *OrderClient, consigneeId string) {
	allure.Step(allure.Description("set consignee id"),
		allure.Action(func() {
			consigneeIdResponse, err := orderClient.client.SetConsigneeId(orderClient.ctx, &orderv1.SetConsigneeIdRequest{
				ConsigneeId: consigneeId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			checkConsigneeId := consigneeIdResponse.GetOrder()
			assert.NotEmpty(t, checkConsigneeId)
		}),
	)
}

func SetBonusApprove(t *testing.T, orderClient *OrderClient, token string) {
	allure.Step(allure.Description("set bonus approve"),
		allure.Action(func() {
			bonusApproveResponse, err := orderClient.client.SetBonusApprove(orderClient.ctx, &orderv1.SetBonusApproveRequest{
				VerifyToken: &verifyv1.Verification{Token: token},
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			checkBonusApprove := bonusApproveResponse.GetOrder().GetModules().GetBonusUse().GetBonusApprove()
			assert.Equal(t, true, checkBonusApprove)
		}),
	)
}

func SetBonus(t *testing.T, orderClient *OrderClient, bonus int32) {
	allure.Step(allure.Description("set bonus"),
		allure.Action(func() {
			bonusResponse, err := orderClient.client.SetBonus(orderClient.ctx, &orderv1.SetBonusRequest{
				UseBonus: bonus,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			checkBonus := bonusResponse.GetOrder().GetModules().GetBonusUse().GetCanUse()
			assert.Equal(t, true, checkBonus)
		}),
	)
}

func SetBoolOption(t *testing.T, orderClient *OrderClient, value bool, id string) {
	allure.Step(allure.Description("set bool option"),
		allure.Action(func() {
			boolOptionResponse, err := orderClient.client.SetBoolOption(orderClient.ctx, &orderv1.SetBoolOptionRequest{
				Value: value,
				Id:    id,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			checkBoolOption := boolOptionResponse.GetOrder().GetModules().GetBoolOptions()
			assert.NotEmpty(t, checkBoolOption)
		}),
	)
}

func GetPaymentUrl(t *testing.T, orderHistoryClient *OrderHistoryClient, orderId string) string {
	var paymentUrl string

	allure.Step(allure.Description("get payment url"),
		allure.Action(func() {
			urlResponse, err := orderHistoryClient.client.GetPaymentUrl(orderHistoryClient.ctx, &orderhistoryv1.GetPaymentUrlRequest{
				Id: orderId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			paymentUrl = urlResponse.GetUrl()
			assert.NotEmpty(t, paymentUrl)
		}),
	)

	return paymentUrl
}
