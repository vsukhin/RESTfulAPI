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

	apicompany, err := helpers.LoadCompany(dtocompany, r, companycoderepository, companyaddressrepository, companybankrepository,
		companyemployeerepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apicompany)
}

// post /api/v1.0/organisations/
func CreateCompany(errors binding.Errors, viewcompany models.ViewCompany, r render.Render,
	unitrepository services.UnitRepository, companytyperepository services.CompanyTypeRepository, companyrepository services.CompanyRepository,
	companycoderepository services.CompanyCodeRepository, companyaddressrepository services.CompanyAddressRepository,
	companybankrepository services.CompanyBankRepository, companyemployeerepository services.CompanyEmployeeRepository,
	companyclassrepository services.CompanyClassRepository, addresstyperepository services.AddressTypeRepository,
	session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
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
	dtocompany.Active = true

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
		*companycodes, *companyaddresses, viewcompany.CompanyBanks, viewcompany.CompanyStaff, dtocompany.VAT, dtocompany.Locked, !dtocompany.Active))
}

// put /api/v1.0/organisations/:orgid/
func UpdateCompany(errors binding.Errors, viewcompany models.ViewCompany, r render.Render, params martini.Params,
	unitrepository services.UnitRepository, companytyperepository services.CompanyTypeRepository, companyrepository services.CompanyRepository,
	companycoderepository services.CompanyCodeRepository, companyaddressrepository services.CompanyAddressRepository,
	companybankrepository services.CompanyBankRepository, companyemployeerepository services.CompanyEmployeeRepository,
	companyclassrepository services.CompanyClassRepository, addresstyperepository services.AddressTypeRepository,
	session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
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

	if dtocompany.Locked {
		log.Error("Company is locked %v", dtocompany.ID)
		r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_CHANGES_DENIED,
			Message: config.Localization[session.Language].Errors.Api.Data_Changes_Denied})
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
		*companycodes, *companyaddresses, viewcompany.CompanyBanks, viewcompany.CompanyStaff, dtocompany.VAT, dtocompany.Locked, !dtocompany.Active))
}

// delete /api/v1.0/organisations/:orgid/
func DeleteCompany(r render.Render, params martini.Params, companyrepository services.CompanyRepository, session *models.DtoSession) {
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

	if dtocompany.Locked {
		log.Error("Company is locked %v", dtocompany.ID)
		r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_CHANGES_DENIED,
			Message: config.Localization[session.Language].Errors.Api.Data_Changes_Denied})
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

// get /api/v1.0/unit/contract/
func GetContracts(w http.ResponseWriter, request *http.Request, r render.Render, contractrepository services.ContractRepository,
	appendixrepository services.AppendixRepository, session *models.DtoSession) {
	dtocontracts, err := contractrepository.GetByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	apicontracts := new([]models.ApiContract)
	for _, dtocontract := range *dtocontracts {
		apiappendices, err := appendixrepository.GetByContract(dtocontract.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}
		*apicontracts = append(*apicontracts, *models.NewApiContract(dtocontract.Company_ID, dtocontract.Confirmed, dtocontract.Signed,
			dtocontract.SignedDate, dtocontract.Name, dtocontract.File_ID, *apiappendices))
	}

	helpers.RenderJSONArray(apicontracts, len(*apicontracts), w, r)
}

// patch /api/v1.0/unit/contract/
func UpdateContracts(errors binding.Errors, changecontract models.ChangeContract, w http.ResponseWriter, request *http.Request, r render.Render,
	contractrepository services.ContractRepository, appendixrepository services.AppendixRepository, unitrepository services.UnitRepository,
	companyrepository services.CompanyRepository, sessionrepository services.SessionRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}

	unit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	companies, err := companyrepository.GetByUnit(unit.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	for _, company := range *companies {
		found := false
		for _, viewcontract := range changecontract {
			if viewcontract.Company_ID == company.ID {
				found = true
				break
			}
		}
		if !found {
			log.Error("Can't find company in contract list %v", company.ID)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}
	}
	for _, viewcontract := range changecontract {
		found := false
		for _, company := range *companies {
			if company.ID == viewcontract.Company_ID {
				found = true
				break
			}
		}
		if !found {
			log.Error("Unknown company in contract list %v", viewcontract.Company_ID)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}
	}

	dtocontracts, err := contractrepository.GetByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	newcontracts := new([]models.DtoContract)
	for i := range changecontract {
		found := false
		for j := range *dtocontracts {
			if (*dtocontracts)[j].Company_ID == changecontract[i].Company_ID {
				(*dtocontracts)[j].Confirmed = changecontract[i].Confirmed
				found = true
				break
			}
		}
		if !found {
			token, err := sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
			*newcontracts = append(*newcontracts, *models.NewDtoContract(0, changecontract[i].Company_ID, token, changecontract[i].Confirmed, false,
				time.Time{}, 0, time.Now(), true, []models.DtoAppendix{}))
		}
	}
	*dtocontracts = append(*dtocontracts, *newcontracts...)

	err = contractrepository.SaveAll(dtocontracts, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	apicontracts := new([]models.ApiContract)
	for _, dtocontract := range *dtocontracts {
		apiappendices, err := appendixrepository.GetByContract(dtocontract.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}
		*apicontracts = append(*apicontracts, *models.NewApiContract(dtocontract.Company_ID, dtocontract.Confirmed, dtocontract.Signed,
			dtocontract.SignedDate, dtocontract.Name, dtocontract.File_ID, *apiappendices))
	}

	helpers.RenderJSONArray(apicontracts, len(*apicontracts), w, r)
}
