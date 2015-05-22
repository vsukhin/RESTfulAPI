package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

func CheckMobileOperator(mobileoperator_id int, r render.Render, mobileoperatorrepository services.MobileOperatorRepository,
	language string) (dtomobileoperator *models.DtoMobileOperator, err error) {
	dtomobileoperator, err = mobileoperatorrepository.Get(mobileoperator_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtomobileoperator.Active {
		log.Error("Mobile operator is not active %v", dtomobileoperator.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Not active mobile operator")
	}

	return dtomobileoperator, nil
}
