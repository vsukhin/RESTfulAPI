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
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
	"types"
)

//  В принципе сессия уже проверена и продлена, так как вызов сделан под авторизацией
// get /api/v1.0/session/:token
func KeepSession(request *http.Request, r render.Render, params martini.Params, sessionservice *services.SessionService, session *models.DtoSession) {

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

//  В принципе сессия уже проверена, так как вызов сделан под авторизацией
// get /api/v1.0/ping/:token
func Ping(request *http.Request, r render.Render, params martini.Params, sessionservice *services.SessionService, session *models.DtoSession) {

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// post /api/v1.0/user/session/
func CreateSession(errors binding.Errors, viewsession models.ViewSession, r render.Render, params martini.Params,
	userservice *services.UserService, sessionservice *services.SessionService, captchaservice *services.CaptchaService) {
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	if captchaservice.Check(viewsession.CaptchaHash, viewsession.CaptchaValue, r) != nil {
		return
	}

	var language string
	if viewsession.Language != "" {
		language = strings.ToLower(viewsession.Language)
	} else {
		language = config.Configuration.Server.DefaultLanguage
	}

	token, err := sessionservice.GenerateToken(helpers.TOKEN_LENGTH)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return
	}

	var user *models.DtoUser
	viewsession.Login = strings.ToLower(viewsession.Login)
	user, err = userservice.FindByLogin(viewsession.Login)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_LOGIN_OR_PASSWORD_WRONG,
			Message: config.Localization[language].Errors.Api.Login_Or_Password_Wrong})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(viewsession.Password)) != nil {
		log.Error("Can't compare password hashes")
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_LOGIN_OR_PASSWORD_WRONG,
			Message: config.Localization[language].Errors.Api.Login_Or_Password_Wrong})
		return
	}

	if !user.Active || !user.Confirmed {
		log.Error("User is not active or confirmed %v", user.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_USER_BLOCKED,
			Message: config.Localization[language].Errors.Api.User_Blocked})
		return
	}

	dtosession := models.NewDtoSession(token, user.ID, user.Roles, time.Now(), language)

	err = sessionservice.Create(dtosession, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return
	}

	user.LastLogin = dtosession.LastActivity
	err = userservice.Update(user, true, false)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiSession(dtosession.LastActivity.Add(config.Configuration.Server.SessionTimeout), dtosession.AccessToken))
}

// delete /api/v1.0/session/:token
func DeleteSession(r render.Render, params martini.Params, sessionservice *services.SessionService, session *models.DtoSession) {
	err := sessionservice.Delete(session.AccessToken, true)
	if err != nil {
		middlewares.GeneratingSessionErrorResponse(r, session.AccessToken)
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
