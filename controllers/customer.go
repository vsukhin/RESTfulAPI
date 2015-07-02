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
	"strings"
	"time"
	"types"
)

// get /api/v1.0/customers/services/
func GetAvailableFacilities(w http.ResponseWriter, r render.Render, facilityrepository services.FacilityRepository, session *models.DtoSession) {
	facilities, err := facilityrepository.GetAllAvailable()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(facilities, len(*facilities), w, r)
}

// options /api/v1.0/projects/
func GetMetaProjects(r render.Render, projectrepository services.ProjectRepository, session *models.DtoSession) {
	project, err := projectrepository.GetMeta(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, project)
}

// get /api/v1.0/projects/
func GetAllProjects(w http.ResponseWriter, request *http.Request, r render.Render, projectrepository services.ProjectRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.ProjectLongSearch), nil, request, r, session.Language)
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
	sorts, err = helpers.GetOrderArray(new(models.ProjectLongSearch), request, r, session.Language)
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

	projects, err := projectrepository.GetByUser(session.UserID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(projects, len(*projects), w, r)
}

// get /api/v1.0/projects/onthego/
func GetActiveProjects(request *http.Request, w http.ResponseWriter, r render.Render, projectrepository services.ProjectRepository, session *models.DtoSession) {
	helpers.GetProjects(session.UserID, true, request, w, r, projectrepository, session.Language)
}

// get /api/v1.0/projects/wasarchived/
func GetArchiveProjects(request *http.Request, w http.ResponseWriter, r render.Render, projectrepository services.ProjectRepository, session *models.DtoSession) {
	helpers.GetProjects(session.UserID, false, request, w, r, projectrepository, session.Language)
}

// post /api/v1.0/projects/
func CreateProject(errors binding.Errors, viewproject models.ViewProject, r render.Render, projectrepository services.ProjectRepository,
	unitrepository services.UnitRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewproject, errors, r, session.Language) != nil {
		return
	}

	unit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	dtoproject := new(models.DtoProject)
	dtoproject.Unit_ID = unit.ID
	dtoproject.Name = viewproject.Name
	dtoproject.Active = true
	dtoproject.Created = time.Now()

	err = projectrepository.Create(dtoproject)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongProject(dtoproject.ID, dtoproject.Name, !dtoproject.Active, dtoproject.Created))
}

// get /api/v1.0/projects/:prid/
func GetProject(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	session *models.DtoSession) {
	dtoproject, err := helpers.CheckProject(r, params, projectrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongProject(dtoproject.ID, dtoproject.Name, !dtoproject.Active, dtoproject.Created))
}

// put /api/v1.0/projects/:prid/
func UpdateProject(errors binding.Errors, viewproject models.ViewUpdateProject, r render.Render, params martini.Params,
	projectrepository services.ProjectRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewproject, errors, r, session.Language) != nil {
		return
	}
	dtoproject, err := helpers.CheckProject(r, params, projectrepository, session.Language)
	if err != nil {
		return
	}

	dtoproject.Name = viewproject.Name
	dtoproject.Active = !viewproject.Archive
	err = projectrepository.Update(dtoproject)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongProject(dtoproject.ID, dtoproject.Name, !dtoproject.Active, dtoproject.Created))
}

// delete /api/v1.0/projects/:prid/
func DeleteProject(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	session *models.DtoSession) {
	dtoproject, err := helpers.CheckProject(r, params, projectrepository, session.Language)
	if err != nil {
		return
	}
	if !dtoproject.Active {
		log.Error("Project is not active %v", dtoproject.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	has, err := projectrepository.HasNotCompletedOrder(dtoproject.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if has {
		log.Error("Project has not completed order")
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	has, err = projectrepository.HasNotPaidOrder(dtoproject.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if has {
		log.Error("Project has not paid order")
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	err = projectrepository.Deactivate(dtoproject)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// options /api/v1.0/smsfrom/
func GetMetaSMSSenders(r render.Render, smssenderrepository services.SMSSenderRepository, session *models.DtoSession) {
	smssender, err := smssenderrepository.GetMeta(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, smssender)
}

// get /api/v1.0/smsfrom/
func GetSMSSenders(w http.ResponseWriter, request *http.Request, r render.Render, smssenderrepository services.SMSSenderRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.SMSSenderSearch), nil, request, r, session.Language)
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
	sorts, err = helpers.GetOrderArray(new(models.SMSSenderSearch), request, r, session.Language)
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

	smssenders, err := smssenderrepository.GetByUser(session.UserID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(smssenders, len(*smssenders), w, r)
}

// get /api/v1.0/smsfrom/:frmid/
func GetSMSSender(r render.Render, params martini.Params, smssenderrepository services.SMSSenderRepository,
	session *models.DtoSession) {
	dtosmssender, err := helpers.CheckSMSSender(r, params, smssenderrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiMiddleSMSSender(dtosmssender.ID, dtosmssender.Name, dtosmssender.Registered, !dtosmssender.Active))
}

// post /api/v1.0/smsfrom/
func CreateSMSSender(errors binding.Errors, viewsmssender models.ViewSMSSender, r render.Render, params martini.Params,
	smssenderrepository services.SMSSenderRepository, unitrepository services.UnitRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewsmssender, errors, r, session.Language) != nil {
		return
	}

	found, err := smssenderrepository.Exists(viewsmssender.Name)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if found {
		r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_SMSSENDER_INUSE,
			Message: config.Localization[session.Language].Errors.Api.SMSSender_InUse})
		return
	}

	unit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	dtosmssender := new(models.DtoSMSSender)
	dtosmssender.Unit_ID = unit.ID
	dtosmssender.Name = viewsmssender.Name
	dtosmssender.Created = time.Now()
	dtosmssender.Registered = false
	dtosmssender.Withdraw = false
	dtosmssender.Withdrawn = false
	dtosmssender.Active = true

	err = smssenderrepository.Create(dtosmssender)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortSMSSender(dtosmssender.ID, dtosmssender.Name, dtosmssender.Registered))
}

// patch /api/v1.0/smsfrom/:frmid/
func UpdateSMSSender(errors binding.Errors, viewsmssender models.ViewSMSSender, r render.Render, params martini.Params,
	smssenderrepository services.SMSSenderRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewsmssender, errors, r, session.Language) != nil {
		return
	}
	dtosmssender, err := helpers.CheckSMSSender(r, params, smssenderrepository, session.Language)
	if err != nil {
		return
	}

	found, err := smssenderrepository.Exists(viewsmssender.Name)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if found && dtosmssender.Name != viewsmssender.Name {
		r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_SMSSENDER_INUSE,
			Message: config.Localization[session.Language].Errors.Api.SMSSender_InUse})
		return
	}

	dtosmssender.Name = viewsmssender.Name
	err = smssenderrepository.Update(dtosmssender)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortSMSSender(dtosmssender.ID, dtosmssender.Name, dtosmssender.Registered))
}

// delete //api/v1.0/smsfrom/:frmid/
func DeleteSMSSender(r render.Render, params martini.Params, smssenderrepository services.SMSSenderRepository,
	session *models.DtoSession) {
	dtosmssender, err := helpers.CheckSMSSender(r, params, smssenderrepository, session.Language)
	if err != nil {
		return
	}
	_, err = helpers.IsSMSSenderActive(dtosmssender.ID, r, smssenderrepository, session.Language)
	if err != nil {
		return
	}

	dtosmssender.Withdraw = true
	err = smssenderrepository.Update(dtosmssender)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// options /api/v1.0/reports/
func GetMetaReports(r render.Render, reportrepository services.ReportRepository, session *models.DtoSession) {
	report, err := reportrepository.GetMeta(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, report)
}

// post /api/v1.0/aggregates/reports/
func CreateReport(errors binding.Errors, viewreport models.ViewReport, r render.Render,
	unitrepository services.UnitRepository, projectrepository services.ProjectRepository,
	orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	complexstatusrepository services.ComplexStatusRepository, reportrepository services.ReportRepository,
	session *models.DtoSession) {
	if helpers.CheckValidation(&viewreport, errors, r, session.Language) != nil {
		return
	}

	unit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	dtoreport := new(models.DtoReport)
	dtoreport.Unit_ID = unit.ID
	dtoreport.User_ID = session.UserID
	dtoreport.Created = time.Now()
	dtoreport.Active = true

	for _, apiperiod := range viewreport.Periods {
		var begin, end time.Time
		if apiperiod.Begin != "" {
			begin, err = models.ParseDate(apiperiod.Begin)
			if err != nil {
				log.Error("Can't parse begin date %v", apiperiod.Begin)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
		}
		if apiperiod.End != "" {
			end, err = models.ParseDate(apiperiod.End)
			if err != nil {
				log.Error("Can't parse end date %v", apiperiod.End)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
		}
		if !begin.IsZero() && !end.IsZero() && begin.Sub(end) > 0 {
			log.Error("Date begin can't be bigger than date end %v", apiperiod.Begin, apiperiod.End)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}
		dtoreport.Periods = append(dtoreport.Periods, *models.NewDtoReportPeriod(0, dtoreport.ID, begin, end))
	}

	for _, apiproject := range viewreport.Projects {
		_, err = helpers.CheckProjectValidity(apiproject.Project_ID, r, projectrepository, session.Language)
		if err != nil {
			return
		}
		err = helpers.CheckProjectAccess(apiproject.Project_ID, session.UserID, r, projectrepository, session.Language)
		if err != nil {
			return
		}
		dtoreport.Projects = append(dtoreport.Projects, *models.NewDtoReportProject(dtoreport.ID, apiproject.Project_ID))
	}

	for _, apiorder := range viewreport.Orders {
		_, err = helpers.CheckOrderValidity(apiorder.Order_ID, r, orderrepository, session.Language)
		if err != nil {
			return
		}
		err = helpers.CheckOrderAccess(apiorder.Order_ID, session.UserID, r, orderrepository, session.Language)
		if err != nil {
			return
		}
		dtoreport.Orders = append(dtoreport.Orders, *models.NewDtoReportOrder(dtoreport.ID, apiorder.Order_ID))
	}

	var budgeted models.BudgetedBy
	switch strings.ToLower(viewreport.Budgeted) {
	case "":
		budgeted = models.TYPE_BUDGETEDBY_UNKNOWN
	case models.TYPE_BUDGETEDBY_FACILITY_VALUE:
		budgeted = models.TYPE_BUDGETEDBY_FACILITY
	case models.TYPE_BUDGETEDBY_COMPLEX_STATUS_VALUE:
		budgeted = models.TYPE_BUDGETEDBY_COMPLEX_STATUS
	case models.TYPE_BUDGETEDBY_SUPPLIER_VALUE:
		budgeted = models.TYPE_BUDGETEDBY_SUPPLIER
	default:
		log.Error("Unknown budgeted type %v", viewreport.Budgeted)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	dtoreport.Budgeted = budgeted

	for _, apifacility := range viewreport.Facilities {
		err = helpers.CheckFacilityValidity(apifacility.Facility_ID, r, facilityrepository, session.Language)
		if err != nil {
			return
		}
		dtoreport.Facilities = append(dtoreport.Facilities, *models.NewDtoReportFacility(dtoreport.ID, apifacility.Facility_ID))
	}

	for _, apicomplexstatus := range viewreport.ComplexStatuses {
		_, err = helpers.CheckComplexStatus(apicomplexstatus.ComplexStatus_ID, r, complexstatusrepository, session.Language)
		if err != nil {
			return
		}
		dtoreport.ComplexStatuses = append(dtoreport.ComplexStatuses, *models.NewDtoReportComplexStatus(dtoreport.ID, apicomplexstatus.ComplexStatus_ID))
	}

	for _, apisupplier := range viewreport.Suppliers {
		err = helpers.CheckUnitValidity(apisupplier.Supplier_ID, session.Language, r, unitrepository)
		if err != nil {
			return
		}
		dtoreport.Suppliers = append(dtoreport.Suppliers, *models.NewDtoReportSupplier(dtoreport.ID, apisupplier.Supplier_ID))
	}

	if viewreport.Settings.Field != "" || viewreport.Settings.Order != "" {
		var valid bool
		apiorderreport := new(models.ApiOrderReport)
		valid, err = apiorderreport.Check(viewreport.Settings.Field)
		if !valid || err != nil {
			log.Error("Unknown field name %v", viewreport.Settings.Field)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}
		if strings.ToLower(viewreport.Settings.Order) != helpers.PARAM_SORT_ASC &&
			strings.ToLower(viewreport.Settings.Order) != helpers.PARAM_SORT_DESC {
			log.Error("Unknown sort operation %v", viewreport.Settings.Order)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}
	}

	if viewreport.Settings.Page == 0 && viewreport.Settings.Count > 0 {
		log.Error("Page number can't be %v", viewreport.Settings.Page)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	dtoreport.Settings = *models.NewDtoReportSettings(dtoreport.ID, viewreport.Settings.Field, viewreport.Settings.Order,
		viewreport.Settings.Page, viewreport.Settings.Count)

	err = reportrepository.Create(dtoreport, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiReport(dtoreport.ID, fmt.Sprintf("%v", dtoreport.Created), dtoreport.Unit_ID, dtoreport.User_ID,
		viewreport.Periods, viewreport.Projects, viewreport.Orders, viewreport.Budgeted,
		viewreport.Facilities, viewreport.ComplexStatuses, viewreport.Suppliers,
		viewreport.Settings))
}

// options /api/v1.0/reports/aggregates/:aggregateId/
func GetReport(r render.Render, params martini.Params, reportrepository services.ReportRepository,
	reportperiodrepository services.ReportPeriodRepository, reportprojectrepository services.ReportProjectRepository,
	reportorderrepository services.ReportOrderRepository, reportfacilityrepository services.ReportFacilityRepository,
	reportcomplexstatusrepository services.ReportComplexStatusRepository, reportsupplierrepository services.ReportSupplierRepository,
	reportsettingsrepository services.ReportSettingsRepository, session *models.DtoSession) {
	dtoreport, err := helpers.CheckReport(r, params, reportrepository, session.Language)
	if err != nil {
		return
	}

	apireport, err := helpers.FillReport(dtoreport, r, reportperiodrepository, reportprojectrepository, reportorderrepository,
		reportfacilityrepository, reportcomplexstatusrepository, reportsupplierrepository, reportsettingsrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apireport)
}

// get /api/v1.0/reports/aggregates/:aggregateId/
func GetComplexReport(r render.Render, params martini.Params, reportrepository services.ReportRepository,
	complexreportrepository services.ComplexReportRepository, reportperiodrepository services.ReportPeriodRepository,
	reportprojectrepository services.ReportProjectRepository, reportorderrepository services.ReportOrderRepository,
	reportfacilityrepository services.ReportFacilityRepository, reportcomplexstatusrepository services.ReportComplexStatusRepository,
	reportsupplierrepository services.ReportSupplierRepository, reportsettingsrepository services.ReportSettingsRepository, session *models.DtoSession) {
	dtoreport, err := helpers.CheckReport(r, params, reportrepository, session.Language)
	if err != nil {
		return
	}

	apireport, err := helpers.FillReport(dtoreport, r, reportperiodrepository, reportprojectrepository, reportorderrepository,
		reportfacilityrepository, reportcomplexstatusrepository, reportsupplierrepository, reportsettingsrepository, session.Language)
	if err != nil {
		return
	}

	complexreport, err := complexreportrepository.Get(session.UserID, apireport)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, complexreport)
}
