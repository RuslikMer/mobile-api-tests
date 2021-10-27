package acceptancetest

// TestUserCredentials - Структура данных
type TestUserCredentials struct {
	Password    string
	NewPassword string
	Email       string
	PhoneNumber string
	Ip          string
	Code        string
}

// UserData - Хранилище данных
var UserData = TestUserCredentials{
	Email:       "", // Почта для авторизации по почте/паролю
	Password:    "password",                  // пароль для авторизации по почте/паролю, телефону/паролю
	NewPassword: "NewPassword",               // новый пароль при смене в лк
	PhoneNumber: "",               // используется для всех авторизаций
	Ip:          "",               // ip пользователя
}

// UserRegistrationData - Регистрационные данные пользователя
var UserRegistrationData = TestUserCredentials{
	PhoneNumber: "",                     // Телефон для регистрации
	Email:       "", // Почта для регистрации
	Password:    "",                        // Пароль для регистрации
	Code:        "",                            // Проверочный код
}
