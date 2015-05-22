package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"image/jpeg"
	"net/http"
	"strings"
	"time"
	"types"
)

const (
	CAPTCHA_LENGTH  = 6
	CAPTCHA_WIDTH   = 180
	CAPTCHA_HEIGHT  = 80
	CAPTCHA_QUALITY = 10

	NEWS_NUMBER  = 10
	NEWS_VERSION = "2.0"

	SUBSCRIPTION_METHOD_NAME = "/api/v1.0/subscriptions/news/:email/"
)

// get /api/v1.0/captcha/native/
func GetCaptcha(r render.Render, captcharepository services.CaptchaRepository, sessionrepository services.SessionRepository) {
	token, err := sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	digits := captcha.RandomDigits(CAPTCHA_LENGTH)
	value := ""
	for _, d := range digits {
		value += fmt.Sprintf("%v", d)
	}
	image := captcha.NewImage("", digits, CAPTCHA_WIDTH, CAPTCHA_HEIGHT)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, image, &jpeg.Options{Quality: CAPTCHA_QUALITY})
	if err != nil {
		log.Error("Can't convert image to jpeg format %v", err)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	dtocaptcha := models.NewDtoCaptcha(token, buf.Bytes(), value, time.Now(), false)

	err = captcharepository.Create(dtocaptcha)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}
	apicaptcha := models.NewApiCaptcha(dtocaptcha.Hash, base64.StdEncoding.EncodeToString(dtocaptcha.Image))

	r.JSON(http.StatusOK, apicaptcha)
}

// post /api/v1.0/emails/confirm/
func ConfirmEmail(errors binding.Errors, confirm models.EmailConfirm, request *http.Request, r render.Render,
	emailrepository services.EmailRepository, sessionrepository services.SessionRepository, userrepository services.UserRepository,
	templaterepository services.TemplateRepository) {
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}

	email, err := emailrepository.FindByCode(confirm.ConfirmationToken)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_CONFIRMATION_CODE_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Confirmation_Code_Wrong})
		return
	}

	user, err := userrepository.Get(email.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	if !user.Active {
		log.Error("User is not active %v", user.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_USER_BLOCKED,
			Message: config.Localization[user.Language].Errors.Api.User_Blocked})
		return
	}

	if email.Primary {
		if email.Code == user.Code {
			for index, _ := range *user.Emails {
				if (*user.Emails)[index].Email == email.Email {
					(*user.Emails)[index].Code = ""
					(*user.Emails)[index].Confirmed = true
				}
			}

			token, err := sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
				return
			}
			user.Confirmed = true
			user.Code = token

			err = userrepository.Update(user, false, true)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
				return
			}

			for _, confEmail := range *user.Emails {
				if confEmail.Confirmed {
					if helpers.SendPassword(user.Language, &confEmail, user, request, r, emailrepository, templaterepository) != nil {
						return
					}
				}
			}
		} else {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_CONFIRMATION_CODE_WRONG,
				Message: config.Localization[user.Language].Errors.Api.Confirmation_Code_Wrong})
			return
		}
	} else {
		email.Code = ""
		email.Confirmed = true
		err = emailrepository.Update(email, nil)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
			return
		}
	}

	r.JSON(http.StatusAccepted, types.ResponseOK{Message: config.Localization[user.Language].Messages.OK})
}

// get /api/v1.0/
func HomePageTemplate(w http.ResponseWriter, templaterepository services.TemplateRepository) {
	err := templaterepository.GenerateHTML("homepage", w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// get /api/v1.0/classification/contacts/
func GetAvailableContacts(w http.ResponseWriter, r render.Render, classifierrepository services.ClassifierRepository, session *models.DtoSession) {
	classifiers, err := classifierrepository.GetAllAvailable()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(classifiers, len(*classifiers), w, r)
}

// options /api/v1.0/services/
// options /api/v1.0/suppliers/services/
func GetFacilities(w http.ResponseWriter, r render.Render, facilityrepository services.FacilityRepository, session *models.DtoSession) {
	facilities, err := facilityrepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(facilities, len(*facilities), w, r)
}

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

// get /api/v1.0/services/periods/
func GetPeriods(w http.ResponseWriter, r render.Render, periodrepository services.PeriodRepository, session *models.DtoSession) {
	periods, err := periodrepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(periods, len(*periods), w, r)
}

// get /api/v1.0/services/events/
func GetEvents(w http.ResponseWriter, r render.Render, eventrepository services.EventRepository, session *models.DtoSession) {
	events, err := eventrepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(events, len(*events), w, r)
}

// get /api/v1.0/classification/mobileoperators/
func GetMobileOperators(w http.ResponseWriter, r render.Render, mobileoperatorrepository services.MobileOperatorRepository, session *models.DtoSession) {
	mobileoperators, err := mobileoperatorrepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(mobileoperators, len(*mobileoperators), w, r)
}

// get /api/v1.0/classification/services/
func GetFacilityTypes(w http.ResponseWriter, r render.Render, facilitytyperepository services.FacilityTypeRepository, session *models.DtoSession) {
	facilitytypes, err := facilitytyperepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(facilitytypes, len(*facilitytypes), w, r)
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
func GetCompanyClasses(w http.ResponseWriter, request *http.Request, r render.Render, params martini.Params,
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

// get /api/v1.0/classification/addresses/
func GetAddressTypes(w http.ResponseWriter, r render.Render, addresstyperepository services.AddressTypeRepository, session *models.DtoSession) {
	addresstypes, err := addresstyperepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(addresstypes, len(*addresstypes), w, r)
}

// /api/v1.0/services/suppliers/sms/
func GetSMSSuppliers(w http.ResponseWriter, r render.Render, supplierfacilityrepository services.SupplierFacilityRepository, session *models.DtoSession) {
	smsfacilities, err := supplierfacilityrepository.GetAll(models.SERVICE_TYPE_SMS)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(smsfacilities, len(*smsfacilities), w, r)
}

// /api/v1.0/services/suppliers/hlr/
func GetHLRSuppliers(w http.ResponseWriter, r render.Render, supplierfacilityrepository services.SupplierFacilityRepository, session *models.DtoSession) {
	hlrfacilities, err := supplierfacilityrepository.GetAll(models.SERVICE_TYPE_HLR)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(hlrfacilities, len(*hlrfacilities), w, r)
}

// /api/v1.0/services/suppliers/recognize/
func GetRecognizeSuppliers(w http.ResponseWriter, r render.Render, supplierfacilityrepository services.SupplierFacilityRepository, session *models.DtoSession) {
	recognizefacilities, err := supplierfacilityrepository.GetAll(models.SERVICE_TYPE_RECOGNIZE)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(recognizefacilities, len(*recognizefacilities), w, r)
}

// /api/v1.0/services/suppliers/verification/
func GetVerifySuppliers(w http.ResponseWriter, r render.Render, supplierfacilityrepository services.SupplierFacilityRepository, session *models.DtoSession) {
	verifyfacilities, err := supplierfacilityrepository.GetAll(models.SERVICE_TYPE_VERIFY)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(verifyfacilities, len(*verifyfacilities), w, r)
}

// get /api/v1.0/classification/orderstatuses/
func GetComplexStatuses(w http.ResponseWriter, r render.Render, complexstatusrepository services.ComplexStatusRepository, session *models.DtoSession) {
	complexstatuses, err := complexstatusrepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(complexstatuses, len(*complexstatuses), w, r)
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
		*companycodes, *companyaddresses, *companybanks, *companystaff, !dtocompany.Active))
}

// post /api/v1.0/organisations/
func CreateCompany(errors binding.Errors, viewcompany models.ViewCompany, r render.Render, params martini.Params,
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
		*companycodes, *companyaddresses, viewcompany.CompanyBanks, viewcompany.CompanyStaff, !dtocompany.Active))
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
		*companycodes, *companyaddresses, viewcompany.CompanyBanks, viewcompany.CompanyStaff, !dtocompany.Active))
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
