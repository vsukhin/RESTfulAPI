package middlewares

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
	"time"
	"types"
)

func RequireAdminRights(r render.Render, session *models.DtoSession) {
	if !IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		log.Error("There is no required role for user %v", session.UserID)
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
}

func RequireSupplierRights(r render.Render, session *models.DtoSession) {
	if !IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER,
		models.USER_ROLE_SUPPLIER}) {
		log.Error("There is no required role for user %v", session.UserID)
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
}

func RequireUserRights(r render.Render, session *models.DtoSession) {
	if !IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER,
		models.USER_ROLE_SUPPLIER, models.USER_ROLE_CUSTOMER}) {
		log.Error("There is no required role for user %v", session.UserID)
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
}

func RequireTableRights(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	session *models.DtoSession) {
	if !IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		var allowed bool = false
		if IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_SUPPLIER, models.USER_ROLE_CUSTOMER}) {
			names := []string{helpers.PARAM_NAME_TABLE_ID, helpers.PARAM_NAME_TEMPORABLE_TABLE_ID}
			for _, name := range names {
				param := params[name]
				if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
					tableid, err := strconv.ParseInt(param, 0, 64)
					if err == nil {
						allowed, err = customertablerepository.CheckUserAccess(session.UserID, tableid)
						if err == nil {
							if allowed {
								return
							} else {
								log.Error("Table %v is not accessible for user %v", tableid, session.UserID)
							}
						}
					} else {
						log.Error("Can't convert parameter %v with value %v", err, param)
					}
				} else {
					log.Error("Parameter value is wrong %v", param)
				}
			}
		} else {
			log.Error("There is no required role for user %v", session.UserID)
		}
		if !allowed {
			r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
				Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
			return
		}
	}
}

func RequireOrderRights(r render.Render, params martini.Params, orderrepository services.OrderRepository, session *models.DtoSession) {
	if !IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		var allowed bool = false
		if IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_SUPPLIER}) {
			param := params[helpers.PARAM_NAME_ORDER_ID]
			if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
				orderid, err := strconv.ParseInt(param, 0, 64)
				if err == nil {
					allowed, err = orderrepository.CheckSupplierAccess(session.UserID, orderid)
					if err == nil {
						if allowed {
							return
						} else {
							log.Error("Order %v is not accessible for supplier %v", orderid, session.UserID)
						}
					}
				} else {
					log.Error("Can't convert parameter %v with value %v", err, param)
				}
			} else {
				log.Error("Parameter value is wrong %v", param)
			}
		} else {
			log.Error("There is no required role for user %v", session.UserID)
		}
		if !allowed {
			r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
				Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
			return
		}
	}
}

func RequireMessageRights(r render.Render, params martini.Params, orderrepository services.OrderRepository, session *models.DtoSession) {
	if !IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		var allowed bool = false
		if IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_SUPPLIER, models.USER_ROLE_CUSTOMER}) {
			param := params[helpers.PARAM_NAME_ORDER_ID]
			if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
				orderid, err := strconv.ParseInt(param, 0, 64)
				if err == nil {
					allowed, err = orderrepository.CheckUserAccess(session.UserID, orderid)
					if err == nil {
						if allowed {
							return
						} else {
							log.Error("Order %v is not accessible for user %v", orderid, session.UserID)
						}
					}
				} else {
					log.Error("Can't convert parameter %v with value %v", err, param)
				}
			} else {
				log.Error("Parameter value is wrong %v", param)
			}
		} else {
			log.Error("There is no required role for user %v", session.UserID)
		}
		if !allowed {
			r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
				Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
			return
		}
	}
}

func RequireProjectRights(r render.Render, params martini.Params, projectrepository services.ProjectRepository, session *models.DtoSession) {
	if !IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER}) {
		var allowed bool = false
		if IsUserRoleAllowed(session.Roles, []models.UserRole{models.USER_ROLE_CUSTOMER}) {
			param := params[helpers.PARAM_NAME_PROJECT_ID]
			if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
				projectid, err := strconv.ParseInt(param, 0, 64)
				if err == nil {
					allowed, err = projectrepository.CheckCustomerAccess(session.UserID, projectid)
					if err == nil {
						if allowed {
							return
						} else {
							log.Error("Project %v is not accessible for customer %v", projectid, session.UserID)
						}
					}
				} else {
					log.Error("Can't convert parameter %v with value %v", err, param)
				}
			} else {
				log.Error("Parameter value is wrong %v", param)
			}
		} else {
			log.Error("There is no required role for user %v", session.UserID)
		}
		if !allowed {
			r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
				Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
			return
		}
	}
}

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

func RequireSession(request *http.Request, r render.Render, sessionrepository services.SessionRepository,
	context martini.Context, params martini.Params, updateSession bool, takeParamFromURI bool) {
	session, token, err := sessionrepository.GetAndSaveSession(request, r, params, updateSession, takeParamFromURI, false)
	if err != nil {
		GeneratingSessionErrorResponse(r, token)
	}
	context.Map(session)
}

func RequireSessionCheckWithRoute(request *http.Request, r render.Render, sessionrepository services.SessionRepository,
	context martini.Context, params martini.Params) {
	RequireSession(request, r, sessionrepository, context, params, false, true)
}

func RequireSessionCheckWithoutRoute(request *http.Request, r render.Render, sessionrepository services.SessionRepository,
	context martini.Context, params martini.Params) {
	RequireSession(request, r, sessionrepository, context, params, false, false)
}

func RequireSessionKeepWithRoute(request *http.Request, r render.Render, sessionrepository services.SessionRepository,
	context martini.Context, params martini.Params) {
	RequireSession(request, r, sessionrepository, context, params, true, true)
}

func RequireSessionKeepWithoutRoute(request *http.Request, r render.Render, sessionrepository services.SessionRepository,
	context martini.Context, params martini.Params) {
	RequireSession(request, r, sessionrepository, context, params, true, false)
}

func UtcNow() time.Time {
	return time.Now().Truncate(time.Millisecond).UTC()
}
