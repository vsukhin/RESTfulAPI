package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

func CheckComplexStatus(complexstatus_id int, r render.Render, complexstatusrepository services.ComplexStatusRepository,
	language string) (dtocomplexstatus *models.DtoComplexStatus, err error) {
	dtocomplexstatus, err = complexstatusrepository.Get(complexstatus_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return dtocomplexstatus, nil
}
