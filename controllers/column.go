package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
	"time"
	"types"
)

// get /api/v1.0/tables/fieldtypes/
func GetColumnTypes(w http.ResponseWriter, request *http.Request, r render.Render, columntyperepository services.ColumnTypeRepository,
	session *models.DtoSession) {
	query := ""
	default_filter := false
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.ColumnTypeSearch), nil, request, r, session.Language)
	if err != nil {
		return
	}
	if len(*filters) != 0 {
		var masks []string
		for _, filter := range *filters {
			var exps []string
			for _, field := range filter.Fields {
				if field == "private" {
					default_filter = true
				}
				if field == "private" && ((filter.Op == "=" && filter.Value == "true") || (filter.Op == "!=" && filter.Value == "false") ||
					(filter.Op == "like" && filter.Value == "true")) {
					exps = append(exps, "`private` = false", "`private` = true")
				} else {
					exps = append(exps, "`"+field+"` "+filter.Op+" "+filter.Value)
				}
			}
			masks = append(masks, "("+strings.Join(exps, " or ")+")")
		}
		query += " and "
		query += strings.Join(masks, " and ")
	}
	if !default_filter {
		query += " and (`private` = false)"
	}

	var sorts *[]models.OrderExp
	sorts, err = helpers.GetOrderArray(new(models.ColumnTypeSearch), request, r, session.Language)
	if err != nil {
		return
	}
	if len(*sorts) != 0 {
		var orders []string
		for _, sort := range *sorts {
			orders = append(orders, " `"+sort.Field+"` "+sort.Order)
		}
		query += " order by"
		query += strings.Join(orders, ",")
	}

	var limit string
	limit, err = helpers.GetLimitQuery(request, r, session.Language)
	if err != nil {
		return
	}
	query += limit

	columntypes, err := columntyperepository.GetAll(query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(columntypes, len(*columntypes), w, r)
}

// post /api/v1.0/tables/:tid/field/
func CreateTableColumn(errors binding.Errors, viewtablecolumn models.ViewApiTableColumn, r render.Render, params martini.Params,
	customertablerepository services.CustomerTableRepository, columntyperepository services.ColumnTypeRepository,
	tablecolumnrepository services.TableColumnRepository, session *models.DtoSession) {
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
	if helpers.IsColumnTypeActive(r, columntyperepository, viewtablecolumn.TypeID, session.Language) != nil {
		return
	}

	dtotablecolumn := new(models.DtoTableColumn)
	dtotablecolumn.Name = viewtablecolumn.Name
	dtotablecolumn.Created = time.Now()
	if viewtablecolumn.Position == 0 {
		dtotablecolumn.Position, err = tablecolumnrepository.GetDefaultPosition(tableid)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}
		dtotablecolumn.Position++
	} else {
		dtotablecolumn.Position = viewtablecolumn.Position
	}
	dtotablecolumn.Column_Type_ID = viewtablecolumn.TypeID
	dtotablecolumn.Customer_Table_ID = tableid
	dtotablecolumn.Prebuilt = false
	dtotablecolumn.FieldNum, err = helpers.FindFreeColumn(tableid, 0, r, tablecolumnrepository, session.Language)
	if err != nil {
		return
	}
	dtotablecolumn.Active = true
	dtotablecolumn.Edition = 0

	err = tablecolumnrepository.Create(dtotablecolumn, nil)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiTableColumn(dtotablecolumn.ID, dtotablecolumn.Name, dtotablecolumn.Column_Type_ID, dtotablecolumn.Position))
}

// get /api/v1.0/tables/:tid/field/
func GetTableColumns(w http.ResponseWriter, r render.Render, params martini.Params, tablecolumnrepository services.TableColumnRepository,
	customertablerepository services.CustomerTableRepository, session *models.DtoSession) {
	tableid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_TABLE_ID], session.Language)
	if err != nil {
		return
	}
	_, err = helpers.IsTableAvailable(r, customertablerepository, tableid, session.Language)
	if err != nil {
		return
	}

	tablecolumns, err := tablecolumnrepository.GetByTable(tableid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	apicolumns := new([]models.ApiTableColumn)
	for _, tablecolumn := range *tablecolumns {
		*apicolumns = append(*apicolumns, *models.NewApiTableColumn(tablecolumn.ID, tablecolumn.Name, tablecolumn.Column_Type_ID, tablecolumn.Position))
	}

	helpers.RenderJSONArray(apicolumns, len(*apicolumns), w, r)
}

// get /api/v1.0/tables/:tid/field/:cid/
func GetTableColumn(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository, session *models.DtoSession) {
	dtotablecolumn, err := helpers.CheckTableColumn(r, params, columntyperepository, customertablerepository, tablecolumnrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewViewApiTableColumn(dtotablecolumn.Name, dtotablecolumn.Column_Type_ID, dtotablecolumn.Position))
}

// put /api/v1.0/tables/:tid/field/:cid/
func UpdateTableColumn(errors binding.Errors, viewtablecolumn models.ViewApiTableColumn, r render.Render, params martini.Params,
	customertablerepository services.CustomerTableRepository, columntyperepository services.ColumnTypeRepository,
	tablecolumnrepository services.TableColumnRepository, tablerowrepository services.TableRowRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	oldtablecolumn, err := helpers.CheckTableColumn(r, params, columntyperepository, customertablerepository, tablecolumnrepository, session.Language)
	if err != nil {
		return
	}
	if oldtablecolumn.Prebuilt {
		log.Error("Can't update prebuilt column %v", oldtablecolumn.ID)
		r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_CHANGES_DENIED,
			Message: config.Localization[session.Language].Errors.Api.Data_Changes_Denied})
		return
	}
	if helpers.IsColumnTypeActive(r, columntyperepository, viewtablecolumn.TypeID, session.Language) != nil {
		return
	}

	newtablecolumn := new(models.DtoTableColumn)
	*newtablecolumn = *oldtablecolumn

	newtablecolumn.Name = viewtablecolumn.Name
	if viewtablecolumn.Position == 0 {
		newtablecolumn.Position, err = tablecolumnrepository.GetDefaultPosition(newtablecolumn.Customer_Table_ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}
		newtablecolumn.Position++
	} else {
		newtablecolumn.Position = viewtablecolumn.Position
	}
	newtablecolumn.Column_Type_ID = viewtablecolumn.TypeID
	newtablecolumn.Created = time.Now()
	newtablecolumn.Edition += 1

	oldtablecolumn.ID = 0
	oldtablecolumn.Active = false
	oldtablecolumn.Original_ID = newtablecolumn.ID

	err = tablecolumnrepository.Update(newtablecolumn, oldtablecolumn, false, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	go helpers.CheckTableColumnCells(newtablecolumn, columntyperepository, tablerowrepository)

	r.JSON(http.StatusOK, models.NewViewApiTableColumn(newtablecolumn.Name, newtablecolumn.Column_Type_ID, newtablecolumn.Position))
}

// delete /api/v1.0/tables/:tid/field/:cid/
func DeleteTableColumn(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository, session *models.DtoSession) {
	dtotablecolumn, err := helpers.CheckTableColumn(r, params, columntyperepository, customertablerepository, tablecolumnrepository, session.Language)
	if err != nil {
		return
	}
	if dtotablecolumn.Prebuilt {
		log.Error("Can't delete prebuilt column %v", dtotablecolumn.ID)
		r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_CHANGES_DENIED,
			Message: config.Localization[session.Language].Errors.Api.Data_Changes_Denied})
		return
	}

	err = tablecolumnrepository.Deactivate(dtotablecolumn)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// put /api/v1.0/tables/:tid/sequence/
func UpdateOrderTableColumn(errors binding.Errors, viewordertablecolumns models.ViewApiOrderTableColumns, w http.ResponseWriter, r render.Render,
	params martini.Params, customertablerepository services.CustomerTableRepository, tablecolumnrepository services.TableColumnRepository,
	session *models.DtoSession) {
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
	dtotablecolumns, err := helpers.CheckColumnSet(viewordertablecolumns, tableid, r, tablecolumnrepository, session.Language)
	if err != nil {
		return
	}

	for i, _ := range *dtotablecolumns {
		for _, viewordertablecolumn := range viewordertablecolumns {
			if (*dtotablecolumns)[i].ID == viewordertablecolumn.ID {
				if viewordertablecolumn.Position == 0 {
					(*dtotablecolumns)[i].Position, err = tablecolumnrepository.GetDefaultPosition(tableid)
					if err != nil {
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
							Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
						return
					}
					(*dtotablecolumns)[i].Position++
				} else {
					(*dtotablecolumns)[i].Position = viewordertablecolumn.Position
				}
			}
		}
	}

	err = tablecolumnrepository.UpdateAll(dtotablecolumns)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	helpers.RenderJSONArray(viewordertablecolumns, len(viewordertablecolumns), w, r)
}
