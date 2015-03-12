package services

import (
	"github.com/coopernurse/gorp"
	logging "github.com/op/go-logging"
)

type Repository struct {
	DbContext *gorp.DbMap // объект ORM для работы с бд
	Table     string      // имя таблицы репозитория в бд
}

var (
	log *logging.Logger
)

func init() {
	log = logging.MustGetLogger("services")
}

// Конструктор создания объекта репозитория
func NewRepository(dbmap *gorp.DbMap, table string) *Repository {
	return &Repository{
		DbContext: dbmap,
		Table:     table,
	}
}
