/* Helpers package provides supporting methods for controller functions */

package helpers

import (
	"application/config"
	"application/models"
	"errors"
	"fmt"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	logging "github.com/op/go-logging"
	"net/http"
	"strconv"
	"types"
)

const (
	PARAM_LENGTH_MAX = 255
	TOKEN_LENGTH     = 64
)

var (
	log config.Logger = logging.MustGetLogger("helpers")
)

func InitLogger(logger config.Logger) {
	log = logger
}

func CheckValidation(object interface{}, binerrs binding.Errors, r render.Render, language string) error {
	fielderrors, errcode := models.ConvertErrors(language, binerrs)
	switch errcode {
	case types.TYPE_ERROR_LANGUAGE_NOTSUPPORTED:
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_LANGUAGE_NOTSUPPORTED,
			Message: config.Localization[language].Errors.Api.Language_NotSupported})
		return errors.New("Wrong language")
	case types.TYPE_ERROR_NONE:
		if len(*fielderrors) > 0 {
			for i, _ := range *fielderrors {
				value, found := models.GetFieldValue((*fielderrors)[i].Field, object)
				if found {
					(*fielderrors)[i].FieldValue = fmt.Sprintf("%v", value)
				}
			}
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong, Errors: *fielderrors})
			return errors.New("Wrong fields")
		}
	}

	return nil
}

func CheckParameterInt(r render.Render, param string, language string) (value int64, err error) {
	if param == "" || len(param) > PARAM_LENGTH_MAX {
		log.Error("Parameter is too long or too short %v", param)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return 0, errors.New("Length is wrong")
	}

	value, err = strconv.ParseInt(param, 0, 64)
	if err != nil {
		log.Error("Can't convert to number %v with value %v", err, param)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return 0, errors.New("Value is wrong")
	}

	return value, nil
}

func RenderJSONArray(object interface{}, length int, w http.ResponseWriter, r render.Render) {
	if length != 0 {
		r.JSON(http.StatusOK, object)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}
}
