package helpers

import (
	"application/config"
	"application/services"
	"errors"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

func CheckFacility(facilityid int64, r render.Render, facilityrepository services.FacilityRepository,
	language string) (err error) {
	dtofacility, err := facilityrepository.Get(facilityid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}
	if !dtofacility.Active {
		log.Error("Service is not active %v", dtofacility.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return errors.New("Service not active")
	}

	return nil
}

func CheckFacilityValidity(facilityid int64, r render.Render, facilityrepository services.FacilityRepository,
	language string) (err error) {
	_, err = facilityrepository.Get(facilityid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}

	return nil
}

func CheckFacilityAlias(facilityid int64, alias string, r render.Render, facilityrepository services.FacilityRepository,
	language string) (err error) {
	dtofacility, err := facilityrepository.Get(facilityid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}
	if dtofacility.Alias != alias {
		log.Error("Order service is not macthed to the service method %v", facilityid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return errors.New("Wrong service alias")
	}

	return nil
}
