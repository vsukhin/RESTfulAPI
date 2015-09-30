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
	PARAM_NAME_UNIT_USER_ID = "uid"
)

func CheckUnitUser(admin_id int64, r render.Render, params martini.Params, userrepository services.UserRepository,
	language string) (dtouser *models.DtoUser, err error) {
	user_id, err := CheckParameterInt(r, params[PARAM_NAME_UNIT_USER_ID], language)
	if err != nil {
		return nil, err
	}

	dtouser, err = userrepository.Get(user_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	dtoadmin, err := userrepository.Get(admin_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	if dtouser.UnitID != dtoadmin.UnitID {
		log.Error("User unit %v and admin unit %v don't match", dtouser.UnitID, dtoadmin.UnitID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Not matched units")
	}

	return dtouser, nil
}

func CheckPrimaryUnitUserEmail(emails []models.ViewEmail, language string, r render.Render) (count int, err error) {
	count = 0
	for _, checkEmail := range emails {
		if checkEmail.Primary {
			count++
		}
	}
	if count > 1 {
		log.Error("Only one primary email is allowed")
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return 0, errors.New("Wrong primary emails amount")
	}

	return count, nil
}

func CheckPrimaryUnitUserMobilePhone(mobilephones []models.ViewMobilePhone, language string, r render.Render) (count int, err error) {
	count = 0
	for _, checkMobilePhone := range mobilephones {
		if checkMobilePhone.Primary {
			count++
		}
	}
	if count > 1 {
		log.Error("Only one primary mobile phone is allowed")
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return 0, errors.New("Wrong primary mobile phones amount")
	}

	return count, nil
}
