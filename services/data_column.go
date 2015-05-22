package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type DataColumnRepository interface {
	Get(order_id int64, table_column_id int64) (datacolumn *models.DtoDataColumn, err error)
	GetByOrder(order_id int64) (datacolumns *[]models.ApiDataColumn, err error)
	Create(dtodatacolumn *models.DtoDataColumn, trans *gorp.Transaction) (err error)
	DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error)
}

type DataColumnService struct {
	*Repository
}

func NewDataColumnService(repository *Repository) *DataColumnService {
	repository.DbContext.AddTableWithName(models.DtoDataColumn{}, repository.Table).SetKeys(false, "order_id", "table_column_id")
	return &DataColumnService{Repository: repository}
}

func (datacolumnservice *DataColumnService) Get(order_id int64, table_column_id int64) (datacolumn *models.DtoDataColumn, err error) {
	datacolumn = new(models.DtoDataColumn)
	err = datacolumnservice.DbContext.SelectOne(datacolumn, "select * from "+datacolumnservice.Table+
		" where order_id = ? and table_column_id = ?", order_id, table_column_id)
	if err != nil {
		log.Error("Error during getting data column object from database %v with value %v, %v", err, order_id, table_column_id)
		return nil, err
	}

	return datacolumn, nil
}

func (datacolumnservice *DataColumnService) GetByOrder(order_id int64) (datacolumns *[]models.ApiDataColumn, err error) {
	datacolumns = new([]models.ApiDataColumn)
	_, err = datacolumnservice.DbContext.Select(datacolumns,
		"select d.table_column_id, c.name, c.column_type_id, c.position from "+datacolumnservice.Table+
			" d inner join table_columns c on d.table_column_id = c.id where d.order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all data column object from database %v with value %v", err, order_id)
		return nil, err
	}

	return datacolumns, nil
}

func (datacolumnservice *DataColumnService) Create(dtodatacolumn *models.DtoDataColumn, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtodatacolumn)
	} else {
		err = datacolumnservice.DbContext.Insert(dtodatacolumn)
	}
	if err != nil {
		log.Error("Error during creating data column object in database %v", err)
		return err
	}

	return nil
}

func (datacolumnservice *DataColumnService) DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+datacolumnservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = datacolumnservice.DbContext.Exec("delete from "+datacolumnservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting data column objects for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
