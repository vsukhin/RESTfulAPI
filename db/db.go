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

	TABLE_USERS            = "users"
	TABLE_SESSIONS         = "sessions"
	TABLE_GROUPS           = "groups"
	TABLE_FILES            = "files"
	TABLE_EMAILS           = "emails"
	TABLE_CAPTCHAS         = "captchas"
	TABLE_UNITS            = "units"
	TABLE_CUSTOMER_TABLES  = "customer_tables"
	TABLE_TABLE_TYPES      = "table_types"
	TABLE_COLUMN_TYPES     = "column_types"
	TABLE_TABLE_COLUMNS    = "table_columns"
	TABLE_TABLE_DATA       = "table_data"
	TABLE_FACILITIES       = "services"
	TABLE_PRICE_PROPERTIES = "price_properties"
	TABLE_DATA_FORMATS     = "data_formats"
	TABLE_VIRTUAL_DIRS     = "virtual_dirs"
	TABLE_IMPORT_STEPS     = "import_steps"
	TABLE_ORDERS           = "orders"
	TABLE_MESSAGES         = "messages"
	TABLE_ORDER_STATUSES   = "order_statuses"
	TABLE_STATUSES         = "statuses"
	TABLE_PROJECTS         = "projects"
	TABLE_CLASSIFIERS      = "classifiers"
	TABLE_MOBILE_PHONES    = "mobile_phones"
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
		cfgMySql.Socket = config.Configuration.MySql[MYSQL_CONFIG_NUMBER].Socker
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
			return errors.New("Uknown protocol")
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
