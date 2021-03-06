package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"net/http"
	"strconv"
	"strings"
	"time"
	"types"

	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"lib"

	"github.com/dchest/captcha"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

const (
	CAPTCHA_LENGTH  = 6
	CAPTCHA_WIDTH   = 180
	CAPTCHA_HEIGHT  = 80
	CAPTCHA_QUALITY = 10
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

	lib.NetHttp.SetNoCache(r.Header())
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
		for index, _ := range *user.Emails {
			if (*user.Emails)[index].Email == email.Email {
				(*user.Emails)[index].Code = ""
				(*user.Emails)[index].Confirmed = true
			}
		}
		sendconfirmation := !user.Confirmed
		user.Confirmed = true

		err = userrepository.Update(user, false, true)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
			return
		}

		if sendconfirmation {
			for _, useremail := range *user.Emails {
				if useremail.Confirmed {
					if helpers.SendConfirmation(user.Language, &useremail, request, r, emailrepository, templaterepository) != nil {
						return
					}
				}
			}
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

// get /api/v1.0/services/suppliers/
func GetSuppliers(w http.ResponseWriter, request *http.Request, r render.Render, supplierfacilityrepository services.SupplierFacilityRepository,
	unitrepository services.UnitRepository, projectrepository services.ProjectRepository, orderrepository services.OrderRepository,
	session *models.DtoSession) {
	dtounit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	query := ""
	var filters *[]models.FilterExp
	filters, err = helpers.GetFilterArray(new(models.SupplierFacilitySearch), nil, request, r, session.Language)
	if err != nil {
		return
	}
	if len(*filters) != 0 {
		var masks []string
		for _, filter := range *filters {
			var exps []string
			for _, field := range filter.Fields {
				switch field {
				case "project":
					project_id, err := strconv.ParseInt(filter.Value, 0, 64)
					if err != nil {
						log.Error("Can't convert to number %v with value %v", err, filter.Value)
						r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
							Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
						return
					}
					err = helpers.CheckProjectAccess(project_id, session.UserID, r, projectrepository, session.Language)
					if err != nil {
						return
					}
					exps = append(exps, fmt.Sprintf("f.supplier_id in (select supplier_id from orders where project_id %v %v)",
						filter.Op, project_id))
				case "order":
					order_id, err := strconv.ParseInt(filter.Value, 0, 64)
					if err != nil {
						log.Error("Can't convert to number %v with value %v", err, filter.Value)
						r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
							Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
						return
					}
					err = helpers.CheckOrderAccess(order_id, session.UserID, r, orderrepository, session.Language)
					if err != nil {
						return
					}
					exps = append(exps, fmt.Sprintf("f.supplier_id in (select supplier_id from orders where id %v %v)",
						filter.Op, order_id))
				default:
					exps = append(exps, field+" "+filter.Op+" "+filter.Value)
				}
			}
			masks = append(masks, "("+strings.Join(exps, " or ")+")")
		}
		query += " and "
		query += strings.Join(masks, " and ")
	}

	var sorts *[]models.OrderExp
	sorts, err = helpers.GetOrderArray(new(models.SupplierFacilitySort), request, r, session.Language)
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

	facilities, err := supplierfacilityrepository.GetByUnit(dtounit.ID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(facilities, len(*facilities), w, r)
}

// get /api/v1.0/services/suppliers/sms/
func GetSMSSuppliers(w http.ResponseWriter, r render.Render, supplierfacilityrepository services.SupplierFacilityRepository, session *models.DtoSession) {
	smsfacilities, err := supplierfacilityrepository.GetByAlias(models.SERVICE_TYPE_SMS)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(smsfacilities, len(*smsfacilities), w, r)
}

// get /api/v1.0/services/suppliers/header/
func GetHeaderSuppliers(w http.ResponseWriter, r render.Render, supplierfacilityrepository services.SupplierFacilityRepository, session *models.DtoSession) {
	headerfacilities, err := supplierfacilityrepository.GetByAlias(models.SERVICE_TYPE_HEADER)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(headerfacilities, len(*headerfacilities), w, r)
}

// get /api/v1.0/services/suppliers/hlr/
func GetHLRSuppliers(w http.ResponseWriter, r render.Render, supplierfacilityrepository services.SupplierFacilityRepository, session *models.DtoSession) {
	hlrfacilities, err := supplierfacilityrepository.GetByAlias(models.SERVICE_TYPE_HLR)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(hlrfacilities, len(*hlrfacilities), w, r)
}

// get /api/v1.0/services/suppliers/recognize/
func GetRecognizeSuppliers(w http.ResponseWriter, r render.Render, supplierfacilityrepository services.SupplierFacilityRepository, session *models.DtoSession) {
	recognizefacilities, err := supplierfacilityrepository.GetByAlias(models.SERVICE_TYPE_RECOGNIZE)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(recognizefacilities, len(*recognizefacilities), w, r)
}

// get /api/v1.0/services/suppliers/verification/
func GetVerifySuppliers(w http.ResponseWriter, r render.Render, supplierfacilityrepository services.SupplierFacilityRepository, session *models.DtoSession) {
	verifyfacilities, err := supplierfacilityrepository.GetByAlias(models.SERVICE_TYPE_VERIFY)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(verifyfacilities, len(*verifyfacilities), w, r)
}

// get /api/v1.0/services/suppliers/sms/price/
func GetSMSPrices(w http.ResponseWriter, r render.Render, pricerepository services.PriceRepository, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository, mobileoperatorrepository services.MobileOperatorRepository, session *models.DtoSession) {
	smshlrprices, err := helpers.GetSMSHLRPrices(models.SERVICE_TYPE_SMS, 0, r, pricerepository, tablecolumnrepository,
		tablerowrepository, mobileoperatorrepository, session.Language, false)
	if err != nil {
		return
	}

	helpers.RenderJSONArray(smshlrprices, len(*smshlrprices), w, r)
}

// get /api/v1.0/services/suppliers/hlr/price/
func GetHLRPrices(w http.ResponseWriter, r render.Render, pricerepository services.PriceRepository, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository, mobileoperatorrepository services.MobileOperatorRepository, session *models.DtoSession) {
	smshlrprices, err := helpers.GetSMSHLRPrices(models.SERVICE_TYPE_HLR, 0, r, pricerepository, tablecolumnrepository,
		tablerowrepository, mobileoperatorrepository, session.Language, false)
	if err != nil {
		return
	}

	helpers.RenderJSONArray(smshlrprices, len(*smshlrprices), w, r)
}

// get /api/v1.0/classification/recognizeproducts/
func GetRecognizeProducts(w http.ResponseWriter, r render.Render, recognizeproductrepository services.RecognizeProductRepository, session *models.DtoSession) {
	recognizeproducts, err := recognizeproductrepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(recognizeproducts, len(*recognizeproducts), w, r)
}

// get /api/v1.0/classification/verificationproducts/
func GetVerifyProducts(w http.ResponseWriter, r render.Render, verifyproductrepository services.VerifyProductRepository, session *models.DtoSession) {
	verifyproducts, err := verifyproductrepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(verifyproducts, len(*verifyproducts), w, r)
}

// get /api/v1.0/classification/header/
func GetHeaderProducts(w http.ResponseWriter, r render.Render, headerproductrepository services.HeaderProductRepository, session *models.DtoSession) {
	headerproducts, err := headerproductrepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(headerproducts, len(*headerproducts), w, r)
}

// get /api/v1.0/services/suppliers/header/price/
func GetHeaderPrices(w http.ResponseWriter, r render.Render, pricerepository services.PriceRepository, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository, headerproductrepository services.HeaderProductRepository, session *models.DtoSession) {
	headerprices, err := helpers.GetHeaderPrices(models.SERVICE_TYPE_HEADER, 0, r, pricerepository, tablecolumnrepository, tablerowrepository,
		headerproductrepository, session.Language, false)
	if err != nil {
		return
	}

	helpers.RenderJSONArray(headerprices, len(*headerprices), w, r)
}

// get /api/v1.0/services/suppliers/recognize/price/
func GetRecognizePrices(w http.ResponseWriter, r render.Render, pricerepository services.PriceRepository, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository, recognizeproductrepository services.RecognizeProductRepository, session *models.DtoSession) {
	recognizeprices, err := helpers.GetRecognizePrices(models.SERVICE_TYPE_RECOGNIZE, 0, r, pricerepository, tablecolumnrepository,
		tablerowrepository, recognizeproductrepository, session.Language, false)
	if err != nil {
		return
	}

	helpers.RenderJSONArray(recognizeprices, len(*recognizeprices), w, r)
}

// get /api/v1.0/services/suppliers/verification/price/
func GetVerifyPrices(w http.ResponseWriter, r render.Render, pricerepository services.PriceRepository, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository, verifyproductrepository services.VerifyProductRepository, session *models.DtoSession) {
	verifyprices, err := helpers.GetVerifyPrices(models.SERVICE_TYPE_VERIFY, 0, r, pricerepository, tablecolumnrepository, tablerowrepository,
		verifyproductrepository, session.Language, false)
	if err != nil {
		return
	}

	helpers.RenderJSONArray(verifyprices, len(*verifyprices), w, r)
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

// get /api/v1.0/classification/tariffplans/
func GetTariffPlans(w http.ResponseWriter, r render.Render, tariffplanrepository services.TariffPlanRepository, session *models.DtoSession) {
	tariffplans, err := tariffplanrepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(tariffplans, len(*tariffplans), w, r)
}

// post /api/v1.0/sayhello/
func CreateFeedback(errors binding.Errors, viewfeedback models.ViewFeedback, request *http.Request, r render.Render,
	feedbackrepository services.FeedbackRepository, captcharepository services.CaptchaRepository,
	sessionrepository services.SessionRepository, emailrepository services.EmailRepository, templaterepository services.TemplateRepository,
	accesslogrepository services.AccessLogRepository) {
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	var user_id int64 = 0
	session, _, err := sessionrepository.GetAndSaveSession(request, r, nil, false, false, true)
	if err == nil {
		user_id = session.UserID
	}
	if user_id == 0 {
		if viewfeedback.CaptchaHash == "" {
			log.Error("Captcha required for submitting feedback")
			r.JSON(helpers.HTTP_STATUS_CAPTCHA_REQUIRED, types.Error{Code: types.TYPE_ERROR_CAPTCHA_REQUIRED,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Captcha_Required})
			return
		}
	}
	if helpers.Check(viewfeedback.CaptchaHash, viewfeedback.CaptchaValue, r, captcharepository) != nil {
		return
	}

	dtofeedback := new(models.DtoFeedback)
	dtofeedback.User_ID = user_id
	dtofeedback.Name = viewfeedback.Name
	dtofeedback.Email = strings.ToLower(viewfeedback.Email)
	dtofeedback.Message = viewfeedback.Message

	dtoaccesslog, err := helpers.CreateAccessLog(request.RequestURI, request, r, accesslogrepository, config.Configuration.Server.DefaultLanguage)
	if err != nil {
		return
	}
	dtofeedback.AccessLog_ID = dtoaccesslog.ID

	err = feedbackrepository.Create(dtofeedback)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	buf, err := templaterepository.GenerateText(models.NewDtoHTMLTemplate(dtofeedback.Message, config.Configuration.Server.DefaultLanguage),
		services.TEMPLATE_FEEDBACK, services.TEMPLATE_DIRECTORY_EMAILS, "")
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	subject := config.Localization[config.Configuration.Server.DefaultLanguage].Messages.FeedbackSubject
	err = emailrepository.SendHTML(config.Configuration.Mail.Receiver, subject, buf.String(), "", dtofeedback.Name+"<"+dtofeedback.Email+">")
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[config.Configuration.Server.DefaultLanguage].Messages.OK})
}
