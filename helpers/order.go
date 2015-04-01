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
	PARAM_NAME_ORDER = "oid"
)

func CheckOrder(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	language string) (dtoorder *models.DtoOrder, err error) {
	orderid, err := CheckParameterInt(r, params[PARAM_NAME_ORDER], language)
	if err != nil {
		return
	}
	dtoorder, err = orderrepository.Get(orderid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	return dtoorder, nil
}
