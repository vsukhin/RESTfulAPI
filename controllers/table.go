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
	"strings"
	"time"
	"types"
)

// options /api/v1.0/tables/
func GetTableTypes(r render.Render, tabletypeservice *services.TableTypeService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	tabletypes, err := tabletypeservice.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, tabletypes)
}

// get /api/v1.0/tables/
func GetUnitTables(request *http.Request, r render.Render, customertableservice *services.CustomerTableService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	var err error
	query := ""

	var filters *[]models.FilterExp
	filters, err = helpers.GetFilterArray(new(models.TableSearch), nil, request, r, session.Language)
	if err != nil {
		return
	}

	if len(*filters) != 0 {
		var masks []string
		for _, filter := range *filters {
			var exps []string
			for _, field := range filter.Fields {
				exps = append(exps, field+" "+filter.Op+" "+filter.Value)
			}
			masks = append(masks, "("+strings.Join(exps, " or ")+")")
		}
		query += " and "
		query += strings.Join(masks, " and ")
	}

	var sorts *[]models.OrderExp
	sorts, err = helpers.GetOrderArray(new(models.TableSearch), request, r, session.Language)
	if err != nil {
		return
	}
	if len(*sorts) != 0 {
		var orders []string
		for _, sort := range *sorts {
			orders = append(orders, " "+sort.Field+" "+sort.Order)
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

	customertables, err := customertableservice.GetByUnit(query, session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, customertables)
}

// get /api/v1.0/tables/:tid/
func GetTable(r render.Render, params martini.Params, customertableservice *services.CustomerTableService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	dtocustomertable, err := helpers.CheckTable(r, params, customertableservice, session.Language)
	if err != nil {
		return
	}

	var customertablemeta *models.ApiLongCustomerTable
	customertablemeta, err = customertableservice.GetEx(dtocustomertable.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, customertablemeta)
}

// post /api/v1.0/tables/
func CreateTable(errors binding.Errors, viewcustomertable models.ViewShortCustomerTable, r render.Render,
	userservice *services.UserService, customertableservice *services.CustomerTableService,
	tabletypeservice *services.TableTypeService, unitservice *services.UnitService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}

	unitid, typeid, err := helpers.CheckCustomerTableParameters(r, viewcustomertable.UnitID, models.TABLE_TYPE_DEFAULT, session.UserID, session.Language,
		userservice, unitservice, tabletypeservice)
	if err != nil {
		return
	}

	dtocustomertable := new(models.DtoCustomerTable)
	dtocustomertable.Name = viewcustomertable.Name
	dtocustomertable.Created = time.Now()
	dtocustomertable.TypeID = typeid
	dtocustomertable.UnitID = unitid
	dtocustomertable.Active = true
	dtocustomertable.Permanent = true

	err = customertableservice.Create(dtocustomertable)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongCustomerTable(dtocustomertable.ID, dtocustomertable.Name, models.TABLE_TYPE_DEFAULT, dtocustomertable.UnitID))
}

// put /api/v1.0/tables/:tid/
func UpdateTable(errors binding.Errors, viewcustomertable models.ViewLongCustomerTable, r render.Render, params martini.Params,
	userservice *services.UserService, customertableservice *services.CustomerTableService, unitservice *services.UnitService,
	tabletypeservice *services.TableTypeService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	dtocustomertable, err := helpers.CheckTable(r, params, customertableservice, session.Language)
	if err != nil {
		return
	}

	unitid, typeid, err := helpers.CheckCustomerTableParameters(r, viewcustomertable.UnitID, viewcustomertable.Type, session.UserID, session.Language,
		userservice, unitservice, tabletypeservice)
	if err != nil {
		return
	}

	if dtocustomertable.TypeID != typeid {
		if viewcustomertable.Type != models.TABLE_TYPE_DEFAULT {
			log.Error("Can change table type to %v", viewcustomertable.Type)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}
	}

	dtocustomertable.Name = viewcustomertable.Name
	dtocustomertable.TypeID = typeid
	dtocustomertable.UnitID = unitid

	err = customertableservice.Update(dtocustomertable)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortCustomerTable(dtocustomertable.Name, viewcustomertable.Type, dtocustomertable.UnitID))

}

// delete /api/v1.0/tables/:tid/
func DeleteTable(r render.Render, params martini.Params, customertableservice *services.CustomerTableService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	dtocustomertable, err := helpers.CheckTable(r, params, customertableservice, session.Language)
	if err != nil {
		return
	}

	err = customertableservice.Deactivate(dtocustomertable)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// options /api/v1.0/tables/:tid/data/
func GetTableMetaData(r render.Render, params martini.Params, customertableservice *services.CustomerTableService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	dtocustomertable, err := helpers.CheckTable(r, params, customertableservice, session.Language)
	if err != nil {
		return
	}

	customertablemeta := new(models.ApiMetaCustomerTable)
	customertablemeta, err = customertableservice.GetMeta(dtocustomertable.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, customertablemeta)
}

// put /api/v1.0/tables/:tid/price/
func UpdatePriceTable(errors binding.Errors, viewpriceproperties models.ViewApiPriceProperties, r render.Render, params martini.Params,
	customertableservice *services.CustomerTableService, pricepropertiesservice *services.PricePropertiesService,
	facilityservice *services.FacilityService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	dtocustomertable, err := helpers.CheckTable(r, params, customertableservice, session.Language)
	if err != nil {
		return
	}

	found := false
	found, err = pricepropertiesservice.Exists(dtocustomertable.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	if found {
		log.Error("Customer table is already linked to price table %v", dtocustomertable.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	var facility *models.DtoFacility
	facility, err = facilityservice.Get(viewpriceproperties.Service_ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	if !facility.Active {
		log.Error("Facility is not active %v", facility.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	if viewpriceproperties.After_ID != 0 {
		_, err = helpers.IsTableAvailable(r, customertableservice, viewpriceproperties.After_ID, session.Language)
		if err != nil {
			return
		}
		var priceproperties *models.DtoPriceProperties
		priceproperties, err = pricepropertiesservice.Get(viewpriceproperties.After_ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}
		if priceproperties.Service_ID != viewpriceproperties.Service_ID {
			log.Error("Service is not the same for after price %v", viewpriceproperties.After_ID)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}
	}

	if !viewpriceproperties.Begin.IsZero() && viewpriceproperties.Begin.Sub(time.Now()) < 0 {
		log.Error("Begin date is in the past %v", viewpriceproperties.Begin)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	if !viewpriceproperties.End.IsZero() && viewpriceproperties.End.Sub(time.Now()) < 0 {
		log.Error("End date is in the past %v", viewpriceproperties.End)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	if !viewpriceproperties.Begin.IsZero() && !viewpriceproperties.End.IsZero() &&
		viewpriceproperties.Begin.Sub(viewpriceproperties.End) > 0 {
		log.Error("Begin date can't be bigger than end date %v", viewpriceproperties.Begin, viewpriceproperties.End)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	dtopriceproperties := models.NewDtoPriceProperties(dtocustomertable.ID, viewpriceproperties.Service_ID,
		viewpriceproperties.After_ID, viewpriceproperties.Begin, viewpriceproperties.End, time.Now())
	err = pricepropertiesservice.Create(dtopriceproperties, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, viewpriceproperties)
}