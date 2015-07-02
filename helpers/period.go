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

func CheckPeriod(period_id int, r render.Render, periodrepository services.PeriodRepository,
	language string) (dtoperiod *models.DtoPeriod, err error) {
	dtoperiod, err = periodrepository.Get(period_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtoperiod.Active {
		log.Error("Period is not active %v", dtoperiod.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Not active period")
	}

	return dtoperiod, nil
}
