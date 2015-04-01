package services

import (
	"application/models"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"strings"
)

const (
	MAX_COLUMNS_PER_SELECT = 77
)

type TableRowRepository interface {
	Get(id int64) (tablerow *models.DtoTableRow, err error)
	GetAll(startquery string, endquery string, tableid int64, tablecolumns *[]models.DtoTableColumn) (apitablerows *[]models.ApiInfoTableRow, err error)
	GetDefaultPosition(tableid int64) (value int64, err error)
	Create(tablerow *models.DtoTableRow, inTrans bool) (err error)
	Update(newtablerow *models.DtoTableRow, oldtablerow *models.DtoTableRow, briefly bool, inTrans bool) (err error)
	Deactivate(tablerow *models.DtoTableRow, inTrans bool) (err error)
	Delete(tablerow *models.DtoTableRow, inTrans bool) (err error)
	GetValidation(offset int64, count int64, tableid int64, tablecolumns *[]models.DtoTableColumn) (dtotablerows *[]models.DtoTableRow, err error)
	SaveValidation(tablerows *[]models.DtoTableRow, tablecolumns *[]models.DtoTableColumn) (err error)
}

type TableRowService struct {
	*Repository
}

func NewTableRowService(repository *Repository) *TableRowService {
	repository.DbContext.AddTableWithName(models.DtoTableRow{}, repository.Table).SetKeys(true, "id")
	return &TableRowService{Repository: repository}
}

func (tablerowservice *TableRowService) Get(id int64) (tablerow *models.DtoTableRow, err error) {
	tablerow = new(models.DtoTableRow)
	query := "select id, customer_table_id, created, active, wrong, position, edition, original_id,"
	for i := 0; i < MAX_COLUMNS_PER_SELECT; i++ {
		query += fmt.Sprintf(" field%v, valid%v, checked%v", i+1, i+1, i+1)
		if i != MAX_COLUMNS_PER_SELECT-1 {
			query += ","
		}
	}
	err = tablerowservice.DbContext.SelectOne(tablerow, query+" from "+tablerowservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting table row object from database %v with value %v", err, id)
		return nil, err
	}

	query = "select"
	for i := MAX_COLUMNS_PER_SELECT; i < models.MAX_COLUMN_NUMBER; i++ {
		query += fmt.Sprintf(" field%v, valid%v, checked%v", i+1, i+1, i+1)
		if i != models.MAX_COLUMN_NUMBER-1 {
			query += ","
		}
	}
	temptablerow := new(models.DtoTableRow)
	err = tablerowservice.DbContext.SelectOne(temptablerow, query+" from "+tablerowservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting table row object from database %v with value %v", err, id)
		return nil, err
	}

	for i := MAX_COLUMNS_PER_SELECT; i < models.MAX_COLUMN_NUMBER; i++ {
		for _, column := range []string{"field", "valid", "checked"} {
			data, found := models.GetDbTagValue(fmt.Sprintf(column+"%v", i+1), temptablerow)
			var stringvalue string
			var boolvalue bool
			var ok bool
			if column == "field" {
				stringvalue, ok = data.(string)
			} else {
				boolvalue, ok = data.(bool)
			}
			if !found || !ok {
				log.Error("Can't find "+column+" for table column %v", i+1)
				return nil, errors.New("Error " + column)
			}
			if column == "field" {
				found = models.SetDbTagValue(fmt.Sprintf(column+"%v", i+1), tablerow, stringvalue)
			} else {
				found = models.SetDbTagValue(fmt.Sprintf(column+"%v", i+1), tablerow, boolvalue)
			}
			if !found {
				log.Error("Can't find "+column+" for table column %v", i+1)
				return nil, errors.New("Error " + column)
			}
		}
	}

	return tablerow, nil
}

func (tablerowservice *TableRowService) GetAll(startquery string, endquery string, tableid int64,
	tablecolumns *[]models.DtoTableColumn) (apitablerows *[]models.ApiInfoTableRow, err error) {
	dtotablerows := new([]models.DtoTableRow)
	query := "id"
	for _, tablecolumn := range *tablecolumns {
		query += fmt.Sprintf(", field%v, valid%v", tablecolumn.FieldNum, tablecolumn.FieldNum)
	}
	_, err = tablerowservice.DbContext.Select(dtotablerows,
		"select "+query+" from "+tablerowservice.Table+" where"+startquery+" customer_table_id = ? and active = 1"+endquery, tableid)
	if err != nil {
		log.Error("Error during getting all table row object from database %v with value %v", err, tableid)
		return nil, err
	}
	apitablerows = new([]models.ApiInfoTableRow)
	for _, dtotablerow := range *dtotablerows {
		apitablerow := new(models.ApiInfoTableRow)
		apitablerow.ID = dtotablerow.ID
		cells, err := dtotablerow.TableRowToApiTableCells(tablecolumns)
		if err != nil {
			log.Error("Error during getting all table row object from database %v with value %v", err, tableid)
			return nil, err
		}
		valid := true
		for _, cell := range *cells {
			valid = valid && cell.Valid
		}
		apitablerow.Valid = valid
		apitablerow.Cells = *cells

		*apitablerows = append(*apitablerows, *apitablerow)
	}

	return apitablerows, nil
}

func (tablerowservice *TableRowService) GetValidation(offset int64, count int64, tableid int64,
	tablecolumns *[]models.DtoTableColumn) (dtotablerows *[]models.DtoTableRow, err error) {
	dtotablerows = new([]models.DtoTableRow)
	query := "id"
	for _, tablecolumn := range *tablecolumns {
		query += fmt.Sprintf(", field%v", tablecolumn.FieldNum)
	}
	_, err = tablerowservice.DbContext.Select(dtotablerows,
		"select "+query+" from "+tablerowservice.Table+" where customer_table_id = ? and active = 1 and position >= ? order by position asc limit "+
			fmt.Sprintf("%v", count), tableid, offset)
	if err != nil {
		log.Error("Error during getting validation table row object from database %v with value %v", err, tableid)
		return nil, err
	}

	return dtotablerows, nil
}

func (tablerowservice *TableRowService) SaveValidation(tablerows *[]models.DtoTableRow, tablecolumns *[]models.DtoTableColumn) (err error) {
	if len(*tablerows) == 0 || len(*tablecolumns) == 0 {
		return nil
	}
	elements := new([]string)
	for i, tablerow := range *tablerows {
		query := ""
		if i == 0 {
			query += fmt.Sprintf("select %v as id", tablerow.ID)
		} else {
			query += fmt.Sprintf("select %v", tablerow.ID)
		}
		for _, tablecolumn := range *tablecolumns {
			tablecell, err := tablerow.TableRowToDtoTableCell(&tablecolumn)
			if err != nil {
				log.Error("Error during setting validation table row object in database %v", err)
				return err
			}
			valid := ""
			if tablecell.Valid {
				valid = "1"
			} else {
				valid = "0"
			}
			if i == 0 {
				query += ", " + valid + fmt.Sprintf(" as valid%v", tablecolumn.FieldNum)
			} else {
				query += ", " + valid
			}
		}
		*elements = append(*elements, query)
	}
	components := new([]string)
	for _, tablecolumn := range *tablecolumns {
		query := fmt.Sprintf("t.checked%v = 1, t.valid%v = d.valid%v", tablecolumn.FieldNum, tablecolumn.FieldNum, tablecolumn.FieldNum)
		*components = append(*components, query)
	}
	_, err = tablerowservice.DbContext.Exec("update " + tablerowservice.Table + " t, (" + strings.Join(*elements, " union all ") + ") d set " +
		strings.Join(*components, " ,") + " where t.id = d.id")
	if err != nil {
		log.Error("Error during setting validation table row object in database %v", err)
		return err
	}

	return nil
}

func (tablerowservice *TableRowService) GetDefaultPosition(tableid int64) (value int64, err error) {
	value, err = tablerowservice.DbContext.SelectInt("select count(*) from "+tablerowservice.Table+" where customer_table_id = ? and active = 1", tableid)
	if err != nil {
		log.Error("Error during getting default row position from database %v with value %v", err, tableid)
		return 0, err
	}

	return value, nil
}

func (tablerowservice *TableRowService) Create(tablerow *models.DtoTableRow, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = tablerowservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating row object in database %v", err)
			return err
		}
	}

	_, err = tablerowservice.DbContext.Exec("update "+tablerowservice.Table+" set position = position + 1"+
		" where customer_table_id = ? and position >= ? and active = 1", tablerow.Customer_Table_ID, tablerow.Position)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating row object in database %v", err)
		return err
	}

	err = tablerowservice.DbContext.Insert(tablerow)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating row object in database %v", err)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during creating row object in database %v", err)
			return err
		}
	}

	return nil
}

func (tablerowservice *TableRowService) Update(newtablerow *models.DtoTableRow, oldtablerow *models.DtoTableRow, briefly bool, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = tablerowservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating row object in database %v", err)
			return err
		}
	}

	if !briefly {
		err = tablerowservice.DbContext.Insert(oldtablerow)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during updating table row object in database %v", err)
			return err
		}
	}

	_, err = tablerowservice.DbContext.Update(newtablerow)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating row object in database %v with value %v", err, newtablerow.ID)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during updating row object in database %v", err)
			return err
		}
	}

	return nil
}

func (tablerowservice *TableRowService) Deactivate(tablerow *models.DtoTableRow, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = tablerowservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during deactivating row object in database %v", err)
			return err
		}
	}

	_, err = tablerowservice.DbContext.Exec("update "+tablerowservice.Table+" set position = position - 1"+
		" where customer_table_id = ? and position > ? and active = 1", tablerow.Customer_Table_ID, tablerow.Position)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deactivating row object in database %v with value %v", err, tablerow.ID)
		return err
	}

	_, err = tablerowservice.DbContext.Exec("update "+tablerowservice.Table+" set active = 0 where id = ?", tablerow.ID)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deactivating row object from database %v with value %v", err, tablerow.ID)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during deactivating row object in database %v", err)
			return err
		}
	}

	return nil
}

func (tablerowservice *TableRowService) Delete(tablerow *models.DtoTableRow, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = tablerowservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during deleting row object in database %v", err)
			return err
		}
	}

	_, err = tablerowservice.DbContext.Exec("update "+tablerowservice.Table+" set position = position - 1"+
		" where customer_table_id = ? and position > ? and active = 1", tablerow.Customer_Table_ID, tablerow.Position)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting row object in database %v with value %v", err, tablerow.ID)
		return err
	}

	_, err = tablerowservice.DbContext.Exec("delete from "+tablerowservice.Table+" where id = ?", tablerow.ID)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting row object from database %v with value %v", err, tablerow.ID)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during deleting row object in database %v", err)
			return err
		}
	}

	return nil
}
