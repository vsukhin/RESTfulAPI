package services

import (
	"application/config"
	"application/models"
	"github.com/coopernurse/gorp"
	"time"
)

type CustomerTableService struct {
	TableColumnService *TableColumnService
	TableRowService    *TableRowService
	*Repository
}

func NewCustomerTableService(repository *Repository) *CustomerTableService {
	repository.DbContext.AddTableWithName(models.DtoCustomerTable{}, repository.Table).SetKeys(true, "id")
	return &CustomerTableService{Repository: repository}
}

func (customertableservice *CustomerTableService) ImportData(customertable *models.DtoCustomerTable,
	dtotablecolumns *[]models.DtoTableColumn, dtotablerows *[]models.DtoTableRow, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		log.Info("Starting transaction %v", time.Now())
		trans, err = customertableservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during importing customer table object in database %v at %v", err, time.Now())
			return err
		}
	}

	log.Info("Inserting table object %v", time.Now())
	err = customertableservice.DbContext.Insert(customertable)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during importing customer table object in database %v", err)
		return err
	}

	log.Info("Starting inserting table columns %v", time.Now())
	for i, _ := range *dtotablecolumns {
		(*dtotablecolumns)[i].Customer_Table_ID = customertable.ID
		err = customertableservice.TableColumnService.Create(&(*dtotablecolumns)[i])
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
				log.Info("Cancelling table columns insertion %v", time.Now())
			}
			return err
		}
	}
	log.Info("Ending inserting table columns %v", time.Now())

	log.Info("Starting inserting table rows %v", time.Now())
	for i, _ := range *dtotablerows {
		(*dtotablerows)[i].Customer_Table_ID = customertable.ID
		for j, _ := range *(*dtotablerows)[i].Cells {
			(*(*dtotablerows)[i].Cells)[j].Table_Column_ID = (*dtotablecolumns)[j].ID
		}
		err = customertableservice.TableRowService.Create(&(*dtotablerows)[i], false)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
				log.Info("Cancelling table rows insertion %v", time.Now())
			}
			return err
		}
		if i%1000 == 0 {
			log.Info("Continue inserting table rows %v at position %v", time.Now(), i)
		}
	}
	log.Info("Ending inserting table rows %v", time.Now())

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

	_, err = customertableservice.DbContext.Update(customertable)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating imported customer table object in database %v", err)
		return err
	}

	err = customertableservice.TableColumnService.UpdateBriefly(dtotablecolumns, false)
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

func (customertableservice *CustomerTableService) ClearExpiredTables() {
	for {
		tables, err := customertableservice.GetExpired(config.Configuration.Server.TableTimeout)
		if err == nil {
			for _, table := range *tables {
				err = customertableservice.Deactivate(&table)
			}
		}
		time.Sleep(time.Minute)
	}
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
		"select c.id, c.name, t.name as type, c.unit_id from "+customertableservice.Table+" c left join table_types t on c.type_id = t.id where c.id = ?", id)
	if err != nil {
		log.Error("Error during getting extended customer table object from database %v with value %v", err, id)
		return nil, err
	}

	return customertable, nil
}

func (customertableservice *CustomerTableService) GetMeta(id int64) (customertable *models.ApiMetaCustomerTable, err error) {
	customertable = new(models.ApiMetaCustomerTable)
	customertable.NumOfRows, err = customertableservice.DbContext.SelectInt(
		"select count(*) from table_rows where active = 1 and customer_table_id = ?", id)
	if err != nil {
		log.Error("Error during getting meta customer table object from database %v with value %v", err, id)
		return nil, err
	}

	customertable.NumOfCols, err = customertableservice.DbContext.SelectInt(
		"select count(*) from table_columns where active = 1 and customer_table_id = ? "+
			"and (column_type_id = 0 or column_type_id in (select id from column_types where active = 1))", id)
	if err != nil {
		log.Error("Error during getting meta customer table object from database %v with value %v with value %v", err, id)
		return nil, err
	}

	var notchecked int64
	notchecked, err = customertableservice.DbContext.SelectInt(
		"select count(*) from table_cells where active = 1 and checked = 0 and table_column_id in "+
			"(select id from table_columns where customer_table_id = ? and active = 1"+
			" and (column_type_id = 0 or column_type_id in (select id from column_types where active = 1)))"+
			"and table_row_id in (select id from table_rows where active = 1 and customer_table_id = ?)", id, id)
	if err != nil {
		log.Error("Error during getting meta customer table object from database %v with value %v", err, id)
		return nil, err
	}
	customertable.Checked = notchecked == 0

	var numofvalid int64
	numofvalid, err = customertableservice.DbContext.SelectInt(
		"select count(*) from table_cells where active = 1 and valid = 1 and table_column_id in "+
			"(select id from table_columns where customer_table_id = ? and active = 1"+
			" and (column_type_id = 0 or column_type_id in (select id from column_types where active = 1)))"+
			"and table_row_id in (select id from table_rows where active = 1 and customer_table_id = ?)", id, id)
	if err != nil {
		log.Error("Error during getting meta customer table object from database %v with value %v", err, id)
		return nil, err
	}
	customertable.QaulityPer = 100
	if customertable.NumOfRows != 0 && customertable.NumOfCols != 0 {
		customertable.QaulityPer = byte(100 * numofvalid / (customertable.NumOfRows * customertable.NumOfCols))
	}

	customertable.NumOfWrongRows, err = customertableservice.DbContext.SelectInt(
		"select count(*) from table_rows where active = 1 and customer_table_id = ? and id in "+
			"(select table_row_id from table_cells where active = 1 and valid = 0 and table_column_id in "+
			"(select id from table_columns where active = 1 and customer_table_id = ?"+
			" and (column_type_id = 0 or column_type_id in (select id from column_types where active = 1))))", id, id)
	if err != nil {
		log.Error("Error during getting meta customer table object from database %v with value %v", err, id)
		return nil, err
	}

	return customertable, nil
}

func (customertableservice *CustomerTableService) GetByUnit(filter string, userid int64) (customertables *[]models.ApiLongCustomerTable, err error) {
	customertables = new([]models.ApiLongCustomerTable)
	_, err = customertableservice.DbContext.Select(customertables,
		"select c.id, c.name, t.name as type, c.unit_id from "+customertableservice.Table+
			" c left join table_types t on c.type_id = t.id where c.active = 1 and c.permanent = 1 and"+
			" c.unit_id = (select unit_id from users where id = ?)"+filter, userid)
	if err != nil {
		log.Error("Error during getting unit customer table object from database %v with value %v", err, userid)
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
		log.Error("Error during updating customer table object in database %v", err)
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
