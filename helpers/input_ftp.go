package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

func CheckInputFtp(orderid int64, r render.Render, inputftprepository services.InputFtpRepository,
	language string) (dtoinputftp *models.DtoInputFtp, err error) {
	found, err := inputftprepository.Exists(orderid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if found {
		dtoinputftp, err = inputftprepository.Get(orderid)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
	} else {
		dtoinputftp = new(models.DtoInputFtp)
		dtoinputftp.Order_ID = orderid
	}

	return dtoinputftp, nil
}
