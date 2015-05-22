package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"sort"
	"types"
)

const (
	PARAM_NAME_COLUMN_ID = "cid"
)

func CheckColumnValidity(tableid int64, columnid int64, r render.Render, columntyperepository services.ColumnTypeRepository,
	tablecolumnrepository services.TableColumnRepository, language string) (dtotablecolumn *models.DtoTableColumn, err error) {
	dtotablecolumn, err = tablecolumnrepository.Get(columnid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtotablecolumn.Active {
		log.Error("Table column is not active %v", columnid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Column is not active")
	}

	err = IsColumnTypeActive(r, columntyperepository, dtotablecolumn.Column_Type_ID, language)
	if err != nil {
		return nil, err
	}

	if dtotablecolumn.Customer_Table_ID != tableid {
		log.Error("Column %v doesn't belong table %v", columnid, tableid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Not matched table and column")
	}

	return dtotablecolumn, nil
}

func CheckTableColumn(r render.Render, params martini.Params, columntyperepository services.ColumnTypeRepository,
	customertablerepository services.CustomerTableRepository, tablecolumnrepository services.TableColumnRepository,
	language string) (dtotablecolumn *models.DtoTableColumn, err error) {
	var tableid int64
	var columnid int64

	tableid, err = CheckParameterInt(r, params[PARAM_NAME_TABLE_ID], language)
	if err != nil {
		return nil, err
	}
	columnid, err = CheckParameterInt(r, params[PARAM_NAME_COLUMN_ID], language)
	if err != nil {
		return nil, err
	}
	_, err = IsTableAvailable(r, customertablerepository, tableid, language)
	if err != nil {
		return nil, err
	}

	return CheckColumnValidity(tableid, columnid, r, columntyperepository, tablecolumnrepository, language)
}

func IsColumnTypeActive(r render.Render, columntyperepository services.ColumnTypeRepository, typeid int, language string) (err error) {
	dtocolumntype, err := columntyperepository.Get(typeid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return errors.New("Column is not found")
	}
	if !dtocolumntype.Active {
		log.Error("Column type is not active %v", typeid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return errors.New("Colimn is not active")
	}

	return nil
}

func CheckColumnSet(ids models.IDs, tableid int64, r render.Render, tablecolumnrepository services.TableColumnRepository,
	language string) (tablecolumns *[]models.DtoTableColumn, err error) {
	tablecolumns, err = tablecolumnrepository.GetByTable(tableid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	for _, tablecolumn := range *tablecolumns {
		found := false
		for _, id := range ids.GetIDs() {
			if tablecolumn.ID == id {
				found = true
				break
			}
		}
		if !found {
			log.Error("Can't found column %v", tablecolumn.ID)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, errors.New("Column not found")
		}
	}

	for _, id := range ids.GetIDs() {
		found := false
		for _, tablecolumn := range *tablecolumns {
			if id == tablecolumn.ID {
				found = true
				break
			}
		}
		if !found {
			log.Error("Can't found column %v for table", id)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, errors.New("Column not found")
		}
	}

	return tablecolumns, nil
}

func FindFreeColumn(tableid int64, r render.Render, tablecolumnrepository services.TableColumnRepository,
	language string) (fieldnum byte, err error) {
	tablecolumns, err := tablecolumnrepository.GetByTable(tableid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return 0, err
	}

	allcolumns := make(map[int]bool)
	for i := 0; i < models.MAX_COLUMN_NUMBER; i++ {
		allcolumns[i] = true
	}
	for _, tablecolumn := range *tablecolumns {
		allcolumns[int(tablecolumn.FieldNum)-1] = false
	}

	var keys []int
	for i, avialable := range allcolumns {
		if avialable {
			keys = append(keys, i)
		}
	}
	sort.Ints(keys)

	found := false
	for _, number := range keys {
		fieldnum = byte(number) + 1
		found = true
		break
	}
	if !found {
		log.Error("Can't find free column for table %v", tableid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return 0, errors.New("No free column")
	}

	return fieldnum, nil
}
