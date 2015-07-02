package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
	"types"
)

const (
	PARAM_NAME_PROJECT_ID = "prid"
)

func CheckProject(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	language string) (dtoproject *models.DtoProject, err error) {
	project_id, err := CheckParameterInt(r, params[PARAM_NAME_PROJECT_ID], language)
	if err != nil {
		return nil, err
	}

	dtoproject, err = projectrepository.Get(project_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return dtoproject, nil
}

func CheckProjectValidity(project_id int64, r render.Render, projectrepository services.ProjectRepository,
	language string) (dtoproject *models.DtoProject, err error) {
	dtoproject, err = projectrepository.Get(project_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return dtoproject, nil
}

func CheckProjectAccess(project_id int64, user_id int64, r render.Render, projectrepository services.ProjectRepository,
	language string) (err error) {
	allowed, err := projectrepository.CheckCustomerAccess(user_id, project_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}
	if !allowed {
		log.Error("Project %v is not accessible for user %v", project_id, user_id)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return errors.New("Project not accessible")
	}

	return nil
}

func CheckProjectOrder(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	orderrepository services.OrderRepository, language string) (dtoproject *models.DtoProject, dtoorder *models.DtoOrder, err error) {
	dtoproject, err = CheckProject(r, params, projectrepository, language)
	if err != nil {
		return nil, nil, err
	}
	dtoorder, err = CheckOrder(r, params, orderrepository, language)
	if err != nil {
		return nil, nil, err
	}
	if dtoproject.ID != dtoorder.Project_ID {
		log.Error("Order %v doesn't belong to project %v", dtoorder.ID, dtoproject.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, errors.New("Not matched project and order")
	}

	return dtoproject, dtoorder, nil
}

func GetProjects(userid int64, active bool, request *http.Request, w http.ResponseWriter, r render.Render, projectrepository services.ProjectRepository,
	language string) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := GetFilterArray(new(models.ProjectShortSearch), nil, request, r, language)
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
	sorts, err = GetOrderArray(new(models.ProjectShortSearch), request, r, language)
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
	limit, err = GetLimitQuery(request, r, language)
	if err != nil {
		return
	}
	query += limit

	projects, err := projectrepository.GetByUserWithStatus(userid, active, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return
	}

	RenderJSONArray(projects, len(*projects), w, r)
}
