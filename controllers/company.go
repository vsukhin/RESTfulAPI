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

// options /api/v1.0/organisations/
func GetMetaCompanies(r render.Render, companyrepository services.CompanyRepository, session *models.DtoSession) {
	company, err := companyrepository.GetMeta(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, company)
}

// get /api/v1.0/organisations/
func GetCompanies(w http.ResponseWriter, request *http.Request, r render.Render, companyrepository services.CompanyRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.CompanySearch), nil, request, r, session.Language)
	if err != nil {
		return
	}
	if len(*filters) != 0 {
		var masks []string
		for _, filter := range *filters {
			var exps []string
			for _, field := range filter.Fields {
				exps = append(exps, "`"+field+"` "+filter.Op+" "+filter.Value)
			}
			masks = append(masks, "("+strings.Join(exps, " or ")+")")
		}
		query += " and "
		query += strings.Join(masks, " and ")
	}

	var sorts *[]models.OrderExp
	sorts, err = helpers.GetOrderArray(new(models.CompanySearch), request, r, session.Language)
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

	companies, err := companyrepository.GetByUser(session.UserID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(companies, len(*companies), w, r)
}

// get /api/v1.0/classification/legalformorganisation/
func GetCompanyTypes(w http.ResponseWriter, r render.Render, companytyperepository services.CompanyTypeRepository, session *models.DtoSession) {
	companytypes, err := companytyperepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(companytypes, len(*companytypes), w, r)
}

// get /api/v1.0/classification/organisationClasses/
func GetCompanyClasses(w http.ResponseWriter, request *http.Request, r render.Render,
	companyclassrepository services.CompanyClassRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.CompanyClassSearch), nil, request, r, session.Language)
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
	sorts, err = helpers.GetOrderArray(new(models.CompanyClassSearch), request, r, session.Language)
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

	companyclasses, err := companyclassrepository.GetAll(query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(companyclasses, len(*companyclasses), w, r)
}

// get /api/v1.0/organisations/:orgid/
func GetCompany(r render.Render, params martini.Params, companyrepository services.CompanyRepository,
	companycoderepository services.CompanyCodeRepository, companyaddressrepository services.CompanyAddressRepository,
	companybankrepository services.CompanyBankRepository, companyemployeerepository services.CompanyEmployeeRepository,
	session *models.DtoSession) {
	dtocompany, err := helpers.CheckCompany(r, params, companyrepository, session.Language)
	if err != nil {
		return
	}

	companycodes, err := companycoderepository.GetByCompany(dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	companyaddresses, err := companyaddressrepository.GetByCompany(dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	companybanks, err := companybankrepository.GetByCompany(dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	companystaff, err := companyemployeerepository.GetByCompany(dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiMiddleCompany(dtocompany.Primary, dtocompany.Company_Type_ID,
		dtocompany.FullName_Rus, dtocompany.FullName_Eng, dtocompany.ShortName_Rus, dtocompany.ShortName_Eng, dtocompany.Resident,
		*companycodes, *companyaddresses, *companybanks, *companystaff, dtocompany.VAT, !dtocompany.Active))
}

// post /api/v1.0/organisations/
func CreateCompany(errors binding.Errors, viewcompany models.ViewCompany, r render.Render,
	unitrepository services.UnitRepository, companytyperepository services.CompanyTypeRepository, companyrepository services.CompanyRepository,
	companycoderepository services.CompanyCodeRepository, companyaddressrepository services.CompanyAddressRepository,
	companybankrepository services.CompanyBankRepository, companyemployeerepository services.CompanyEmployeeRepository,
	companyclassrepository services.CompanyClassRepository, addresstyperepository services.AddressTypeRepository,
	session *models.DtoSession) {
	if helpers.CheckValidation(&viewcompany, errors, r, session.Language) != nil {
		return
	}
	unit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	dtocompany := new(models.DtoCompany)
	dtocompany.Unit_ID = unit.ID
	dtocompany.Created = time.Now()

	err = helpers.FillCompany(&viewcompany, dtocompany, r, companytyperepository, companyclassrepository, addresstyperepository, session.Language)
	if err != nil {
		return
	}

	err = companyrepository.Create(dtocompany, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	companycodes, err := companycoderepository.GetByCompany(dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	companyaddresses, err := companyaddressrepository.GetByCompany(dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongCompany(dtocompany.ID, dtocompany.Primary, dtocompany.Company_Type_ID,
		dtocompany.FullName_Rus, dtocompany.FullName_Eng, dtocompany.ShortName_Rus, dtocompany.ShortName_Eng, dtocompany.Resident,
		*companycodes, *companyaddresses, viewcompany.CompanyBanks, viewcompany.CompanyStaff, dtocompany.VAT, !dtocompany.Active))
}

// put /api/v1.0/organisations/:orgid/
func UpdateCompany(errors binding.Errors, viewcompany models.ViewCompany, r render.Render, params martini.Params,
	unitrepository services.UnitRepository, companytyperepository services.CompanyTypeRepository, companyrepository services.CompanyRepository,
	companycoderepository services.CompanyCodeRepository, companyaddressrepository services.CompanyAddressRepository,
	companybankrepository services.CompanyBankRepository, companyemployeerepository services.CompanyEmployeeRepository,
	companyclassrepository services.CompanyClassRepository, addresstyperepository services.AddressTypeRepository,
	session *models.DtoSession) {
	if helpers.CheckValidation(&viewcompany, errors, r, session.Language) != nil {
		return
	}
	dtocompany, err := helpers.CheckCompany(r, params, companyrepository, session.Language)
	if err != nil {
		return
	}

	err = helpers.FillCompany(&viewcompany, dtocompany, r, companytyperepository, companyclassrepository, addresstyperepository, session.Language)
	if err != nil {
		return
	}

	err = companyrepository.Update(dtocompany, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	companycodes, err := companycoderepository.GetByCompany(dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	companyaddresses, err := companyaddressrepository.GetByCompany(dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiMiddleCompany(dtocompany.Primary, dtocompany.Company_Type_ID,
		dtocompany.FullName_Rus, dtocompany.FullName_Eng, dtocompany.ShortName_Rus, dtocompany.ShortName_Eng, dtocompany.Resident,
		*companycodes, *companyaddresses, viewcompany.CompanyBanks, viewcompany.CompanyStaff, dtocompany.VAT, !dtocompany.Active))
}

// delete /api/v1.0/organisations/:orgid/
func DeleteCompany(r render.Render, params martini.Params, companyrepository services.CompanyRepository,
	session *models.DtoSession) {
	dtocompany, err := helpers.CheckCompany(r, params, companyrepository, session.Language)
	if err != nil {
		return
	}
	if !dtocompany.Active {
		log.Error("Company is not active %v", dtocompany.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	err = companyrepository.Deactivate(dtocompany)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
