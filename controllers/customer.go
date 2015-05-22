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

// get /api/v1.0/customers/services/
func GetAvailableFacilities(w http.ResponseWriter, r render.Render, facilityrepository services.FacilityRepository, session *models.DtoSession) {
	facilities, err := facilityrepository.GetAllAvailable()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(facilities, len(*facilities), w, r)
}

// options /api/v1.0/projects/
func GetMetaProjects(r render.Render, projectrepository services.ProjectRepository, session *models.DtoSession) {
	project, err := projectrepository.GetMeta(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, project)
}

// get /api/v1.0/projects/
func GetAllProjects(w http.ResponseWriter, request *http.Request, r render.Render, projectrepository services.ProjectRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.ProjectLongSearch), nil, request, r, session.Language)
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
	sorts, err = helpers.GetOrderArray(new(models.ProjectLongSearch), request, r, session.Language)
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

	projects, err := projectrepository.GetByUser(session.UserID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(projects, len(*projects), w, r)
}

// get /api/v1.0/projects/onthego/
func GetActiveProjects(request *http.Request, w http.ResponseWriter, r render.Render, projectrepository services.ProjectRepository, session *models.DtoSession) {
	helpers.GetProjects(session.UserID, true, request, w, r, projectrepository, session.Language)
}

// get /api/v1.0/projects/wasarchived/
func GetArchiveProjects(request *http.Request, w http.ResponseWriter, r render.Render, projectrepository services.ProjectRepository, session *models.DtoSession) {
	helpers.GetProjects(session.UserID, false, request, w, r, projectrepository, session.Language)
}

// post /api/v1.0/projects/
func CreateProject(errors binding.Errors, viewproject models.ViewProject, r render.Render, projectrepository services.ProjectRepository,
	unitrepository services.UnitRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}

	unit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	dtoproject := new(models.DtoProject)
	dtoproject.Unit_ID = unit.ID
	dtoproject.Name = viewproject.Name
	dtoproject.Active = true
	dtoproject.Created = time.Now()

	err = projectrepository.Create(dtoproject)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongProject(dtoproject.ID, dtoproject.Name, !dtoproject.Active, dtoproject.Created))
}

// get /api/v1.0/projects/:prid/
func GetProject(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	session *models.DtoSession) {
	dtoproject, err := helpers.CheckProject(r, params, projectrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongProject(dtoproject.ID, dtoproject.Name, !dtoproject.Active, dtoproject.Created))
}

// put /api/v1.0/projects/:prid/
func UpdateProject(errors binding.Errors, viewproject models.ViewUpdateProject, r render.Render, params martini.Params,
	projectrepository services.ProjectRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	dtoproject, err := helpers.CheckProject(r, params, projectrepository, session.Language)
	if err != nil {
		return
	}

	dtoproject.Name = viewproject.Name
	dtoproject.Active = !viewproject.Archive
	err = projectrepository.Update(dtoproject)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongProject(dtoproject.ID, dtoproject.Name, !dtoproject.Active, dtoproject.Created))
}

// delete /api/v1.0/projects/:prid/
func DeleteProject(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	session *models.DtoSession) {
	dtoproject, err := helpers.CheckProject(r, params, projectrepository, session.Language)
	if err != nil {
		return
	}
	if !dtoproject.Active {
		log.Error("Project is not active %v", dtoproject.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	has, err := projectrepository.HasNotCompletedOrder(dtoproject.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if has {
		log.Error("Project has not completed order")
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	has, err = projectrepository.HasNotPaidOrder(dtoproject.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if has {
		log.Error("Project has not paid order")
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	err = projectrepository.Deactivate(dtoproject)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// options /api/v1.0/smsfrom/
func GetMetaSMSSenders(r render.Render, smssenderrepository services.SMSSenderRepository, session *models.DtoSession) {
	smssender, err := smssenderrepository.GetMeta(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, smssender)
}

// get /api/v1.0/smsfrom/
func GetSMSSenders(w http.ResponseWriter, request *http.Request, r render.Render, smssenderrepository services.SMSSenderRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.SMSSenderSearch), nil, request, r, session.Language)
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
	sorts, err = helpers.GetOrderArray(new(models.SMSSenderSearch), request, r, session.Language)
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

	smssenders, err := smssenderrepository.GetByUser(session.UserID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(smssenders, len(*smssenders), w, r)
}

// get /api/v1.0/smsfrom/:frmid/
func GetSMSSender(r render.Render, params martini.Params, smssenderrepository services.SMSSenderRepository,
	session *models.DtoSession) {
	dtosmssender, err := helpers.CheckSMSSender(r, params, smssenderrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiMiddleSMSSender(dtosmssender.ID, dtosmssender.Name, dtosmssender.Registered, !dtosmssender.Active))
}

// post /api/v1.0/smsfrom/
func CreateSMSSender(errors binding.Errors, viewsmssender models.ViewSMSSender, r render.Render, params martini.Params,
	smssenderrepository services.SMSSenderRepository, unitrepository services.UnitRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}

	found, err := smssenderrepository.Exists(viewsmssender.Name)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if found {
		r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_SMSSENDER_INUSE,
			Message: config.Localization[session.Language].Errors.Api.SMSSender_InUse})
		return
	}

	unit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	dtosmssender := new(models.DtoSMSSender)
	dtosmssender.Unit_ID = unit.ID
	dtosmssender.Name = viewsmssender.Name
	dtosmssender.Created = time.Now()
	dtosmssender.Registered = false
	dtosmssender.Withdraw = false
	dtosmssender.Withdrawn = false
	dtosmssender.Active = true

	err = smssenderrepository.Create(dtosmssender)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortSMSSender(dtosmssender.ID, dtosmssender.Name, dtosmssender.Registered))
}

// patch /api/v1.0/smsfrom/:frmid/
func UpdateSMSSender(errors binding.Errors, viewsmssender models.ViewSMSSender, r render.Render, params martini.Params,
	smssenderrepository services.SMSSenderRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	dtosmssender, err := helpers.CheckSMSSender(r, params, smssenderrepository, session.Language)
	if err != nil {
		return
	}

	found, err := smssenderrepository.Exists(viewsmssender.Name)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if found && dtosmssender.Name != viewsmssender.Name {
		r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_SMSSENDER_INUSE,
			Message: config.Localization[session.Language].Errors.Api.SMSSender_InUse})
		return
	}

	dtosmssender.Name = viewsmssender.Name
	err = smssenderrepository.Update(dtosmssender)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortSMSSender(dtosmssender.ID, dtosmssender.Name, dtosmssender.Registered))
}

// delete //api/v1.0/smsfrom/:frmid/
func DeleteSMSSender(r render.Render, params martini.Params, smssenderrepository services.SMSSenderRepository,
	session *models.DtoSession) {
	dtosmssender, err := helpers.CheckSMSSender(r, params, smssenderrepository, session.Language)
	if err != nil {
		return
	}
	_, err = helpers.IsSMSSenderActive(dtosmssender.ID, r, smssenderrepository, session.Language)
	if err != nil {
		return
	}

	dtosmssender.Withdraw = true
	err = smssenderrepository.Update(dtosmssender)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
