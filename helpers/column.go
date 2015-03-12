package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

const (
	PARAM_NAME_COLUMN_ID = "cid"
)

func CheckColumnValidity(tableid int64, columnid int64, r render.Render, columntypeservice *services.ColumnTypeService,
	tablecolumnservice *services.TableColumnService, language string) (dtotablecolumn *models.DtoTableColumn, err error) {
	dtotablecolumn, err = tablecolumnservice.Get(columnid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}
	if !dtotablecolumn.Active {
		log.Error("Table column is not active %v", columnid)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Column is not active")
	}
	if dtotablecolumn.Column_Type_ID != 0 {
		err = IsColumnTypeActive(r, columntypeservice, dtotablecolumn.Column_Type_ID, language)
		if err != nil {
			return nil, err
		}
	}
	if dtotablecolumn.Customer_Table_ID != tableid {
		log.Error("Column %v doesn't belong table %v", columnid, tableid)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Non matched table and column")
	}

	return dtotablecolumn, nil
}

func CheckTableColumn(r render.Render, params martini.Params, columntypeservice *services.ColumnTypeService,
	customertableservice *services.CustomerTableService, tablecolumnservice *services.TableColumnService,
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
	_, err = IsTableAvailable(r, customertableservice, tableid, language)
	if err != nil {
		return nil, err
	}

	return CheckColumnValidity(tableid, columnid, r, columntypeservice, tablecolumnservice, language)
}

func IsColumnTypeActive(r render.Render, columntypeservice *services.ColumnTypeService, typeid int64, language string) (err error) {
	var dtocolumntype *models.DtoColumnType
	dtocolumntype, err = columntypeservice.Get(typeid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return errors.New("Column is not found")
	}
	if !dtocolumntype.Active {
		log.Error("Column type is not active %v", typeid)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return errors.New("Colimn is not active")
	}

	return nil
}

func CheckColumnSet(ids models.IDs, tableid int64, r render.Render, tablecolumnservice *services.TableColumnService,
	language string) (err error) {

	var tablecolumns *[]models.ApiTableColumn
	tablecolumns, err = tablecolumnservice.GetByTable(tableid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return err
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
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return errors.New("Column not found")
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
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return errors.New("Column not found")
		}
	}

	return nil
}
