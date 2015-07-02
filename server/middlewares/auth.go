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

func IsAdmin(roles []models.UserRole) bool {
	return IsUserRoleAllowed(roles, []models.UserRole{models.USER_ROLE_ADMINISTRATOR, models.USER_ROLE_DEVELOPER})
}

func IsSupplier(roles []models.UserRole) bool {
	return IsUserRoleAllowed(roles, []models.UserRole{models.USER_ROLE_SUPPLIER})
}

func IsCustomer(roles []models.UserRole) bool {
	return IsUserRoleAllowed(roles, []models.UserRole{models.USER_ROLE_CUSTOMER})
}

func IsUser(roles []models.UserRole) bool {
	return IsUserRoleAllowed(roles, []models.UserRole{models.USER_ROLE_SUPPLIER, models.USER_ROLE_CUSTOMER})
}

func RequireAdminRights(r render.Render, session *models.DtoSession) {
	if !IsAdmin(session.Roles) {
		log.Error("There is no required role for user %v", session.UserID)
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
}

func RequireSupplierRights(r render.Render, session *models.DtoSession) {
	if !(IsAdmin(session.Roles) || IsSupplier(session.Roles)) {
		log.Error("There is no required role for user %v", session.UserID)
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
}

func RequireCustomerRights(r render.Render, session *models.DtoSession) {
	if !(IsAdmin(session.Roles) || IsCustomer(session.Roles)) {
		log.Error("There is no required role for user %v", session.UserID)
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
}

func RequireUserRights(r render.Render, session *models.DtoSession) {
	if !(IsAdmin(session.Roles) || IsUser(session.Roles)) {
		log.Error("There is no required role for user %v", session.UserID)
		r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
			Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
		return
	}
}

func RequireTableRights(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	session *models.DtoSession) {
	if !IsAdmin(session.Roles) {
		var allowed bool = false
		if IsUser(session.Roles) {
			param := params[helpers.PARAM_NAME_TABLE_ID]
			if param == "" {
				param = params[helpers.PARAM_NAME_TEMPORABLE_TABLE_ID]
			}
			if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
				tableid, err := strconv.ParseInt(param, 0, 64)
				if err == nil {
					allowed, err = customertablerepository.CheckUserAccess(session.UserID, tableid)
					if err == nil {
						if !allowed {
							log.Error("Table %v is not accessible for user %v", tableid, session.UserID)
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

func RequireEditableTable(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	pricepropertiesrepository services.PricePropertiesRepository, session *models.DtoSession) {
	if !IsAdmin(session.Roles) {
		var allowed bool = false
		param := params[helpers.PARAM_NAME_TABLE_ID]
		if param == "" {
			param = params[helpers.PARAM_NAME_TEMPORABLE_TABLE_ID]
		}
		if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
			tableid, err := strconv.ParseInt(param, 0, 64)
			if err == nil {
				dtocustomertable, err := customertablerepository.Get(tableid)
				if err == nil {
					switch dtocustomertable.TypeID {
					case models.TABLE_TYPE_READONLY:
						log.Error("Can't edit readonly table %v", dtocustomertable.ID)
					case models.TABLE_TYPE_HIDDEN_READONLY:
						log.Error("Can't edit readonly table %v", dtocustomertable.ID)
					case models.TABLE_TYPE_PRICE:
						found, err := pricepropertiesrepository.Exists(dtocustomertable.ID)
						if err == nil {
							if found {
								priceproperties, err := pricepropertiesrepository.Get(dtocustomertable.ID)
								if err == nil {
									if priceproperties.Published {
										log.Error("Can't edit published price list %v", dtocustomertable.ID)
									} else {
										allowed = true
									}
								}
							} else {
								allowed = true
							}
						}
					default:
						allowed = true
					}
				}
			} else {
				log.Error("Can't convert parameter %v with value %v", err, param)
			}
		} else {
			log.Error("Parameter value is wrong %v", param)
		}
		if !allowed {
			r.JSON(http.StatusForbidden, types.Error{Code: types.TYPE_ERROR_METHOD_NOTALLOWED,
				Message: config.Localization[session.Language].Errors.Api.Method_NotAllowed})
			return
		}
	}
}

func RequireOrderRights(r render.Render, params martini.Params, orderrepository services.OrderRepository, session *models.DtoSession) {
	if !IsAdmin(session.Roles) {
		var allowed bool = false
		if IsSupplier(session.Roles) {
			param := params[helpers.PARAM_NAME_ORDER_ID]
			if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
				orderid, err := strconv.ParseInt(param, 0, 64)
				if err == nil {
					allowed, err = orderrepository.CheckSupplierAccess(session.UserID, orderid)
					if err == nil {
						if !allowed {
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
	if !IsAdmin(session.Roles) {
		var allowed bool = false
		if IsUser(session.Roles) {
			param := params[helpers.PARAM_NAME_ORDER_ID]
			if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
				orderid, err := strconv.ParseInt(param, 0, 64)
				if err == nil {
					allowed, err = orderrepository.CheckUserAccess(session.UserID, orderid)
					if err == nil {
						if !allowed {
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
	if !IsAdmin(session.Roles) {
		var allowed bool = false
		if IsCustomer(session.Roles) {
			param := params[helpers.PARAM_NAME_PROJECT_ID]
			if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
				projectid, err := strconv.ParseInt(param, 0, 64)
				if err == nil {
					allowed, err = projectrepository.CheckCustomerAccess(session.UserID, projectid)
					if err == nil {
						if !allowed {
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

func RequireSMSSenderRights(r render.Render, params martini.Params, smssenderrepository services.SMSSenderRepository, session *models.DtoSession) {
	if !IsAdmin(session.Roles) {
		var allowed bool = false
		if IsCustomer(session.Roles) {
			param := params[helpers.PARAM_NAME_SMSSENDER_ID]
			if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
				smssender_id, err := strconv.ParseInt(param, 0, 64)
				if err == nil {
					allowed, err = smssenderrepository.CheckCustomerAccess(session.UserID, smssender_id)
					if err == nil {
						if !allowed {
							log.Error("SMSFrom %v is not accessible for customer %v", smssender_id, session.UserID)
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

func RequireCompanyRights(r render.Render, params martini.Params, companyrepository services.CompanyRepository, session *models.DtoSession) {
	if !IsAdmin(session.Roles) {
		var allowed bool = false
		if IsUser(session.Roles) {
			param := params[helpers.PARAM_NAME_COMPANY_ID]
			if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
				companyid, err := strconv.ParseInt(param, 0, 64)
				if err == nil {
					allowed, err = companyrepository.CheckUserAccess(session.UserID, companyid)
					if err == nil {
						if !allowed {
							log.Error("Company %v is not accessible for user %v", companyid, session.UserID)
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

func RequireInvoiceRights(r render.Render, params martini.Params, invoicerepository services.InvoiceRepository, session *models.DtoSession) {
	if !IsAdmin(session.Roles) {
		var allowed bool = false
		if IsCustomer(session.Roles) {
			param := params[helpers.PARAM_NAME_INVOICE_ID]
			if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
				invoiceid, err := strconv.ParseInt(param, 0, 64)
				if err == nil {
					allowed, err = invoicerepository.CheckUserAccess(session.UserID, invoiceid)
					if err == nil {
						if !allowed {
							log.Error("Invoice %v is not accessible for user %v", invoiceid, session.UserID)
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

func RequireUnitRights(r render.Render, userrepository services.UserRepository, session *models.DtoSession) {
	if !IsAdmin(session.Roles) {
		var allowed bool = false
		if IsUser(session.Roles) {
			var err error
			allowed, err = userrepository.CheckUnitAccess(session.UserID)
			if err == nil {
				if !allowed {
					log.Error("Unit administration is not accessible for user %v", session.UserID)
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

func RequireReportAccessRights(r render.Render, userrepository services.UserRepository, session *models.DtoSession) {
	if !IsAdmin(session.Roles) {
		var allowed bool = false
		if IsCustomer(session.Roles) {
			var err error
			allowed, err = userrepository.CheckReportAccess(session.UserID)
			if err == nil {
				if !allowed {
					log.Error("Reports are not accessible for user %v", session.UserID)
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

func RequireReportRights(r render.Render, params martini.Params, userrepository services.UserRepository,
	reportrepository services.ReportRepository, session *models.DtoSession) {
	if !IsAdmin(session.Roles) {
		var allowed bool = false
		if IsCustomer(session.Roles) {
			param := params[helpers.PARAM_NAME_REPORT_ID]
			if param != "" && len(param) <= helpers.PARAM_LENGTH_MAX {
				reportid, err := strconv.ParseInt(param, 0, 64)
				if err == nil {
					allowed, err = userrepository.CheckReportAccess(session.UserID)
					if err == nil {
						if allowed {
							allowed, err = reportrepository.CheckCustomerAccess(session.UserID, reportid)
							if err == nil {
								if !allowed {
									log.Error("Report %v is not accessible for user %v", reportid, session.UserID)
								}
							}
						} else {
							log.Error("Reports are not accessible for user %v", session.UserID)
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
