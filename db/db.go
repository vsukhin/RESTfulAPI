/* DB package provides methods and data structures responsible for exchanging data with sql server */

package db

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/ziutek/mymysql/godrv"

	"application/config"
	"fmt"
	logging "github.com/op/go-logging"
)

const (
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
	TABLE_TABLE_ROWS       = "table_data"
	TABLE_TABLE_CELLS      = "table_cells"
	TABLE_FACILITIES       = "services"
	TABLE_PRICE_PROPERTIES = "price_properties"
	TABLE_DATA_FORMATS     = "data_formats"
	TABLE_VIRTUAL_DIRS     = "virtual_dirs"
)

var (
	log   *logging.Logger
	DbMap *gorp.DbMap
)

func init() {
	log = logging.MustGetLogger("db")

	// If mysql exists
	if len(config.Configuration.MySql) > 0 {
		cfgMySql := new(config.MysqlConfiguration)

		cfgMySql.Driver = config.Configuration.MySql[0].Driver
		cfgMySql.Host = config.Configuration.MySql[0].Host
		cfgMySql.Port = config.Configuration.MySql[0].Port
		cfgMySql.Type = config.Configuration.MySql[0].Type
		cfgMySql.Socket = config.Configuration.MySql[0].Socker
		cfgMySql.Name = config.Configuration.MySql[0].Name
		cfgMySql.Login = config.Configuration.MySql[0].Login
		cfgMySql.Password = config.Configuration.MySql[0].Password
		cfgMySql.Charset = config.Configuration.MySql[0].Charset

		connect := ""
		switch cfgMySql.Type {
		case "tcp":
			connect = "tcp:" + fmt.Sprintf("%s:%d", cfgMySql.Host, cfgMySql.Port)
		case "socket":
			connect = "unix:" + cfgMySql.Socket
		default:
			log.Error("Unknown type of connection protocol %s", cfgMySql.Type)
		}

		if connect != "" {
			connect += "*" + cfgMySql.Name + "/" + cfgMySql.Login + "/" + cfgMySql.Password

			db, err := sql.Open(cfgMySql.Driver, connect)
			if err == nil {
				DbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", cfgMySql.Charset}}
			} else {
				log.Error("MySQL connection error: %v", err)
			}
		}
	} else {
		log.Error("MySQL configuration is empty!")
	}
}
