package services

import (
	"application/models"
)

type ColumnTypeRepository interface {
	Get(id int64) (columntype *models.DtoColumnType, err error)
	GetAll() (columntypes *[]models.ApiColumnType, err error)
	GetByTable(tableid int64) (columntypes map[int64]models.DtoColumnType, err error)
}

type ColumnTypeService struct {
	*Repository
}

func NewColumnTypeService(repository *Repository) *ColumnTypeService {
	repository.DbContext.AddTableWithName(models.DtoColumnType{}, repository.Table).SetKeys(true, "id")
	return &ColumnTypeService{
		repository,
	}
}

func (columntypeservice *ColumnTypeService) Get(id int64) (columntype *models.DtoColumnType, err error) {
	columntype = new(models.DtoColumnType)
	err = columntypeservice.DbContext.SelectOne(columntype, "select * from "+columntypeservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting column type object from database %v with value %v", err, id)
		return nil, err
	}

	return columntype, nil
}

func (columntypeservice *ColumnTypeService) GetAll() (columntypes *[]models.ApiColumnType, err error) {
	columntypes = new([]models.ApiColumnType)
	_, err = columntypeservice.DbContext.Select(columntypes,
		"select id, name, description, required, `regexp`, "+
			"case horAlignmentHead when 1 then 'left' when 2 then 'center' when 3 then 'right' end as alignmentHead, "+
			"case horAlignmentBody when 1 then 'left' when 2 then 'center' when 3 then 'right' end as alignmentBody from "+
			columntypeservice.Table+" where active = 1")
	if err != nil {
		log.Error("Error during getting all column type object from database %v", err)
		return nil, err
	}

	return columntypes, nil
}

func (columntypeservice *ColumnTypeService) GetByTable(tableid int64) (columntypes map[int64]models.DtoColumnType, err error) {
	columntypes = make(map[int64]models.DtoColumnType)

	tempcolumntypes := new([]models.DtoColumnType)
	_, err = columntypeservice.DbContext.Select(tempcolumntypes,
		"select c.* from "+columntypeservice.Table+" c inner join table_columns t on c.id = t.column_type_id where c.active = 1"+
			" and t.active = 1 and t.customer_table_id = ? order by t.position asc", tableid)
	if err != nil {
		log.Error("Error during getting all column type object from database %v with value %v", err, tableid)
		return nil, err
	}
	for _, columntype := range *tempcolumntypes {
		columntypes[columntype.ID] = columntype
	}

	return columntypes, nil
}
