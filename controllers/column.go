package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/server/middlewares"
	"application/services"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"time"
	"types"
)

// get /api/v1.0/tables/fieldtypes/
func GetColumnTypes(r render.Render, columntypeservice *services.ColumnTypeService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}

	columntypes, err := columntypeservice.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, columntypes)
}

// post /api/v1.0/tables/:tid/field/
func CreateTableColumn(errors binding.Errors, viewtablecolumn models.ViewApiTableColumn, r render.Render, params martini.Params,
	customertableservice *services.CustomerTableService, columntypeservice *services.ColumnTypeService,
	tablecolumnservice *services.TableColumnService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	tableid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_TABLE_ID], session.Language)
	if err != nil {
		return
	}
	_, err = helpers.IsTableAvailable(r, customertableservice, tableid, session.Language)
	if err != nil {
		return
	}

	var typeid int64 = 0
	if viewtablecolumn.TypeID != 0 {
		if helpers.IsColumnTypeActive(r, columntypeservice, viewtablecolumn.TypeID, session.Language) != nil {
			return
		}
		typeid = viewtablecolumn.TypeID
	}

	dtotablecolumn := new(models.DtoTableColumn)
	dtotablecolumn.Name = viewtablecolumn.Name
	dtotablecolumn.Created = time.Now()
	dtotablecolumn.Position = viewtablecolumn.Position
	dtotablecolumn.Column_Type_ID = typeid
	dtotablecolumn.Customer_Table_ID = tableid
	dtotablecolumn.Prebuilt = false
	dtotablecolumn.Active = true
	dtotablecolumn.Edition = 0

	err = tablecolumnservice.Create(dtotablecolumn)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiTableColumn(dtotablecolumn.ID, dtotablecolumn.Name, dtotablecolumn.Column_Type_ID, dtotablecolumn.Position))
}

// get /api/v1.0/tables/:tid/field/
func GetTableColumns(r render.Render, params martini.Params, tablecolumnservice *services.TableColumnService,
	customertableservice *services.CustomerTableService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	tableid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_TABLE_ID], session.Language)
	if err != nil {
		return
	}
	_, err = helpers.IsTableAvailable(r, customertableservice, tableid, session.Language)
	if err != nil {
		return
	}

	var tablecolumns *[]models.ApiTableColumn
	tablecolumns, err = tablecolumnservice.GetByTable(tableid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, tablecolumns)
}

// get /api/v1.0/tables/:tid/field/:cid/
func GetTableColumn(r render.Render, params martini.Params, customertableservice *services.CustomerTableService,
	columntypeservice *services.ColumnTypeService, tablecolumnservice *services.TableColumnService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	dtotablecolumn, err := helpers.CheckTableColumn(r, params, columntypeservice, customertableservice, tablecolumnservice, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewViewApiTableColumn(dtotablecolumn.Name, dtotablecolumn.Column_Type_ID, dtotablecolumn.Position))
}

// put /api/v1.0/tables/:tid/field/:cid/
func UpdateTableColumn(errors binding.Errors, viewtablecolumn models.ViewApiTableColumn, r render.Render, params martini.Params,
	customertableservice *services.CustomerTableService, columntypeservice *services.ColumnTypeService,
	tablecolumnservice *services.TableColumnService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	oldtablecolumn, err := helpers.CheckTableColumn(r, params, columntypeservice, customertableservice, tablecolumnservice, session.Language)
	if err != nil {
		return
	}
	if oldtablecolumn.Prebuilt {
		log.Error("Can't update prebuilt column %v", oldtablecolumn.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	var typeid int64 = 0
	if viewtablecolumn.TypeID != 0 {
		if helpers.IsColumnTypeActive(r, columntypeservice, viewtablecolumn.TypeID, session.Language) != nil {
			return
		}
		typeid = viewtablecolumn.TypeID
	}
	newtablecolumn := new(models.DtoTableColumn)
	*newtablecolumn = *oldtablecolumn

	newtablecolumn.Name = viewtablecolumn.Name
	newtablecolumn.Position = viewtablecolumn.Position
	newtablecolumn.Column_Type_ID = typeid
	newtablecolumn.Created = time.Now()
	newtablecolumn.Edition += 1

	oldtablecolumn.ID = 0
	oldtablecolumn.Active = false
	oldtablecolumn.Original_ID = newtablecolumn.ID

	err = tablecolumnservice.Update(newtablecolumn, oldtablecolumn, false, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewViewApiTableColumn(newtablecolumn.Name, newtablecolumn.Column_Type_ID, newtablecolumn.Position))
}

// delete /api/v1.0/tables/:tid/field/:cid/
func DeleteTableColumn(r render.Render, params martini.Params, customertableservice *services.CustomerTableService,
	columntypeservice *services.ColumnTypeService, tablecolumnservice *services.TableColumnService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	dtotablecolumn, err := helpers.CheckTableColumn(r, params, columntypeservice, customertableservice, tablecolumnservice, session.Language)
	if err != nil {
		return
	}
	if dtotablecolumn.Prebuilt {
		log.Error("Can't delete prebuilt column %v", dtotablecolumn.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	err = tablecolumnservice.Deactivate(dtotablecolumn)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// put /api/v1.0/tables/:tid/sequence/
func UpdateOrderTableColumn(errors binding.Errors, viewordertablecolumns models.ViewApiOrderTableColumns, r render.Render, params martini.Params,
	customertableservice *services.CustomerTableService, columntypeservice *services.ColumnTypeService,
	tablecolumnservice *services.TableColumnService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	tableid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_TABLE_ID], session.Language)
	if err != nil {
		return
	}
	_, err = helpers.IsTableAvailable(r, customertableservice, tableid, session.Language)
	if err != nil {
		return
	}
	if helpers.CheckColumnSet(viewordertablecolumns, tableid, r, tablecolumnservice, session.Language) != nil {
		return
	}

	var dtotablecolumns *[]models.DtoTableColumn
	for _, viewordertablecolumn := range viewordertablecolumns {
		dtotablecolumn, err := tablecolumnservice.Get(viewordertablecolumn.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}
		dtotablecolumn.Position = viewordertablecolumn.Position
		*dtotablecolumns = append(*dtotablecolumns, *dtotablecolumn)
	}

	err = tablecolumnservice.UpdateBriefly(dtotablecolumns, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, viewordertablecolumns)
}
