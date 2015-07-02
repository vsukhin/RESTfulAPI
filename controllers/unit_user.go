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

// options /api/v1.0/user/administration/users/
func GetMetaUnitUser(r render.Render, userrepository services.UserRepository, session *models.DtoSession) {
	metaunituser, err := userrepository.GetMetaByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, metaunituser)
}

// get /api/v1.0/user/administration/users/
func GetUnitUsers(w http.ResponseWriter, request *http.Request, r render.Render,
	userrepository services.UserRepository, session *models.DtoSession) {
	query := ""

	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.SearchUnitUser), nil, request, r, session.Language)
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
	sorts, err = helpers.GetOrderArray(new(models.SearchUnitUser), request, r, session.Language)
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

	users, err := userrepository.GetAllByUser(session.UserID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(users, len(*users), w, r)
}

// post /api/v1.0/user/administration/users/
func CreateUnitUser(errors binding.Errors, user models.ViewShortUnitUser, request *http.Request, r render.Render,
	userrepository services.UserRepository, emailrepository services.EmailRepository, sessionrepository services.SessionRepository,
	unitrepository services.UnitRepository, templaterepository services.TemplateRepository, grouprepository services.GroupRepository,
	classifierrepository services.ClassifierRepository, mobilephonerepository services.MobilePhoneRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&user, errors, r, session.Language) != nil {
		return
	}

	dtouser := new(models.DtoUser)
	dtouser.Created = time.Now()
	dtouser.LastLogin = dtouser.Created
	dtouser.Active = true
	dtouser.Confirmed = false
	dtouser.Surname = user.Surname
	dtouser.Name = user.Name
	dtouser.MiddleName = user.MiddleName
	dtouser.WorkPhone = user.WorkPhone
	dtouser.JobTitle = user.JobTitle
	dtouser.Language = strings.ToLower(user.Language)
	dtouser.Code = ""
	dtouser.Password = ""
	dtouser.ReportAccess = true
	dtouser.CaptchaRequired = false

	roles, err := grouprepository.GetDefault()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	dtouser.Roles = *roles
	dtouser.Creator_ID = session.UserID

	unit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	dtouser.UnitID = unit.ID
	dtouser.UnitAdmin = user.UnitAdmin

	if helpers.CheckPrimaryUnitUserEmail(user.Emails, session.Language, r) != nil {
		return
	}
	dtouser.Emails = new([]models.DtoEmail)
	if helpers.CheckPrimaryUnitUserMobilePhone(user.MobilePhones, session.Language, r) != nil {
		return
	}
	dtouser.MobilePhones = new([]models.DtoMobilePhone)

	for _, updEmail := range user.Emails {
		updEmail.Email = strings.ToLower(updEmail.Email)
		emailExists, err := helpers.CheckEmailAvailability(updEmail.Email, session.Language, r, emailrepository)
		if err != nil {
			return
		}
		code := ""
		classifier, err := helpers.CheckClassifier(updEmail.Classifier_ID, r, classifierrepository, session.Language)
		if err != nil {
			return
		}

		code, err = sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}
		if updEmail.Primary {
			dtouser.Code = code
			dtouser.Password = code
		}
		*dtouser.Emails = append(*dtouser.Emails, models.DtoEmail{
			Email:         updEmail.Email,
			Created:       dtouser.LastLogin,
			Primary:       updEmail.Primary,
			Confirmed:     false,
			Subscription:  false,
			Code:          code,
			Language:      strings.ToLower(updEmail.Language),
			Exists:        emailExists,
			Classifier_ID: classifier.ID,
		})
	}

	for _, updMobilePhone := range user.MobilePhones {
		phoneExists, err := helpers.CheckMobilePhoneAvailability(updMobilePhone.Phone, session.Language, r, mobilephonerepository)
		if err != nil {
			return
		}
		classifier, err := helpers.CheckClassifier(updMobilePhone.Classifier_ID, r, classifierrepository, session.Language)
		if err != nil {
			return
		}

		*dtouser.MobilePhones = append(*dtouser.MobilePhones, models.DtoMobilePhone{
			Phone:         updMobilePhone.Phone,
			Created:       dtouser.LastLogin,
			Primary:       updMobilePhone.Primary,
			Confirmed:     false,
			Subscription:  false,
			Code:          "",
			Language:      strings.ToLower(updMobilePhone.Language),
			Exists:        phoneExists,
			Classifier_ID: classifier.ID,
		})
	}

	err = userrepository.Create(dtouser, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	if helpers.SendConfirmations(dtouser, session, request, r, emailrepository, templaterepository, true) != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiUserTiny(dtouser.ID))
}

// get /api/v1.0/user/administration/users/:uid/
func GetUnitUser(r render.Render, params martini.Params, userrepository services.UserRepository, session *models.DtoSession) {
	user, err := helpers.CheckUnitUser(session.UserID, r, params, userrepository, session.Language)
	if err != nil {
		return
	}

	apiuser := models.NewApiShortUnitUser(user.UnitAdmin, user.Active, user.Confirmed, user.Surname, user.Name, user.MiddleName,
		user.WorkPhone, user.JobTitle, user.Language, *new([]models.ViewApiEmail), *new([]models.ViewApiMobilePhone))
	for _, email := range *user.Emails {
		apiuser.Emails = append(apiuser.Emails, *models.NewViewApiEmail(email.Email, email.Primary, email.Confirmed, /*email.Subscription,*/
			email.Language, email.Classifier_ID))
	}
	for _, mobilephone := range *user.MobilePhones {
		apiuser.MobilePhones = append(apiuser.MobilePhones, *models.NewViewApiMobilePhone(mobilephone.Phone, mobilephone.Primary,
			mobilephone.Confirmed /*mobilephone.Subscription,*/, mobilephone.Language, mobilephone.Classifier_ID))
	}

	r.JSON(http.StatusOK, apiuser)
}

// put /api/v1.0/user/administration/users/:uid/
func UpdateUnitUser(errors binding.Errors, user models.ViewLongUnitUser, request *http.Request, r render.Render, params martini.Params,
	userrepository services.UserRepository, emailrepository services.EmailRepository, sessionrepository services.SessionRepository,
	unitrepository services.UnitRepository, templaterepository services.TemplateRepository, grouprepository services.GroupRepository,
	classifierrepository services.ClassifierRepository, mobilephonerepository services.MobilePhoneRepository,
	session *models.DtoSession) {
	if helpers.CheckValidation(&user, errors, r, session.Language) != nil {
		return
	}
	dtouser, err := helpers.CheckUnitUser(session.UserID, r, params, userrepository, session.Language)
	if err != nil {
		return
	}

	dtouser.Surname = user.Surname
	dtouser.Name = user.Name
	dtouser.MiddleName = user.MiddleName
	dtouser.WorkPhone = user.WorkPhone
	dtouser.JobTitle = user.JobTitle
	dtouser.Language = strings.ToLower(user.Language)

	if helpers.CheckUserRoles(user.Roles, session.Language, r, grouprepository) != nil {
		return
	}
	if middlewares.IsAdmin(user.Roles) {
		log.Error("Can't give administrator role to user at unit administration")
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	dtouser.Roles = user.Roles
	dtouser.UnitAdmin = user.UnitAdmin

	if helpers.CheckPrimaryUnitUserEmail(user.Emails, session.Language, r) != nil {
		return
	}
	if helpers.CheckPrimaryUnitUserMobilePhone(user.MobilePhones, session.Language, r) != nil {
		return
	}

	arrInEmails := user.Emails
	arrOutEmails := new([]models.DtoEmail)

	var updEmail models.ViewEmail
	var curEmail models.DtoEmail

	for _, updEmail = range arrInEmails {
		updEmail.Email = strings.ToLower(updEmail.Email)
		found := false
		code := ""
		classifier, err := helpers.CheckClassifier(updEmail.Classifier_ID, r, classifierrepository, session.Language)
		if err != nil {
			return
		}
		for _, curEmail = range *dtouser.Emails {
			if updEmail.Email == curEmail.Email {
				found = true
				break
			}
		}

		if !found {
			var emailExists bool
			emailExists, err = helpers.CheckEmailAvailability(updEmail.Email, session.Language, r, emailrepository)
			if err != nil {
				return
			}

			*arrOutEmails = append(*arrOutEmails, models.DtoEmail{
				Email:         updEmail.Email,
				UserID:        dtouser.ID,
				Created:       time.Now(),
				Primary:       updEmail.Primary,
				Confirmed:     false,
				Subscription:  false,
				Code:          code,
				Language:      strings.ToLower(updEmail.Language),
				Exists:        emailExists,
				Classifier_ID: classifier.ID,
			})
		} else {
			curEmail.Primary = updEmail.Primary
			curEmail.Language = strings.ToLower(updEmail.Language)
			curEmail.Classifier_ID = classifier.ID

			*arrOutEmails = append(*arrOutEmails, curEmail)
		}

		if !found {
			code, err = sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}

			(*arrOutEmails)[len(*arrOutEmails)-1].Code = code
		}

		if updEmail.Primary {
			if !found || !curEmail.Confirmed {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
		}
	}

	arrInMobilePhones := user.MobilePhones
	arrOutMobilePhones := new([]models.DtoMobilePhone)

	var updMobilePhone models.ViewMobilePhone
	var curMobilePhone models.DtoMobilePhone

	for _, updMobilePhone = range arrInMobilePhones {
		found := false
		classifier, err := helpers.CheckClassifier(updMobilePhone.Classifier_ID, r, classifierrepository, session.Language)
		if err != nil {
			return
		}
		for _, curMobilePhone = range *dtouser.MobilePhones {
			if updMobilePhone.Phone == curMobilePhone.Phone {
				found = true
				break
			}
		}

		if !found {
			var phoneExists bool
			phoneExists, err = helpers.CheckMobilePhoneAvailability(updMobilePhone.Phone, session.Language, r, mobilephonerepository)
			if err != nil {
				return
			}

			*arrOutMobilePhones = append(*arrOutMobilePhones, models.DtoMobilePhone{
				Phone:         updMobilePhone.Phone,
				UserID:        dtouser.ID,
				Created:       time.Now(),
				Primary:       updMobilePhone.Primary,
				Confirmed:     false,
				Subscription:  false,
				Code:          "",
				Language:      strings.ToLower(updMobilePhone.Language),
				Exists:        phoneExists,
				Classifier_ID: classifier.ID,
			})
		} else {
			curMobilePhone.Primary = updMobilePhone.Primary
			curMobilePhone.Language = strings.ToLower(updMobilePhone.Language)
			curMobilePhone.Classifier_ID = classifier.ID

			*arrOutMobilePhones = append(*arrOutMobilePhones, curMobilePhone)
		}
	}

	dtouser.Emails = arrOutEmails
	dtouser.MobilePhones = arrOutMobilePhones
	err = userrepository.Update(dtouser, false, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	if helpers.SendConfirmations(dtouser, session, request, r, emailrepository, templaterepository, false) != nil {
		return
	}

	apiuser := models.NewApiLongUnitUser(dtouser.UnitAdmin, dtouser.Active, dtouser.Confirmed, dtouser.Surname, dtouser.Name, dtouser.MiddleName,
		dtouser.WorkPhone, dtouser.JobTitle, dtouser.Language, dtouser.Roles, *new([]models.ViewApiEmail), *new([]models.ViewApiMobilePhone))
	for _, email := range *dtouser.Emails {
		apiuser.Emails = append(apiuser.Emails, *models.NewViewApiEmail(email.Email, email.Primary, email.Confirmed, /*email.Subscription,*/
			email.Language, email.Classifier_ID))
	}
	for _, mobilephone := range *dtouser.MobilePhones {
		apiuser.MobilePhones = append(apiuser.MobilePhones, *models.NewViewApiMobilePhone(mobilephone.Phone, mobilephone.Primary,
			mobilephone.Confirmed /*mobilephone.Subscription,*/, mobilephone.Language, mobilephone.Classifier_ID))
	}

	r.JSON(http.StatusOK, apiuser)
}

// delete /api/v1.0/user/administration/users/:uid/
func DeleteUnitUser(r render.Render, params martini.Params, userrepository services.UserRepository, session *models.DtoSession) {
	dtouser, err := helpers.CheckUnitUser(session.UserID, r, params, userrepository, session.Language)
	if err != nil {
		return
	}

	err = userrepository.Delete(dtouser.ID, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
