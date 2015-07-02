package administration

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

// options /api/v1.0/administration/units/
func GetUnitMetaData(r render.Render, unitrepository services.UnitRepository, session *models.DtoSession) {
	unitmeta, err := unitrepository.GetMeta()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, unitmeta)
}

// get /api/v1.0/administration/units/
func GetUnits(w http.ResponseWriter, request *http.Request, r render.Render, unitrepository services.UnitRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.UnitSearch), nil, request, r, session.Language)
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
		query += " where "
		query += strings.Join(masks, " and ")
	}

	var sorts *[]models.OrderExp
	sorts, err = helpers.GetOrderArray(new(models.UnitSearch), request, r, session.Language)
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

	units, err := unitrepository.GetAll(query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(units, len(*units), w, r)
}

// post /api/v1.0/administration/units/
func CreateUnit(errors binding.Errors, viewunit models.ViewShortUnit, r render.Render, unitrepository services.UnitRepository,
	session *models.DtoSession) {
	if helpers.CheckValidation(&viewunit, errors, r, session.Language) != nil {
		return
	}

	dtounit := new(models.DtoUnit)
	dtounit.Name = viewunit.Name
	dtounit.Active = true
	dtounit.Subscribed = false
	dtounit.Paid = false
	dtounit.Created = time.Now()

	err := unitrepository.Create(dtounit, nil)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongUnit(dtounit.ID, dtounit.Created, dtounit.Name, !dtounit.Active))
}

// get /api/v1.0/administration/units/:unitId/
func GetUnit(r render.Render, params martini.Params, unitrepository services.UnitRepository,
	session *models.DtoSession) {
	dtounit, err := helpers.CheckUnit(r, params, unitrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongUnit(dtounit.ID, dtounit.Created, dtounit.Name, !dtounit.Active))
}

// put /api/v1.0/administration/units/:unitId/
func UpdateUnit(errors binding.Errors, viewunit models.ViewLongUnit, r render.Render, params martini.Params,
	unitrepository services.UnitRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewunit, errors, r, session.Language) != nil {
		return
	}
	dtounit, err := helpers.CheckUnit(r, params, unitrepository, session.Language)
	if err != nil {
		return
	}

	dtounit.Name = viewunit.Name
	dtounit.Active = !viewunit.Deleted
	err = unitrepository.Update(dtounit)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongUnit(dtounit.ID, dtounit.Created, dtounit.Name, !dtounit.Active))
}

// delete /api/v1.0/administration/units/:unitId/
func DeleteUnit(r render.Render, params martini.Params, unitrepository services.UnitRepository, userrepository services.UserRepository,
	customertablerepository services.CustomerTableRepository, projectrepository services.ProjectRepository, orderrepository services.OrderRepository,
	facilityrepository services.FacilityRepository, companyrepository services.CompanyRepository, smssenderrepository services.SMSSenderRepository,
	invoicerepository services.InvoiceRepository, session *models.DtoSession) {
	dtounit, err := helpers.CheckUnit(r, params, unitrepository, session.Language)
	if err != nil {
		return
	}
	if !dtounit.Active {
		log.Error("Unit is not active %v", dtounit.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	users, tables, projects, orders, facilities, companies, smssenders, invoices, err := helpers.GetUnitDependences(dtounit.ID, r, userrepository,
		customertablerepository, projectrepository, orderrepository, facilityrepository, companyrepository, smssenderrepository, invoicerepository,
		session.Language)
	if err != nil {
		return
	}

	if len(*users) != 0 || len(*tables) != 0 || len(*projects) != 0 || len(*orders) != 0 || len(*facilities) != 0 || len(*companies) != 0 ||
		len(*smssenders) != 0 || len(*invoices) != 0 {
		log.Error("Can't delete unit %v with %v users, %v tables, %v projects, %v orders, %v facilities, %v organisations, %v smsfroms, %v invoices",
			dtounit.ID, len(*users), len(*tables), len(*projects), len(*orders), len(*facilities), len(*companies), len(*smssenders), len(*invoices))
		r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_DELETE_DENIED,
			Message: config.Localization[session.Language].Errors.Api.Data_Delete_Denied})
		return
	}

	err = unitrepository.Deactivate(dtounit)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// options /api/v1.0/administration/units/:unitId/dependences/
func GetUnitDependences(r render.Render, params martini.Params, unitrepository services.UnitRepository, userrepository services.UserRepository,
	customertablerepository services.CustomerTableRepository, projectrepository services.ProjectRepository, orderrepository services.OrderRepository,
	facilityrepository services.FacilityRepository, companyrepository services.CompanyRepository, smssenderrepository services.SMSSenderRepository,
	invoicerepository services.InvoiceRepository, session *models.DtoSession) {
	dtounit, err := helpers.CheckUnit(r, params, unitrepository, session.Language)
	if err != nil {
		return
	}
	users, tables, projects, orders, facilities, companies, smssenders, invoices, err := helpers.GetUnitDependences(dtounit.ID, r, userrepository,
		customertablerepository, projectrepository, orderrepository, facilityrepository, companyrepository, smssenderrepository, invoicerepository,
		session.Language)
	if err != nil {
		return
	}

	unitmeta := new(models.ApiLongMetaUnit)
	unitmeta.NumOfUsers = int64(len(*users))
	unitmeta.NumOfTables = int64(len(*tables))
	unitmeta.NumOfProjects = int64(len(*projects))
	unitmeta.NumOfOrders = int64(len(*orders))
	unitmeta.NumOfFacilities = int64(len(*facilities))
	unitmeta.NumOfCompanies = int64(len(*companies))
	unitmeta.NumOfSMSSenders = int64(len(*smssenders))
	unitmeta.NumOfInvoices = int64(len(*invoices))

	r.JSON(http.StatusOK, unitmeta)
}

// get /api/v1.0/administration/units/:unitId/users/
func GetUnitUsers(w http.ResponseWriter, r render.Render, params martini.Params, unitrepository services.UnitRepository,
	userrepository services.UserRepository, session *models.DtoSession) {
	dtounit, err := helpers.CheckUnit(r, params, unitrepository, session.Language)
	if err != nil {
		return
	}

	users, err := userrepository.GetByUnit(dtounit.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(users, len(*users), w, r)
}

// get /api/v1.0/administration/units/:unitId/tables/
func GetUnitTables(w http.ResponseWriter, r render.Render, params martini.Params, unitrepository services.UnitRepository,
	customertablerepository services.CustomerTableRepository, session *models.DtoSession) {
	dtounit, err := helpers.CheckUnit(r, params, unitrepository, session.Language)
	if err != nil {
		return
	}

	tables, err := customertablerepository.GetByUnit(dtounit.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(tables, len(*tables), w, r)
}

// get /api/v1.0/administration/units/:unitId/projects/
func GetUnitProjects(w http.ResponseWriter, r render.Render, params martini.Params, unitrepository services.UnitRepository,
	projectrepository services.ProjectRepository, session *models.DtoSession) {
	dtounit, err := helpers.CheckUnit(r, params, unitrepository, session.Language)
	if err != nil {
		return
	}

	projects, err := projectrepository.GetByUnit(dtounit.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(projects, len(*projects), w, r)
}

// get /api/v1.0/administration/units/:unitId/orders/
func GetUnitOrders(w http.ResponseWriter, r render.Render, params martini.Params, unitrepository services.UnitRepository,
	orderrepository services.OrderRepository, session *models.DtoSession) {
	dtounit, err := helpers.CheckUnit(r, params, unitrepository, session.Language)
	if err != nil {
		return
	}

	orders, err := orderrepository.GetByUnit(dtounit.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(orders, len(*orders), w, r)
}

// get /api/v1.0/administration/units/:unitId/services/
func GetUnitFacilities(w http.ResponseWriter, r render.Render, params martini.Params, unitrepository services.UnitRepository,
	facilityrepository services.FacilityRepository, session *models.DtoSession) {
	dtounit, err := helpers.CheckUnit(r, params, unitrepository, session.Language)
	if err != nil {
		return
	}

	facilities, err := facilityrepository.GetByUnit(dtounit.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(facilities, len(*facilities), w, r)
}

// get /api/v1.0/administration/units/:unitId/organisations/
func GetUnitCompanies(w http.ResponseWriter, r render.Render, params martini.Params, unitrepository services.UnitRepository,
	companyrepository services.CompanyRepository, session *models.DtoSession) {
	dtounit, err := helpers.CheckUnit(r, params, unitrepository, session.Language)
	if err != nil {
		return
	}

	companies, err := companyrepository.GetByUnit(dtounit.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(companies, len(*companies), w, r)
}

// get /api/v1.0/administration/units/:unitId/smsfroms/
func GetUnitSMSSenders(w http.ResponseWriter, r render.Render, params martini.Params, unitrepository services.UnitRepository,
	smssenderrepository services.SMSSenderRepository, session *models.DtoSession) {
	dtounit, err := helpers.CheckUnit(r, params, unitrepository, session.Language)
	if err != nil {
		return
	}

	smssenders, err := smssenderrepository.GetByUnit(dtounit.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(smssenders, len(*smssenders), w, r)
}

// get /api/v1.0/administration/units/:unitId/invoices/
func GetUnitInvoices(w http.ResponseWriter, request *http.Request, r render.Render, params martini.Params,
	unitrepository services.UnitRepository, invoicerepository services.InvoiceRepository, session *models.DtoSession) {
	dtounit, err := helpers.CheckUnit(r, params, unitrepository, session.Language)
	if err != nil {
		return
	}
	query := ""
	var filters *[]models.FilterExp
	filters, err = helpers.GetFilterArray(new(models.InvoiceSearch), nil, request, r, session.Language)
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
	sorts, err = helpers.GetOrderArray(new(models.InvoiceSearch), request, r, session.Language)
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

	invoices, err := invoicerepository.GetByUnit(dtounit.ID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(invoices, len(*invoices), w, r)
}

// options /api/v1.0/administration/orders/
func GetOrderMetaData(r render.Render, orderrepository services.OrderRepository, session *models.DtoSession) {
	ordermeta, err := orderrepository.GetFullMeta()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, ordermeta)
}

// get /api/v1.0/administration/orders/
func GetOrders(w http.ResponseWriter, request *http.Request, r render.Render, orderrepository services.OrderRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.OrderAdminSearch), nil, request, r, session.Language)
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
		query += " where "
		query += strings.Join(masks, " and ")
	}

	var sorts *[]models.OrderExp
	sorts, err = helpers.GetOrderArray(new(models.OrderAdminSearch), request, r, session.Language)
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

	orders, err := orderrepository.GetAll(query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(orders, len(*orders), w, r)
}

// get /api/v1.0/administration/orders/:oid/
func GetOrder(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	orderstatusrepository services.OrderStatusRepository, session *models.DtoSession) {
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}

	dtoorderstatuses, err := orderstatusrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiFullOrderFromDto(dtoorder, dtoorderstatuses))
}

// put /api/v1.0/administration/orders/:oid/
func UpdateOrder(errors binding.Errors, vieworder models.ViewFullOrder, r render.Render, params martini.Params,
	orderrepository services.OrderRepository, unitrepository services.UnitRepository,
	userrepository services.UserRepository, facilityrepository services.FacilityRepository,
	projectrepository services.ProjectRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&vieworder, errors, r, session.Language) != nil {
		return
	}
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}

	apiorder, err := helpers.UpdateFullOrder(dtoorder, &vieworder, r, params, orderrepository, unitrepository, facilityrepository,
		userrepository, projectrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apiorder)
}

// delete /api/v1.0/administration/orders/:oid/
func DeleteOrder(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	orderstatusrepository services.OrderStatusRepository, session *models.DtoSession) {
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}

	orderstatus := models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_DEL, true, "", time.Now())
	err = orderstatusrepository.Save(orderstatus, nil)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
