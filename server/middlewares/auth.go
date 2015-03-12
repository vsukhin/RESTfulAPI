package middlewares

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"

	"application/config"
	"application/models"
	"application/services"
	"time"
	"types"
)

func IsUserRoleAllowed(existingRoles []models.UserRole, requiredRoles []models.UserRole) bool {
	for _, requiredRole := range requiredRoles {
		for _, existingRole := range existingRoles {
			if requiredRole == existingRole {
				return true
			}
		}
	}

	return false
}

func GeneratingSessionErrorResponse(r render.Render, token string) {
	if token == "" {
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_TOKEN_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Token_Wrong})
	} else {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_TOKEN_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Token_Wrong})
	}
	return
}

func RequireSession(request *http.Request, r render.Render, sessionservice *services.SessionService,
	context martini.Context, params martini.Params, updateSession bool, takeParamFromURI bool) {
	session, token, err := sessionservice.GetAndSaveSession(request, r, params, updateSession, takeParamFromURI, false)
	if err != nil {
		GeneratingSessionErrorResponse(r, token)
	}
	context.Map(session)
}

func RequireSessionCheckWithRoute(request *http.Request, r render.Render, sessionservice *services.SessionService,
	context martini.Context, params martini.Params) {
	RequireSession(request, r, sessionservice, context, params, false, true)
}

func RequireSessionCheckWithoutRoute(request *http.Request, r render.Render, sessionservice *services.SessionService,
	context martini.Context, params martini.Params) {
	RequireSession(request, r, sessionservice, context, params, false, false)
}

func RequireSessionKeepWithRoute(request *http.Request, r render.Render, sessionservice *services.SessionService,
	context martini.Context, params martini.Params) {
	RequireSession(request, r, sessionservice, context, params, true, true)
}

func RequireSessionKeepWithoutRoute(request *http.Request, r render.Render, sessionservice *services.SessionService,
	context martini.Context, params martini.Params) {
	RequireSession(request, r, sessionservice, context, params, true, false)
}

func UtcNow() time.Time {
	return time.Now().Truncate(time.Millisecond).UTC()
}
