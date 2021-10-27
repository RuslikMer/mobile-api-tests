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

func TestAddDeliveryAddressToProfileAndDelete(t *testing.T) {
	allure.Test(t, allure.Description("add delivery address to profile and delete"),
		allure.Action(func() {
			//задаю адрес доставки
			addressId := pages.AddNewAddress(t, acceptancetest.AddressClient, "Москва", "Арбатская площадь")

			//удаляю адрес доставки по id
			pages.DeleteAddressById(t, acceptancetest.AddressClient, []string{addressId.GetId()})
		}),
	)
}

func TestAddAvatarToProfileAndDelete(t *testing.T) {
	allure.Test(t, allure.Description("add avatar to profile and delete"),
		allure.Action(func() {
			//добавляю аватар
			pages.UploadAvatar(t, acceptancetest.ProfileClient, "../upload/test.jpeg")

			//удаляю аватар
			pages.DeleteAvatar(t, acceptancetest.ProfileClient)
		}),
	)
}

func TestAddNewContactAndDelete(t *testing.T) {
	allure.Test(t, allure.Description("add new contact and delete"),
		allure.Action(func() {
			//добавляю нового получателя
			contact := pages.AddNewContact(t, acceptancetest.ContactClient, "Тестирование", "Тестирование")

			//получаю контактные данные получателя по его id
			contactInfo := pages.GetContactById(t, acceptancetest.ContactClient, contact.Id)

			//обновляю контактные данные получателя
			pages.UpdateContactById(t, acceptancetest.ContactClient, contactInfo.Id, "Тестирование2", "Тестирование2")

			//удаляю получателя
			pages.DeleteContactById(t, acceptancetest.ContactClient, []string{contactInfo.Id})
		}),
	)
}

func TestUpdateSavedAddress(t *testing.T) {
	allure.Test(t, allure.Description("update saved address"),
		allure.Action(func() {
			//получаю все адреса
			allAddress := pages.GetAllAddresses(t, acceptancetest.AddressClient)

			//ищу адрес по его Id
			addressInfo := pages.GetAddressById(t, acceptancetest.AddressClient, allAddress[0].Id)

			//изменяю данные адреса
			pages.ChangeAddressData(t, acceptancetest.ProfileClient, acceptancetest.AddressClient, addressInfo)
		}),
	)
}

func TestUpdateUserData(t *testing.T) {
	allure.Test(t, allure.Description("update user data"),
		allure.Action(func() {
			//обновляю данные пользователя
			pages.ChangeUserData(t, acceptancetest.ProfileClient)
		}),
	)
}

func TestChangePassword(t *testing.T) {
	allure.Test(t, allure.Description("change password"),
		allure.Action(func() {
			//меняю пароль
			pages.ChangePassword(t, acceptancetest.ProfileClient, acceptancetest.UserData.Password, acceptancetest.UserData.NewPassword)

			//получаю гостевой токен
			token, _ := initializer.GetGuestToken()

			//задаю клиента для новой сессии
			initializer.SetClients(initializer.GetGuestContext(token))

			//авторизовываюсь по новому паролю
			authResponse := pages.AuthByEmail(t, acceptancetest.AuthClient, acceptancetest.UserData.Email, acceptancetest.UserData.NewPassword)

			//задаю авторизованного клиента
			initializer.SetClients(initializer.GetGuestContext(authResponse.GetAccessToken()))

			//меняю пароль обратно
			pages.ChangePassword(t, acceptancetest.ProfileClient, acceptancetest.UserData.NewPassword, acceptancetest.UserData.Password)

			//получаю гостевой токен
			token, _ = initializer.GetGuestToken()

			//задаю клиента для новой сессии
			initializer.SetClients(initializer.GetGuestContext(token))

			//авторизовываюсь по старому паролю
			authResponse = pages.AuthByEmail(t, acceptancetest.AuthClient, acceptancetest.UserData.Email, acceptancetest.UserData.Password)

			//задаю авторизованного клиента
			initializer.SetClients(initializer.GetGuestContext(authResponse.GetAccessToken()))

			//выхожу из лк
			pages.Logout(t, acceptancetest.AuthClient)
		}),
	)
}
