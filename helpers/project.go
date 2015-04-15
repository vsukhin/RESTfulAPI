package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
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
		return nil, nil, errors.New("Non matched project and order")
	}

	return dtoproject, dtoorder, nil
}
