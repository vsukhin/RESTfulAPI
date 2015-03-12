package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"

	"application/controllers"
	"application/models"
	"application/server/middlewares"
)

func routes() martini.Router {
	var route martini.Router

	route = martini.NewRouter()

	route.Group("/api/v1.0/session", func(a martini.Router) {
		// Проверка токена с продлением +
		a.Get("/:token", middlewares.RequireSessionKeepWithRoute, controllers.KeepSession)
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, controllers.KeepSession)
		// Завершение сеанса пользователя +
		a.Delete("/:token", middlewares.RequireSessionKeepWithRoute, controllers.DeleteSession)
		a.Delete("/", middlewares.RequireSessionKeepWithoutRoute, controllers.DeleteSession)
	})

	route.Group("/api/v1.0/files", func(a martini.Router) {
		// Загрузка файла на сервер +
		a.Post("/", middlewares.RequireSessionKeepWithoutRoute, binding.MultipartForm(models.ViewFile{}), controllers.UploadFile)
		// Отображение картинки по ключу +
		a.Get("/:key/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetFile)
		// Удаление файла на сервере +
		a.Delete("/:key/", middlewares.RequireSessionKeepWithoutRoute, controllers.DeleteFile)
	})

	route.Group("/api/v1.0", func(a martini.Router) {
		// Шаблон домашней страницы +
		a.Get("/", controllers.HomePageTemplate)
		// Проверка доступности сервера +
		a.Get("/ping/:token", middlewares.RequireSessionCheckWithRoute, controllers.Ping)
		a.Get("/ping/", middlewares.RequireSessionCheckWithoutRoute, controllers.Ping)
		// Запрос картинки с капчей +
		a.Get("/captcha/native/", controllers.GetCaptcha)
		// Подтверждение email пользователя +
		a.Post("/emails/confirm/", binding.Json(models.EmailConfirm{}), controllers.ConfirmEmail)
		// Загрузка картинок ?
		a.Get("/images/:type/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetImage)
	})

	route.Group("/api/v1.0/user", func(a martini.Router) {
		// Аутентификация пользователя +
		a.Post("/session/", binding.Json(models.ViewSession{}), controllers.CreateSession)
		// Получение списка групп пользователей +
		a.Get("/groups/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetGroups)
		// Получение информации о пользователе +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetUserInfo)
		// Изменение информации о пользователе +
		a.Put("/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ChangeUser{}), controllers.UpdateUserInfo)
		// Получение информации о e-mail пользователя +
		a.Get("/emails/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetUserEmails)
		// Изменение информации о e-mail пользователя +
		a.Put("/emails/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.UpdateEmails{}), controllers.UpdateUserEmails)
	})

	route.Group("/api/v1.0/users", func(a martini.Router) {
		// Регистрация пользователя +
		a.Post("/register/:token", binding.Json(models.ViewUser{}), controllers.Register)
		a.Post("/register/", binding.Json(models.ViewUser{}), controllers.Register)
		// Восстановление пароля пользователя +
		a.Post("/password/", binding.Json(models.ViewUser{}), controllers.RestorePassword)
		// Смена пароля при восстановлении +
		a.Put("/password/:code/", binding.Json(models.PasswordUpdate{}), controllers.UpdatePassword)
		a.Put("/password/", binding.Json(models.PasswordUpdate{}), controllers.UpdatePassword)
		// Отказ от восстановления пароля +
		a.Delete("/password/:code/", controllers.DeletePasswordRestoring)
	})

	route.Group("/api/v1.0/administration", func(a martini.Router) {
		// Получение списка всех групп доступа +
		a.Get("/groups/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetGroupsInfo)
		// Получение информации о данных списка пользователей +
		a.Options("/users/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetUserMetaData)
		// Получить список пользователей +
		a.Get("/users/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetUsers) //массив ошибок
		// Добавление нового пользователя +
		a.Post("/users/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewApiUserFull{}), controllers.CreateUser)
		// Получение подробной информации о пользователе +
		a.Get("/users/:userid/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetUserFullInfo)
		// Изменение пользователя +
		a.Put("/users/:userid/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewApiUserFull{}), controllers.UpdateUser)
		// Удаление пользователя +
		a.Delete("/users/:userid/", middlewares.RequireSessionKeepWithoutRoute, controllers.DeleteUser)
	})

	route.Group("/api/v1.0/services", func(a martini.Router) {
		// Получение списка услуг поставщиков +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetFacilities)
	})

	route.Group("/api/v1.0/supplier", func(a martini.Router) {
		// Получение списка услуг поставщиков +
		a.Options("/services/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetFacilities)
	})

	route.Group("/api/v1.0/tables", func(a martini.Router) {
		// Получение списка типов таблиц +
		a.Options("/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetTableTypes)
		// Создание таблицы +
		a.Post("/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewShortCustomerTable{}), controllers.CreateTable)
		// Получение списка таблиц +
		a.Get("/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetUnitTables)
		// Получение списка типов колонок +
		a.Get("/fieldtypes/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetColumnTypes)
		// Импорт таблицы ?
		a.Post("/import/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewImportTable{}), controllers.ImportDataFromFile)
		// Получение списка колонок импортируемой таблицы ?
		a.Get("/import/:tmpid/columns/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetImportDataColumns)
		// Сохранение списка импортируемых колонок и присвоение типа колонкам ?
		a.Put("/import/:tmpid/columns/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewImportColumns{}), controllers.UpdateImportDataColumns)
		// Получение информации об экспорте данных +
		a.Options("/export/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetExportDataMeta)
		//  Получение экспортируемых данных данных +
		a.Get("/export/:token/:fid/", controllers.GetExportData)
		// Получение таблицы +
		a.Get("/:tid/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetTable)
		// Изменение таблицы +
		a.Put("/:tid/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewLongCustomerTable{}), controllers.UpdateTable)
		// Удаление таблицы +
		a.Delete("/:tid/", middlewares.RequireSessionKeepWithoutRoute, controllers.DeleteTable)
		// Создание колонки в таблице +
		a.Post("/:tid/field/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewApiTableColumn{}), controllers.CreateTableColumn)
		// Получение списка колонок таблицы в порядке отображения +
		a.Get("/:tid/field/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetTableColumns)
		// Получение колонки таблицы +
		a.Get("/:tid/field/:cid/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetTableColumn)
		// Изменение колонки в таблице +
		a.Put("/:tid/field/:cid/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewApiTableColumn{}), controllers.UpdateTableColumn)
		// Удаление колонки в таблице +
		a.Delete("/:tid/field/:cid/", middlewares.RequireSessionKeepWithoutRoute, controllers.DeleteTableColumn)
		// Изменение порядка отображения колонки +
		a.Put("/:tid/sequence/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewApiOrderTableColumns{}), controllers.UpdateOrderTableColumn)
		// Получение информации о данных в таблице +
		a.Options("/:tid/data/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetTableMetaData)
		// Получение данных таблицы +
		a.Get("/:tid/data/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetTableData)
		// Получение строки данных таблицы +
		a.Get("/:tid/data/:rid/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetTableRow)
		// Внесение строки данных в таблицу +
		a.Post("/:tid/data/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewApiTableRow{}), controllers.CreateTableRow)
		// Изменение строки данных в таблице +
		a.Put("/:tid/data/:rid/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewApiTableRow{}), controllers.UpdateTableRow)
		// Удаление строки данных из таблицы +
		a.Delete("/:tid/data/:rid/", middlewares.RequireSessionKeepWithoutRoute, controllers.DeleteTableRow)
		// Получение информации о данных в ячейке таблицы +
		a.Options("/:tid/cell/:rid/:cid/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetTableMetaCell)
		// Получение данных ячейки таблицы +
		a.Get("/:tid/cell/:rid/:cid/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetTableCell)
		// Изменение данных ячейки таблицы +
		a.Put("/:tid/cell/:rid/:cid/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewTableCell{}), controllers.UpdateTableCell)
		// Удаление данных ячейки таблицы +
		a.Delete("/:tid/cell/:rid/:cid/", middlewares.RequireSessionKeepWithoutRoute, controllers.DeleteTableCell)
		// Экспорт таблицы ?
		a.Get("/:tid/export/", middlewares.RequireSessionKeepWithoutRoute, controllers.ExportDataToFile)
		// Проверка статуса готовности экспортируемого файла +
		a.Options("/:tid/export/:fid/", middlewares.RequireSessionKeepWithoutRoute, controllers.GetExportDataStatus)
		// Изменение настроек для таблицы являющейся прайс-листом  +
		a.Put("/:tid/price/", middlewares.RequireSessionKeepWithoutRoute, binding.Json(models.ViewApiPriceProperties{}), controllers.UpdatePriceTable)
	})

	route.NotFound(middlewares.Default)
	return route
}
