package services

import (
	"application/models"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"strconv"
	"strings"
)

type TableColumnRepository interface {
	Get(id int64) (tablecolumn *models.DtoTableColumn, err error)
	Check(field string) (valid bool, err error)
	Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error)
	GetAllFields(parameter interface{}) (fields *[]string)
	GetByTable(tableid int64) (tablecolumns *[]models.DtoTableColumn, err error)
	GetDefaultPosition(tableid int64) (position int64, err error)
	Create(tablecolumn *models.DtoTableColumn, trans *gorp.Transaction) (err error)
	CreateAll(tablecolumns *[]models.DtoTableColumn) (err error)
	Update(newtablecolumn *models.DtoTableColumn, oldtablecolumn *models.DtoTableColumn, briefly bool, inTrans bool) (err error)
	UpdateAll(tablecolumns *[]models.DtoTableColumn) (err error)
	UpdateBriefly(tablecolumns *[]models.DtoTableColumn, trans *gorp.Transaction) (err error)
	Deactivate(tablecolumn *models.DtoTableColumn) (err error)
	Delete(id int64) (err error)
}

type TableColumnService struct {
	*Repository
}

func NewTableColumnService(repository *Repository) *TableColumnService {
	repository.DbContext.AddTableWithName(models.DtoTableColumn{}, repository.Table).SetKeys(true, "id")
	return &TableColumnService{Repository: repository}
}

func (tablecolumnservice *TableColumnService) Get(id int64) (tablecolumn *models.DtoTableColumn, err error) {
	tablecolumn = new(models.DtoTableColumn)
	err = tablecolumnservice.DbContext.SelectOne(tablecolumn, "select * from "+tablecolumnservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting table column object from database %v with value %v", err, id)
		return nil, err
	}

	return tablecolumn, nil
}

func (tablecolumnservice *TableColumnService) Check(field string) (valid bool, err error) {
	var count int64
	var value int64

	value, err = strconv.ParseInt(field, 0, 64)
	if err != nil {
		log.Error("Can't convert to number %v with value %v", err, field)
		return false, err
	}

	count, err = tablecolumnservice.DbContext.SelectInt("select count(*) from "+tablecolumnservice.Table+
		" t left join column_types c on t.column_type_id = c.id where t.active = 1 and t.id = ?"+
		" and (c.active = 1)", value)
	if err != nil {
		log.Error("Error during checking table column object from database %v with value %v", err, value)
		return false, err
	}

	return count != 0, nil
}

func (tablecolumnservice *TableColumnService) Extract(infield string, invalue string) (outfield string, outvalue string,
	errField error, errValue error) {

	valid, err := tablecolumnservice.Check(infield)
	if !valid || err != nil {
		errField = errors.New("Unknown field")
		return "", "", errField, nil
	}
	outfield = infield

	if strings.Contains(invalue, "'") {
		errValue = errors.New("Wrong field value")
		return "", "", nil, errValue
	}
	outvalue = "'" + invalue + "'"

	return outfield, outvalue, nil, nil
}

func (tablecolumnservice *TableColumnService) GetAllFields(parameter interface{}) (fields *[]string) {
	fields = new([]string)
	tableid, ok := parameter.(int64)
	if ok {
		tablecolumns, err := tablecolumnservice.GetByTable(tableid)
		if err == nil {
			for _, column := range *tablecolumns {
				*fields = append(*fields, fmt.Sprintf("%v", column.ID))
			}
		}
	}

	return fields
}

func (tablecolumnservice *TableColumnService) GetByTable(tableid int64) (tablecolumns *[]models.DtoTableColumn, err error) {
	tablecolumns = new([]models.DtoTableColumn)
	_, err = tablecolumnservice.DbContext.Select(tablecolumns,
		"select t.* from "+tablecolumnservice.Table+
			" t left join column_types c on t.column_type_id = c.id where t.active = 1 "+
			"and (c.active = 1) and t.customer_table_id = ? order by position asc",
		tableid)
	if err != nil {
		log.Error("Error during getting all table column object from database %v with value %v", err, tableid)
		return nil, err
	}

	return tablecolumns, nil
}

func (tablecolumnservice *TableColumnService) GetDefaultPosition(tableid int64) (position int64, err error) {
	position, err = tablecolumnservice.DbContext.SelectInt("select max(position) from "+tablecolumnservice.Table+
		" t left join column_types c on t.column_type_id = c.id where t.active = 1 "+
		"and (c.active = 1) and t.customer_table_id = ?",
		tableid)
	if err != nil {
		log.Error("Error during getting all table column object from database %v with value %v", err, tableid)
		return 0, err
	}

	return position, nil
}

func (tablecolumnservice *TableColumnService) Create(tablecolumn *models.DtoTableColumn, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(tablecolumn)
	} else {
		err = tablecolumnservice.DbContext.Insert(tablecolumn)
	}
	if err != nil {
		log.Error("Error during creating table column object in database %v", err)
		return err
	}

	return nil
}

func (tablecolumnservice *TableColumnService) CreateAll(tablecolumns *[]models.DtoTableColumn) (err error) {
	var trans *gorp.Transaction

	trans, err = tablecolumnservice.DbContext.Begin()
	if err != nil {
		log.Error("Error during creating table columnn object in database %v", err)
		return err
	}

	for _, tablecolumn := range *tablecolumns {
		err = trans.Insert(&tablecolumn)
		if err != nil {
			_ = trans.Rollback()
			log.Error("Error during creating table column object in database %v", err)
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		log.Error("Error during creating table column object in database %v", err)
		return err
	}

	return nil
}

func (tablecolumnservice *TableColumnService) Update(newtablecolumn *models.DtoTableColumn,
	oldtablecolumn *models.DtoTableColumn, briefly bool, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = tablecolumnservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating table columnn object in database %v", err)
			return err
		}
	}

	if !briefly {
		if inTrans {
			err = trans.Insert(oldtablecolumn)
		} else {
			err = tablecolumnservice.DbContext.Insert(oldtablecolumn)
		}
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during updating table column object in database %v", err)
			return err
		}

		if newtablecolumn.Column_Type_ID != oldtablecolumn.Column_Type_ID {
			if inTrans {
				_, err = trans.Exec("update table_data set checked"+fmt.Sprintf("%v", newtablecolumn.FieldNum)+
					" = 0, valid"+fmt.Sprintf("%v", newtablecolumn.FieldNum)+" = 0 where customer_table_id = ?",
					newtablecolumn.Customer_Table_ID)
			} else {
				_, err = tablecolumnservice.DbContext.Exec("update table_data set checked"+fmt.Sprintf("%v", newtablecolumn.FieldNum)+
					" = 0, valid"+fmt.Sprintf("%v", newtablecolumn.FieldNum)+" = 0 where customer_table_id = ?",
					newtablecolumn.Customer_Table_ID)
			}
			if err != nil {
				if inTrans {
					_ = trans.Rollback()
				}
				log.Error("Error during updating table column object in database %v with value %v for %v",
					err, newtablecolumn.FieldNum, newtablecolumn.Customer_Table_ID)
				return err
			}
		}
	}

	if inTrans {
		_, err = trans.Update(newtablecolumn)
	} else {
		_, err = tablecolumnservice.DbContext.Update(newtablecolumn)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating table column object in database %v", err)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during updating table column object in database %v", err)
			return err
		}
	}

	return nil
}

func (tablecolumnservice *TableColumnService) UpdateAll(tablecolumns *[]models.DtoTableColumn) (err error) {
	var trans *gorp.Transaction

	trans, err = tablecolumnservice.DbContext.Begin()
	if err != nil {
		log.Error("Error during updating table columnn object in database %v", err)
		return err
	}

	for _, tablecolumn := range *tablecolumns {
		_, err = trans.Update(&tablecolumn)
		if err != nil {
			_ = trans.Rollback()
			log.Error("Error during briefly updating table column object in database %v", err)
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		log.Error("Error during updating table column object in database %v", err)
		return err
	}

	return nil
}

func (tablecolumnservice *TableColumnService) UpdateBriefly(tablecolumns *[]models.DtoTableColumn, trans *gorp.Transaction) (err error) {
	for _, tablecolumn := range *tablecolumns {
		if trans != nil {
			_, err = trans.Update(&tablecolumn)
		} else {
			_, err = tablecolumnservice.DbContext.Update(&tablecolumn)
		}
		if err != nil {
			log.Error("Error during briefly updating table column object in database %v", err)
			return err
		}
	}

	return nil
}

func (tablecolumnservice *TableColumnService) Deactivate(tablecolumn *models.DtoTableColumn) (err error) {
	_, err = tablecolumnservice.DbContext.Exec("update "+tablecolumnservice.Table+" set active = 0 where id = ?", tablecolumn.ID)
	if err != nil {
		log.Error("Error during deactivating table column object from database %v with value %v", err, tablecolumn.ID)
		return err
	}

	return nil
}

func (tablecolumnservice *TableColumnService) Delete(id int64) (err error) {
	_, err = tablecolumnservice.DbContext.Exec("delete from "+tablecolumnservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during deleting table column object in database %v with value %v", err, id)
		return err
	}

	return nil
}
