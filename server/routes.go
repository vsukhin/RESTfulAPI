package server

import (
	"application/controllers"
	"application/controllers/administration"
	"application/models"
	"application/server/middlewares"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func PrintRoutes(router martini.Router) {
	routes := router.All()
	fmt.Printf("%-3v|%-8v|%-50v|%v\n", "#", "method", "pattern", "description")
	for i, routeinfo := range routes {
		fmt.Printf("%-3v|%-8v|%-50v|%v", i+1, routeinfo.Method(), routeinfo.Pattern(), routeinfo.GetName()+"\n")
	}
}

func Routes() martini.Router {
	var router martini.Router

	router = martini.NewRouter()

	router.Group("/api/v1.0/session", func(a martini.Router) {
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
		// Аутентификация пользователя +
		a.Post("/session/", binding.Json(models.ViewSession{}), controllers.CreateSession).
			Name("Аутентификация пользователя")
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

	router.Group("/api/v1.0/classification/contacts", func(a martini.Router) {
		// Получение справочника классификации контактов  +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetAvailableContacts).
			Name("Получение справочника классификации контактов")
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
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireSupplierRights, controllers.GetFacilities).
			Name("Получение списка услуг поставщиков")
	})

	router.Group("/api/v1.0/suppliers/services", func(a martini.Router) {
		// Получение списка услуг поставщиков +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireSupplierRights, controllers.GetFacilities).
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
		// Получение расширенной информации о заказе по оказываемому сервису
		a.Get("/:oid/services/:sid/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireOrderRights, controllers.GetOrderInfo).
			Name("Получение расширенной информации о заказе по оказываемому сервису")
		// Внесение изменений в информацию о заказе по оказываемому сервису
		a.Put("/:oid/services/:sid/", middlewares.RequireSessionKeepWithoutRoute,
			middlewares.RequireOrderRights, controllers.UpdateOrderInfo).
			Name("Внесение изменений в информацию о заказе по оказываемому сервису")
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
			binding.Json(models.ViewImportColumns{}), controllers.UpdateImportDataColumns).
			Name("Сохранение списка импортируемых колонок и присвоение типа колонкам")
		// Получение информации об экспорте данных +
		a.Options("/export/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetExportDataMeta).
			Name("Получение информации об экспорте данных")
		// Получение экспортируемых данных данных +
		a.Get("/export/:token/:fid/", controllers.GetExportData).
			Name("Получение экспортируемых данных данных")
		// Получение таблицы +
		a.Get("/:tid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTable).
			Name("Получение таблицы")
		// Изменение таблицы +
		a.Put("/:tid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			binding.Json(models.ViewLongCustomerTable{}), controllers.UpdateTable).
			Name("Изменение таблицы")
		// Удаление таблицы +
		a.Delete("/:tid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.DeleteTable).
			Name("Удаление таблицы")
		// Создание колонки в таблице +
		a.Post("/:tid/field/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			binding.Json(models.ViewApiTableColumn{}), controllers.CreateTableColumn).
			Name("Создание колонки в таблице")
		// Получение списка колонок таблицы в порядке отображения +
		a.Get("/:tid/field/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTableColumns).
			Name("Получение списка колонок таблицы в порядке отображения")
		// Получение колонки таблицы +
		a.Get("/:tid/field/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTableColumn).
			Name("Получение колонки таблицы")
		// Изменение колонки в таблице +
		a.Put("/:tid/field/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			binding.Json(models.ViewApiTableColumn{}), controllers.UpdateTableColumn).
			Name("Изменение колонки в таблице")
		// Удаление колонки в таблице +
		a.Delete("/:tid/field/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.DeleteTableColumn).
			Name("Удаление колонки в таблице")
		// Изменение порядка отображения колонки +
		a.Put("/:tid/sequence/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			binding.Json(models.ViewApiOrderTableColumns{}), controllers.UpdateOrderTableColumn).
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
			binding.Json(models.ViewApiTableRow{}), controllers.CreateTableRow).
			Name("Внесение строки данных в таблицу")
		// Изменение строки данных в таблице +
		a.Put("/:tid/data/:rid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			binding.Json(models.ViewApiTableRow{}), controllers.UpdateTableRow).
			Name("Изменение строки данных в таблице")
		// Удаление строки данных из таблицы +
		a.Delete("/:tid/data/:rid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.DeleteTableRow).
			Name("Удаление строки данных из таблицы")
		// Получение информации о данных в ячейке таблицы +
		a.Options("/:tid/cell/:rid/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTableMetaCell).
			Name("Получение информации о данных в ячейке таблицы")
		// Получение данных ячейки таблицы +
		a.Get("/:tid/cell/:rid/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetTableCell).
			Name("Получение данных ячейки таблицы")
		// Изменение данных ячейки таблицы +
		a.Put("/:tid/cell/:rid/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
			binding.Json(models.ViewTableCell{}), controllers.UpdateTableCell).
			Name("Изменение данных ячейки таблицы")
		// Удаление данных ячейки таблицы +
		a.Delete("/:tid/cell/:rid/:cid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.DeleteTableCell).
			Name("Удаление данных ячейки таблицы")
		// Экспорт таблицы +
		a.Get("/:tid/export/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.ExportDataToFile).
			Name("Экспорт таблицы") //массив ошибок
		// Проверка статуса готовности экспортируемого файла +
		a.Options("/:tid/export/:fid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights, controllers.GetExportDataStatus).
			Name("Проверка статуса готовности экспортируемого файла")
		// Изменение настроек для таблицы являющейся прайс-листом  +
		a.Put("/:tid/price/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireTableRights,
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
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetAvailableFacilities).
			Name("Получение списка доступных услуг")
	})

	router.Group("/api/v1.0/projects", func(a martini.Router) {
		// Получение общей информации о проектах +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetMetaProjects).
			Name("Получение общей информации о проектах")
		// Получение списка проектов +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights, controllers.GetProjects).
			Name("Получение списка проектов")
		// Создание проекта +
		a.Post("/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireUserRights,
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
		a.Put("/:prid/orders/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireProjectRights,
			binding.Json(models.ViewLongOrder{}), controllers.UpdateProjectOrder).
			Name("Изменение информации о заказе проекта")
		// Удаление заказа проекта +
		a.Delete("/:prid/orders/:oid/", middlewares.RequireSessionKeepWithoutRoute, middlewares.RequireProjectRights, controllers.DeleteProjectOrder).
			Name("Удаление заказа проекта")
	})

	router.NotFound(middlewares.Default)
	return router
}
