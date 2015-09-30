package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type WorkTableRepository interface {
	Get(order_id int64, customer_table_id int64) (worktable *models.DtoWorkTable, err error)
	GetByOrder(order_id int64) (worktables *[]models.ApiWorkTable, err error)
	Create(dtoworktable *models.DtoWorkTable, trans *gorp.Transaction) (err error)
	DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error)
}

type WorkTableService struct {
	*Repository
}

func NewWorkTableService(repository *Repository) *WorkTableService {
	repository.DbContext.AddTableWithName(models.DtoWorkTable{}, repository.Table).SetKeys(false, "order_id", "customer_table_id")
	return &WorkTableService{Repository: repository}
}

func (worktableservice *WorkTableService) Get(order_id int64, customer_table_id int64) (worktable *models.DtoWorkTable, err error) {
	worktable = new(models.DtoWorkTable)
	err = worktableservice.DbContext.SelectOne(worktable, "select * from "+worktableservice.Table+
		" where order_id = ? and customer_table_id = ?", order_id, customer_table_id)
	if err != nil {
		log.Error("Error during getting work table object from database %v with value %v, %v", err, order_id, customer_table_id)
		return nil, err
	}

	return worktable, nil
}

func (worktableservice *WorkTableService) GetByOrder(order_id int64) (worktables *[]models.ApiWorkTable, err error) {
	worktables = new([]models.ApiWorkTable)
	_, err = worktableservice.DbContext.Select(worktables,
		"select r.customer_table_id, t.created, t.type_id from "+worktableservice.Table+
			" r inner join customer_tables t on r.customer_table_id = t.id where t.active = 1 and t.permanent = 1 and r.order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all work table object from database %v with value %v", err, order_id)
		return nil, err
	}

	return worktables, nil
}

func (worktableservice *WorkTableService) Create(dtoworktable *models.DtoWorkTable, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoworktable)
	} else {
		err = worktableservice.DbContext.Insert(dtoworktable)
	}
	if err != nil {
		log.Error("Error during creating work table object in database %v", err)
		return err
	}

	return nil
}

func (worktableservice *WorkTableService) DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+worktableservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = worktableservice.DbContext.Exec("delete from "+worktableservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting work table objects for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
