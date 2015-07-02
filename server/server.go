package server

import (
	"net/http"
	"os"
	"runtime"

	"application/communication"
	"application/config"
	"application/db"
	"application/services"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	userservice                    *services.UserService
	sessionservice                 *services.SessionService
	groupservice                   *services.GroupService
	fileservice                    *services.FileService
	emailservice                   *services.EmailService
	captchaservice                 *services.CaptchaService
	unitservice                    *services.UnitService
	templateservice                *services.TemplateService
	customertableservice           *services.CustomerTableService
	tabletypeservice               *services.TableTypeService
	columntypeservice              *services.ColumnTypeService
	tablecolumnservice             *services.TableColumnService
	tablerowservice                *services.TableRowService
	facilityservice                *services.FacilityService
	pricepropertiesservice         *services.PricePropertiesService
	dataformatservice              *services.DataFormatService
	virtualdirservice              *services.VirtualDirService
	importstepservice              *services.ImportStepService
	orderservice                   *services.OrderService
	statusservice                  *services.StatusService
	orderstatusservice             *services.OrderStatusService
	messageservice                 *services.MessageService
	projectservice                 *services.ProjectService
	classifierservice              *services.ClassifierService
	mobilephoneservice             *services.MobilePhoneService
	companyservice                 *services.CompanyService
	periodservice                  *services.PeriodService
	eventservice                   *services.EventService
	mobileoperatorservice          *services.MobileOperatorService
	newsservice                    *services.NewsService
	subscriptionservice            *services.SubscriptionService
	requestservice                 *services.RequestService
	smssenderservice               *services.SMSSenderService
	facilitytypeservice            *services.FacilityTypeService
	smsfacilityservice             *services.SMSFacilityService
	hlrfacilityservice             *services.HLRFacilityService
	recognizefacilityservice       *services.RecognizeFacilityService
	verifyfacilityservice          *services.VerifyFacilityService
	mobileoperatoroperationservice *services.MobileOperatorOperationService
	resulttableservice             *services.ResultTableService
	worktableservice               *services.WorkTableService
	datacolumnservice              *services.DataColumnService
	inputfieldservice              *services.InputFieldService
	inputfileservice               *services.InputFileService
	supplierrequestservice         *services.SupplierRequestService
	inputftpservice                *services.InputFtpService
	companytypeservice             *services.CompanyTypeService
	companyclassservice            *services.CompanyClassService
	addresstypeservice             *services.AddressTypeService
	companycodeservice             *services.CompanyCodeService
	companyaddressservice          *services.CompanyAddressService
	companybankservice             *services.CompanyBankService
	companyemployeeservice         *services.CompanyEmployeeService
	facilitytableservice           *services.FacilityTableService
	smstableservice                *services.SMSTableService
	hlrtableservice                *services.HLRTableService
	verifytableservice             *services.VerifyTableService
	supplierfacilityservice        *services.SupplierFacilityService
	complexstatusservice           *services.ComplexStatusService
	invoiceservice                 *services.InvoiceService
	invoiceitemservice             *services.InvoiceItemService
	priceservice                   *services.PriceService
	recognizeproductservice        *services.RecognizeProductService
	verifyproductservice           *services.VerifyProductService
	inputproductservice            *services.InputProductService
	dataproductservice             *services.DataProductService
	feedbackservice                *services.FeedbackService
	deviceservice                  *services.DeviceService
	smsperiodservice               *services.SMSPeriodService
	smseventservice                *services.SMSEventService
	reportservice                  *services.ReportService
	reportperiodservice            *services.ReportPeriodService
	reportprojectservice           *services.ReportProjectService
	reportorderservice             *services.ReportOrderService
	reportfacilityservice          *services.ReportFacilityService
	reportcomplexstatusservice     *services.ReportComplexStatusService
	reportsupplierservice          *services.ReportSupplierService
	reportsettingsservice          *services.ReportSettingsService
	complexreportservice           *services.ComplexReportService
	transactionservice             *services.TransactionService
	operationservice               *services.OperationService
	transactiontypeservice         *services.TransactionTypeService
	operationtypeservice           *services.OperationTypeService
	orderinvoiceservice            *services.OrderInvoiceService
)

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
	if err = communication.Init(communication.ModeClient); err != nil {
		return
	}
	if config.InitDbCassandra() != nil {
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
	priceservice = services.NewPriceService(services.NewRepository(db.DbMap, ""))
	recognizeproductservice = services.NewRecognizeProductService(services.NewRepository(db.DbMap, db.TABLE_RECOGNIZE_PRODUCTS))
	verifyproductservice = services.NewVerifyProductService(services.NewRepository(db.DbMap, db.TABLE_VERIFY_PRODUCTS))
	inputproductservice = services.NewInputProductService(services.NewRepository(db.DbMap, db.TABLE_INPUT_PRODUCTS))
	dataproductservice = services.NewDataProductService(services.NewRepository(db.DbMap, db.TABLE_DATA_PRODUCTS))
	feedbackservice = services.NewFeedbackService(services.NewRepository(db.DbMap, db.TABLE_FEEDBACK))
	deviceservice = services.NewDeviceService(services.NewRepository(db.DbMap, db.TABLE_DEVICES))
	smsperiodservice = services.NewSMSPeriodService(services.NewRepository(db.DbMap, db.TABLE_SMS_PERIODS))
	smseventservice = services.NewSMSEventService(services.NewRepository(db.DbMap, db.TABLE_SMS_EVENTS))
	reportservice = services.NewReportService(services.NewRepository(db.DbMap, db.TABLE_REPORTS))
	reportperiodservice = services.NewReportPeriodService(services.NewRepository(db.DbMap, db.TABLE_REPORT_PERIODS))
	reportprojectservice = services.NewReportProjectService(services.NewRepository(db.DbMap, db.TABLE_REPORT_PROJECTS))
	reportorderservice = services.NewReportOrderService(services.NewRepository(db.DbMap, db.TABLE_REPORT_ORDERS))
	reportfacilityservice = services.NewReportFacilityService(services.NewRepository(db.DbMap, db.TABLE_REPORT_FACILITIES))
	reportcomplexstatusservice = services.NewReportComplexStatusService(services.NewRepository(db.DbMap, db.TABLE_REPORT_COMPLEX_STATUSES))
	reportsupplierservice = services.NewReportSupplierService(services.NewRepository(db.DbMap, db.TABLE_REPORT_SUPPLIERS))
	reportsettingsservice = services.NewReportSettingsService(services.NewRepository(db.DbMap, db.TABLE_REPORT_SETTINGS))
	complexreportservice = services.NewComplexReportService(services.NewRepository(db.DbMap, ""))
	transactionservice = services.NewTransactionService(services.NewRepository(db.DbMap, db.TABLE_TRANSACTIONS))
	operationservice = services.NewOperationService(services.NewRepository(db.DbMap, db.TABLE_OPERATIONS))
	transactiontypeservice = services.NewTransactionTypeService(services.NewRepository(db.DbMap, db.TABLE_TRANSACTION_TYPES))
	operationtypeservice = services.NewOperationTypeService(services.NewRepository(db.DbMap, db.TABLE_OPERATION_TYPES))
	orderinvoiceservice = services.NewOrderInvoiceService(services.NewRepository(db.DbMap, db.TABLE_ORDER_INVOICES))

	userservice.SessionRepository = sessionservice
	userservice.EmailRepository = emailservice
	userservice.GroupRepository = groupservice
	userservice.UnitRepository = unitservice
	userservice.MessageRepository = messageservice
	userservice.MobilePhoneRepository = mobilephoneservice
	userservice.DeviceRepository = deviceservice

	sessionservice.GroupRepository = groupservice

	customertableservice.TableColumnRepository = tablecolumnservice
	customertableservice.TableRowRepository = tablerowservice

	orderservice.OrderStatusRepository = orderstatusservice

	smsfacilityservice.MobileOperatorOperationRepository = mobileoperatoroperationservice
	smsfacilityservice.SMSPeriodRepository = smsperiodservice
	smsfacilityservice.SMSEventRepository = smseventservice
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
	recognizefacilityservice.InputProductRepository = inputproductservice

	verifyfacilityservice.DataColumnRepository = datacolumnservice
	verifyfacilityservice.ResultTableRepository = resulttableservice
	verifyfacilityservice.WorkTableRepository = worktableservice
	verifyfacilityservice.DataProductRepository = dataproductservice

	companyservice.CompanyCodeRepository = companycodeservice
	companyservice.CompanyAddressRepository = companyaddressservice
	companyservice.CompanyBankRepository = companybankservice
	companyservice.CompanyEmployeeRepository = companyemployeeservice

	smstableservice.FacilityTableRepository = facilitytableservice

	hlrtableservice.FacilityTableRepository = facilitytableservice

	verifytableservice.FacilityTableRepository = facilitytableservice

	invoiceservice.InvoiceItemRepository = invoiceitemservice
	invoiceservice.TransactionRepository = transactionservice
	invoiceservice.OperationRepository = operationservice
	invoiceservice.OrderInvoiceRepository = orderinvoiceservice
	invoiceservice.OrderStatusRepository = orderstatusservice

	reportservice.UserRepository = userservice
	reportservice.ReportPeriodRepository = reportperiodservice
	reportservice.ReportProjectRepository = reportprojectservice
	reportservice.ReportOrderRepository = reportorderservice
	reportservice.ReportFacilityRepository = reportfacilityservice
	reportservice.ReportComplexStatusRepository = reportcomplexstatusservice
	reportservice.ReportSupplierRepository = reportsupplierservice
	reportservice.ReportSettingsRepository = reportsettingsservice

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
		context.Map(priceservice)
		context.Map(recognizeproductservice)
		context.Map(verifyproductservice)
		context.Map(inputproductservice)
		context.Map(dataproductservice)
		context.Map(feedbackservice)
		context.Map(deviceservice)
		context.Map(smsperiodservice)
		context.Map(smseventservice)
		context.Map(reportservice)
		context.Map(reportperiodservice)
		context.Map(reportprojectservice)
		context.Map(reportorderservice)
		context.Map(reportfacilityservice)
		context.Map(reportcomplexstatusservice)
		context.Map(reportsupplierservice)
		context.Map(reportsettingsservice)
		context.Map(complexreportservice)
		context.Map(transactionservice)
		context.Map(operationservice)
		context.Map(transactiontypeservice)
		context.Map(operationtypeservice)
		context.Map(orderinvoiceservice)
	}
}
