/* DB package provides methods and data structures responsible for exchanging data with sql server */

package db

import (
	"application/config"
	"database/sql"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	logging "github.com/op/go-logging"
	_ "github.com/ziutek/mymysql/godrv"
)

const (
	MYSQL_CONFIG_NUMBER = 0

	TABLE_USERS                      = "users"
	TABLE_SESSIONS                   = "sessions"
	TABLE_GROUPS                     = "groups"
	TABLE_FILES                      = "files"
	TABLE_EMAILS                     = "emails"
	TABLE_CAPTCHAS                   = "captchas"
	TABLE_UNITS                      = "units"
	TABLE_CUSTOMER_TABLES            = "customer_tables"
	TABLE_TABLE_TYPES                = "table_types"
	TABLE_COLUMN_TYPES               = "column_types"
	TABLE_TABLE_COLUMNS              = "table_columns"
	TABLE_TABLE_DATA                 = "table_data"
	TABLE_FACILITIES                 = "services"
	TABLE_PRICE_PROPERTIES           = "price_properties"
	TABLE_DATA_FORMATS               = "data_formats"
	TABLE_VIRTUAL_DIRS               = "virtual_dirs"
	TABLE_IMPORT_STEPS               = "import_steps"
	TABLE_ORDERS                     = "orders"
	TABLE_MESSAGES                   = "messages"
	TABLE_ORDER_STATUSES             = "order_statuses"
	TABLE_STATUSES                   = "statuses"
	TABLE_PROJECTS                   = "projects"
	TABLE_CLASSIFIERS                = "classifiers"
	TABLE_MOBILE_PHONES              = "mobile_phones"
	TABLE_COMPANIES                  = "companies"
	TABLE_PERIODS                    = "periods"
	TABLE_EVENTS                     = "events"
	TABLE_MOBILE_OPERATORS           = "mobile_operators"
	TABLE_NEWS                       = "news"
	TABLE_SUBSCRIPTIONS              = "subscriptions"
	TABLE_REQUESTS                   = "requests"
	TABLE_SMS_SENDERS                = "sms_senders"
	TABLE_FACILITY_TYPES             = "service_types"
	TABLE_SMS_FACILITIES             = "sms_services"
	TABLE_HLR_FACILITIES             = "hlr_services"
	TABLE_RECOGNIZE_FACILITIES       = "input_services"
	TABLE_VERIFY_FACILITIES          = "data_services"
	TABLE_MOBILE_OPERATOR_OPERATIONS = "mobile_operator_operations"
	TABLE_RESULT_TABLES              = "result_tables"
	TABLE_WORK_TABLES                = "work_tables"
	TABLE_DATA_COLUMNS               = "data_columns"
	TABLE_INPUT_FIELDS               = "input_fields"
	TABLE_INPUT_FILES                = "input_files"
	TABLE_SUPPLIER_REQUESTS          = "supplier_requests"
	TABLE_INPUT_FTPS                 = "input_ftps"
	TABLE_COMPANY_TYPES              = "company_types"
	TABLE_COMPANY_CLASSES            = "company_classes"
	TABLE_ADDRESS_TYPES              = "address_types"
	TABLE_COMPANY_CODES              = "company_codes"
	TABLE_COMPANY_ADDRESSES          = "company_addresses"
	TABLE_COMPANY_BANKS              = "company_banks"
	TABLE_COMPANY_STAFF              = "company_staff"
	TABLE_SUPPLIER_FACILITIES        = "supplier_services"
	TABLE_COMPLEX_STATUSES           = "complex_statuses"
	TABLE_INVOICES                   = "invoices"
	TABLE_INVOICE_ITEMS              = "invoice_items"
	TABLE_RECOGNIZE_PRODUCTS         = "input_products"
	TABLE_VERIFY_PRODUCTS            = "data_products"
	TABLE_INPUT_PRODUCTS             = "input_order_products"
	TABLE_DATA_PRODUCTS              = "data_order_products"
	TABLE_FEEDBACK                   = "feedback"
	TABLE_DEVICES                    = "devices"
	TABLE_SMS_PERIODS                = "sms_periods"
	TABLE_SMS_EVENTS                 = "sms_events"
	TABLE_REPORTS                    = "reports"
	TABLE_REPORT_PERIODS             = "report_periods"
	TABLE_REPORT_PROJECTS            = "report_projects"
	TABLE_REPORT_ORDERS              = "report_orders"
	TABLE_REPORT_FACILITIES          = "report_services"
	TABLE_REPORT_COMPLEX_STATUSES    = "report_complex_statuses"
	TABLE_REPORT_SUPPLIERS           = "report_suppliers"
	TABLE_REPORT_SETTINGS            = "report_settings"
	TABLE_TRANSACTIONS               = "transactions"
	TABLE_OPERATIONS                 = "operations"
	TABLE_TRANSACTION_TYPES          = "transaction_types"
	TABLE_OPERATION_TYPES            = "operation_types"
	TABLE_ORDER_INVOICES             = "order_invoices"
	TABLE_ACCESS_LOG                 = "access_log"
	TABLE_DOCUMENT_TYPES             = "document_types"
	TABLE_CONTRACTS                  = "contracts"
	TABLE_APPENDICES                 = "appendices"
	TABLE_DOCUMENTS                  = "documents"
	TABLE_DATA_ENCODINGS             = "data_encodings"
	TABLE_HEADER_FACILITIES          = "header_services"
	TABLE_TARIFF_PLANS               = "tariff_plans"
	TABLE_PAYMENTS                   = "payments"
	TABLE_HEADER_PRODUCTS            = "header_products"
)

var (
	DbMap *gorp.DbMap
	log   config.Logger = logging.MustGetLogger("db")
)

func InitLogger(logger config.Logger) {
	log = logger
}

func InitDB() (err error) {
	// If mysql exists
	if len(config.Configuration.MySql) > 0 {
		cfgMySql := new(config.MysqlConfiguration)

		cfgMySql.Driver = config.Configuration.MySql[MYSQL_CONFIG_NUMBER].Driver
		cfgMySql.Host = config.Configuration.MySql[MYSQL_CONFIG_NUMBER].Host
		cfgMySql.Port = config.Configuration.MySql[MYSQL_CONFIG_NUMBER].Port
		cfgMySql.Type = config.Configuration.MySql[MYSQL_CONFIG_NUMBER].Type
		cfgMySql.Socket = config.Configuration.MySql[MYSQL_CONFIG_NUMBER].Socket
		cfgMySql.Name = config.Configuration.MySql[MYSQL_CONFIG_NUMBER].Name
		cfgMySql.Login = config.Configuration.MySql[MYSQL_CONFIG_NUMBER].Login
		cfgMySql.Password = config.Configuration.MySql[MYSQL_CONFIG_NUMBER].Password
		cfgMySql.Charset = config.Configuration.MySql[MYSQL_CONFIG_NUMBER].Charset

		connect := ""
		switch cfgMySql.Type {
		case "tcp":
			connect = "tcp:" + fmt.Sprintf("%s:%d", cfgMySql.Host, cfgMySql.Port)
		case "socket":
			connect = "unix:" + cfgMySql.Socket
		default:
			log.Error("Unknown type of connection protocol %s", cfgMySql.Type)
			return errors.New("Unknown protocol")
		}

		if connect != "" {
			connect += "*" + cfgMySql.Name + "/" + cfgMySql.Login + "/" + cfgMySql.Password

			db, err := sql.Open(cfgMySql.Driver, connect)
			if err == nil {
				DbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", cfgMySql.Charset}}
			} else {
				log.Error("MySQL connection error: %v", err)
				return err
			}
		}
	} else {
		log.Error("MySQL configuration is empty!")
		return errors.New("Empty configuration")
	}

	return nil
}

func ShutdownDB() {
	if DbMap != nil {
		err := DbMap.Db.Close()
		if err != nil {
			log.Error("MySQL close connection error: %v", err)
		}
	}
}
