package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/server/middlewares"
	"application/services"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"image/jpeg"
	"net/http"
	"time"
	"types"
)

const (
	CAPTCHA_LENGTH  = 6
	CAPTCHA_WIDTH   = 240
	CAPTCHA_HEIGHT  = 80
	CAPTCHA_QUALITY = 10
)

// get /api/v1.0/captcha/native/
func GetCaptcha(r render.Render, captchaservice *services.CaptchaService, sessionservice *services.SessionService) {
	token, err := sessionservice.GenerateToken(helpers.TOKEN_LENGTH)
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

	err = captchaservice.Create(dtocaptcha)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}
	apicaptcha := models.NewApiCaptcha(dtocaptcha.Hash, base64.StdEncoding.EncodeToString(dtocaptcha.Image))

	r.JSON(http.StatusOK, apicaptcha)
}

// post /api/v1.0/emails/confirm/
func ConfirmEmail(errors binding.Errors, confirm models.EmailConfirm, request *http.Request, r render.Render, emailservice *services.EmailService,
	sessionservice *services.SessionService, userservice *services.UserService, templateservice *services.TemplateService) {
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}

	email, err := emailservice.FindByCode(confirm.ConfirmationToken)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_CONFIRMATION_CODE_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Confirmation_Code_Wrong})
		return
	}

	var user *models.DtoUser
	user, err = userservice.Get(email.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_CONFIRMATION_CODE_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Confirmation_Code_Wrong})
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

			var token string
			token, err = sessionservice.GenerateToken(helpers.TOKEN_LENGTH)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
				return
			}
			user.Confirmed = true
			user.Code = token

			err = userservice.Update(user, false, true)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
				return
			}

			for _, confEmail := range *user.Emails {
				if confEmail.Confirmed {
					if helpers.SendPassword(user.Language, &confEmail, user, request, r, emailservice, templateservice) != nil {
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
		err = emailservice.Update(email)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[user.Language].Errors.Api.Data_Wrong})
			return
		}
	}

	r.JSON(http.StatusAccepted, types.ResponseOK{Message: config.Localization[user.Language].Messages.OK})
}

// options /api/v1.0/services/
// options /api/v1.0/supplier/services/
func GetFacilities(r render.Render, facilityservice *services.FacilityService, session *models.DtoSession) {
	if !middlewares.IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}

	facilities, err := facilityservice.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, facilities)
}

// get /api/v1.0/
func HomePageTemplate(w http.ResponseWriter, templateservice *services.TemplateService) {
	err := templateservice.GenerateHTML("homepage", w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
