/* Services package provides methods and data structures for database layer implementation */

package services

import (
	"application/config"
	"database/sql"
	"github.com/coopernurse/gorp"
	logging "github.com/op/go-logging"
)

type DbMap interface {
	AddTableWithName(i interface{}, name string) *gorp.TableMap
	Begin() (*gorp.Transaction, error)
	Get(i interface{}, keys ...interface{}) (interface{}, error)
	Insert(list ...interface{}) error
	Update(list ...interface{}) (int64, error)
	Delete(list ...interface{}) (int64, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	SelectInt(query string, args ...interface{}) (int64, error)
	SelectStr(query string, args ...interface{}) (string, error)
	SelectOne(holder interface{}, query string, args ...interface{}) error
}

type Repository struct {
	DbContext DbMap  // объект ORM для работы с бд
	Table     string // имя таблицы репозитория в бд
}

var (
	log config.Logger = logging.MustGetLogger("services")
)

func InitLogger(logger config.Logger) {
	log = logger
}

// Конструктор создания объекта репозитория
func NewRepository(dbmap DbMap, table string) *Repository {
	return &Repository{
		DbContext: dbmap,
		Table:     table,
	}
}
