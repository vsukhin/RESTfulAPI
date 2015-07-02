package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

const (
	PARAM_NAME_REPORT_ID = "aggregateId"
)

func CheckReport(r render.Render, params martini.Params, reportrepository services.ReportRepository,
	language string) (dtoreport *models.DtoReport, err error) {
	report_id, err := CheckParameterInt(r, params[PARAM_NAME_REPORT_ID], language)
	if err != nil {
		return nil, err
	}

	dtoreport, err = reportrepository.Get(report_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtoreport.Active {
		log.Error("Report is not active %v", dtoreport.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Not active report")
	}

	return dtoreport, nil
}

func FillReport(dtoreport *models.DtoReport, r render.Render, reportperiodrepository services.ReportPeriodRepository,
	reportprojectrepository services.ReportProjectRepository, reportorderrepository services.ReportOrderRepository,
	reportfacilityrepository services.ReportFacilityRepository, reportcomplexstatusrepository services.ReportComplexStatusRepository,
	reportsupplierrepository services.ReportSupplierRepository, reportsettingsrepository services.ReportSettingsRepository,
	language string) (apireport *models.ApiReport, err error) {
	periods, err := reportperiodrepository.GetByReport(dtoreport.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	projects, err := reportprojectrepository.GetByReport(dtoreport.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	orders, err := reportorderrepository.GetByReport(dtoreport.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	facilities, err := reportfacilityrepository.GetByReport(dtoreport.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	complexstatuses, err := reportcomplexstatusrepository.GetByReport(dtoreport.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	suppliers, err := reportsupplierrepository.GetByReport(dtoreport.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	settings, err := reportsettingsrepository.GetByReport(dtoreport.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	budgeted := ""
	switch dtoreport.Budgeted {
	case models.TYPE_BUDGETEDBY_UNKNOWN:
	case models.TYPE_BUDGETEDBY_FACILITY:
		budgeted = models.TYPE_BUDGETEDBY_FACILITY_VALUE
	case models.TYPE_BUDGETEDBY_COMPLEX_STATUS:
		budgeted = models.TYPE_BUDGETEDBY_COMPLEX_STATUS_VALUE
	case models.TYPE_BUDGETEDBY_SUPPLIER:
		budgeted = models.TYPE_BUDGETEDBY_SUPPLIER_VALUE
	default:
		log.Error("Unknown budgeted type %v", dtoreport.Budgeted)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, errors.New("Wrong budgeted type")
	}

	return models.NewApiReport(dtoreport.ID, fmt.Sprintf("%v", dtoreport.Created), dtoreport.Unit_ID, dtoreport.User_ID,
		*periods, *projects, *orders, budgeted, *facilities, *complexstatuses, *suppliers,
		*models.NewViewApiReportSettings(settings.Field, settings.Order, settings.Page, settings.Count)), nil
}
