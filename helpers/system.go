package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/martini-contrib/render"
	"net"
	"net/http"
	"time"
	"types"
)

const (
	PARAMETER_NAME_UNSUBSCRIBE_CODE   = "unsubscribeCode"
	PARAMETER_NAME_SUBSCRIBТION_EMAIL = "email"
)

func CheckFrequence(method string, timeout time.Duration, request *http.Request, r render.Render, requestrepository services.RequestRepository,
	language string) (err error) {
	host := request.Header.Get(REQUEST_HEADER_X_FORWARDED_FOR)
	if host == "" {
		host, _, err = net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			log.Error("Can't detect ip address %v from %v", err, request.RemoteAddr)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return err
		}
	}
	exists, err := requestrepository.Exists(host, method)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}
	var hits int64 = 0
	if exists {
		request, err := requestrepository.Get(host, method)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return err
		}
		if time.Now().Sub(request.LastUpdated) < timeout {
			log.Error("Requests are too frequent for method %v", request.Method)
			r.JSON(http.StatusServiceUnavailable, types.Error{Code: types.TYPE_ERROR_REQUEST_TOOFREQUENT,
				Message: config.Localization[language].Errors.Api.Request_Too_Often})
			return errors.New("Frequent requests")
		}
		hits = request.Hits
	}
	dtorequest := models.NewDtoRequest(host, method, time.Now(), hits+1)
	err = requestrepository.Save(dtorequest)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return err
	}

	return nil
}
