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
	grouprepository services.GroupRepository, classifierrepository services.ClassifierRepository,
	accesslogrepository services.AccessLogRepository) {
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	if helpers.Check(viewuser.CaptchaHash, viewuser.CaptchaValue, r, captcharepository) != nil {
		return
	}

	viewuser.Login = strings.ToLower(viewuser.Login)
	emailExists, err := emailrepository.Exists(viewuser.Login)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	if emailExists {
		email, err := emailrepository.Get(viewuser.Login)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
			return
		}
		if !email.Confirmed {
			if viewuser.CaptchaHash == "" {
				r.JSON(helpers.HTTP_STATUS_CAPTCHA_REQUIRED, types.Error{Code: types.TYPE_ERROR_CAPTCHA_REQUIRED,
					Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Captcha_Required})
				return
			}
		} else {
			if viewuser.CaptchaHash == "" {
				r.JSON(helpers.HTTP_STATUS_CAPTCHA_REQUIRED, types.Error{Code: types.TYPE_ERROR_CAPTCHA_REQUIRED,
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
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	dtouser := new(models.DtoUser)
	dtouser.Created = time.Now()
	dtouser.LastLogin = dtouser.Created
	dtouser.Active = true
	dtouser.Confirmed = false
	dtouser.Surname = models.USER_SURNAME_DEFAULT
	dtouser.Name = models.USER_NAME_DEFAULT
	dtouser.Language = config.Configuration.Server.DefaultLanguage
	dtouser.ReportAccess = true
	dtouser.CaptchaRequired = false
	dtouser.NewsBlocked = false
	dtouser.Roles = *roles

	session, token, err := sessionrepository.GetAndSaveSession(request, r, params, false, true, true)
	if err == nil {
		dtouser.Creator_ID = session.UserID

		unit, err := unitrepository.FindByUser(session.UserID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
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

	classifier, err := helpers.CheckClassifier(models.CLASSIFIER_TYPE_MAIN, r, classifierrepository, config.Configuration.Server.DefaultLanguage)
	if err != nil {
		return
	}
	dtoemail := new(models.DtoEmail)
	dtoemail.Email = viewuser.Login
	dtoemail.Created = dtouser.LastLogin
	dtoemail.Primary = true
	dtoemail.Confirmed = false
	dtoemail.Subscription = false
	dtoemail.Code = dtouser.Code
	dtoemail.Language = dtouser.Language
	dtoemail.Exists = emailExists
	dtoemail.Classifier_ID = classifier.ID

	dtouser.Emails = &[]models.DtoEmail{*dtoemail}
	dtouser.MobilePhones = new([]models.DtoMobilePhone)

	dtoaccesslog, err := helpers.CreateAccessLog(dtouser.Code, request, r, accesslogrepository, config.Configuration.Server.DefaultLanguage)
	if err != nil {
		return
	}
	dtouser.Subscr_AccessLog_ID = dtoaccesslog.ID

	err = userrepository.Create(dtouser, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	if helpers.SendPasswordRegistration(config.Configuration.Server.DefaultLanguage, dtoemail, dtouser, request, r, emailrepository, templaterepository) != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiSession(dtouser.LastLogin, ""))
}

// post /api/v1.0/users/password/
func RestorePassword(errors binding.Errors, viewuser models.ViewUser, request *http.Request, r render.Render,
	emailrepository services.EmailRepository, userrepository services.UserRepository, sessionrepository services.SessionRepository,
	captcharepository services.CaptchaRepository, templaterepository services.TemplateRepository) {
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}

	viewuser.Login = strings.ToLower(viewuser.Login)
	dtouser, err := userrepository.FindByLogin(viewuser.Login)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	if dtouser.CaptchaRequired && viewuser.CaptchaHash == "" {
		r.JSON(helpers.HTTP_STATUS_CAPTCHA_REQUIRED, types.Error{Code: types.TYPE_ERROR_CAPTCHA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Captcha_Wrong})
		return
	}
	if helpers.Check(viewuser.CaptchaHash, viewuser.CaptchaValue, r, captcharepository) != nil {
		dtouser.CaptchaRequired = true
		err = userrepository.Update(dtouser, true, false)
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

	dtouser.CaptchaRequired = false
	err = userrepository.Update(dtouser, true, false)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[dtouser.Language].Errors.Api.Data_Wrong})
		return
	}

	for _, confEmail := range *dtouser.Emails {
		if confEmail.Confirmed {
			if helpers.SendPasswordRecovery(dtouser.Language, &confEmail, dtouser, request, r, emailrepository, templaterepository) != nil {
				return
			}
		}
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[dtouser.Language].Messages.OK})
}

// put /api/v1.0/users/password/:code/
func UpdatePassword(errors binding.Errors, password models.PasswordUpdate, request *http.Request, r render.Render, params martini.Params,
	userrepository services.UserRepository, sessionrepository services.SessionRepository, emailrepository services.EmailRepository,
	templaterepository services.TemplateRepository, accesslogrepository services.AccessLogRepository) {
	code := params[helpers.PARAMETER_NAME_CODE]
	if len(code) > helpers.PARAM_LENGTH_MAX {
		log.Error("Wrong parameter length %v", code)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	for _, errBind := range errors {
		for _, field := range errBind.FieldNames {
			if field == helpers.PARAMETER_NAME_CODE {
				if (code == "") || (len(code) > helpers.PARAM_LENGTH_MAX) {
					r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
					return
				}
			} else {
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
				return
			}
		}

	}

	if code != "" {
		password.Code = code
	}

	if len(password.Value) < PASSWORD_LENGTH_MIN {
		log.Error("Password is too simple %v", password.Value)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_PASSWORD_TOOSIMPLE,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Password_Too_Simple})
		return
	}

	user, err := userrepository.FindByCode(password.Code)
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

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password.Value), PASSWORD_LEVEL_ENCRYPTION)
	if err != nil {
		log.Error("Password generating error %v", err)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
		return
	}

	sendconfirmation := !user.Confirmed
	user.Confirmed = true
	user.Password = string(hash[:])
	user.Code = ""
	user.LastLogin = time.Now()

	for i, _ := range *user.Emails {
		if (*user.Emails)[i].Code == password.Code {
			(*user.Emails)[i].Confirmed = true
			(*user.Emails)[i].Code = ""
		}
	}

	if sendconfirmation {
		dtoaccesslog, err := helpers.CreateAccessLog(password.Code, request, r, accesslogrepository, user.Language)
		if err != nil {
			return
		}
		user.Conf_AccessLog_ID = dtoaccesslog.ID
	}

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

// head /api/v1.0/users/password/:code/
func CheckPasswordRestoring(r render.Render, params martini.Params, userrepository services.UserRepository) {
	code := params[helpers.PARAMETER_NAME_CODE]
	if code == "" || len(code) > helpers.PARAM_LENGTH_MAX {
		log.Error("Wrong parameter length %v", code)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_CONFIRMATION_CODE_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Confirmation_Code_Wrong})
		return
	}

	user, err := userrepository.FindByCode(code)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_CONFIRMATION_CODE_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Confirmation_Code_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[user.Language].Messages.OK})
}

// delete /api/v1.0/users/password/:code/
func DeletePasswordRestoring(r render.Render, params martini.Params, userrepository services.UserRepository) {
	code := params[helpers.PARAMETER_NAME_CODE]
	if code == "" || len(code) > helpers.PARAM_LENGTH_MAX {
		log.Error("Wrong parameter length %v", code)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_CONFIRMATION_CODE_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Confirmation_Code_Wrong})
		return
	}

	user, err := userrepository.FindByCode(code)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_CONFIRMATION_CODE_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Confirmation_Code_Wrong})
		return
	}

	if !user.Active {
		log.Error("User is not active%v", user.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_USER_BLOCKED,
			Message: config.Localization[user.Language].Errors.Api.User_Blocked})
		return
	}

	if user.Confirmed {
		user.Code = ""

		err = userrepository.Update(user, true, false)
	} else {
		err = userrepository.Delete(user.ID, true)
	}
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[user.Language].Messages.OK})
}

// get /api/v1.0/user/groups/
func GetGroups(w http.ResponseWriter, r render.Render, grouprepository services.GroupRepository, session *models.DtoSession) {
	groups, err := grouprepository.GetByUserExt(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(groups, len(*groups), w, r)
}

// get /api/v1.0/user/
func GetUserInfo(r render.Render, userrepository services.UserRepository, session *models.DtoSession) {
	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiUserLong(user.ID, user.UnitID, user.UnitAdmin, user.Active,
		user.Confirmed, user.Surname, user.Name, user.MiddleName, user.WorkPhone, user.JobTitle, user.Language))
}

// patch /api/v1.0/user/
func UpdateUserInfo(errors binding.Errors, changeuser models.ChangeUser, r render.Render,
	userrepository services.UserRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	user.Surname = changeuser.Surname
	user.Name = changeuser.Name
	if user.Surname == "" {
		user.Surname = models.USER_SURNAME_DEFAULT
	}
	if user.Name == "" {
		user.Name = models.USER_NAME_DEFAULT
	}
	user.MiddleName = changeuser.MiddleName
	user.WorkPhone = changeuser.WorkPhone
	user.JobTitle = changeuser.JobTitle
	if changeuser.Language != "" {
		user.Language = strings.ToLower(changeuser.Language)
	} else {
		user.Language = session.Language
	}

	err = userrepository.UpdateProfile(user)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiUserLong(user.ID, user.UnitID, user.UnitAdmin, user.Active,
		user.Confirmed, user.Surname, user.Name, user.MiddleName, user.WorkPhone, user.JobTitle, user.Language))
}

// get /api/v1.0/user/emails/
func GetUserEmails(w http.ResponseWriter, r render.Render, userrepository services.UserRepository, session *models.DtoSession) {
	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	emails := new([]models.ViewApiEmail)
	for _, email := range *user.Emails {
		*emails = append(*emails, *models.NewViewApiEmail(email.Email, email.Primary, email.Confirmed, /*email.Subscription,*/
			email.Language, email.Classifier_ID))
	}

	helpers.RenderJSONArray(emails, len(*emails), w, r)
}

// put /api/v1.0/user/emails/
func UpdateUserEmails(w http.ResponseWriter, errors binding.Errors, updateemails models.UpdateEmails, request *http.Request, r render.Render,
	sessionrepository services.SessionRepository, emailrepository services.EmailRepository, userrepository services.UserRepository,
	templaterepository services.TemplateRepository, classifierrepository services.ClassifierRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	count := 0
	for _, checkEmail := range updateemails {
		if checkEmail.Primary {
			count++
		}
	}
	if count > 1 {
		log.Error("Only one primary email is allowed")
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_PRIMARY_EMAIL_NOTSINGLE,
			Message: config.Localization[session.Language].Errors.Api.PrimaryEmail_NotSingle})
		return
	}

	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	count_old := 0
	for _, checkEmail := range *user.Emails {
		if checkEmail.Primary {
			count_old++
		}
	}
	if count_old > count {
		log.Error("Can't delete primary email")
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_PRIMARY_EMAIL_NOTSINGLE,
			Message: config.Localization[session.Language].Errors.Api.PrimaryEmail_NotSingle})
		return
	}

	arrEmails := new([]models.DtoEmail)

	var updEmail models.ViewEmail
	var curEmail models.DtoEmail

	for _, updEmail = range updateemails {
		updEmail.Email = strings.ToLower(updEmail.Email)
		found := false
		code := ""
		for _, curEmail = range *user.Emails {
			if updEmail.Email == curEmail.Email {
				found = true
				break
			}
		}
		classifier, err := helpers.CheckClassifier(updEmail.Classifier_ID, r, classifierrepository, session.Language)
		if err != nil {
			return
		}

		if !found {
			var emailExists bool
			emailExists, err = helpers.CheckEmailAvailability(updEmail.Email, session.Language, r, emailrepository)
			if err != nil {
				return
			}

			*arrEmails = append(*arrEmails, models.DtoEmail{
				Email:         updEmail.Email,
				UserID:        user.ID,
				Classifier_ID: classifier.ID,
				Created:       time.Now(),
				Primary:       updEmail.Primary,
				Confirmed:     false,
				//				Subscription:  updEmail.Subscription,
				Code:     code,
				Language: strings.ToLower(updEmail.Language),
				Exists:   emailExists,
			})
		} else {
			curEmail.Primary = updEmail.Primary
			//			curEmail.Subscription = updEmail.Subscription
			curEmail.Language = strings.ToLower(updEmail.Language)
			curEmail.Classifier_ID = classifier.ID

			*arrEmails = append(*arrEmails, curEmail)
		}

		if !found {
			code, err = sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}

			(*arrEmails)[len(*arrEmails)-1].Code = code
		}

		if updEmail.Primary {
			if !found || !curEmail.Confirmed {
				log.Error("Primary email is not confirmed %v", updEmail.Email)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_PRIMARY_EMAIL_NOTCONFIRMED,
					Message: config.Localization[session.Language].Errors.Api.PrimaryEmail_NotConfirmed})
				return
			}
		}
	}

	user.Emails = arrEmails
	err = userrepository.UpdateEmails(user, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	if helpers.SendConfirmations(user, session, request, r, emailrepository, templaterepository, false) != nil {
		return
	}

	emails := new([]models.ViewApiEmail)
	for _, email := range *user.Emails {
		*emails = append(*emails, *models.NewViewApiEmail(email.Email, email.Primary, email.Confirmed, /*email.Subscription,*/
			email.Language, email.Classifier_ID))
	}

	helpers.RenderJSONArray(emails, len(*emails), w, r)
}

// patch /api/v1.0/user/password/
func ChangePassword(errors binding.Errors, changepassword models.ChangePassword, r render.Render,
	userrepository services.UserRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	if changepassword.NewPassword != changepassword.ConfirmPassword {
		log.Error("New %v and confirm %v passwords don't match each other", changepassword.NewPassword, changepassword.ConfirmPassword)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}
	if len([]rune(changepassword.NewPassword)) < PASSWORD_LENGTH_MIN {
		log.Error("Password is too simple %v", changepassword.NewPassword)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_PASSWORD_TOOSIMPLE,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Password_Too_Simple})
		return
	}

	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changepassword.OldPassword)) != nil {
		log.Error("Can't compare password hashes")
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_LOGIN_OR_PASSWORD_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Login_Or_Password_Wrong})
		return
	}

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(changepassword.NewPassword), PASSWORD_LEVEL_ENCRYPTION)
	if err != nil {
		log.Error("Password generating error %v", err)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
		return
	}

	user.Password = string(hash[:])
	err = userrepository.UpdatePassword(user)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// get /api/v1.0/user/mobilephones/
func GetUserMobilePhones(w http.ResponseWriter, r render.Render, userrepository services.UserRepository, session *models.DtoSession) {
	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	phones := new([]models.ViewApiMobilePhone)
	for _, phone := range *user.MobilePhones {
		*phones = append(*phones, *models.NewViewApiMobilePhone(phone.Phone, phone.Primary, phone.Confirmed, /*phone.Subscription,*/
			phone.Language, phone.Classifier_ID))
	}

	helpers.RenderJSONArray(phones, len(*phones), w, r)
}

// put /api/v1.0/user/mobilephones/
func UpdateUserMobilePhones(w http.ResponseWriter, errors binding.Errors, updatephones models.UpdateMobilePhones, request *http.Request, r render.Render,
	mobilephonerepository services.MobilePhoneRepository, userrepository services.UserRepository,
	classifierrepository services.ClassifierRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	count := 0
	for _, checkMobilePhone := range updatephones {
		if checkMobilePhone.Primary {
			count++
		}
	}
	if count > 1 {
		log.Error("Only one primary mobile phone is allowed")
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_PRIMARY_MOBILEPHONE_NOTSINGLE,
			Message: config.Localization[session.Language].Errors.Api.PrimaryMobilePhone_NotSingle})
		return
	}

	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	count_old := 0
	for _, checkMobilePhone := range *user.MobilePhones {
		if checkMobilePhone.Primary {
			count_old++
		}
	}
	if count_old > count {
		log.Error("Can't delete primary mobile phone")
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_PRIMARY_MOBILEPHONE_NOTSINGLE,
			Message: config.Localization[session.Language].Errors.Api.PrimaryMobilePhone_NotSingle})
		return
	}

	arrMobilePhones := new([]models.DtoMobilePhone)

	var updMobilePhone models.ViewMobilePhone
	var curMobilePhone models.DtoMobilePhone

	for _, updMobilePhone = range updatephones {
		found := false
		for _, curMobilePhone = range *user.MobilePhones {
			if updMobilePhone.Phone == curMobilePhone.Phone {
				found = true
				break
			}
		}
		classifier, err := helpers.CheckClassifier(updMobilePhone.Classifier_ID, r, classifierrepository, session.Language)
		if err != nil {
			return
		}

		if !found {
			var phoneExists bool
			phoneExists, err = helpers.CheckMobilePhoneAvailability(updMobilePhone.Phone, session.Language, r, mobilephonerepository)
			if err != nil {
				return
			}

			*arrMobilePhones = append(*arrMobilePhones, models.DtoMobilePhone{
				Phone:         updMobilePhone.Phone,
				UserID:        user.ID,
				Classifier_ID: classifier.ID,
				Created:       time.Now(),
				Primary:       updMobilePhone.Primary,
				Confirmed:     false,
				//				Subscription:  updMobilePhone.Subscription,
				Code:     "",
				Language: strings.ToLower(updMobilePhone.Language),
				Exists:   phoneExists,
			})
		} else {
			curMobilePhone.Primary = updMobilePhone.Primary
			//			curMobilePhone.Subscription = updMobilePhone.Subscription
			curMobilePhone.Language = strings.ToLower(updMobilePhone.Language)
			curMobilePhone.Classifier_ID = classifier.ID

			*arrMobilePhones = append(*arrMobilePhones, curMobilePhone)
		}

		if updMobilePhone.Primary {
			if !found || !curMobilePhone.Confirmed {
				log.Error("Primary mobile phone is not confirmed %v", updMobilePhone.Phone)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_PRIMARY_MOBILEPHONE_NOTCONFIRMED,
					Message: config.Localization[session.Language].Errors.Api.PrimaryMobilePhone_NotConfirmed})
				return
			}
		}
	}

	user.MobilePhones = arrMobilePhones
	err = userrepository.UpdateMobilePhones(user, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	phones := new([]models.ViewApiMobilePhone)
	for _, phone := range *user.MobilePhones {
		*phones = append(*phones, *models.NewViewApiMobilePhone(phone.Phone, phone.Primary, phone.Confirmed, /*phone.Subscription,*/
			phone.Language, phone.Classifier_ID))
	}

	helpers.RenderJSONArray(phones, len(*phones), w, r)
}

// get /api/v1.0/unit/
func GetUserUnit(r render.Render, userrepository services.UserRepository, unitrepository services.UnitRepository, session *models.DtoSession) {
	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	unit, err := unitrepository.Get(user.UnitID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiFullUnit(unit.ID, unit.Name, unit.Created, unit.Subscribed, unit.Paid, unit.Begin_Paid, unit.End_Paid))
}

// patch /api/v1.0/unit/
func UpdateUserUnit(r render.Render, errors binding.Errors, viewunit models.ViewShortUnit, userrepository services.UserRepository,
	unitrepository services.UnitRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	user, err := userrepository.Get(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if !user.UnitAdmin {
		log.Error("User %v is not admin for unit %v", user.ID, user.UnitID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	unit, err := unitrepository.Get(user.UnitID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	unit.Name = viewunit.Name
	err = unitrepository.Update(unit)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiFullUnit(unit.ID, unit.Name, unit.Created, unit.Subscribed, unit.Paid, unit.Begin_Paid, unit.End_Paid))
}
