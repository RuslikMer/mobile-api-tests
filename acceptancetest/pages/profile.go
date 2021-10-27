package pages

import (
	"bytes"
	authv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/auth/v1"
	overallv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/overall/v1"
	addressv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/profile/address/v1"
	contactv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/profile/contact/v1"
	profilev1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/profile/v1"
	sessionv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/session/v1"
	verifyv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/verify/v1"
	"context"
	"github.com/dailymotion/allure-go"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"image"
	"image/jpeg"
	"log"
	"os"
	"testing"
	"time"
)

type AddressClient struct {
	client addressv1.AddressAPIClient
	ctx    context.Context
}

func NewAddressClient(ctx context.Context, connection *grpc.ClientConn) *AddressClient {
	return &AddressClient{client: addressv1.NewAddressAPIClient(connection), ctx: ctx}
}

type ContactClient struct {
	client contactv1.ContactAPIClient
	ctx    context.Context
}

func NewContactClient(ctx context.Context, connection *grpc.ClientConn) *ContactClient {
	return &ContactClient{client: contactv1.NewContactAPIClient(connection), ctx: ctx}
}

type ProfileClient struct {
	client profilev1.ProfileAPIClient
	ctx    context.Context
}

func NewProfileClient(ctx context.Context, connection *grpc.ClientConn) *ProfileClient {
	return &ProfileClient{client: profilev1.NewProfileAPIClient(connection), ctx: ctx}
}

func GetAllAddresses(t *testing.T, addressClient *AddressClient) []*addressv1.DeliveryAddress {
	var addresses []*addressv1.DeliveryAddress

	allure.Step(allure.Description("get all delivery addresses"),
		allure.Action(func() {
			addressesResponse, err := addressClient.client.GetAll(addressClient.ctx, &addressv1.GetAllRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			addresses = addressesResponse.GetAddresses()
			assert.NotNil(t, addresses)
		}),
	)

	return addresses
}

func GetAddressById(t *testing.T, addressClient *AddressClient, addressId string) *addressv1.DeliveryAddress {
	var address *addressv1.DeliveryAddress

	allure.Step(allure.Description("get delivery address by id"),
		allure.Action(func() {
			addressResponse, err := addressClient.client.Get(addressClient.ctx, &addressv1.GetRequest{
				Id: addressId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			address = addressResponse.GetAddress()
			assert.NotNil(t, address)
		}),
	)

	return address
}

func AddNewAddress(t *testing.T, addressClient *AddressClient, city string, street string) *addressv1.DeliveryAddress {
	var address *addressv1.DeliveryAddress

	allure.Step(allure.Description("add delivery address"),
		allure.Action(func() {
			cityKladrId, streetKladrId := FindAddress(t, addressClient, city, street)
			addResponse, err := addressClient.client.Add(addressClient.ctx, &addressv1.AddRequest{
				CityKladrId: cityKladrId,
				Street: &addressv1.AddRequest_StreetKladrId{
					StreetKladrId: streetKladrId,
				},
				House:    "1",
				Corpus:   "",
				Building: "",
				Flat:     "",
				Porch:    "",
				Floor:    "",
				Comment:  "",
			})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			address = addResponse.GetAddress()
			assert.NotNil(t, address)
		}),
	)

	return address
}

func UpdateAddressById(t *testing.T, addressClient *AddressClient, addressId string, city string, street string) *addressv1.DeliveryAddress {
	var address *addressv1.DeliveryAddress

	allure.Step(allure.Description("update delivery address by id"),
		allure.Action(func() {
			cityKladrId, streetKladrId := FindAddress(t, addressClient, city, street)
			updateResponse, err := addressClient.client.Update(addressClient.ctx, &addressv1.UpdateRequest{
				Id: addressId,
				Fields: &addressv1.UpdateRequest_Fields{
					CityKladrId: &wrapperspb.StringValue{Value: cityKladrId},
					Street: &addressv1.UpdateRequest_Fields_StreetKladrId{
						StreetKladrId: &wrapperspb.StringValue{Value: streetKladrId},
					},
					House:    &wrapperspb.StringValue{Value: "1"},
					Corpus:   &wrapperspb.StringValue{Value: ""},
					Building: &wrapperspb.StringValue{Value: ""},
					Flat:     &wrapperspb.StringValue{Value: ""},
					Porch:    &wrapperspb.StringValue{Value: ""},
					Floor:    &wrapperspb.StringValue{Value: ""},
					Comment:  &wrapperspb.StringValue{Value: ""},
				},
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			address = updateResponse.GetAddress()
			assert.NotNil(t, address)
		}),
	)

	return address
}

func DeleteAddressById(t *testing.T, addressClient *AddressClient, addressIds []string) {
	allure.Step(allure.Description("delete delivery address by id"),
		allure.Action(func() {
			_, err := addressClient.client.Delete(addressClient.ctx, &addressv1.DeleteRequest{
				Ids: addressIds,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
		}),
	)
}

func FindAddress(t *testing.T, addressClient *AddressClient, city string, street string) (string, string) {
	var cityKladrId string
	var streetKladrId string

	allure.Step(allure.Description("find city and street fo delivery address"),
		allure.Action(func() {
			citiesResponse, err := addressClient.client.FindCities(addressClient.ctx, &addressv1.FindCitiesRequest{
				Name: city,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			cityKladrId = citiesResponse.GetCities()[0].GetKladrId()
			assert.NotEmpty(t, cityKladrId)
			streetsResponse, err := addressClient.client.FindStreets(addressClient.ctx, &addressv1.FindStreetsRequest{
				Name:        street,
				CityKladrId: cityKladrId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			streetKladrId = streetsResponse.GetStreets()[0].GetKladrId()
			assert.NotEmpty(t, streetKladrId)
		}),
	)

	return cityKladrId, streetKladrId
}

func GetAllContacts(t *testing.T, contactClient *ContactClient) []*contactv1.Contact {
	var contacts []*contactv1.Contact

	allure.Step(allure.Description("get all contacts"),
		allure.Action(func() {
			contactsResponse, err := contactClient.client.GetAll(contactClient.ctx, &contactv1.GetAllRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			contacts = contactsResponse.GetContacts()
			assert.NotNil(t, contacts)
		}),
	)

	return contacts
}

func GetContactById(t *testing.T, contactClient *ContactClient, contactId string) *contactv1.Contact {
	var contact *contactv1.Contact

	allure.Step(allure.Description("get contacts by id"),
		allure.Action(func() {
			contactResponse, err := contactClient.client.Get(contactClient.ctx, &contactv1.GetRequest{
				Id: contactId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			contact = contactResponse.GetContact()
			assert.NotNil(t, contact)
		}),
	)

	return contact
}

func AddNewContact(t *testing.T, contactClient *ContactClient, firstName string, lastName string) *contactv1.Contact {
	var contact *contactv1.Contact

	allure.Step(allure.Description("add new contact"),
		allure.Action(func() {
			contactResponse, err := contactClient.client.Add(contactClient.ctx, &contactv1.AddRequest{
				FirstName:       firstName,
				LastName:        lastName,
				Phone:           "79323232323",
				AdditionalPhone: "",
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			contact = contactResponse.GetContact()
			assert.NotEmpty(t, contact)
		}),
	)

	return contact
}

func UpdateContactById(t *testing.T, contactClient *ContactClient, contactId string, firstName string, lastName string) *contactv1.Contact {
	var contact *contactv1.Contact

	allure.Step(allure.Description("update contact by id"),
		allure.Action(func() {
			contactResponse, err := contactClient.client.Update(contactClient.ctx, &contactv1.UpdateRequest{
				Id:              contactId,
				FirstName:       &wrapperspb.StringValue{Value: firstName},
				LastName:        &wrapperspb.StringValue{Value: lastName},
				Phone:           &wrapperspb.StringValue{Value: "79323232323"},
				AdditionalPhone: &wrapperspb.StringValue{Value: ""},
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			contact = contactResponse.GetContact()
			assert.NotEmpty(t, contact)
		}),
	)

	return contact
}

func DeleteContactById(t *testing.T, contactClient *ContactClient, contactsIds []string) {
	allure.Step(allure.Description("delete contact by id"),
		allure.Action(func() {
			_, err := contactClient.client.Delete(contactClient.ctx, &contactv1.DeleteRequest{
				Ids: contactsIds,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
		}),
	)
}

func IsEmailBusy(t *testing.T, profileClient *ProfileClient, email string) bool {
	var busy bool

	allure.Step(allure.Description("is email busy"),
		allure.Action(func() {
			emailResponse, err := profileClient.client.IsEmailBusy(profileClient.ctx, &profilev1.IsEmailBusyRequest{
				Email: email,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			busy = emailResponse.GetBusy()
		}),
	)

	return busy
}

func IsPhoneBusy(t *testing.T, profileClient *ProfileClient, phone string) bool {
	var busy bool

	allure.Step(allure.Description("is phone busy"),
		allure.Action(func() {
			phoneResponse, err := profileClient.client.IsPhoneBusy(profileClient.ctx, &profilev1.IsPhoneBusyRequest{
				Phone: phone,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			busy = phoneResponse.GetBusy()
		}),
	)

	return busy
}

func Register(t *testing.T, profileClient *ProfileClient, phone string, email string, token string) *authv1.AuthResponse {
	var authResponse *authv1.AuthResponse

	allure.Step(allure.Description("register new user"),
		allure.Action(func() {
			registerResponse, err := profileClient.client.Register(profileClient.ctx, &profilev1.RegisterRequest{
				Verification: &verifyv1.Verification{Token: token},
				Phone:        phone,
				Email:        email,
				Password:     "password",
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			authResponse = registerResponse.GetResponse()
			assert.NotEmpty(t, authResponse)
		}),
	)

	return authResponse
}

func GetCurrentUser(t *testing.T, profileClient *ProfileClient) *sessionv1.MyProfile {
	var user *sessionv1.MyProfile

	allure.Step(allure.Description("get current user info"),
		allure.Action(func() {
			userResponse, err := profileClient.client.GetCurrent(profileClient.ctx, &profilev1.GetCurrentRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			user = userResponse.GetProfile()
			assert.NotEmpty(t, user)
		}),
	)

	return user
}

func ChangeCity(t *testing.T, profileClient *ProfileClient, cityId string) {
	allure.Step(allure.Description("change city"),
		allure.Action(func() {
			if cityId == "" {
				cityId = "msk_cl"
			}

			changeCityResponse, err := profileClient.client.ChangeCity(profileClient.ctx, &profilev1.ChangeCityRequest{
				CityId: cityId,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			changeCityId := changeCityResponse.GetCity().GetId()
			assert.Equal(t, cityId, changeCityId)
		}),
	)
}

func RestorePasswordByEmail(t *testing.T, profileClient *ProfileClient, email string) {
	allure.Step(allure.Description("restore password by email"),
		allure.Action(func() {
			_, err := profileClient.client.RestorePasswordByEmail(profileClient.ctx, &profilev1.RestorePasswordByEmailRequest{
				Email: email,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
		}),
	)
}

func UpdateUserData(t *testing.T, profileClient *ProfileClient, nickname string, firstName string, lastname string) *sessionv1.MyProfile {
	var user *sessionv1.MyProfile

	allure.Step(allure.Description("update user data"),
		allure.Action(func() {
			updateResponse, err := profileClient.client.Update(profileClient.ctx, &profilev1.UpdateRequest{
				Nickname:  &wrapperspb.StringValue{Value: nickname},
				Firstname: &wrapperspb.StringValue{Value: firstName},
				Lastname:  &wrapperspb.StringValue{Value: lastname},
				BirthdateTime: &timestamp.Timestamp{
					Seconds: time.Date(1995, time.February, 4, 0, 0, 0, 0, time.UTC).Unix(),
				},
				Gender: &profilev1.UpdateRequest_GenderValue{Value: overallv1.Gender_GENDER_MALE},
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			user = updateResponse.GetProfile()
			assert.NotEmpty(t, user)
		}),
	)

	return user
}

func ChangePassword(t *testing.T, profileClient *ProfileClient, oldPassword string, newPassword string) string {
	var authToken string

	allure.Step(allure.Description("change password"),
		allure.Action(func() {
			changePasswordResponse, err := profileClient.client.ChangePassword(profileClient.ctx, &profilev1.ChangePasswordRequest{
				OldPassword: oldPassword,
				NewPassword: newPassword,
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			authToken = changePasswordResponse.GetResponse().GetAccessToken()
			assert.NotEmpty(t, authToken)
		}),
	)

	return authToken
}

func UploadAvatar(t *testing.T, profileClient *ProfileClient, filePath string) {
	allure.Step(allure.Description("upload avatar"),
		allure.Action(func() {
			f, err := os.Open(filePath)
			if err != nil {
				log.Println(err)
				assert.Nil(t, err)
			}

			image, _, err := image.Decode(f)
			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, image, nil)
			bytes := buf.Bytes()

			uploadResponse, err := profileClient.client.UploadAvatar(profileClient.ctx, &profilev1.UploadAvatarRequest{
				File: bytes,
			})

			if err != nil {
				log.Println(err)
				assert.Nil(t, err)
			}

			avatarUrl := uploadResponse.GetProfile().GetAvatarUrl()
			assert.NotEmpty(t, avatarUrl)
		}),
	)
}

func DeleteAvatar(t *testing.T, profileClient *ProfileClient) {
	allure.Step(allure.Description("delete avatar"),
		allure.Action(func() {
			deleteResponse, err := profileClient.client.DeleteAvatar(profileClient.ctx, &profilev1.DeleteAvatarRequest{})
			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			avatarUrl := deleteResponse.GetProfile().GetAvatarUrl()
			assert.Empty(t, avatarUrl)
		}),
	)
}

func VerifyPhone(t *testing.T, profileClient *ProfileClient, token string) {
	allure.Step(allure.Description("verify phone"),
		allure.Action(func() {
			verifyResponse, err := profileClient.client.VerifyPhone(profileClient.ctx, &profilev1.VerifyPhoneRequest{
				VerifyToken: &verifyv1.Verification{Token: token},
			})

			if err != nil {
				log.Println(err)
			}

			assert.Nil(t, err)
			verify := verifyResponse.GetProfile().GetIsPhoneConfirmed()
			assert.True(t, verify)
		}),
	)
}

func ChangeAddressData(t *testing.T, profileClient *ProfileClient, addressClient *AddressClient, addressInfo *addressv1.DeliveryAddress) {
	var newCityName, newStreetName, cityId string

	allure.Step(allure.Description("change address data"),
		allure.Action(func() {
			if addressInfo.StreetName == "Молодёжная улица" && addressInfo.CityName == "Москва" {
				newCityName = "Казань"
				newStreetName = "Центральная"
				cityId = "kzn_cl"
			} else {
				newCityName = "Москва"
				newStreetName = "Молодёжная"
				cityId = "msk_cl"
			}

			//меняю город
			ChangeCity(t, profileClient, cityId)

			//обновляю данные адреса
			updateAddress := UpdateAddressById(t, addressClient, addressInfo.Id, newCityName, newStreetName)

			//проверяю, что город изменился
			assert.NotEqual(t, addressInfo.CityName, updateAddress.CityName)

			//проверяю, что улица изменилась
			assert.NotEqual(t, addressInfo.StreetName, updateAddress.StreetName)
		}),
	)
}

func ChangeUserData(t *testing.T, profileClient *ProfileClient) {
	var newNickName, newFirstName, newLastName string

	allure.Step(allure.Description("change user data"),
		allure.Action(func() {
			//получаю данные пользователя
			userData := GetCurrentUser(t, profileClient)
			if userData.Nickname == "TestNickName" && userData.FirstName == "TestFirstName" && userData.LastName == "TestLastName" {
				newNickName = "TestUpdateNickName"
				newFirstName = "TestUpdateFirstName"
				newLastName = "TestUpdateLastName"
			} else {
				newNickName = "TestNickName"
				newFirstName = "TestFirstName"
				newLastName = "TestLastName"
			}

			//обновляю данные пользователя
			updateUserData := UpdateUserData(t, profileClient, newNickName, newFirstName, newLastName)

			//проверяю, что никнейм изменился
			assert.NotEqual(t, userData.Nickname, updateUserData.Nickname)

			//проверяю, что имя изменилось
			assert.NotEqual(t, userData.FirstName, updateUserData.FirstName)

			//проверяю, что фамилия изменилась
			assert.NotEqual(t, userData.LastName, updateUserData.LastName)
		}),
	)
}
