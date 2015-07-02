package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/martini-contrib/render"
	"types"
)

const (
	HTTP_STATUS_CAPTCHA_REQUIRED = 449
)

func Check(hash string, value string, r render.Render, captcharepository services.CaptchaRepository) (err error) {
	if hash != "" {
		if value != "" {
			var captcha *models.DtoCaptcha
			captcha, err = captcharepository.Get(hash)
			if err != nil {
				r.JSON(HTTP_STATUS_CAPTCHA_REQUIRED, types.Error{Code: types.TYPE_ERROR_CAPTCHA_WRONG,
					Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Captcha_Wrong})
				return err
			}
			if captcha.InUse {
				r.JSON(HTTP_STATUS_CAPTCHA_REQUIRED, types.Error{Code: types.TYPE_ERROR_CAPTCHA_REQUIRED,
					Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Captcha_Required})
				return errors.New("Captcha is required")
			}
			if value != captcha.Value {
				r.JSON(HTTP_STATUS_CAPTCHA_REQUIRED, types.Error{Code: types.TYPE_ERROR_CAPTCHA_WRONG,
					Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Captcha_Wrong})
				return errors.New("Not matched captchas")
			} else {
				captcha.InUse = true
				err = captcharepository.Update(captcha)
				if err != nil {
					r.JSON(HTTP_STATUS_CAPTCHA_REQUIRED, types.Error{Code: types.TYPE_ERROR_CAPTCHA_WRONG,
						Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Captcha_Wrong})
					return err
				}
			}
		} else {
			r.JSON(HTTP_STATUS_CAPTCHA_REQUIRED, types.Error{Code: types.TYPE_ERROR_CAPTCHA_REQUIRED,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Captcha_Required})
			return errors.New("Captcha is required")
		}
	}

	return nil
}
