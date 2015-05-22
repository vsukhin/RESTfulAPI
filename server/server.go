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
var companyservice *services.CompanyService
var periodservice *services.PeriodService
var eventservice *services.EventService
var mobileoperatorservice *services.MobileOperatorService
var newsservice *services.NewsService
var subscriptionservice *services.SubscriptionService
var requestservice *services.RequestService
var smssenderservice *services.SMSSenderService
var facilitytypeservice *services.FacilityTypeService
var smsfacilityservice *services.SMSFacilityService
var hlrfacilityservice *services.HLRFacilityService
var recognizefacilityservice *services.RecognizeFacilityService
var verifyfacilityservice *services.VerifyFacilityService
var mobileoperatoroperationservice *services.MobileOperatorOperationService
var resulttableservice *services.ResultTableService
var worktableservice *services.WorkTableService
var datacolumnservice *services.DataColumnService
var inputfieldservice *services.InputFieldService
var inputfileservice *services.InputFileService
var supplierrequestservice *services.SupplierRequestService
var inputftpservice *services.InputFtpService
var companytypeservice *services.CompanyTypeService
var companyclassservice *services.CompanyClassService
var addresstypeservice *services.AddressTypeService
var companycodeservice *services.CompanyCodeService
var companyaddressservice *services.CompanyAddressService
var companybankservice *services.CompanyBankService
var companyemployeeservice *services.CompanyEmployeeService
var facilitytableservice *services.FacilityTableService
var smstableservice *services.SMSTableService
var hlrtableservice *services.HLRTableService
var verifytableservice *services.VerifyTableService
var supplierfacilityservice *services.SupplierFacilityService
var complexstatusservice *services.ComplexStatusService
var invoiceservice *services.InvoiceService
var invoiceitemservice *services.InvoiceItemService

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
	companyservice = services.NewCompanyService(services.NewRepository(db.DbMap, db.TABLE_COMPANIES))
	periodservice = services.NewPeriodService(services.NewRepository(db.DbMap, db.TABLE_PERIODS))
	eventservice = services.NewEventService(services.NewRepository(db.DbMap, db.TABLE_EVENTS))
	mobileoperatorservice = services.NewMobileOperatorService(services.NewRepository(db.DbMap, db.TABLE_MOBILE_OPERATORS))
	newsservice = services.NewNewsService(services.NewRepository(db.DbMap, db.TABLE_NEWS))
	subscriptionservice = services.NewSubscriptionService(services.NewRepository(db.DbMap, db.TABLE_SUBSCRIPTIONS))
	requestservice = services.NewRequestService(services.NewRepository(db.DbMap, db.TABLE_REQUESTS))
	smssenderservice = services.NewSMSSenderService(services.NewRepository(db.DbMap, db.TABLE_SMS_SENDERS))
	facilitytypeservice = services.NewFacilityTypeService(services.NewRepository(db.DbMap, db.TABLE_FACILITY_TYPES))
	smsfacilityservice = services.NewSMSFacilityService(services.NewRepository(db.DbMap, db.TABLE_SMS_FACILITIES))
	hlrfacilityservice = services.NewHLRFacilityService(services.NewRepository(db.DbMap, db.TABLE_HLR_FACILITIES))
	recognizefacilityservice = services.NewRecognizeFacilityService(services.NewRepository(db.DbMap, db.TABLE_RECOGNIZE_FACILITIES))
	verifyfacilityservice = services.NewVerifyFacilityService(services.NewRepository(db.DbMap, db.TABLE_VERIFY_FACILITIES))
	mobileoperatoroperationservice = services.NewMobileOperatorOperationService(services.NewRepository(db.DbMap, db.TABLE_MOBILE_OPERATOR_OPERATIONS))
	resulttableservice = services.NewResultTableService(services.NewRepository(db.DbMap, db.TABLE_RESULT_TABLES))
	worktableservice = services.NewWorkTableService(services.NewRepository(db.DbMap, db.TABLE_WORK_TABLES))
	datacolumnservice = services.NewDataColumnService(services.NewRepository(db.DbMap, db.TABLE_DATA_COLUMNS))
	inputfieldservice = services.NewInputFieldService(services.NewRepository(db.DbMap, db.TABLE_INPUT_FIELDS))
	inputfileservice = services.NewInputFileService(services.NewRepository(db.DbMap, db.TABLE_INPUT_FILES))
	supplierrequestservice = services.NewSupplierRequestService(services.NewRepository(db.DbMap, db.TABLE_SUPPLIER_REQUESTS))
	inputftpservice = services.NewInputFtpService(services.NewRepository(db.DbMap, db.TABLE_INPUT_FTPS))
	companytypeservice = services.NewCompanyTypeService(services.NewRepository(db.DbMap, db.TABLE_COMPANY_TYPES))
	companyclassservice = services.NewCompanyClassService(services.NewRepository(db.DbMap, db.TABLE_COMPANY_CLASSES))
	addresstypeservice = services.NewAddressTypeService(services.NewRepository(db.DbMap, db.TABLE_ADDRESS_TYPES))
	companycodeservice = services.NewCompanyCodeService(services.NewRepository(db.DbMap, db.TABLE_COMPANY_CODES))
	companyaddressservice = services.NewCompanyAddressService(services.NewRepository(db.DbMap, db.TABLE_COMPANY_ADDRESSES))
	companybankservice = services.NewCompanyBankService(services.NewRepository(db.DbMap, db.TABLE_COMPANY_BANKS))
	companyemployeeservice = services.NewCompanyEmployeeService(services.NewRepository(db.DbMap, db.TABLE_COMPANY_STAFF))
	facilitytableservice = services.NewFacilityTableService(services.NewRepository(db.DbMap, db.TABLE_TABLE_COLUMNS))
	smstableservice = services.NewSMSTableService(services.NewRepository(db.DbMap, db.TABLE_CUSTOMER_TABLES))
	hlrtableservice = services.NewHLRTableService(services.NewRepository(db.DbMap, db.TABLE_CUSTOMER_TABLES))
	verifytableservice = services.NewVerifyTableService(services.NewRepository(db.DbMap, db.TABLE_CUSTOMER_TABLES))
	supplierfacilityservice = services.NewSupplierFacilityService(services.NewRepository(db.DbMap, db.TABLE_SUPPLIER_FACILITIES))
	complexstatusservice = services.NewComplexStatusService(services.NewRepository(db.DbMap, db.TABLE_COMPLEX_STATUSES))
	invoiceservice = services.NewInvoiceService(services.NewRepository(db.DbMap, db.TABLE_INVOICES))
	invoiceitemservice = services.NewInvoiceItemService(services.NewRepository(db.DbMap, db.TABLE_INVOICE_ITEMS))

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

	smsfacilityservice.MobileOperatorOperationRepository = mobileoperatoroperationservice
	smsfacilityservice.ResultTableRepository = resulttableservice
	smsfacilityservice.WorkTableRepository = worktableservice

	hlrfacilityservice.MobileOperatorOperationRepository = mobileoperatoroperationservice
	hlrfacilityservice.ResultTableRepository = resulttableservice
	hlrfacilityservice.WorkTableRepository = worktableservice

	recognizefacilityservice.InputFieldRepository = inputfieldservice
	recognizefacilityservice.InputFileRepository = inputfileservice
	recognizefacilityservice.SupplierRequestRepository = supplierrequestservice
	recognizefacilityservice.InputFtpRepository = inputftpservice
	recognizefacilityservice.ResultTableRepository = resulttableservice
	recognizefacilityservice.WorkTableRepository = worktableservice

	verifyfacilityservice.DataColumnRepository = datacolumnservice
	verifyfacilityservice.ResultTableRepository = resulttableservice
	verifyfacilityservice.WorkTableRepository = worktableservice

	companyservice.CompanyCodeRepository = companycodeservice
	companyservice.CompanyAddressRepository = companyaddressservice
	companyservice.CompanyBankRepository = companybankservice
	companyservice.CompanyEmployeeRepository = companyemployeeservice

	smstableservice.FacilityTableRepository = facilitytableservice

	hlrtableservice.FacilityTableRepository = facilitytableservice

	verifytableservice.FacilityTableRepository = facilitytableservice

	invoiceservice.InvoiceItemRepository = invoiceitemservice

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
		},
		martini.StaticOptions{
			Exclude: "/subscriptions/",
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
		context.Map(companyservice)
		context.Map(periodservice)
		context.Map(eventservice)
		context.Map(mobileoperatorservice)
		context.Map(newsservice)
		context.Map(subscriptionservice)
		context.Map(requestservice)
		context.Map(smssenderservice)
		context.Map(facilitytypeservice)
		context.Map(smsfacilityservice)
		context.Map(hlrfacilityservice)
		context.Map(recognizefacilityservice)
		context.Map(verifyfacilityservice)
		context.Map(mobileoperatoroperationservice)
		context.Map(resulttableservice)
		context.Map(worktableservice)
		context.Map(datacolumnservice)
		context.Map(inputfieldservice)
		context.Map(inputfileservice)
		context.Map(supplierrequestservice)
		context.Map(inputftpservice)
		context.Map(companytypeservice)
		context.Map(companyclassservice)
		context.Map(addresstypeservice)
		context.Map(companycodeservice)
		context.Map(companyaddressservice)
		context.Map(companybankservice)
		context.Map(companyemployeeservice)
		context.Map(facilitytableservice)
		context.Map(smstableservice)
		context.Map(hlrtableservice)
		context.Map(verifytableservice)
		context.Map(supplierfacilityservice)
		context.Map(complexstatusservice)
		context.Map(invoiceservice)
		context.Map(invoiceitemservice)
	}
}
