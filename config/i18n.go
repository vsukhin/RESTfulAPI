package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Resource struct {
	Messages struct {
		User_Email          string `yaml:"User_Email"`          // Email пользователя
		Greetings           string `yaml:"Greetings"`           // Приветствие письма
		Signature           string `yaml:"Signature"`           // Подпись письма
		PasswordCode        string `yaml:"PasswordCode"`        // Код подтверждения пароля
		EmailCode           string `yaml:"EmailCode"`           // Код подтверждения email
		RegistrationSubject string `yaml:"RegistrationSubject"` // Подтверждение регистрации
		PasswordSubject     string `yaml:"PasswordSubject"`     // Подтверждение смены пароля
		EmailSubject        string `yaml:"EmailSubject"`        // Подтверждение email
		OK                  string `yaml:"OK"`                  // Все в порядке
		NewsHeader          string `yaml:"NewsHeader"`          // Заголовок новостей
		SubscriptionSubject string `yaml:"SubscriptionSubject"` // Подтверждение подписки
		SubscribeCode       string `yaml:"SubscribeCode"`       // Код для подписки на новости
		UnsubscribeCode     string `yaml:"UnsubscribeCode"`     // Код для отписки от новостей

	} `yaml:"Messages"` // Общая информация

	Errors struct {
		Binding struct {
			Field_Empty    string `yaml:"Field_Empty"`    // Ошибка незаполненного поля
			Field_Short    string `yaml:"Field_Short"`    // Ошибка короткого поля
			Field_Long     string `yaml:"Field_Long"`     // Ошибка длинного поля
			Language_Wrong string `yaml:"Language_Wrong"` // Ошибка неверного языка
			Field_Regexp   string `yaml:"Field_Regexp"`   // Ошибка неверной маски
		} `yaml:"Binding"` // Ошибки привязки

		Api struct {
			Captcha_Required        string `yaml:"Captcha_Required"`        // Ошибка незаполненной капчи
			Login_Or_Password_Wrong string `yaml:"Login_Or_Password_Wrong"` // Ошибка неверного логина или пароля
			Captcha_Wrong           string `yaml:"Captcha_Wrong"`           // Ошибка неверной капчи
			Language_NotSupported   string `yaml:"Language_NotSupported"`   // Ошибка неподдерживаемого языка
			User_Blocked            string `yaml:"User_Blocked"`            // Ошибка блокировки пользователя
			Token_Wrong             string `yaml:"Token_Wrong"`             // Ошибка неверного токена
			Data_Wrong              string `yaml:"Data_Wrong"`              // Ошибка неверных данных
			Ip_Or_Net_Blocked       string `yaml:"Ip_Or_Net_Blocked"`       // Ошибка заблокированного ip или сети
			Password_Too_Simple     string `yaml:"Password_Too_Simple"`     // Ошибка слишком простого пароля
			Confirmation_Code_Wrong string `yaml:"Confirmation_Code_Wrong"` // Ошибка неверного кода подтверждения
			Request_Too_Often       string `yaml:"Request_Too_Often"`       // Ошибка слишком частых запросов
			Method_NotAllowed       string `yaml:"Method_NotAllowed"`       // Ошибка неразрешенного метода
			Object_NotExist         string `yaml:"Object_NotExist"`         // Ошибка несуществующего объекта
			Equipment_Info_Wrong    string `yaml:"Equipment_Info_Wrong"`    // Ошибка неверной информации об оборудовании
			File_NotImage           string `yaml:"File_NotImage"`           // Ошибка файла неявляющегося картинкой
			Email_InUse             string `yaml:"Email_InUse"`             // Ошибка уже используемого email
			Data_Changes_Denied     string `yaml:"Data_Changes_Denied"`     // Ошибка изменения данных
			Data_Delete_Denied      string `yaml:"Data_Delete_Denied"`      // Ошибка удаления данных
			MobilePhone_InUse       string `yaml:"MobilePhone_InUse"`       // Ошибка уже используемого мобильного телефона
			Token_Hash_Wrong        string `yaml:"Token_Hash_Wrong"`        // Ошибка декодирования хэша токена
			Device_Wrong            string `yaml:"Device_Wrong"`            // Ошибка невернoго устройства
			SMSSender_InUse         string `yaml:"SMSSender_InUse"`         // Ошибка уже используемого отправителя
		} `yaml:"Api"` // Ошибки API

	} `yaml:"Errors"` // Сообщения об ошибках
}

var (
	Localization map[string]Resource
)

func InitI18n() (err error) {
	Localization = make(map[string]Resource)
	for _, lang := range Configuration.Server.AvailableLanguages {
		resourcefile, err := ioutil.ReadFile(filepath.Join(Configuration.Server.ResourceStorage, lang+".yml"))
		if err != nil {
			logger.Fatalf("Can't read data from resource file: %v", err)
			return err
		}
		resourcedata := new(Resource)
		if err = yaml.Unmarshal(resourcefile, resourcedata); err != nil {
			logger.Fatalf("Can't unmarshal data from yaml to resource structure: %v", err)
			return err
		} else {
			Localization[lang] = *resourcedata
		}
	}

	return nil
}
