package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type ResultTableRepository interface {
	Get(order_id int64, customer_table_id int64) (resulttable *models.DtoResultTable, err error)
	GetByOrder(order_id int64) (resulttables *[]models.ApiResultTable, err error)
	Create(dtoresulttable *models.DtoResultTable, trans *gorp.Transaction) (err error)
	DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error)
}

type ResultTableService struct {
	*Repository
}

func NewResultTableService(repository *Repository) *ResultTableService {
	repository.DbContext.AddTableWithName(models.DtoResultTable{}, repository.Table).SetKeys(false, "order_id", "customer_table_id")
	return &ResultTableService{Repository: repository}
}

func (resulttableservice *ResultTableService) Get(order_id int64, customer_table_id int64) (resulttable *models.DtoResultTable, err error) {
	resulttable = new(models.DtoResultTable)
	err = resulttableservice.DbContext.SelectOne(resulttable, "select * from "+resulttableservice.Table+
		" where order_id = ? and customer_table_id = ?", order_id, customer_table_id)
	if err != nil {
		log.Error("Error during getting result table object from database %v with value %v, %v", err, order_id, customer_table_id)
		return nil, err
	}

	return resulttable, nil
}

func (resulttableservice *ResultTableService) GetByOrder(order_id int64) (resulttables *[]models.ApiResultTable, err error) {
	resulttables = new([]models.ApiResultTable)
	_, err = resulttableservice.DbContext.Select(resulttables,
		"select r.customer_table_id, t.created, t.type_id from "+resulttableservice.Table+
			" r inner join customer_tables t on r.customer_table_id = t.id where r.order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all result table object from database %v with value %v", err, order_id)
		return nil, err
	}

	return resulttables, nil
}

func (resulttableservice *ResultTableService) Create(dtoresulttable *models.DtoResultTable, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoresulttable)
	} else {
		err = resulttableservice.DbContext.Insert(dtoresulttable)
	}
	if err != nil {
		log.Error("Error during creating resiult table object in database %v", err)
		return err
	}

	return nil
}

func (resulttableservice *ResultTableService) DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+resulttableservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = resulttableservice.DbContext.Exec("delete from "+resulttableservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting result table objects for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
