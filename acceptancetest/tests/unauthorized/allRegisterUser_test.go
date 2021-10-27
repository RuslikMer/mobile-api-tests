// +build acceptancetest

package unauthorized

import (
	"code.citik.ru/mic/mobile-api/acceptancetest"
	verifyv1 "code.citik.ru/mic/mobile-api/acceptancetest/grpcclient/gen/citilink/mobileapi/verify/v1"
	"code.citik.ru/mic/mobile-api/acceptancetest/pages"
	"context"
	"github.com/dailymotion/allure-go"
	"testing"
)

func init() {
	initializer = acceptancetest.NewInitializer(context.Background())
	initializer.GlobalInit(false)
}

func TestVerifyPhone(t *testing.T) {
	allure.Test(t, allure.Description("verify phone to profile"),
		allure.Action(func() {
			//проверяю свободен ли номер телефона
			pages.IsPhoneBusy(t, acceptancetest.ProfileClient, acceptancetest.UserRegistrationData.PhoneNumber)

			//проверяю свободен ли e-mail
			pages.IsEmailBusy(t, acceptancetest.ProfileClient, acceptancetest.UserRegistrationData.Email)

			//получаю токен
			requestId := pages.SendSms(t, acceptancetest.VerifyClient, verifyv1.SendSmsRequest_OPERATION_REGISTRATION, acceptancetest.UserRegistrationData.PhoneNumber)

			//верифицирую введённый код перед регистрацией
			token := pages.SmsCheck(t, acceptancetest.VerifyClient, requestId, acceptancetest.UserRegistrationData.Code)

			//регистрирую пользователя
			pages.Register(t, acceptancetest.ProfileClient, acceptancetest.UserRegistrationData.PhoneNumber, acceptancetest.UserRegistrationData.Email, token)

			//отправляю подтверждающее смс на телефон из ЛК
			verificationRequestId := pages.SendVerificationSms(t, acceptancetest.VerifyClient, verifyv1.SendVerificationSmsRequest_OPERATION_VERIFY_PHONE)

			//верифицирую введенный код в ЛК
			verificationToken := pages.CheckSmsVerification(t, acceptancetest.VerifyClient, verificationRequestId, acceptancetest.UserRegistrationData.Code)

			//подтверждение номера телефона, который указан в профиле как неподтвержденный
			pages.VerifyPhone(t, acceptancetest.ProfileClient, verificationToken)

			//получаю идентификатор пользователя
			userId := pages.GetCurrentUser(t, acceptancetest.ProfileClient).GetId()

			//удаляю пользователя
			acceptancetest.DeleteUserById(t, userId)
		}),
	)
}
