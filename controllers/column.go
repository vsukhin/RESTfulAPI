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
	"time"
	"types"
)

// get /api/v1.0/tables/fieldtypes/
func GetColumnTypes(w http.ResponseWriter, r render.Render, columntyperepository services.ColumnTypeRepository, session *models.DtoSession) {
	columntypes, err := columntyperepository.GetAll()
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
	if helpers.CheckValidation(&viewtablecolumn, errors, r, session.Language) != nil {
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
	dtotablecolumn.Position = viewtablecolumn.Position
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
	tablecolumnrepository services.TableColumnRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewtablecolumn, errors, r, session.Language) != nil {
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
	newtablecolumn.Position = viewtablecolumn.Position
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
	if helpers.CheckValidation(&viewordertablecolumns, errors, r, session.Language) != nil {
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
				(*dtotablecolumns)[i].Position = viewordertablecolumn.Position
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
