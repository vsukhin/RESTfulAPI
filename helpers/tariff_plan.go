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

func CheckTariffPlan(tariff_plan_id int, r render.Render, tariffplanrepository services.TariffPlanRepository,
	language string) (dtotariffplan *models.DtoTariffPlan, err error) {
	dtotariffplan, err = tariffplanrepository.Get(tariff_plan_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtotariffplan.Active {
		log.Error("Tariff plan is not active %v", dtotariffplan.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Tariff plan not active")
	}

	return dtotariffplan, nil
}
