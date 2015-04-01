package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
	"time"
	"types"
)

// get /api/v1.0/administration/users/
func GetUsers(request *http.Request, r render.Render, params martini.Params,
	userrepository services.UserRepository, session *models.DtoSession) {
	var err error
	query := ""

	var filters *[]models.FilterExp
	filters, err = helpers.GetFilterArray(new(models.UserSearch), nil, request, r, session.Language)
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
	sorts, err = helpers.GetOrderArray(new(models.UserSearch), request, r, session.Language)
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

	users, err := userrepository.GetAll(query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, users)
}

// post /api/v1.0/administration/users/
func CreateUser(errors binding.Errors, user models.ViewApiUserFull, request *http.Request, r render.Render, params martini.Params,
	userrepository services.UserRepository, emailrepository services.EmailRepository, sessionrepository services.SessionRepository,
	unitrepository services.UnitRepository, templaterepository services.TemplateRepository,
	grouprepository services.GroupRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}

	dtouser := new(models.DtoUser)
	dtouser.Created = time.Now()
	dtouser.LastLogin = dtouser.Created
	dtouser.Active = user.Active
	dtouser.Confirmed = user.Confirmed
	dtouser.Name = user.Name
	dtouser.Language = strings.ToLower(user.Language)
	dtouser.Code = ""
	dtouser.Password = ""

	if helpers.CheckUserRoles(user.Roles, session.Language, r, grouprepository) != nil {
		return
	}
	dtouser.Roles = user.Roles

	if user.Creator_ID != 0 {
		if helpers.CheckUser(user.Creator_ID, session.Language, r, userrepository) != nil {
			return
		}
	}
	dtouser.Creator_ID = user.Creator_ID

	if user.Unit_ID != 0 {
		if helpers.CheckUnit(user.Unit_ID, session.Language, r, unitrepository) != nil {
			return
		}
	}
	dtouser.UnitID = user.Unit_ID
	dtouser.UnitAdmin = user.UnitAdmin

	if helpers.CheckPrimaryEmail(&user, session.Language, r) != nil {
		return
	}
	dtouser.Emails = new([]models.DtoEmail)

	for _, apiEmail := range user.Emails {
		apiEmail.Email = strings.ToLower(apiEmail.Email)
		emailExists, err := helpers.CheckEmailAvailability(apiEmail.Email, session.Language, r, emailrepository)
		if err != nil {
			return
		}
		code := ""

		if !apiEmail.Confirmed {
			code, err = sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
			if apiEmail.Primary {
				dtouser.Code = code
				dtouser.Password = code
			}
		}
		*dtouser.Emails = append(*dtouser.Emails, models.DtoEmail{
			Email:        apiEmail.Email,
			Created:      dtouser.LastLogin,
			Primary:      apiEmail.Primary,
			Confirmed:    apiEmail.Confirmed,
			Subscription: apiEmail.Subscription,
			Code:         code,
			Language:     apiEmail.Language,
			Exists:       emailExists,
		})
	}

	err := userrepository.Create(dtouser, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	if helpers.SendConfirmations(dtouser, session, request, r, emailrepository, templaterepository) != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiUserTiny(dtouser.ID))
}

// put /api/v1.0/administration/users/:userId/
func UpdateUser(errors binding.Errors, user models.ViewApiUserFull, request *http.Request, r render.Render, params martini.Params,
	userrepository services.UserRepository, emailrepository services.EmailRepository, sessionrepository services.SessionRepository,
	unitrepository services.UnitRepository, templaterepository services.TemplateRepository, grouprepository services.GroupRepository,
	session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	userid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_USER_ID], session.Language)
	if err != nil {
		return
	}

	dtouser, err := userrepository.Get(userid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	dtouser.Active = user.Active
	dtouser.Confirmed = user.Confirmed
	dtouser.Name = user.Name
	dtouser.Language = strings.ToLower(user.Language)

	if helpers.CheckUserRoles(user.Roles, session.Language, r, grouprepository) != nil {
		return
	}
	dtouser.Roles = user.Roles

	if user.Creator_ID != 0 {
		if helpers.CheckUser(user.Creator_ID, session.Language, r, userrepository) != nil {
			return
		}
	}
	dtouser.Creator_ID = user.Creator_ID

	if user.Unit_ID != 0 {
		if helpers.CheckUnit(user.Unit_ID, session.Language, r, unitrepository) != nil {
			return
		}
	}
	dtouser.UnitID = user.Unit_ID
	dtouser.UnitAdmin = user.UnitAdmin

	if helpers.CheckPrimaryEmail(&user, session.Language, r) != nil {
		return
	}

	arrInEmails := user.Emails
	arrOutEmails := new([]models.DtoEmail)

	var updEmail models.ViewApiEmail
	var curEmail models.DtoEmail

	for _, updEmail = range arrInEmails {
		updEmail.Email = strings.ToLower(updEmail.Email)
		found := false
		code := ""
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
				Email:        updEmail.Email,
				UserID:       dtouser.ID,
				Created:      time.Now(),
				Primary:      updEmail.Primary,
				Confirmed:    updEmail.Confirmed,
				Subscription: updEmail.Subscription,
				Code:         code,
				Language:     updEmail.Language,
				Exists:       emailExists,
			})
		} else {
			curEmail.Primary = updEmail.Primary
			curEmail.Confirmed = updEmail.Confirmed
			curEmail.Subscription = updEmail.Subscription
			curEmail.Code = code
			curEmail.Language = updEmail.Language

			*arrOutEmails = append(*arrOutEmails, curEmail)
		}

		if !updEmail.Confirmed {
			code, err = sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}

			(*arrOutEmails)[len(*arrOutEmails)-1].Code = code

			if updEmail.Primary {
				dtouser.Code = code
				dtouser.Password = code
			}
		}
	}

	dtouser.Emails = arrOutEmails
	err = userrepository.Update(dtouser, false, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	if helpers.SendConfirmations(dtouser, session, request, r, emailrepository, templaterepository) != nil {
		return
	}

	user.Unit_ID = dtouser.UnitID
	r.JSON(http.StatusOK, user)
}

// delete /api/v1.0/administration/users/:userId/
func DeleteUser(r render.Render, params martini.Params, userrepository services.UserRepository, session *models.DtoSession) {
	userid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_USER_ID], session.Language)
	if err != nil {
		return
	}

	user, err := userrepository.Get(userid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	err = userrepository.Delete(user.ID, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// options /api/v1.0/administration/users/
func GetUserMetaData(r render.Render, userrepository services.UserRepository, session *models.DtoSession) {
	usermeta, err := userrepository.GetMeta()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, usermeta)
}

// get /api/v1.0/administration/users/:userId/
func GetUserFullInfo(r render.Render, params martini.Params, userrepository services.UserRepository, session *models.DtoSession) {
	userid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_USER_ID], session.Language)
	if err != nil {
		return
	}

	user, err := userrepository.Get(userid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	userfull := models.NewViewApiUserFull(user.Creator_ID, user.UnitID, user.UnitAdmin, user.Active, user.Confirmed,
		user.Name, user.Language, user.Roles, *new([]models.ViewApiEmail))
	for _, email := range *user.Emails {
		userfull.Emails = append(userfull.Emails, *models.NewViewApiEmail(email.Email, email.Primary, email.Confirmed, email.Subscription, email.Language))
	}

	r.JSON(http.StatusOK, userfull)
}

// get /api/v1.0/administration/groups/
func GetGroupsInfo(r render.Render, grouprepository services.GroupRepository, session *models.DtoSession) {
	groups, err := grouprepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, groups)
}
