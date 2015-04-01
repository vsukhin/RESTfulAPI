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
	PARAM_NAME_ROW_ID = "rid"
	PARAM_QUERY_ROW   = "rowId"
)

func CheckRowValidity(tableid int64, rowid int64, r render.Render, tablerowrepository services.TableRowRepository,
	language string) (dtotablerow *models.DtoTableRow, err error) {
	dtotablerow, err = tablerowrepository.Get(rowid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}
	if !dtotablerow.Active {
		log.Error("Table row is not active %v", rowid)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Row is not active")
	}
	if dtotablerow.Customer_Table_ID != tableid {
		log.Error("Row %v doesn't belong table %v", rowid, tableid)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Non matched table and row")
	}

	return dtotablerow, nil
}

func CheckTableRow(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	tablerowrepository services.TableRowRepository, language string) (dtotablerow *models.DtoTableRow, err error) {
	var tableid int64
	var rowid int64

	tableid, err = CheckParameterInt(r, params[PARAM_NAME_TABLE_ID], language)
	if err != nil {
		return nil, err
	}
	rowid, err = CheckParameterInt(r, params[PARAM_NAME_ROW_ID], language)
	if err != nil {
		return nil, err
	}
	_, err = IsTableAvailable(r, customertablerepository, tableid, language)
	if err != nil {
		return nil, err
	}

	return CheckRowValidity(tableid, rowid, r, tablerowrepository, language)
}
