package services

import (
	"application/models"
)

type CaptchaRepository interface {
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
