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

// get /api/v1.0/classification/contacts/
func GetAvailableContacts(r render.Render, classifierrepository services.ClassifierRepository, session *models.DtoSession) {
	classifiers, err := classifierrepository.GetAllAvailable()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, classifiers)
}

// get /api/v1.0/customers/services/
func GetAvailableFacilities(r render.Render, facilityrepository services.FacilityRepository, session *models.DtoSession) {
	facilities, err := facilityrepository.GetAllAvailable()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, facilities)
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
func GetProjects(request *http.Request, r render.Render, projectrepository services.ProjectRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.ProjectSearch), nil, request, r, session.Language)
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
	sorts, err = helpers.GetOrderArray(new(models.ProjectSearch), request, r, session.Language)
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

	r.JSON(http.StatusOK, projects)
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

// options /api/v1.0/projects/:prid/orders/
func GetMetaProjectOrders(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	orderrepository services.OrderRepository, session *models.DtoSession) {
	dtoproject, err := helpers.CheckProject(r, params, projectrepository, session.Language)
	if err != nil {
		return
	}

	order, err := orderrepository.GetMetaByProject(dtoproject.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, order)
}

// get /api/v1.0/projects/:prid/orders/
func GetProjectOrders(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	orderrepository services.OrderRepository, session *models.DtoSession) {
	dtoproject, err := helpers.CheckProject(r, params, projectrepository, session.Language)
	if err != nil {
		return
	}

	orders, err := orderrepository.GetByProject(dtoproject.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, orders)
}

// post /api/v1.0/projects/:prid/orders/
func CreateProjectOrder(errors binding.Errors, vieworder models.ViewShortOrder, r render.Render, params martini.Params,
	projectrepository services.ProjectRepository, orderrepository services.OrderRepository, unitrepository services.UnitRepository,
	session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	dtoproject, err := helpers.CheckProject(r, params, projectrepository, session.Language)
	if err != nil {
		return
	}
	unit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	dtoorder := new(models.DtoOrder)
	dtoorder.Project_ID = dtoproject.ID
	dtoorder.Creator_ID = session.UserID
	dtoorder.Unit_ID = unit.ID
	dtoorder.Name = vieworder.Name
	dtoorder.Step = 0
	dtoorder.Created = time.Now()

	order := new(models.ViewLongOrder)
	dtoorderstatuses := order.ToOrderStatuses(dtoorder.ID)
	err = orderrepository.Create(dtoorder, dtoorderstatuses, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongOrderFromDto(dtoorder, dtoorderstatuses))
}

// get /api/v1.0/projects/:prid/orders/:oid/
func GetProjectOrder(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	orderrepository services.OrderRepository, orderstatusrepository services.OrderStatusRepository, session *models.DtoSession) {
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
	if err != nil {
		return
	}

	dtoorderstatuses, err := orderstatusrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongOrderFromDto(dtoorder, dtoorderstatuses))
}

// put /api/v1.0/projects/:prid/orders/:oid/
func UpdateProjectOrder(errors binding.Errors, vieworder models.ViewLongOrder, r render.Render, params martini.Params,
	projectrepository services.ProjectRepository, orderrepository services.OrderRepository, unitrepository services.UnitRepository,
	facilityrepository services.FacilityRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
	if err != nil {
		return
	}

	apiorder, err := helpers.UpdateOrder(dtoorder, &vieworder, r, params, orderrepository, unitrepository, facilityrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apiorder)
}

// delete /api/v1.0/projects/:prid/orders/:oid/
func DeleteProjectOrder(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	orderrepository services.OrderRepository, orderstatusrepository services.OrderStatusRepository,
	session *models.DtoSession) {
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
	if err != nil {
		return
	}
	confirmed, err := orderrepository.IsConfirmed(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if confirmed {
		log.Error("Can't delete confirmed order %v", dtoorder.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	orderstatus := models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_DEL, true, "", time.Now())
	err = orderstatusrepository.Save(orderstatus)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
