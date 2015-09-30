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
	PARAM_NAME_SMSSENDER_ID = "hdrid"
)

func CheckSMSSender(r render.Render, params martini.Params, smssenderrepository services.SMSSenderRepository,
	language string) (dtosmssender *models.DtoSMSSender, err error) {
	smssender_id, err := CheckParameterInt(r, params[PARAM_NAME_SMSSENDER_ID], language)
	if err != nil {
		return nil, err
	}

	dtosmssender, err = IsSMSSenderActive(smssender_id, r, smssenderrepository, language)
	if err != nil {
		return nil, err
	}

	return dtosmssender, nil
}

func IsSMSSenderActive(smssender_id int64, r render.Render, smssenderrepository services.SMSSenderRepository,
	language string) (dtosmssender *models.DtoSMSSender, err error) {
	dtosmssender, err = smssenderrepository.Get(smssender_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtosmssender.Active {
		log.Error("SMS sender is not active %v", dtosmssender.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Wrong SMS sender")
	}

	return dtosmssender, nil
}

func IsSMSSenderAccessible(smssender_id int64, user_id int64, r render.Render, smssenderrepository services.SMSSenderRepository,
	language string) (err error) {
	allowed, err := smssenderrepository.CheckCustomerAccess(user_id, smssender_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}
	if !allowed {
		log.Error("SMS sender %v is not accessible for customer %v", smssender_id, user_id)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return errors.New("Not accessible SMS sender")
	}

	return nil
}
