package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
	"strings"
	"time"
	"types"
)

// get /api/v1.0/tables/:tid/data/
func GetTableData(request *http.Request, r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	tablerowrepository services.TableRowRepository, tablecolumnrepository services.TableColumnRepository, session *models.DtoSession) {
	tableid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_TABLE_ID], session.Language)
	if err != nil {
		return
	}
	_, err = helpers.IsTableAvailable(r, customertablerepository, tableid, session.Language)
	if err != nil {
		return
	}

	filters, err := helpers.GetFilterArray(tablecolumnrepository, tableid, request, r, session.Language)
	if err != nil {
		return
	}

	sorts, err := helpers.GetOrderArray(tablecolumnrepository, request, r, session.Language)
	if err != nil {
		return
	}

	tablecolumns, err := tablecolumnrepository.GetByTable(tableid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	startquery := ""
	endquery := ""

	for _, filter := range *filters {
		var exps []string
		for _, field := range filter.Fields {
			found := false
			for _, tablecolumn := range *tablecolumns {
				value, err := strconv.ParseInt(field, 0, 64)
				if err != nil {
					log.Error("Can't convert to number %v filter column with value %v", err, field)
					r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
					return
				}
				if tablecolumn.ID == value {
					found = true
					exps = append(exps, fmt.Sprintf(" field%v ", tablecolumn.FieldNum)+filter.Op+" "+filter.Value)
					break
				}
			}
			if !found {
				log.Error("Filter column %v doesn't belong table %v", field, tableid)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
		}
		startquery += " (" + strings.Join(exps, " or ") + ")" + " and"
	}

	endquery += " order by"
	for _, sort := range *sorts {
		found := false
		for _, tablecolumn := range *tablecolumns {
			value, err := strconv.ParseInt(sort.Field, 0, 64)
			if err != nil {
				log.Error("Can't convert to number %v order column with value %v", err, sort.Field)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
			if tablecolumn.ID == value {
				found = true
				endquery += fmt.Sprintf(" field%v ", tablecolumn.FieldNum) + sort.Order + ","
				break
			}
		}
		if !found {
			log.Error("Order column %v doesn't belong table %v", sort.Field, tableid)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}
	}
	endquery += " position asc"

	var limit string
	limit, err = helpers.GetLimitQuery(request, r, session.Language)
	if err != nil {
		return
	}
	endquery += limit

	tablerows, err := tablerowrepository.GetAll(startquery, endquery, tableid, tablecolumns)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, tablerows)
}

// get /api/v1.0/tables/:tid/data/:rowid/
func GetTableRow(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	tablerowrepository services.TableRowRepository, tablecolumnrepository services.TableColumnRepository, session *models.DtoSession) {
	dtotablerow, err := helpers.CheckTableRow(r, params, customertablerepository, tablerowrepository, session.Language)
	if err != nil {
		return
	}

	tablecolumns, err := tablecolumnrepository.GetByTable(dtotablerow.Customer_Table_ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	tablecells, err := dtotablerow.TableRowToDtoTableCells(tablecolumns)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	cells := new(models.ViewApiTableRow)
	for _, cell := range *tablecells {
		*cells = append(*cells, *models.NewViewApiTableRowCell(cell.Table_Column_ID, cell.Value))
	}

	r.JSON(http.StatusOK, cells)
}

// post /api/v1.0/tables/:tid/data/
func CreateTableRow(request *http.Request, errors binding.Errors, r render.Render, viewtablecells models.ViewApiTableRow, params martini.Params,
	customertablerepository services.CustomerTableRepository, tablerowrepository services.TableRowRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	tableid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_TABLE_ID], session.Language)
	if err != nil {
		return
	}
	_, err = helpers.IsTableAvailable(r, customertablerepository, tableid, session.Language)
	if err != nil {
		return
	}
	tablecolumns, err := helpers.CheckColumnSet(viewtablecells, tableid, r, tablecolumnrepository, session.Language)
	if err != nil {
		return
	}
	row := request.URL.Query().Get(helpers.PARAM_QUERY_ROW)
	var tablerow *models.DtoTableRow
	if row != "" {
		var rowid int64
		rowid, err = helpers.CheckParameterInt(r, row, session.Language)
		if err != nil {
			return
		}
		tablerow, err = helpers.CheckRowValidity(tableid, rowid, r, tablerowrepository, session.Language)
		if err != nil {
			return
		}
	} else {
		tablerow = nil
	}

	valid := true
	dtotablerow := new(models.DtoTableRow)
	dtotablerow.Customer_Table_ID = tableid
	dtotablerow.Created = time.Now()
	dtotablerow.Active = true
	dtotablerow.Wrong = false
	dtotablerow.Edition = 0
	if tablerow != nil {
		dtotablerow.Position = tablerow.Position
	} else {
		dtotablerow.Position, err = tablerowrepository.GetDefaultPosition(tableid)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}
	}
	cells := new([]models.DtoTableCell)
	for _, cell := range viewtablecells {
		dtotablecolumn, err := helpers.CheckColumnValidity(dtotablerow.Customer_Table_ID, cell.Table_Column_ID, r, columntyperepository,
			tablecolumnrepository, session.Language)
		if err != nil {
			return
		}
		dtotablecell := new(models.DtoTableCell)
		dtotablecell.Table_Column_ID = cell.Table_Column_ID
		dtotablecell.Value = cell.Value
		dtotablecell.Checked = true
		dtotablecell.Valid, err = helpers.ValidateCell(dtotablecell.Value, dtotablecolumn, r, columntyperepository, session.Language)
		if err != nil {
			return
		}
		valid = valid && dtotablecell.Valid

		*cells = append(*cells, *dtotablecell)
	}
	if dtotablerow.TableCellsToTableRow(cells, tablecolumns) != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	err = tablerowrepository.Create(dtotablerow, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongTableRow(dtotablerow.ID, valid))
}

// put /api/v1.0/tables/:tid/data/:rid/
func UpdateTableRow(errors binding.Errors, r render.Render, viewtablecells models.ViewApiTableRow, params martini.Params,
	customertablerepository services.CustomerTableRepository, tablerowrepository services.TableRowRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	oldtablerow, err := helpers.CheckTableRow(r, params, customertablerepository, tablerowrepository, session.Language)
	if err != nil {
		return
	}
	tablecolumns, err := helpers.CheckColumnSet(viewtablecells, oldtablerow.Customer_Table_ID, r, tablecolumnrepository, session.Language)
	if err != nil {
		return
	}

	newtablerow := new(models.DtoTableRow)
	*newtablerow = *oldtablerow

	newtablerow.Created = time.Now()
	newtablerow.Edition += 1
	newtablerow.Wrong = false

	oldtablerow.ID = 0
	oldtablerow.Active = false
	oldtablerow.Original_ID = newtablerow.ID

	valid := true
	cells := new([]models.DtoTableCell)
	for _, cell := range viewtablecells {
		dtotablecolumn, err := helpers.CheckColumnValidity(newtablerow.Customer_Table_ID, cell.Table_Column_ID, r, columntyperepository,
			tablecolumnrepository, session.Language)
		if err != nil {
			return
		}
		dtotablecell := new(models.DtoTableCell)
		dtotablecell.Table_Column_ID = cell.Table_Column_ID
		dtotablecell.Value = cell.Value
		dtotablecell.Checked = true
		dtotablecell.Valid, err = helpers.ValidateCell(dtotablecell.Value, dtotablecolumn, r, columntyperepository, session.Language)
		if err != nil {
			return
		}
		valid = valid && dtotablecell.Valid

		*cells = append(*cells, *dtotablecell)
	}
	if newtablerow.TableCellsToTableRow(cells, tablecolumns) != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	err = tablerowrepository.Update(newtablerow, oldtablerow, false, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortTableRow(valid))
}

// delete /api/v1.0/tables/:tid/data/:rid/
func DeleteTableRow(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	tablerowrepository services.TableRowRepository, session *models.DtoSession) {

	dtotablerow, err := helpers.CheckTableRow(r, params, customertablerepository, tablerowrepository, session.Language)
	if err != nil {
		return
	}

	err = tablerowrepository.Deactivate(dtotablerow, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
