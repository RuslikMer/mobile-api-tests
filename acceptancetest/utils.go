package acceptancetest

import (
	authv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/auth/v1"
	sessionv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/session/v1"
	"code.citik.ru/mic/mobile-api/acceptancetest/pages"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"os"
	"time"
)

var OrderClient *pages.OrderClient
var HomepageClient *pages.HomepageClient
var SelfDeliveryClient *pages.SelfDeliveryClient
var OrderHistoryClient *pages.OrderHistoryClient
var CheckoutClient *pages.CheckoutClient
var ProductClient *pages.ProductClient
var CourierDeliveryClient *pages.CourierDeliveryClient
var CatalogClient *pages.ProductFilterClient
var ContactClient *pages.ContactClient
var AddressClient *pages.AddressClient
var ProfileClient *pages.ProfileClient
var CityClient *pages.CityClient
var CategoryClient *pages.CategoryClient
var VerifyClient *pages.VerifyClient
var AuthClient *pages.AuthClient

// Initializer Структура для подготовки окружения для запуска тестов
type Initializer struct {
	connection *grpc.ClientConn
	ctx        context.Context
}

func (i *Initializer) Ctx() context.Context {
	return i.ctx
}

func (i *Initializer) Connection() *grpc.ClientConn {
	return i.connection
}

func NewInitializer(ctx context.Context) *Initializer {
	return &Initializer{ctx: ctx}
}

// GlobalInit общий метод, максимально полный по функционалу
// можно создавать подобного рода методы только из нужных методов Initializer
func (i *Initializer) GlobalInit(auth bool) {
	var token string

	apiEndpoint := os.Getenv("TEST_APP_GRPC_ADDR")
	grpcKeepAliveParams := grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:    6 * time.Minute,
		Timeout: 1 * time.Second,
	})

	conn, err := grpc.Dial(apiEndpoint, grpc.WithInsecure(), grpcKeepAliveParams, grpc.WithBlock())
	if err != nil {
		fmt.Println("error on creating grpc connection")
		os.Exit(1)
	}

	i.connection = conn
	if auth {
		token, err = i.GetAuthToken()
	} else {
		token, err = i.GetGuestToken()
	}

	if err != nil {
		fmt.Print(fmt.Sprintf("err getting token: %s", err))
		os.Exit(1)
	}

	i.SetContext(token, true)
	i.SetClients(i.Ctx())
}

// setContext формирует контекст, путем добавления метадаты
func (i *Initializer) SetContext(accessToken string, isTester bool) {
	md := make([]string, 0)
	md = append(md, os.Getenv("CL_REAL_IP_HEADER_NAME"), UserData.Ip)
	if accessToken != "" {
		md = append(md, os.Getenv("CL_JWT_MD_PARAM"), accessToken)
	}

	if isTester {
		md = append(md, os.Getenv("CL_TESTER_HEADER_NAME"), "true")
	}

	i.ctx = metadata.AppendToOutgoingContext(i.ctx, md...)
}

// GetAuthToken получает новый токен авторизованного пользователя
func (i *Initializer) GetAuthToken() (string, error) {
	token, err := i.GetGuestToken()
	if err != nil {
		fmt.Print(fmt.Sprintf("err getting guest token: %s", err))
		os.Exit(1)
	}

	authRequest := &authv1.AuthByPasswordRequest{
		Credentials: &authv1.AuthByPasswordRequest_Phone{
			Phone: UserData.PhoneNumber,
		},
		Password: UserData.Password,
	}

	// Инициируем клиента для получения сессии авторизованного пользователя
	authClient := authv1.NewAuthAPIClient(i.connection)
	authResponse, err := authClient.AuthByPassword(i.GetGuestContext(token), authRequest)
	if err != nil {
		return "", fmt.Errorf("can't get auth token: %w", err)
	}

	return authResponse.GetResponse().GetAccessToken(), nil
}

// GetGuestToken получает новый гостевой токен
func (i *Initializer) GetGuestToken() (string, error) {
	// Инициируем клиента для получения сессии неавторизованного пользователя
	sessionClient := sessionv1.NewSessionAPIClient(i.connection)

	// Выполняем запрос на получение сессии
	sessionResponse, err := sessionClient.Register(i.ctx, &sessionv1.RegisterRequest{
		Device: &sessionv1.Device{
			Id:         "test",
			Os:         sessionv1.Device_OS_ANDROID,
			AppVersion: "1",
			SystemLang: "test",
			PushToken:  "",
			Timezone:   "test",
		},
	})

	if err != nil {
		fmt.Println("error on session register")
		fmt.Print(err)
		os.Exit(1)
	}

	return sessionResponse.GetAccessToken(), nil
}

func (i *Initializer) GetGuestContext(token string) context.Context {
	// Добавляем в метадату access-token и ip пользователя
	return metadata.AppendToOutgoingContext(
		context.Background(),
		os.Getenv("CL_JWT_MD_PARAM"),
		token,
		os.Getenv("CL_REAL_IP_HEADER_NAME"),
		UserData.Ip,
		os.Getenv("CL_TESTER_HEADER_NAME"),
		"true",
	)
}

// SetClients задает данные текущего клиента
func (i *Initializer) SetClients(ctx context.Context) {
	OrderClient = pages.NewOrderClient(ctx, i.Connection())
	HomepageClient = pages.NewHomepageClient(ctx, i.Connection())
	SelfDeliveryClient = pages.NewSelfDeliveryClient(ctx, i.Connection())
	OrderHistoryClient = pages.NewOrderHistoryClient(ctx, i.Connection())
	CheckoutClient = pages.NewCheckoutClient(ctx, i.Connection())
	ProductClient = pages.NewProductClient(ctx, i.Connection())
	CourierDeliveryClient = pages.NewCourierDeliveryClient(ctx, i.Connection())
	CatalogClient = pages.NewProductFilterClient(ctx, i.Connection())
	ContactClient = pages.NewContactClient(ctx, i.Connection())
	AddressClient = pages.NewAddressClient(ctx, i.Connection())
	ProfileClient = pages.NewProfileClient(ctx, i.Connection())
	CityClient = pages.NewCityClient(ctx, i.Connection())
	CategoryClient = pages.NewCategoryClient(ctx, i.Connection())
	VerifyClient = pages.NewVerifyClient(ctx, i.Connection())
	AuthClient = pages.NewAuthClient(ctx, i.Connection())
}
