package server

import (
	"application/controllers"
	"application/controllers/administration"
	"application/models"
	"application/server/middlewares"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"strconv"
)

func PrintRoutes(router martini.Router) {
	routes := router.All()

	urlLength := 0
	for _, routeinfo := range routes {
		if len(routeinfo.Pattern()) > urlLength {
			urlLength = len(routeinfo.Pattern())
		}
	}

	fmt.Printf("%-4v|%-8v|%-"+strconv.Itoa(urlLength)+"v|%v\n", "#", "method", "pattern", "description")
	for i, routeinfo := range routes {
		fmt.Printf("%-4v|%-8v|%-"+strconv.Itoa(urlLength)+"v|%v", i+1, routeinfo.Method(), routeinfo.Pattern(), routeinfo.GetName()+"\n")
	}
}

func Routes() martini.Router {
	var router martini.Router

	router = martini.NewRouter()

	router.Group("/subscriptions", func(a martini.Router) {
		// Выдача последних новостей в виде ленты RSS +
		a.Get("/news/rss/", controllers.GetNewsRss).
			Name("Выдача последних новостей в виде ленты RSS")
		// Удаление подписки на новости +
		a.Get("/unsubscribe/:unsubscribeCode/", controllers.UnsubscribeFromNews).
			Name("Удаление подписки на новости")
	})

	router.Group("/api/v1.0/session", func(a martini.Router) {
		// Аутентификация пользователя +
		a.Post("/user/", binding.Json(models.ViewSession{}), controllers.CreateSession).
			Name("Аутентификация пользователя")
		// Проверка токена с продлением +
		a.Get("/:token", middlewares.RequireSessionKeepWithRoute, controllers.KeepSession).
			Name("Проверка токена с продлением")
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, controllers.KeepSession).
			Name("Проверка токена с продлением")
		// Завершение сеанса пользователя +
		a.Delete("/:token", middlewares.RequireSessionKeepWithRoute, controllers.DeleteSession).
			Name("Завершение сеанса пользователя")
		a.Delete("/", middlewares.RequireSessionKeepWithoutRoute, controllers.DeleteSession).
			Name("Завершение сеанса пользователя")
	})

	router.Group("/api/v1.0/files", func(a martini.Router) {
		// Загрузка файла на сервер +
		a.Post("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights,
			binding.MultipartForm(models.ViewFile{}), controllers.UploadFile).
			Name("Загрузка файла на сервер")
		// Отображение картинки по ключу +
		a.Get("/:key/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetFile).
			Name("Отображение картинки по ключу")
		// Удаление файла на сервере +
		a.Delete("/:key/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.DeleteFile).
			Name("Удаление файла на сервере")
	})

	router.Group("/api/v1.0", func(a martini.Router) {
		// Шаблон домашней страницы +
		a.Get("/", controllers.HomePageTemplate).
			Name("Шаблон домашней страницы")
		// Проверка доступности сервера +
		a.Get("/ping/:token", middlewares.RequireSessionCheckWithRoute, controllers.Ping).
			Name("Проверка доступности сервера")
		a.Get("/ping/", middlewares.RequireSessionCheckWithoutRoute, controllers.Ping).
			Name("Проверка доступности сервера")
		// Запрос картинки с капчей +
		a.Get("/captcha/native/", controllers.GetCaptcha).
			Name("Запрос картинки с капчей")
		// Подтверждение email пользователя +
		a.Post("/emails/confirm/", binding.Json(models.EmailConfirm{}), controllers.ConfirmEmail).
			Name("Подтверждение email пользователя")
		// Загрузка картинок ?
		a.Get("/images/:type/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetImage).
			Name("Загрузка картинок")
	})

	router.Group("/api/v1.0/user", func(a martini.Router) {
		// Получение списка групп пользователей +
		a.Get("/groups/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetGroups).
			Name("Получение списка групп пользователей")
		// Получение информации о пользователе +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetUserInfo).
			Name("Получение информации о пользователе")
		// Изменение информации о пользователе +
		a.Patch("/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ChangeUser{}), controllers.UpdateUserInfo).
			Name("Изменение информации о пользователе")
		// Получение информации о e-mail пользователя +
		a.Get("/emails/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetUserEmails).
			Name("Получение информации о e-mail пользователя")
		// Изменение информации о e-mail пользователя +
		a.Put("/emails/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.UpdateEmails{}), controllers.UpdateUserEmails).
			Name("Изменение информации о e-mail пользователя")
		// Изменение пароля пользователя на новый +
		a.Patch("/password/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ChangePassword{}), controllers.ChangePassword).
			Name("Изменение пароля пользователя на новый")
		// Получение информации о мобильных телефонах пользователя +
		a.Get("/mobilephones/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetUserMobilePhones).
			Name("Получение информации о мобильных телефонах пользователя")
		// Изменение информации о мобильных телефонах пользователя +
		a.Put("/mobilephones/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.UpdateMobilePhones{}), controllers.UpdateUserMobilePhones).
			Name("Изменение информации о мобильных телефонах пользователя")
	})

	router.Group("/api/v1.0/users", func(a martini.Router) {
		// Регистрация пользователя +
		a.Post("/register/:token", binding.Json(models.ViewUser{}), controllers.Register).
			Name("Регистрация пользователя")
		a.Post("/register/", binding.Json(models.ViewUser{}), controllers.Register).
			Name("Регистрация пользователя")
		// Восстановление пароля пользователя +
		a.Post("/password/", binding.Json(models.ViewUser{}), controllers.RestorePassword).
			Name("Восстановление пароля пользователя")
		// Смена пароля при восстановлении +
		a.Put("/password/:code/", binding.Json(models.PasswordUpdate{}), controllers.UpdatePassword).
			Name("Смена пароля при восстановлении")
		a.Put("/password/", binding.Json(models.PasswordUpdate{}), controllers.UpdatePassword).
			Name("Смена пароля при восстановлении")
		// Отказ от восстановления пароля +
		a.Delete("/password/:code/", controllers.DeletePasswordRestoring).
			Name("Отказ от восстановления пароля")
	})

	router.Group("/api/v1.0/administration/groups", func(a martini.Router) {
		// Получение списка всех групп доступа +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetGroupsInfo).
			Name("Получение списка всех групп доступа")
	})

	router.Group("/api/v1.0/administration/users", func(a martini.Router) {
		// Получение информации о данных списка пользователей +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUserMetaData).
			Name("Получение информации о данных списка пользователей")
		// Получить список пользователей +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUsers).
			Name("Получить список пользователей") //массив ошибок
		// Добавление нового пользователя +
		a.Post("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights,
			binding.Json(models.ViewApiUserFull{}), administration.CreateUser).
			Name("Добавление нового пользователя")
		// Получение подробной информации о пользователе +
		a.Get("/:userid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUserFullInfo).
			Name("Получение подробной информации о пользователе")
		// Изменение пользователя +
		a.Put("/:userid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights,
			binding.Json(models.ViewApiUserFull{}), administration.UpdateUser).
			Name("Изменение пользователя")
		// Удаление пользователя +
		a.Delete("/:userid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.DeleteUser).
			Name("Удаление пользователя")
	})

	router.Group("/api/v1.0/administration/units", func(a martini.Router) {
		// Получение общей информации о объединениях пользователей +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUnitMetaData).
			Name("Получение общей информации о объединениях пользователей")
		// Получение списка объединений всех пользователей +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUnits).
			Name("Получение списка объединений всех пользователей")
		// Создание нового объединения пользователей +
		a.Post("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights,
			binding.Json(models.ViewShortUnit{}), administration.CreateUnit).
			Name("Создание нового объединения пользователей")
		// Получение подробной информации об объединении пользователей +
		a.Get("/:unitId/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUnit).
			Name("Получение подробной информации об объединении пользователей")
		// Изменение информации об объединении пользователей +
		a.Put("/:unitId/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights,
			binding.Json(models.ViewLongUnit{}), administration.UpdateUnit).
			Name("Изменение информации об объединении пользователей")
		// Удаление объединения пользователей +
		a.Delete("/:unitId/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.DeleteUnit).
			Name("Удаление объединения пользователей")
		// Получение количества объектов зависящих от объединения +
		a.Options("/:unitId/dependences/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUnitDependences).
			Name("Получение количества объектов зависящих от объединения")
		// Получение списка пользователей объединения +
		a.Get("/:unitId/users/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUnitUsers).
			Name("Получение списка пользователей объединения")
		// Получение списка таблиц объединения +
		a.Get("/:unitId/tables/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUnitTables).
			Name("Получение списка таблиц объединения")
		// Получение списка проектов заказчиков объединения +
		a.Get("/:unitId/projects/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUnitProjects).
			Name("Получение списка проектов заказчиков объединения")
		// Получение списка заказов объединения +
		a.Get("/:unitId/orders/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUnitOrders).
			Name("Получение списка заказов объединения")
		// Получение списка услуг привязанных к объединению +
		a.Get("/:unitId/services/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUnitFacilities).
			Name("Получение списка услуг привязанных к объединению")
		// Получение списка компаний объединения +
		a.Get("/:unitId/organisations/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUnitCompanies).
			Name("Получение списка компаний объединения")
		// Получение списка зарегистрированных за объединением имён отправителей sms +
		a.Get("/:unitId/smsfroms/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUnitSMSSenders).
			Name("Получение списка зарегистрированных за объединением имён отправителей sms")
		// Получение списка счетов организаций объединения +
		a.Get("/:unitId/invoices/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetUnitInvoices).
			Name("Получение списка счетов организаций объединения")
	})

	router.Group("/api/v1.0/administration/orders", func(a martini.Router) {
		// Общая информация о заказах +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetOrderMetaData).
			Name("Общая информация о заказах")
		// Получение списка всех заказов +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetOrders).
			Name("Получение списка всех заказов")
		// Получение заказа +
		a.Get("/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetOrder).
			Name("Получение заказа")
		// Изменение заказа +
		a.Put("/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights,
			binding.Json(models.ViewFullOrder{}), administration.UpdateOrder).
			Name("Изменение заказа")
		// Удаление заказа +
		a.Delete("/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.DeleteOrder).
			Name("Удаление заказа")
	})

	router.Group("/api/v1.0/classification", func(a martini.Router) {
		// Получение справочника классификации контактов  +
		a.Get("/contacts/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetAvailableContacts).
			Name("Получение справочника классификации контактов")
		// Получение справочника операторов мобильной связи +
		a.Get("/mobileoperators/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetMobileOperators).
			Name("Получение справочника операторов мобильной связи")
		// Получение списка категорий услуг +
		a.Get("/services/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetFacilityTypes).
			Name("Получение списка категорий услуг")
		// Получение справочника правовых форм организаций +
		a.Get("/legalformorganisation/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetCompanyTypes).
			Name("Получение справочника правовых форм организаций")
		// Получение справочника кодов классификации организаций +
		a.Get("/organisationClasses/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetCompanyClasses).
			Name("Получение справочника кодов классификации организаций")
		// Получение справочника типов адресов +
		a.Get("/addresses/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetAddressTypes).
			Name("Получение справочника типов адресов")
		// Справочник статусов заказа +
		a.Get("/orderstatuses/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetComplexStatuses).
			Name("Справочник статусов заказа")
	})

	router.Group("/api/v1.0/administration/classification/contacts", func(a martini.Router) {
		// Получение классификатора контактов полностью +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetClassifiers).
			Name("Получение классификатора контактов полностью")
		// Создание новой записи в классификаторе контактов +
		a.Post("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights,
			binding.Json(models.ViewClassifier{}), administration.CreateClassifier).
			Name("Создание новой записи в классификаторе контактов")
		// Получение одной записи классификатора контактов +
		a.Get("/:id/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.GetClassifier).
			Name("Получение одной записи классификатора контактов")
		// Изменение записи классификатора контактов +
		a.Put("/:id/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights,
			binding.Json(models.ViewUpdateClassifier{}), administration.UpdateClassifier).
			Name("Изменение записи классификатора контактов")
		// Удаление записи классификатора контактов +
		a.Delete("/:id/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireAdminRights, administration.DeleteClassifier).
			Name("Удаление записи классификатора контактов")
	})

	router.Group("/api/v1.0/services", func(a martini.Router) {
		// Получение списка услуг поставщиков +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetFacilities).
			Name("Получение списка услуг поставщиков")
		// Получение справочника периодов +
		a.Get("/periods/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetPeriods).
			Name("Получение справочника периодов")
		// Получение справочника событий +
		a.Get("/events/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetEvents).
			Name("Получение справочника событий")
		// Получение поставщиков услуг оказывающих услугу sms рассылка +
		a.Get("/suppliers/sms/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetSMSSuppliers).
			Name("Получение поставщиков услуг оказывающих услугу sms рассылка")
		// Получение поставщиков услуги «HLR запрос» +
		a.Get("/suppliers/hlr/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetHLRSuppliers).
			Name("Получение поставщиков услуги «HLR запрос»")
		// Получение поставщиков услуги «Ввод данных» +
		a.Get("/suppliers/recognize/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetRecognizeSuppliers).
			Name("Получение поставщиков услуги «Ввод данных»")
		// Получение поставщиков услуги «Верификация базы данных» +
		a.Get("/suppliers/verification/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetVerifySuppliers).
			Name("Получение поставщиков услуги «Верификация базы данных»")
	})

	router.Group("/api/v1.0/suppliers/services", func(a martini.Router) {
		// Получение списка услуг поставщиков +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetFacilities).
			Name("Получение списка услуг поставщиков")
		// Получение текущего списка услуг оказываемых поставщиком услуг +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireSupplierRights, controllers.GetSupplierFacilities).
			Name("Получение текущего списка услуг оказываемых поставщиком услуг")
		// Изменение списка оказываемых услуг поставщиком услуг +
		a.Put("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireSupplierRights,
			binding.Json(models.ViewFacilities{}), controllers.UpdateSupplierFacilities).
			Name("Изменение списка оказываемых услуг поставщиком услуг")
	})

	router.Group("/api/v1.0/suppliers/orders", func(a martini.Router) {
		// Получение общей информации о заказах +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireSupplierRights, controllers.GetMetaOrders).
			Name("Получение общей информации о заказах")
		// Получение списка заказов +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireSupplierRights, controllers.GetOrders).
			Name("Получение списка заказов")
		// Получение полной информации о заказе +
		a.Get("/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireOrderRights, controllers.GetOrder).
			Name("Получение полной информации о заказе")
		// Изменение информации о заказе +
		a.Put("/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireOrderRights,
			binding.Json(models.ViewLongOrder{}), controllers.UpdateOrder).
			Name("Изменение информации о заказе")
		// Отклонение заказа +
		a.Delete("/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireOrderRights, controllers.DeleteOrder).
			Name("Отклонение заказа")
		// Получение расширенной информации заказа - SMS рассылка +
		a.Get("/:oid/service/sms/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireOrderRights, controllers.GetSMSOrder).
			Name("Получение расширенной информации заказа - SMS рассылка")
		// Внесение изменений в расширенную информацию заказа - SMS рассылка +
		a.Put("/:oid/service/sms/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireOrderRights, binding.Json(models.ViewSMSFacility{}), controllers.UpdateSMSOrder).
			Name("Внесение изменений в расширенную информацию заказа - SMS рассылка")
		// Получение расширенной информации заказа - HLR запросы +
		a.Get("/:oid/service/hlr/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireOrderRights, controllers.GetHLROrder).
			Name("Получение расширенной информации заказа - HLR запросы")
		// Внесение изменений в расширенную информацию заказа - HLR запросы +
		a.Put("/:oid/service/hlr/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireOrderRights, binding.Json(models.ViewHLRFacility{}), controllers.UpdateHLROrder).
			Name("Внесение изменений в расширенную информацию заказа - HLR запросы")
		// Получение расширенной информации заказа - Ввод данных +
		a.Get("/:oid/service/recognize/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireOrderRights, controllers.GetRecognizeOrder).
			Name("Получение расширенной информации заказа - Ввод данных")
		// Внесение изменений в расширенную информацию заказа - Ввод данных +
		a.Put("/:oid/service/recognize/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireOrderRights, binding.Json(models.ViewRecognizeFacility{}), controllers.UpdateRecognizeOrder).
			Name("Внесение изменений в расширенную информацию заказа - Ввод данных")
		// Получение расширенной информации заказа - Верификация базы данных +
		a.Get("/:oid/service/verification/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireOrderRights, controllers.GetVerifyOrder).
			Name("Получение расширенной информации заказа - Верификация базы данных")
		// Внесение изменений в расширенную информацию заказа - Верификация базы данных +
		a.Put("/:oid/service/verification/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireOrderRights, binding.Json(models.ViewVerifyFacility{}), controllers.UpdateVerifyOrder).
			Name("Внесение изменений в расширенную информацию заказа - Верификация базы данных")
	})

	router.Group("/api/v1.0/tables", func(a martini.Router) {
		// Получение списка типов таблиц +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetTableTypes).
			Name("Получение списка типов таблиц")
		// Создание таблицы +
		a.Post("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights,
			binding.Json(models.ViewShortCustomerTable{}), controllers.CreateTable).
			Name("Создание таблицы")
		// Получение списка таблиц +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetUnitTables).
			Name("Получение списка таблиц")
		// Получение списка типов колонок +
		a.Get("/fieldtypes/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetColumnTypes).
			Name("Получение списка типов колонок")
		// Импорт таблицы +
		a.Post("/import/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights,
			binding.Json(models.ViewImportTable{}), controllers.ImportDataFromFile).
			Name("Импорт таблицы")
		// Проверка статуса импорта таблицы +
		a.Options("/import/:tmpid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			controllers.GetImportDataStatus).
			Name("Проверка статуса импорта таблицы")
		// Получение списка колонок импортируемой таблицы +
		a.Get("/import/:tmpid/columns/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			controllers.GetImportDataColumns).
			Name("Получение списка колонок импортируемой таблицы")
		// Сохранение списка импортируемых колонок и присвоение типа колонкам +
		a.Put("/import/:tmpid/columns/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			middlewares.RequireEditableTable, binding.Json(models.ViewImportColumns{}), controllers.UpdateImportDataColumns).
			Name("Сохранение списка импортируемых колонок и присвоение типа колонкам")
		// Получение информации об экспорте данных +
		a.Options("/export/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetExportDataMeta).
			Name("Получение информации об экспорте данных")
		// Получение экспортируемых данных данных +
		a.Get("/export/:token/:fid/", controllers.GetExportData).
			Name("Получение экспортируемых данных данных")
		// Получение списка таблиц подходящих под услугу sms рассылка +
		a.Get("/services/sms/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights, controllers.GetSMSTables).
			Name("Получение списка таблиц подходящих под услугу sms рассылка")
		// Получение списка таблиц подходящих под услугу hlr рассылка +
		a.Get("/services/hlr/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights, controllers.GetHLRTables).
			Name("Получение списка таблиц подходящих под услугу hlr рассылка")
		// Получение списка таблиц подходящих под услугу верификация базы данных +
		a.Get("/services/verification/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights, controllers.GetVerifyTables).
			Name("Получение списка таблиц подходящих под услугу верификация базы данных")
		// Получение таблицы +
		a.Get("/:tid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTable).
			Name("Получение таблицы")
		// Изменение таблицы +
		a.Put("/:tid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			middlewares.RequireEditableTable, binding.Json(models.ViewLongCustomerTable{}), controllers.UpdateTable).
			Name("Изменение таблицы")
		// Удаление таблицы +
		a.Delete("/:tid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			middlewares.RequireEditableTable, controllers.DeleteTable).
			Name("Удаление таблицы")
		// Создание колонки в таблице +
		a.Post("/:tid/field/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			middlewares.RequireEditableTable, binding.Json(models.ViewApiTableColumn{}), controllers.CreateTableColumn).
			Name("Создание колонки в таблице")
		// Получение списка колонок таблицы в порядке отображения +
		a.Get("/:tid/field/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTableColumns).
			Name("Получение списка колонок таблицы в порядке отображения")
		// Получение колонки таблицы +
		a.Get("/:tid/field/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTableColumn).
			Name("Получение колонки таблицы")
		// Изменение колонки в таблице +
		a.Put("/:tid/field/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			middlewares.RequireEditableTable, binding.Json(models.ViewApiTableColumn{}), controllers.UpdateTableColumn).
			Name("Изменение колонки в таблице")
		// Удаление колонки в таблице +
		a.Delete("/:tid/field/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			middlewares.RequireEditableTable, controllers.DeleteTableColumn).
			Name("Удаление колонки в таблице")
		// Изменение порядка отображения колонки +
		a.Put("/:tid/sequence/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			middlewares.RequireEditableTable, binding.Json(models.ViewApiOrderTableColumns{}), controllers.UpdateOrderTableColumn).
			Name("Изменение порядка отображения колонки")
		// Получение информации о данных в таблице +
		a.Options("/:tid/data/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTableMetaData).
			Name("Получение информации о данных в таблице")
		// Получение данных таблицы +
		a.Get("/:tid/data/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTableData).
			Name("Получение данных таблицы")
		// Получение строки данных таблицы +
		a.Get("/:tid/data/:rid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTableRow).
			Name("Получение строки данных таблицы")
		// Внесение строки данных в таблицу +
		a.Post("/:tid/data/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			middlewares.RequireEditableTable, binding.Json(models.ViewApiTableRow{}), controllers.CreateTableRow).
			Name("Внесение строки данных в таблицу")
		// Изменение строки данных в таблице +
		a.Put("/:tid/data/:rid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			middlewares.RequireEditableTable, binding.Json(models.ViewApiTableRow{}), controllers.UpdateTableRow).
			Name("Изменение строки данных в таблице")
		// Удаление строки данных из таблицы +
		a.Delete("/:tid/data/:rid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			middlewares.RequireEditableTable, controllers.DeleteTableRow).
			Name("Удаление строки данных из таблицы")
		// Получение информации о данных в ячейке таблицы +
		a.Options("/:tid/cell/:rid/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTableMetaCell).
			Name("Получение информации о данных в ячейке таблицы")
		// Получение данных ячейки таблицы +
		a.Get("/:tid/cell/:rid/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTableCell).
			Name("Получение данных ячейки таблицы")
		// Изменение данных ячейки таблицы +
		a.Put("/:tid/cell/:rid/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			middlewares.RequireEditableTable, binding.Json(models.ViewTableCell{}), controllers.UpdateTableCell).
			Name("Изменение данных ячейки таблицы")
		// Удаление данных ячейки таблицы +
		a.Delete("/:tid/cell/:rid/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			middlewares.RequireEditableTable, controllers.DeleteTableCell).
			Name("Удаление данных ячейки таблицы")
		// Экспорт таблицы +
		a.Get("/:tid/export/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.ExportDataToFile).
			Name("Экспорт таблицы") //массив ошибок
		// Проверка статуса готовности экспортируемого файла +
		a.Options("/:tid/export/:fid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetExportDataStatus).
			Name("Проверка статуса готовности экспортируемого файла")
		// Изменение настроек для таблицы являющейся прайс-листом  +
		a.Put("/:tid/price/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireSupplierRights, middlewares.RequireTableRights,
			binding.Json(models.ViewApiPriceProperties{}), controllers.UpdatePriceTable).
			Name("Изменение настроек для таблицы являющейся прайс-листом")
	})

	router.Group("/api/v1.0/messages/orders", func(a martini.Router) {
		// Получение общей информации о переписке в рамках заказа +
		a.Options("/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireMessageRights, controllers.GetMetaMessages).
			Name("Получение общей информации о переписке в рамках заказа")
		// Получение списка сообщений в рамках заказа +
		a.Get("/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireMessageRights, controllers.GetMessages).
			Name("Получение списка сообщений в рамках заказа")
		// Создание сообщения в рамках заказа +
		a.Post("/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireMessageRights,
			binding.Json(models.ViewMessage{}), controllers.CreateMessage).
			Name("Создание сообщения в рамках заказа")
		// Получение сообщения в рамках заказа +
		a.Get("/:oid/message/:mid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireMessageRights, controllers.GetMessage).
			Name("Получение сообщения в рамках заказа")
		// Пометка сообщения как просмотренное +
		a.Patch("/:oid/message/:mid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireMessageRights, controllers.MarkMessage).
			Name("Пометка сообщения как просмотренное")
		// Пометка всех сообщений заказа как просмотренных +
		a.Patch("/:oid/messages/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireMessageRights, controllers.MarkMessages).
			Name("Пометка всех сообщений заказа как просмотренных")
		// Изменение сообщения в рамках заказа +
		a.Put("/:oid/message/:mid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireMessageRights,
			binding.Json(models.ViewMessage{}), controllers.UpdateMessage).
			Name("Изменение сообщения в рамках заказа")
		// Удаление сообщения в рамках заказа +
		a.Delete("/:oid/message/:mid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireMessageRights, controllers.DeleteMessage).
			Name("Удаление сообщения в рамках заказа")
	})

	router.Group("/api/v1.0/customers/services", func(a martini.Router) {
		// Получение списка доступных услуг +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights, controllers.GetAvailableFacilities).
			Name("Получение списка доступных услуг")
	})

	router.Group("/api/v1.0/customers/invoices", func(a martini.Router) {
		// Получение общей информации о счетах +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights, controllers.GetMetaInvoices).
			Name("Получение общей информации о счетах")
		// Получение списка счетов +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights, controllers.GetInvoices).
			Name("Получение списка счетов")
		// Создание счёта на оплату +
		a.Post("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights,
			binding.Json(models.ViewInvoice{}), controllers.CreateInvoice).
			Name("Создание счёта на оплату")
		// Получение подробной информации о счёте +
		a.Get("/:iid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireInvoiceRights, controllers.GetInvoice).
			Name("Получение подробной информации о счёте")
		// Изменение информации о счёте +
		a.Patch("/:iid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireInvoiceRights,
			binding.Json(models.ViewInvoice{}), controllers.UpdateInvoice).
			Name("Изменение информации о счёте")
		// Отказ от оплаты счёта +
		a.Delete("/:iid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireInvoiceRights, controllers.DeleteInvoice).
			Name("Отказ от оплаты счёта")
	})

	router.Group("/api/v1.0/projects", func(a martini.Router) {
		// Получение общей информации о проектах +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights, controllers.GetMetaProjects).
			Name("Получение общей информации о проектах")
		// Получение списка проектов +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights, controllers.GetAllProjects).
			Name("Получение списка проектов")
		// Получение списка проектов находящихся в работе +
		a.Get("/onthego/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights, controllers.GetActiveProjects).
			Name("Получение списка проектов находящихся в работе")
		// Получение списка проектов находящихся в архиве +
		a.Get("/wasarchived/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights, controllers.GetArchiveProjects).
			Name("Получение списка проектов находящихся в архиве")
		// Создание проекта +
		a.Post("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights,
			binding.Json(models.ViewProject{}), controllers.CreateProject).
			Name("Создание проекта")
		// Получение проекта +
		a.Get("/:prid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireProjectRights, controllers.GetProject).
			Name("Получение проекта")
		// Изменение проекта +
		a.Put("/:prid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireProjectRights,
			binding.Json(models.ViewUpdateProject{}), controllers.UpdateProject).
			Name("Изменение проекта")
		// Перемещение проекта в архив +
		a.Delete("/:prid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireProjectRights, controllers.DeleteProject).
			Name("Перемещение проекта в архив")
		// Получение общей информации о заказах проекта +
		a.Options("/:prid/orders/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireProjectRights, controllers.GetMetaProjectOrders).
			Name("Получение общей информации о заказах проекта")
		// Получение списка заказов проекта +
		a.Get("/:prid/orders/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireProjectRights, controllers.GetProjectOrders).
			Name("Получение списка заказов проекта")
		// Создание нового заказа проекта +
		a.Post("/:prid/orders/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireProjectRights,
			binding.Json(models.ViewShortOrder{}), controllers.CreateProjectOrder).
			Name("Создание нового заказа проекта")
		// Получение полной информации о заказе проекта +
		a.Get("/:prid/orders/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireProjectRights, controllers.GetProjectOrder).
			Name("Получение полной информации о заказе проекта")
		// Изменение информации о заказе проекта +
		a.Patch("/:prid/orders/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireProjectRights,
			binding.Json(models.ViewMiddleOrder{}), controllers.UpdateProjectOrder).
			Name("Изменение информации о заказе проекта")
		// Удаление заказа проекта +
		a.Delete("/:prid/orders/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireProjectRights, controllers.DeleteProjectOrder).
			Name("Удаление заказа проекта")
		// Получение расширенной информации заказа - SMS рассылка +
		a.Get("/:prid/orders/:oid/service/sms/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireProjectRights, controllers.GetProjectSMSOrder).
			Name("Получение расширенной информации заказа - SMS рассылка")
		// Внесение изменений в расширенную информацию заказа - SMS рассылка +
		a.Put("/:prid/orders/:oid/service/sms/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireProjectRights, binding.Json(models.ViewSMSFacility{}), controllers.UpdateProjectSMSOrder).
			Name("Внесение изменений в расширенную информацию заказа - SMS рассылка")
		// Получение расширенной информации заказа - HLR запросы +
		a.Get("/:prid/orders/:oid/service/hlr/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireProjectRights, controllers.GetProjectHLROrder).
			Name("Получение расширенной информации заказа - HLR запросы")
		// Внесение изменений в расширенную информацию заказа - HLR запросы +
		a.Put("/:prid/orders/:oid/service/hlr/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireProjectRights, binding.Json(models.ViewHLRFacility{}), controllers.UpdateProjectHLROrder).
			Name("Внесение изменений в расширенную информацию заказа - HLR запросы")
		// Получение расширенной информации заказа - Ввод данных +
		a.Get("/:prid/orders/:oid/service/recognize/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireProjectRights, controllers.GetProjectRecognizeOrder).
			Name("Получение расширенной информации заказа - Ввод данных")
		// Внесение изменений в расширенную информацию заказа - Ввод данных +
		a.Put("/:prid/orders/:oid/service/recognize/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireProjectRights, binding.Json(models.ViewRecognizeFacility{}), controllers.UpdateProjectRecognizeOrder).
			Name("Внесение изменений в расширенную информацию заказа - Ввод данных")
		// Получение расширенной информации заказа - Верификация базы данных +
		a.Get("/:prid/orders/:oid/service/verification/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireProjectRights, controllers.GetProjectVerifyOrder).
			Name("Получение расширенной информации заказа - Верификация базы данных")
		// Внесение изменений в расширенную информацию заказа - Верификация базы данных +
		a.Put("/:prid/orders/:oid/service/verification/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireProjectRights, binding.Json(models.ViewVerifyFacility{}), controllers.UpdateProjectVerifyOrder).
			Name("Внесение изменений в расширенную информацию заказа - Верификация базы данных")
	})

	router.Group("/api/v1.0/organisations", func(a martini.Router) {
		// Получение общей информации о компаниях объединения +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetMetaCompanies).
			Name("Получение общей информации о компаниях объединения")
		// Получение списка компаний объединения +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetCompanies).
			Name("Получение списка компаний объединения")
		// Получение подробной информации об организации объединения +
		a.Get("/:orgid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCompanyRights, controllers.GetCompany).
			Name("Получение подробной информации об организации объединения")
		// Добавление новой организации объединения +
		a.Post("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights,
			binding.Json(models.ViewCompany{}), controllers.CreateCompany).
			Name("Добавление новой организации объединения")
		// Изменение информации об организации объединения +
		a.Put("/:orgid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCompanyRights,
			binding.Json(models.ViewCompany{}), controllers.UpdateCompany).
			Name("Изменение информации об организации объединения")
		// Удаление организации объединения +
		a.Delete("/:orgid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCompanyRights, controllers.DeleteCompany).
			Name("Удаление организации объединения")
	})

	router.Group("/api/v1.0/subscriptions", func(a martini.Router) {
		// Проверка подписки на новости +
		a.Get("/news/:email/", controllers.GetNewsSubscription).
			Name("Проверка подписки на новости")
		// Подписка на новости +
		a.Post("/news/", binding.Json(models.ViewSubscription{}), controllers.CreateSubscription).
			Name("Подписка на новости")
		// Подтверждение подписки на новостную рассылку +
		a.Patch("/news/", binding.Json(models.SubscriptionConfirm{}), controllers.ConfirmSubscription).
			Name("Подтверждение подписки на новостную рассылку")
	})

	router.Group("/api/v1.0/smsfrom", func(a martini.Router) {
		// Получение общей информации по FROM для SMS +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights, controllers.GetMetaSMSSenders).
			Name("Получение общей информации по FROM для SMS")
		// Получение списка зарегистрированных FROM для SMS +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights, controllers.GetSMSSenders).
			Name("Получение списка зарегистрированных FROM для SMS")
		// Получение FROM отправителя для SMS по id +
		a.Get("/:frmid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireSMSSenderRights, controllers.GetSMSSender).
			Name("Получение FROM отправителя для SMS по id")
		// Регистрация FROM отправителя для SMS +
		a.Post("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireCustomerRights,
			binding.Json(models.ViewSMSSender{}), controllers.CreateSMSSender).
			Name("Регистрация FROM отправителя для SMS")
		// Изменение FROM отправителя для SMS +
		a.Patch("/:frmid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireSMSSenderRights,
			binding.Json(models.ViewSMSSender{}), controllers.UpdateSMSSender).
			Name("Изменение FROM отправителя для SMS")
		// Удаление FROM отправителя для SMS +
		a.Delete("/:frmid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireSMSSenderRights, controllers.DeleteSMSSender).
			Name("Удаление FROM отправителя для SMS")
	})

	router.NotFound(middlewares.Default)
	return router
}
