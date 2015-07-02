package services

import (
	"application/config"
	"application/models"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"path/filepath"
	"strings"
	"time"
)

type CustomerTableRepository interface {
	Copy(srccustomertable *models.DtoCustomerTable, inTrans bool) (destcustomertable *models.DtoCustomerTable, err error)
	ImportDataStructure(dtotablecolumns *[]models.DtoTableColumn, inTrans bool) (err error)
	UpdateImportStructure(customertable *models.DtoCustomerTable, dtotablecolumns *[]models.DtoTableColumn, inTrans bool) (err error)
	ImportData(file *models.DtoFile, dtocustomertable *models.DtoCustomerTable, dataformat models.DataFormat,
		hasheader bool, dtotablecolumns *[]models.DtoTableColumn, version string) (err error)
	ExportData(viewexporttable *models.ViewExportTable, file *models.DtoFile, customertable *models.DtoCustomerTable,
		tablecolumns *[]models.DtoTableColumn, hasheader bool, version string) (err error)
	CheckUserAccess(user_id int64, id int64) (allowed bool, err error)
	Get(id int64) (customertable *models.DtoCustomerTable, err error)
	GetEx(id int64) (customertable *models.ApiLongCustomerTable, err error)
	GetMeta(id int64) (customertable *models.ApiMetaCustomerTable, err error)
	GetByUser(userid int64, filter string, fulllist bool) (customertables *[]models.ApiLongCustomerTable, err error)
	GetByUnit(unitid int64) (customertables *[]models.ApiMiddleCustomerTable, err error)
	GetExpired(timeout time.Duration) (customertables *[]models.DtoCustomerTable, err error)
	Create(customertable *models.DtoCustomerTable) (err error)
	Update(customertable *models.DtoCustomerTable) (err error)
	Deactivate(customertable *models.DtoCustomerTable) (err error)
	Delete(id int64) (err error)
}

type CustomerTableService struct {
	TableColumnRepository TableColumnRepository
	TableRowRepository    TableRowRepository
	*Repository
}

func NewCustomerTableService(repository *Repository) *CustomerTableService {
	repository.DbContext.AddTableWithName(models.DtoCustomerTable{}, repository.Table).SetKeys(true, "id")
	return &CustomerTableService{Repository: repository}
}

func (customertableservice *CustomerTableService) Copy(srccustomertable *models.DtoCustomerTable,
	inTrans bool) (destcustomertable *models.DtoCustomerTable, err error) {
	var trans *gorp.Transaction

	destcustomertable = new(models.DtoCustomerTable)
	*destcustomertable = *srccustomertable
	destcustomertable.ID = 0
	destcustomertable.Created = time.Now()

	if inTrans {
		trans, err = customertableservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating customer table object in database %v", err)
			return nil, err
		}
	}

	if inTrans {
		err = trans.Insert(destcustomertable)
	} else {
		err = customertableservice.DbContext.Insert(destcustomertable)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating customer table object in database %v", err)
		return nil, err
	}

	if inTrans {
		_, err = trans.Exec(
			"insert into table_columns (name, column_type_id, customer_table_id, position, created, prebuilt, active, edition, original_id, fieldnum)"+
				" select name, column_type_id, ?, position, ?, prebuilt, active, 0, 0, fieldnum from table_columns where customer_table_id = ? and active = 1",
			destcustomertable.ID, time.Now(), srccustomertable.ID)
	} else {
		_, err = customertableservice.DbContext.Exec(
			"insert into table_columns (name, column_type_id, customer_table_id, position, created, prebuilt, active, edition, original_id, fieldnum)"+
				" select name, column_type_id, ?, position, ?, prebuilt, active, 0, 0, fieldnum from table_columns where customer_table_id = ? and active = 1",
			destcustomertable.ID, time.Now(), srccustomertable.ID)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating customer table object in database %v", err)
		return nil, err
	}

	query := ""
	for i := 0; i < models.MAX_COLUMN_NUMBER; i++ {
		query += fmt.Sprintf(", field%v, checked%v, valid%v", i+1, i+1, i+1)
	}

	if inTrans {
		_, err = trans.Exec(
			"insert into table_data (customer_table_id, position, created, active, wrong, edition, original_id"+query+
				" select ?, position, ?, active, wrong, 0, 0"+query+" from table_columns where customer_table_id = ? and active = 1",
			destcustomertable.ID, time.Now(), srccustomertable.ID)
	} else {
		_, err = customertableservice.DbContext.Exec(
			"insert into table_data (customer_table_id, position, created, active, wrong, edition, original_id"+query+
				" select ?, position, ?, active, wrong, 0, 0"+query+" from table_columns where customer_table_id = ? and active = 1",
			destcustomertable.ID, time.Now(), srccustomertable.ID)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating customer table object in database %v", err)
		return nil, err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during creating customer table object in database %v", err)
			return nil, err
		}
	}

	return destcustomertable, nil
}

func (customertableservice *CustomerTableService) ImportDataStructure(dtotablecolumns *[]models.DtoTableColumn, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		log.Info("Starting transaction %v", time.Now())
		trans, err = customertableservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during importing customer table object in database %v at %v", err, time.Now())
			return err
		}
	}

	log.Info("Starting inserting table columns %v", time.Now())
	for _, dtotablecolumn := range *dtotablecolumns {
		err = customertableservice.TableColumnRepository.Create(&dtotablecolumn, trans)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
				log.Info("Cancelling inserting table columns %v", time.Now())
			}
			return err
		}
	}
	log.Info("Ending inserting table columns %v", time.Now())

	if inTrans {
		log.Info("Ending transaction %v", time.Now())
		err = trans.Commit()
		if err != nil {
			log.Error("Error during importing customer table object in database %v", err)
			return err
		}
	}

	return nil
}

func (customertableservice *CustomerTableService) UpdateImportStructure(customertable *models.DtoCustomerTable,
	dtotablecolumns *[]models.DtoTableColumn, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = customertableservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating imported customer table object in database %v", err)
			return err
		}
	}

	if inTrans {
		_, err = trans.Update(customertable)
	} else {
		_, err = customertableservice.DbContext.Update(customertable)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating imported customer table object in database %v with value %v", err, customertable.ID)
		return err
	}

	err = customertableservice.TableColumnRepository.UpdateBriefly(dtotablecolumns, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during updating imported customer table object in database %v", err)
			return err
		}
	}

	return nil
}

func (customertableservice *CustomerTableService) ImportData(file *models.DtoFile, dtocustomertable *models.DtoCustomerTable,
	dataformat models.DataFormat, hasheader bool, dtotablecolumns *[]models.DtoTableColumn, version string) (err error) {
	if len(*dtotablecolumns) == 0 {
		log.Error("Can't find any data in customer table object %v for importing with value %v", err, dtocustomertable.ID)
		return errors.New("Empty table")
	}

	query := "load data infile '" + filepath.Join(config.Configuration.FileStorage, file.Path, fmt.Sprintf("%08d", file.ID)+version) +
		"' into table table_data fields terminated by '" +
		string(models.GetDataSeparator(models.DataFormat(dataformat))) +
		"'  escaped by '' lines terminated by '\n'"
	if hasheader {
		query += " ignore 1 lines"
	}
	query += " ("
	for _, tablecolumn := range *dtotablecolumns {
		query += fmt.Sprintf(" field%v", tablecolumn.FieldNum) + ","
	}
	query += " position, wrong) set customer_table_id = " + fmt.Sprintf("%v", dtocustomertable.ID)

	log.Info("Starting inserting table rows %v", time.Now())
	_, err = customertableservice.DbContext.Exec(query)
	if err != nil {
		log.Error("Error during importing customer table object in database %v with value %v", err, dtocustomertable.ID)
		return err
	}
	log.Info("Ending inserting table rows %v", time.Now())

	return nil
}

func (customertableservice *CustomerTableService) ExportData(viewexporttable *models.ViewExportTable, file *models.DtoFile,
	customertable *models.DtoCustomerTable, tablecolumns *[]models.DtoTableColumn, hasheader bool, version string) (err error) {
	if len(*tablecolumns) == 0 {
		log.Error("Can't find any data in customer table object %v for exporting with value %v", err, customertable.ID)
		return errors.New("Empty table")
	}

	log.Info("Starting downloading table rows %v", time.Now())
	query := ""
	if hasheader {
		query = "select"
		for i, tablecolumn := range *tablecolumns {
			columnname := ""
			if !strings.Contains(tablecolumn.Name, "'") {
				columnname = tablecolumn.Name
			}
			query += " '" + columnname + "'"
			if i != len(*tablecolumns)-1 {
				query += ","
			}
		}
		query += " union all "
	}
	query += "( select"
	for i, tablecolumn := range *tablecolumns {
		query += fmt.Sprintf(" field%v", tablecolumn.FieldNum)
		if i != len(*tablecolumns)-1 {
			query += ","
		}
	}
	query += " from table_data where active = 1 and customer_table_id = " + fmt.Sprintf("%v", customertable.ID)
	if viewexporttable.Type != models.EXPORT_DATA_ALL {
		query += " and ("
		for i, tablecolumn := range *tablecolumns {
			if viewexporttable.Type == models.EXPORT_DATA_VALID {
				query += " " + fmt.Sprintf(" valid%v", tablecolumn.FieldNum) + " = 1"
				if i != len(*tablecolumns)-1 {
					query += " and"
				}
			}
			if viewexporttable.Type == models.EXPORT_DATA_INVALID {
				query += " " + fmt.Sprintf(" valid%v", tablecolumn.FieldNum) + " = 0"
				if i != len(*tablecolumns)-1 {
					query += " or"
				}
			}
		}
		query += ")"
	}
	query += " into outfile '" + filepath.Join(config.Configuration.TempDirectory, fmt.Sprintf("%08d", file.ID)+version) + "' fields terminated by '" +
		string(models.GetDataSeparator(models.DataFormat(viewexporttable.Data_Format_ID))) + "' lines terminated by '\n'"
	query += ")"

	_, err = customertableservice.DbContext.Exec(query)
	if err != nil {
		log.Error("Error during exporting customer table object in database %v with value %v", err, customertable.ID)
		return err
	}
	log.Info("Ending downloading table rows %v", time.Now())

	return nil
}

func (customertableservice *CustomerTableService) ClearExpiredTables() {
	for {
		tables, err := customertableservice.GetExpired(config.Configuration.TableTimeout)
		if err == nil {
			for _, table := range *tables {
				err = customertableservice.Deactivate(&table)
			}
		}
		time.Sleep(time.Minute)
	}
}

func (customertableservice *CustomerTableService) CheckUserAccess(user_id int64, id int64) (allowed bool, err error) {
	count, err := customertableservice.DbContext.SelectInt("select count(*) from "+customertableservice.Table+
		" where id = ? and unit_id = (select unit_id from users where id = ?)", id, user_id)
	if err != nil {
		log.Error("Error during checking customer table object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (customertableservice *CustomerTableService) Get(id int64) (customertable *models.DtoCustomerTable, err error) {
	customertable = new(models.DtoCustomerTable)
	err = customertableservice.DbContext.SelectOne(customertable, "select * from "+customertableservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting customer table object from database %v with value %v", err, id)
		return nil, err
	}

	return customertable, nil
}

func (customertableservice *CustomerTableService) GetEx(id int64) (customertable *models.ApiLongCustomerTable, err error) {
	customertable = new(models.ApiLongCustomerTable)
	err = customertableservice.DbContext.SelectOne(customertable,
		"select c.id, c.name, c.type_id as type, c.unit_id from "+customertableservice.Table+
			" c inner join table_types t on c.type_id = t.id where c.id = ?", id)
	if err != nil {
		log.Error("Error during getting extended customer table object from database %v with value %v", err, id)
		return nil, err
	}

	return customertable, nil
}

func (customertableservice *CustomerTableService) GetMeta(id int64) (customertable *models.ApiMetaCustomerTable, err error) {
	customertable = new(models.ApiMetaCustomerTable)
	customertable.NumOfRows, err = customertableservice.DbContext.SelectInt(
		"select count(*) from table_data where active = 1 and customer_table_id = ?", id)
	if err != nil {
		log.Error("Error during getting meta customer table object from database %v with value %v", err, id)
		return nil, err
	}

	tablecolumns, err := customertableservice.TableColumnRepository.GetByTable(id)
	if err != nil {
		return nil, err
	}
	customertable.NumOfCols = int64(len(*tablecolumns))

	query := ""
	for i, tablecolumn := range *tablecolumns {
		query += fmt.Sprintf(" checked%v", tablecolumn.FieldNum) + " = 0"
		if i != len(*tablecolumns)-1 {
			query += " or "
		}
	}
	if query != "" {
		query = "(" + query + ") and"
	}
	var notchecked int64
	notchecked, err = customertableservice.DbContext.SelectInt(
		"select count(*) from table_data where "+query+" active = 1 and customer_table_id = ?", id)
	if err != nil {
		log.Error("Error during getting meta customer table object from database %v with value %v", err, id)
		return nil, err
	}
	customertable.Checked = notchecked == 0

	query = ""
	for i, tablecolumn := range *tablecolumns {
		query += fmt.Sprintf(" valid%v", tablecolumn.FieldNum)
		if i != len(*tablecolumns)-1 {
			query += " + "
		}
	}
	if query != "" {
		query = " sum(" + query + ")"
	} else {
		query = "0"
	}
	var numofvalid int64
	numofvalid, err = customertableservice.DbContext.SelectInt(
		"select "+query+" from table_data where active = 1 and customer_table_id = ?", id)
	if err != nil {
		log.Error("Error during getting meta customer table object from database %v with value %v", err, id)
		return nil, err
	}
	customertable.QaulityPer = 0
	if customertable.NumOfRows != 0 && customertable.NumOfCols != 0 {
		customertable.QaulityPer = byte(100 * numofvalid / (customertable.NumOfRows * customertable.NumOfCols))
	}

	query = ""
	for i, tablecolumn := range *tablecolumns {
		query += fmt.Sprintf(" valid%v", tablecolumn.FieldNum) + " = 0"
		if i != len(*tablecolumns)-1 {
			query += " or "
		}
	}
	if query != "" {
		query = "(" + query + ") and"
	}
	customertable.NumOfWrongRows, err = customertableservice.DbContext.SelectInt(
		"select count(*) from table_data where "+query+" active = 1 and customer_table_id = ?", id)
	if err != nil {
		log.Error("Error during getting meta customer table object from database %v with value %v", err, id)
		return nil, err
	}

	return customertable, nil
}

func (customertableservice *CustomerTableService) GetByUser(userid int64, filter string,
	fulllist bool) (customertables *[]models.ApiLongCustomerTable, err error) {
	customertables = new([]models.ApiLongCustomerTable)
	if !fulllist {
		filter = " and c.type_id != " + fmt.Sprintf("%v", models.TABLE_TYPE_HIDDEN) +
			" and c.type_id != " + fmt.Sprintf("%v", models.TABLE_TYPE_HIDDEN_READONLY) + filter
	}
	_, err = customertableservice.DbContext.Select(customertables,
		"select c.id, c.name, c.type_id as type, c.unit_id from "+customertableservice.Table+
			" c inner join table_types t on c.type_id = t.id where c.active = 1 and c.permanent = 1 and"+
			" c.unit_id = (select unit_id from users where id = ?)"+filter, userid)
	if err != nil {
		log.Error("Error during getting unit customer table object from database %v with value %v", err, userid)
		return nil, err
	}

	return customertables, nil
}

func (customertableservice *CustomerTableService) GetByUnit(unitid int64) (customertables *[]models.ApiMiddleCustomerTable, err error) {
	customertables = new([]models.ApiMiddleCustomerTable)
	_, err = customertableservice.DbContext.Select(customertables,
		"select c.id, c.name, c.type_id as type from "+customertableservice.Table+
			" c inner join table_types t on c.type_id = t.id where c.active = 1 and c.permanent = 1 and c.unit_id = ?", unitid)
	if err != nil {
		log.Error("Error during getting unit customer table object from database %v with value %v", err, unitid)
		return nil, err
	}

	return customertables, nil
}

func (customertableservice *CustomerTableService) GetExpired(timeout time.Duration) (customertables *[]models.DtoCustomerTable, err error) {
	customertables = new([]models.DtoCustomerTable)
	_, err = customertableservice.DbContext.Select(customertables,
		"select * from "+customertableservice.Table+" where permanent = 0 and created < ?", time.Now().Add(-timeout))
	if err != nil {
		log.Error("Error during getting customer table object from database %v with value %v", err, timeout)
		return nil, err
	}

	return customertables, nil
}

func (customertableservice *CustomerTableService) Create(customertable *models.DtoCustomerTable) (err error) {
	err = customertableservice.DbContext.Insert(customertable)
	if err != nil {
		log.Error("Error during creating customer table object in database %v", err)
		return err
	}

	return nil
}

func (customertableservice *CustomerTableService) Update(customertable *models.DtoCustomerTable) (err error) {
	_, err = customertableservice.DbContext.Update(customertable)
	if err != nil {
		log.Error("Error during updating customer table object in database %v with value %v", err, customertable.ID)
		return err
	}

	return nil
}

func (customertableservice *CustomerTableService) Deactivate(customertable *models.DtoCustomerTable) (err error) {
	_, err = customertableservice.DbContext.Exec("update "+customertableservice.Table+" set active = 0 where id = ?", customertable.ID)
	if err != nil {
		log.Error("Error during deactivating customer table object from database %v with value %v", err, customertable.ID)
		return err
	}

	return nil
}

func (customertableservice *CustomerTableService) Delete(id int64) (err error) {
	_, err = customertableservice.DbContext.Exec("delete from "+customertableservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during deleting customer table object in database %v with value %v", err, id)
		return err
	}

	return nil
}
