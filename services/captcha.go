package services

import (
	"application/config"
	"application/models"
	"errors"
	"github.com/martini-contrib/render"
	"types"
)

const (
	HTTP_STATUS_CAPTCHA_REQUIRED = 449
)

type CaptchaRepository interface {
	Check(hash string, value string, r render.Render) (err error)
	Get(hash string) (captcha *models.DtoCaptcha, err error)
	Create(captcha *models.DtoCaptcha) (err error)
	Update(captcha *models.DtoCaptcha) (err error)
	Delete(hash string) (err error)
}

type CaptchaService struct {
	*Repository
}

func NewCaptchaService(repository *Repository) *CaptchaService {
	repository.DbContext.AddTableWithName(models.DtoCaptcha{}, repository.Table).SetKeys(false, "hash")
	return &CaptchaService{
		repository,
	}
}

func (captchaservice *CaptchaService) Check(hash string, value string, r render.Render) (err error) {
	if hash != "" {
		if value != "" {
			var captcha *models.DtoCaptcha
			captcha, err = captchaservice.Get(hash)
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
				err = captchaservice.Update(captcha)
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

func (captchaservice *CaptchaService) Get(hash string) (captcha *models.DtoCaptcha, err error) {
	captcha = new(models.DtoCaptcha)
	err = captchaservice.DbContext.SelectOne(captcha, "select * from "+captchaservice.Table+" where hash = ?", hash)
	if err != nil {
		log.Error("Error during getting captcha object from database %v with value %v", err, hash)
		return nil, err
	}

	return captcha, nil
}

func (captchaservice *CaptchaService) Create(captcha *models.DtoCaptcha) (err error) {
	err = captchaservice.DbContext.Insert(captcha)
	if err != nil {
		log.Error("Error during creating captcha object in database %v", err)
		return err
	}

	return nil
}

func (captchaservice *CaptchaService) Update(captcha *models.DtoCaptcha) (err error) {
	_, err = captchaservice.DbContext.Update(captcha)
	if err != nil {
		log.Error("Error during updating captcha object in database %v with value %v", err, captcha.Hash)
		return err
	}

	return nil
}

func (captchaservice *CaptchaService) Delete(hash string) (err error) {
	_, err = captchaservice.DbContext.Exec("delete from "+captchaservice.Table+" where hash = ?", hash)
	if err != nil {
		log.Error("Error during deleting captcha object in database %v with value %v", err, hash)
		return err
	}

	return nil
}
