package services

import (
	"application/models"
	"fmt"
)

type FacilityTableRepository interface {
	GetColumnsByType(user_id int64, column_type_id int) (tablecolumns *[]models.DtoTableColumn, validcolumns map[int64]byte, err error)
	GetColumnsByTypes(user_id int64, column_type_ids []int) (tablecolumns *[]models.DtoTableColumn, err error)
	GetColumnsByCustomerTable(customertable_id int64, column_type_id int) (
		tablecolumns *[]models.DtoTableColumn, validcolumns map[int64]byte, err error)
}

type FacilityTableService struct {
	*Repository
}

func NewFacilityTableService(repository *Repository) *FacilityTableService {
	return &FacilityTableService{Repository: repository}
}

func (facilitytableservice *FacilityTableService) GetColumnsByTypes(user_id int64,
	column_type_ids []int) (tablecolumns *[]models.DtoTableColumn, err error) {
	tablecolumns = new([]models.DtoTableColumn)
	column_types := ""
	for i, column_type_id := range column_type_ids {
		if i != 0 {
			column_types += " or "
		}
		column_types += "column_type_id = " + fmt.Sprintf("%v", column_type_id)
	}
	_, err = facilitytableservice.DbContext.Select(tablecolumns,
		"select * from "+facilitytableservice.Table+" where customer_table_id in (select id from customer_tables where"+
			" active = 1 and permanent = 1 and unit_id = (select unit_id from users where id = ?)) and active = 1 and ("+column_types+
			") and column_type_id in (select id from column_types where active = 1)",
		user_id)
	if err != nil {
		log.Error("Error during getting all facility tables in database %v with value %v, %v", err, user_id, column_type_ids)
		return nil, err
	}

	return tablecolumns, nil
}

func (facilitytableservice *FacilityTableService) GetColumnsByType(user_id int64,
	column_type_id int) (tablecolumns *[]models.DtoTableColumn, validcolumns map[int64]byte, err error) {
	tablecolumns = new([]models.DtoTableColumn)
	_, err = facilitytableservice.DbContext.Select(tablecolumns,
		"select * from "+facilitytableservice.Table+" where customer_table_id in (select id from customer_tables where"+
			" active = 1 and permanent = 1 and unit_id = (select unit_id from users where id = ?)) and active = 1 and column_type_id = ?"+
			" and column_type_id in (select id from column_types where active = 1)", user_id, column_type_id)
	if err != nil {
		log.Error("Error during getting all facility tables in database %v with value %v, %v", err, user_id, column_type_id)
		return nil, nil, err
	}
	validcolumns = make(map[int64]byte)
	for _, tablecolumn := range *tablecolumns {
		count, err := facilitytableservice.DbContext.SelectInt("select count(*) from table_data where customer_table_id = ?"+
			" and active = 1 and (checked"+fmt.Sprintf("%v", tablecolumn.FieldNum)+" = 0 or valid"+
			fmt.Sprintf("%v", tablecolumn.FieldNum)+" = 0)", tablecolumn.Customer_Table_ID)
		if err != nil {
			log.Error("Error during getting all facility tables in database %v with value %v, %v", err, user_id, column_type_id)
			return nil, nil, err
		}
		validcolumns[tablecolumn.ID] = 0
		if count == 0 {
			validcolumns[tablecolumn.ID]++
		}
	}

	return tablecolumns, validcolumns, nil
}

func (facilitytableservice *FacilityTableService) GetColumnsByCustomerTable(customertable_id int64,
	column_type_id int) (tablecolumns *[]models.DtoTableColumn, validcolumns map[int64]byte, err error) {
	tablecolumns = new([]models.DtoTableColumn)
	_, err = facilitytableservice.DbContext.Select(tablecolumns,
		"select * from "+facilitytableservice.Table+" where customer_table_id = ? and active = 1 and column_type_id = ?"+
			" and column_type_id in (select id from column_types where active = 1)", customertable_id, column_type_id)
	if err != nil {
		log.Error("Error during getting facility table in database %v with value %v, %v", err, customertable_id, column_type_id)
		return nil, nil, err
	}
	validcolumns = make(map[int64]byte)
	for _, tablecolumn := range *tablecolumns {
		count, err := facilitytableservice.DbContext.SelectInt("select count(*) from table_data where customer_table_id = ?"+
			" and active = 1 and (checked"+fmt.Sprintf("%v", tablecolumn.FieldNum)+" = 0 or valid"+
			fmt.Sprintf("%v", tablecolumn.FieldNum)+" = 0)", tablecolumn.Customer_Table_ID)
		if err != nil {
			log.Error("Error during getting facility table in database %v with value %v, %v", err, customertable_id, column_type_id)
			return nil, nil, err
		}
		validcolumns[tablecolumn.ID] = 0
		if count == 0 {
			validcolumns[tablecolumn.ID]++
		}
	}

	return tablecolumns, validcolumns, nil
}
