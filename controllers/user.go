package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
	"types"
)

const (
	PASSWORD_LEVEL_ENCRYPTION = 10
	PASSWORD_LENGTH_MIN       = 8
)

// post /api/v1.0/users/register/:token
func Register(errors binding.Errors, viewuser models.ViewUser, request *http.Request, r render.Render, params martini.Params,
	userrepository services.UserRepository, sessionrepository services.SessionRepository, emailrepository services.EmailRepository,
	unitrepository services.UnitRepository, captcharepository services.CaptchaRepository, templaterepository services.TemplateRepository,
	grouprepository services.GroupRepository) {
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	if captcharepository.Check(viewuser.CaptchaHash, viewuser.CaptchaValue, r) != nil {
		return
	}

	viewuser.Login = strings.ToLower(viewuser.Login)
	emailExists, err := emailrepository.Exists(viewuser.Login)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}
	if emailExists {
		email, err := emailrepository.Get(viewuser.Login)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
			return
		}
		if !email.Confirmed {
			if viewuser.CaptchaHash == "" {
				r.JSON(services.HTTP_STATUS_CAPTCHA_REQUIRED, types.Error{Code: types.TYPE_ERROR_CAPTCHA_REQUIRED,
					Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Captcha_Required})
				return
			}
		} else {
			if viewuser.CaptchaHash == "" {
				r.JSON(services.HTTP_STATUS_CAPTCHA_REQUIRED, types.Error{Code: types.TYPE_ERROR_CAPTCHA_REQUIRED,
					Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Captcha_Required})
				return
			} else {
				r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_EMAIL_INUSE,
					Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Email_InUse})
				return
			}
		}
	}

	roles, err := grouprepository.GetDefault()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	dtouser := new(models.DtoUser)
	dtouser.Created = time.Now()
	dtouser.LastLogin = dtouser.Created
	dtouser.Active = true
	dtouser.Confirmed = false
	dtouser.Name = "Default name for user"
	dtouser.Language = config.Configuration.Server.DefaultLanguage
	dtouser.Roles = *roles

	session, token, err := sessionrepository.GetAndSaveSession(request, r, params, false, true, true)
	if err == nil {
		dtouser.Creator_ID = session.UserID

		unit, err := unitrepository.FindByUser(session.UserID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
			return
		}
		dtouser.UnitID = unit.ID
		dtouser.UnitAdmin = false
	} else {
		dtouser.UnitID = 0
		dtouser.UnitAdmin = true
	}

	token, err = sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}
	dtouser.Code = token
	dtouser.Password = token

	dtoemail := new(models.DtoEmail)
	dtoemail.Email = viewuser.Login
	dtoemail.Created = dtouser.LastLogin
	dtoemail.Primary = true
	dtoemail.Confirmed = false
	dtoemail.Subscription = false
	dtoemail.Code = dtouser.Code
	dtoemail.Language = dtouser.Language
	dtoemail.Exists = emailExists

	dtouser.Emails = &[]models.DtoEmail{*dtoemail}

	err = userrepository.Create(dtouser, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	token, err = sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	session = models.NewDtoSession(token, dtouser.ID, dtouser.Roles, dtouser.LastLogin, dtouser.Language)
	err = sessionrepository.Create(session, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	for _, confEmail := range *dtouser.Emails {
		subject := config.Localization[confEmail.Language].Messages.RegistrationSubject
		buf, err := templaterepository.GenerateText(models.NewDtoTemplate(confEmail.Email, confEmail.Language, request.Host, confEmail.Code),
			services.TEMPLATE_EMAIL, services.TEMPLATE_LAYOUT)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
			return
		}

		err = emailrepository.SendEmail(confEmail.Email, subject, buf.String())
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
			return
		}
	}

	r.JSON(http.StatusOK, models.NewApiSession(session.LastActivity.Add(config.Configuration.Server.SessionTimeout), session.AccessToken))
}

// post /api/v1.0/users/password/
func RestorePassword(errors binding.Errors, viewuser models.ViewUser, request *http.Request, r render.Render,
	emailrepository services.EmailRepository, userrepository services.UserRepository, sessionrepository services.SessionRepository,
	captcharepository services.CaptchaRepository, templaterepository services.TemplateRepository) {
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	if captcharepository.Check(viewuser.CaptchaHash, viewuser.CaptchaValue, r) != nil {
		return
	}

	viewuser.Login = strings.ToLower(viewuser.Login)
	dtouser, err := userrepository.FindByLogin(viewuser.Login)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	if !dtouser.Active || !dtouser.Confirmed {
		log.Error("User is not active or confirmed %v", dtouser.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_USER_BLOCKED,
			Message: config.Localization[dtouser.Language].Errors.Api.User_Blocked})
		return
	}

	token, err := sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[dtouser.Language].Errors.Api.Data_Wrong})
		return
	}
	dtouser.Code = token

	err = userrepository.Update(dtouser, true, false)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[dtouser.Language].Errors.Api.Data_Wrong})
		return
	}

	for _, confEmail := range *dtouser.Emails {
		if confEmail.Confirmed {
			if helpers.SendPassword(dtouser.Language, &confEmail, dtouser, request, r, emailrepository, templaterepository) != nil {
				return
			}
		}
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[dtouser.Language].Messages.OK})
}

// put /api/v1.0/users/password/:code/
func UpdatePassword(errors binding.Errors, password models.PasswordUpdate, r render.Render, params martini.Params,
	userrepository services.UserRepository, sessionrepository services.SessionRepository) {
	code := params[helpers.PARAMETER_NAME_CODE]
	if len(code) > helpers.PARAM_LENGTH_MAX {
		log.Error("Wrong parameter length %v", code)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	for _, errBind := range errors {
		for _, field := range errBind.FieldNames {
			if field == helpers.PARAMETER_NAME_CODE {
				if (code == "") || (len(code) > helpers.PARAM_LENGTH_MAX) {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
					return
				}
			} else {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
				return
			}
		}

	}

	if code != "" {
		password.Code = code
	}

	if len(password.Value) < PASSWORD_LENGTH_MIN {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_PASSWORD_TOOSIMPLE,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Password_Too_Simple})
		return
	}

	user, err := userrepository.FindByCode(password.Code)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	if !user.Active || !user.Confirmed {
		log.Error("User is not active or confirmed %v", user.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_USER_BLOCKED,
			Message: config.Localization[user.Language].Errors.Api.User_Blocked})
		return
	}

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password.Value), PASSWORD_LEVEL_ENCRYPTION)
	if err != nil {
		log.Error("Password generating error %v", err)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
		return
	}

	user.Password = string(hash[:])
	user.Code = ""
	user.LastLogin = time.Now()

	err = userrepository.Update(user, true, false)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
		return
	}

	token, err := sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
		return
	}

	session := models.NewDtoSession(token, user.ID, user.Roles, user.LastLogin, user.Language)
	err = sessionrepository.Create(session, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiSession(session.LastActivity.Add(config.Configuration.Server.SessionTimeout), session.AccessToken))
}

// delete /api/v1.0/users/password/:code/
func DeletePasswordRestoring(r render.Render, params martini.Params, userrepository services.UserRepository) {
	code := params[helpers.PARAMETER_NAME_CODE]
	if code == "" || len(code) > helpers.PARAM_LENGTH_MAX {
		log.Error("Wrong parameter length %v", code)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_CONFIRMATION_CODE_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Confirmation_Code_Wrong})
		return
	}

	user, err := userrepository.FindByCode(code)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_CONFIRMATION_CODE_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Confirmation_Code_Wrong})
		return
	}

	if !user.Active || !user.Confirmed {
		log.Error("User is not active or confirmed %v", user.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_USER_BLOCKED,
			Message: config.Localization[user.Language].Errors.Api.User_Blocked})
		return
	}

	user.Code = ""

	err = userrepository.Update(user, true, false)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[user.Language].Messages.OK})
}

// get /api/v1.0/user/groups/
func GetGroups(r render.Render, grouprepository services.GroupRepository, session *models.DtoSession) {
	groups, err := grouprepository.GetByUserExt(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, groups)
}

// get /api/v1.0/user/
func GetUserInfo(r render.Render, userrepository services.UserRepository, session *models.DtoSession) {
	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiUserLong(user.ID, user.UnitID, user.UnitAdmin, user.Active,
		user.Confirmed, user.Name, user.Language))
}

// put /api/v1.0/user/
func UpdateUserInfo(errors binding.Errors, changeuser models.ChangeUser, r render.Render,
	userrepository services.UserRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	user.Name = changeuser.Name
	if changeuser.Language != "" {
		user.Language = changeuser.Language
	} else {
		user.Language = session.Language
	}

	err = userrepository.Update(user, true, false)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiUserLong(user.ID, user.UnitID, user.UnitAdmin, user.Active,
		user.Confirmed, user.Name, user.Language))
}

// get /api/v1.0/user/emails/
func GetUserEmails(r render.Render, userrepository services.UserRepository, session *models.DtoSession) {
	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	emails := new([]models.ViewApiEmail)
	for _, email := range *user.Emails {
		*emails = append(*emails, *models.NewViewApiEmail(email.Email, email.Primary, email.Confirmed, email.Subscription, email.Language))
	}

	r.JSON(http.StatusOK, emails)
}

// put /api/v1.0/user/emails/
func UpdateUserEmails(errors binding.Errors, updateemails models.UpdateEmails, request *http.Request, r render.Render,
	sessionrepository services.SessionRepository, emailrepository services.EmailRepository, userrepository services.UserRepository,
	templaterepository services.TemplateRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	count := 0
	for _, checkEmail := range updateemails {
		if checkEmail.Primary {
			count++
		}
	}
	if count != 1 {
		log.Error("Only one primary email is allowed")
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	arrEmails := new([]models.DtoEmail)

	var updEmail models.UpdateEmail
	var curEmail models.DtoEmail

	for _, updEmail = range updateemails {
		updEmail.Email = strings.ToLower(updEmail.Email)
		updEmail.Confirmed = false
		found := false
		code := ""
		for _, curEmail = range *user.Emails {
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

			*arrEmails = append(*arrEmails, models.DtoEmail{
				Email:        updEmail.Email,
				UserID:       user.ID,
				Created:      time.Now(),
				Primary:      updEmail.Primary,
				Confirmed:    updEmail.Confirmed,
				Subscription: updEmail.Subscription,
				Code:         "",
				Language:     updEmail.Language,
				Exists:       emailExists,
			})
		} else {
			updEmail.Confirmed = curEmail.Confirmed
			if updEmail.Primary != curEmail.Primary ||
				updEmail.Subscription != curEmail.Subscription ||
				updEmail.Language != curEmail.Language {
				updEmail.Confirmed = false
			}

			curEmail.Primary = updEmail.Primary
			curEmail.Confirmed = updEmail.Confirmed
			curEmail.Subscription = updEmail.Subscription
			curEmail.Language = updEmail.Language

			*arrEmails = append(*arrEmails, curEmail)
		}

		if !updEmail.Confirmed {
			code, err = sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}

			(*arrEmails)[len(*arrEmails)-1].Code = code

			if updEmail.Primary {
				user.Active = true
				user.Confirmed = false

				user.Code = code
				user.Password = code
			}
		}
	}

	user.Emails = arrEmails
	err = userrepository.Update(user, false, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	if helpers.SendConfirmations(user, session, request, r, emailrepository, templaterepository) != nil {
		return
	}

	emails := new([]models.ViewApiEmail)
	for _, email := range *user.Emails {
		*emails = append(*emails, *models.NewViewApiEmail(email.Email, email.Primary, email.Confirmed, email.Subscription, email.Language))
	}

	r.JSON(http.StatusOK, emails)
}
