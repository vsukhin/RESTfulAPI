package services

import (
	"application/db"
	"application/models"
	"database/sql"
	"errors"
	"github.com/coopernurse/gorp"
	"testing"
)

type TestCaptchaDBMap struct {
	Captcha *models.DtoCaptcha
	Err     error
}

func (testCaptchaDBMap *TestCaptchaDBMap) AddTableWithName(i interface{}, name string) *gorp.TableMap {
	return nil
}

func (testCaptchaDBMap *TestCaptchaDBMap) Begin() (*gorp.Transaction, error) {
	return nil, nil
}

func (testCaptchaDBMap *TestCaptchaDBMap) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	return nil, nil
}

func (testCaptchaDBMap *TestCaptchaDBMap) Insert(list ...interface{}) error {
	return testCaptchaDBMap.Err
}

func (testCaptchaDBMap *TestCaptchaDBMap) Update(list ...interface{}) (int64, error) {
	return 0, testCaptchaDBMap.Err
}

func (testCaptchaDBMap *TestCaptchaDBMap) Delete(list ...interface{}) (int64, error) {
	return 0, testCaptchaDBMap.Err
}

func (testCaptchaDBMap *TestCaptchaDBMap) Exec(query string, args ...interface{}) (sql.Result, error) {
	return *new(sql.Result), testCaptchaDBMap.Err
}

func (testCaptchaDBMap *TestCaptchaDBMap) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return nil, testCaptchaDBMap.Err
}

func (testCaptchaDBMap *TestCaptchaDBMap) SelectInt(query string, args ...interface{}) (int64, error) {
	return 0, nil
}

func (testCaptchaDBMap *TestCaptchaDBMap) SelectStr(query string, args ...interface{}) (string, error) {
	return "", nil
}

func (testCaptchaDBMap *TestCaptchaDBMap) SelectOne(holder interface{}, query string, args ...interface{}) error {
	captcha, _ := holder.(*models.DtoCaptcha)
	(*captcha) = (*testCaptchaDBMap.Captcha)
	return testCaptchaDBMap.Err
}

func TestCreateError(t *testing.T) {
	dbmap := new(TestCaptchaDBMap)
	dbmap.Err = errors.New("Captcha error")
	captcha := new(models.DtoCaptcha)
	captchaService := new(CaptchaService)
	captchaService.Repository = NewRepository(dbmap, db.TABLE_CAPTCHAS)

	err := captchaService.Create(captcha)
	if err == nil {
		t.Error("Create should return error")
	}
}

func TestCreateOk(t *testing.T) {
	dbmap := new(TestCaptchaDBMap)
	dbmap.Err = nil
	captcha := new(models.DtoCaptcha)
	captchaService := new(CaptchaService)
	captchaService.Repository = NewRepository(dbmap, db.TABLE_CAPTCHAS)

	err := captchaService.Create(captcha)
	if err != nil {
		t.Error("Create should not return error")
	}
}

func TestUpdateError(t *testing.T) {
	dbmap := new(TestCaptchaDBMap)
	dbmap.Err = errors.New("Captcha error")
	captcha := new(models.DtoCaptcha)
	captchaService := new(CaptchaService)
	captchaService.Repository = NewRepository(dbmap, db.TABLE_CAPTCHAS)

	err := captchaService.Update(captcha)
	if err == nil {
		t.Error("Update should return error")
	}
}

func TestUpdateOk(t *testing.T) {
	dbmap := new(TestCaptchaDBMap)
	dbmap.Err = nil
	captcha := new(models.DtoCaptcha)
	captchaService := new(CaptchaService)
	captchaService.Repository = NewRepository(dbmap, db.TABLE_CAPTCHAS)

	err := captchaService.Update(captcha)
	if err != nil {
		t.Error("Update should not return error")
	}
}

func TestDeleteError(t *testing.T) {
	dbmap := new(TestCaptchaDBMap)
	dbmap.Err = errors.New("Captcha error")
	hash := "12345"
	captchaService := new(CaptchaService)
	captchaService.Repository = NewRepository(dbmap, db.TABLE_CAPTCHAS)

	err := captchaService.Delete(hash)
	if err == nil {
		t.Error("Delete should return error")
	}
}

func TestDeleteOk(t *testing.T) {
	dbmap := new(TestCaptchaDBMap)
	dbmap.Err = nil
	hash := "12345"
	captchaService := new(CaptchaService)
	captchaService.Repository = NewRepository(dbmap, db.TABLE_CAPTCHAS)

	err := captchaService.Delete(hash)
	if err != nil {
		t.Error("Delete should not return error")
	}
}

func TestGetError(t *testing.T) {
	dbmap := new(TestCaptchaDBMap)
	dbmap.Err = errors.New("Captcha error")
	hash := "12345"
	dbmap.Captcha = new(models.DtoCaptcha)
	dbmap.Captcha.Hash = hash
	captchaService := new(CaptchaService)
	captchaService.Repository = NewRepository(dbmap, db.TABLE_CAPTCHAS)

	_, err := captchaService.Get(hash)
	if err == nil {
		t.Error("Get should return error")
	}
}

func TestGetOk(t *testing.T) {
	dbmap := new(TestCaptchaDBMap)
	dbmap.Err = nil
	hash := "12345"
	dbmap.Captcha = new(models.DtoCaptcha)
	dbmap.Captcha.Hash = hash
	captchaService := new(CaptchaService)
	captchaService.Repository = NewRepository(dbmap, db.TABLE_CAPTCHAS)

	captcha, err := captchaService.Get(hash)
	if err != nil {
		t.Error("Get should not return error")
	}
	if captcha.Hash != hash {
		t.Error("Get should return proper captcha")
	}
}
