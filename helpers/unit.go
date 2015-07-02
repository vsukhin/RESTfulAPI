package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

const (
	PARAM_NAME_UNIT_ID = "unitId"
)

func CheckUnit(r render.Render, params martini.Params, unitrepository services.UnitRepository,
	language string) (dtounit *models.DtoUnit, err error) {
	unit_id, err := CheckParameterInt(r, params[PARAM_NAME_UNIT_ID], language)
	if err != nil {
		return nil, err
	}

	dtounit, err = unitrepository.Get(unit_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return dtounit, nil
}

func CheckUnitValidity(unitid int64, language string, r render.Render, unitrepository services.UnitRepository) (err error) {
	_, err = unitrepository.Get(unitid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}

	return nil
}
