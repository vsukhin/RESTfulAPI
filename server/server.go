package server

import (
	"application/config"
	"application/db"
	"application/services"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"os"
	"runtime"
)

var userservice *services.UserService
var sessionservice *services.SessionService
var groupservice *services.GroupService
var fileservice *services.FileService
var emailservice *services.EmailService
var captchaservice *services.CaptchaService
var unitservice *services.UnitService
var templateservice *services.TemplateService
var customertableservice *services.CustomerTableService
var tabletypeservice *services.TableTypeService
var columntypeservice *services.ColumnTypeService
var tablecolumnservice *services.TableColumnService
var tablerowservice *services.TableRowService
var facilityservice *services.FacilityService
var pricepropertiesservice *services.PricePropertiesService
var dataformatservice *services.DataFormatService
var virtualdirservice *services.VirtualDirService
var importstepservice *services.ImportStepService
var orderservice *services.OrderService
var statusservice *services.StatusService
var orderstatusservice *services.OrderStatusService
var messageservice *services.MessageService
var projectservice *services.ProjectService
var classifierservice *services.ClassifierService
var mobilephoneservice *services.MobilePhoneService

func Start() {
	var err error

	// Setting to use the maximum number of sockets and cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	if config.InitConfig() != nil {
		return
	}
	if config.InitLogging() != nil {
		return
	}
	if config.InitI18n() != nil {
		return
	}
	if db.InitDB() != nil {
		return
	}

	// Change working directory
	log.Info("Working directory is: '%s'", config.Configuration.WorkingDirectory)
	err = os.Chdir(config.Configuration.WorkingDirectory)
	if err != nil {
		log.Fatalf("Can't change working directory: %v", err)
		return
	}

	userservice = services.NewUserService(services.NewRepository(db.DbMap, db.TABLE_USERS))
	sessionservice = services.NewSessionService(services.NewRepository(db.DbMap, db.TABLE_SESSIONS))
	groupservice = services.NewGroupService(services.NewRepository(db.DbMap, db.TABLE_GROUPS))
	fileservice = services.NewFileService(services.NewRepository(db.DbMap, db.TABLE_FILES))
	emailservice = services.NewEmailService(services.NewRepository(db.DbMap, db.TABLE_EMAILS))
	captchaservice = services.NewCaptchaService(services.NewRepository(db.DbMap, db.TABLE_CAPTCHAS))
	unitservice = services.NewUnitService(services.NewRepository(db.DbMap, db.TABLE_UNITS))
	templateservice = services.NewTemplateService()
	customertableservice = services.NewCustomerTableService(services.NewRepository(db.DbMap, db.TABLE_CUSTOMER_TABLES))
	tabletypeservice = services.NewTableTypeService(services.NewRepository(db.DbMap, db.TABLE_TABLE_TYPES))
	columntypeservice = services.NewColumnTypeService(services.NewRepository(db.DbMap, db.TABLE_COLUMN_TYPES))
	tablecolumnservice = services.NewTableColumnService(services.NewRepository(db.DbMap, db.TABLE_TABLE_COLUMNS))
	tablerowservice = services.NewTableRowService(services.NewRepository(db.DbMap, db.TABLE_TABLE_DATA))
	facilityservice = services.NewFacilityService(services.NewRepository(db.DbMap, db.TABLE_FACILITIES))
	pricepropertiesservice = services.NewPricePropertiesService(services.NewRepository(db.DbMap, db.TABLE_PRICE_PROPERTIES))
	dataformatservice = services.NewDataFormatService(services.NewRepository(db.DbMap, db.TABLE_DATA_FORMATS))
	virtualdirservice = services.NewVirtualDirService(services.NewRepository(db.DbMap, db.TABLE_VIRTUAL_DIRS))
	importstepservice = services.NewImportStepService(services.NewRepository(db.DbMap, db.TABLE_IMPORT_STEPS))
	orderservice = services.NewOrderService(services.NewRepository(db.DbMap, db.TABLE_ORDERS))
	statusservice = services.NewStatusService(services.NewRepository(db.DbMap, db.TABLE_STATUSES))
	orderstatusservice = services.NewOrderStatusService(services.NewRepository(db.DbMap, db.TABLE_ORDER_STATUSES))
	messageservice = services.NewMessageService(services.NewRepository(db.DbMap, db.TABLE_MESSAGES))
	projectservice = services.NewProjectService(services.NewRepository(db.DbMap, db.TABLE_PROJECTS))
	classifierservice = services.NewClassifierService(services.NewRepository(db.DbMap, db.TABLE_CLASSIFIERS))
	mobilephoneservice = services.NewMobilePhoneService(services.NewRepository(db.DbMap, db.TABLE_MOBILE_PHONES))

	userservice.SessionRepository = sessionservice
	userservice.EmailRepository = emailservice
	userservice.GroupRepository = groupservice
	userservice.UnitRepository = unitservice
	userservice.MessageRepository = messageservice
	userservice.MobilePhoneRepository = mobilephoneservice
	sessionservice.GroupRepository = groupservice
	customertableservice.TableColumnRepository = tablecolumnservice
	customertableservice.TableRowRepository = tablerowservice
	orderservice.OrderStatusRepository = orderstatusservice

	go fileservice.ClearExpiredFiles()
	go customertableservice.ClearExpiredTables()

	routes := Routes()
	mrt := martini.New()

	mrt.Handlers(
		LogRequest,
		bootstrap(),
		martini.Recovery(),
		render.Renderer(render.Options{}),
	)

	// File server
	log.Info("Server DocumentRoot is: '%s'", config.Configuration.Server.DocumentRoot)
	mrt.Use(martini.Static(config.Configuration.Server.DocumentRoot,
		martini.StaticOptions{
			Exclude: "/api/",
		}))

	mrt.MapTo(routes, (*martini.Routes)(nil))
	mrt.Action(routes.Handle)

	if err = http.ListenAndServe(config.Configuration.Server.Address, mrt); err != nil {
		log.Fatalf("Can't launch http server %v", err)
		return
	}
}

func Stop() {
	db.ShutdownDB()
	config.ShutdownLogging()
}

func bootstrap() martini.Handler {
	return func(context martini.Context) {
		context.Map(userservice)
		context.Map(sessionservice)
		context.Map(groupservice)
		context.Map(fileservice)
		context.Map(emailservice)
		context.Map(captchaservice)
		context.Map(unitservice)
		context.Map(templateservice)
		context.Map(customertableservice)
		context.Map(tabletypeservice)
		context.Map(columntypeservice)
		context.Map(tablecolumnservice)
		context.Map(tablerowservice)
		context.Map(facilityservice)
		context.Map(pricepropertiesservice)
		context.Map(dataformatservice)
		context.Map(virtualdirservice)
		context.Map(importstepservice)
		context.Map(orderservice)
		context.Map(statusservice)
		context.Map(orderstatusservice)
		context.Map(messageservice)
		context.Map(projectservice)
		context.Map(classifierservice)
		context.Map(mobilephoneservice)
	}
}
