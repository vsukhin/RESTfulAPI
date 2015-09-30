package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"net/http"
	"time"
	"types"
)

const (
	DEVICE_TTL                   = 5 * time.Minute
	DEVICE_TTL_ADJUSTMENT        = time.Minute
	DEVICE_CODE_MAX              = 9999
	DEVICE_TOKEN_COST            = 7
	METHOD_NAME_NEWDEVICE        = "/api/v1.0/user/devices/link/"
	METHOD_TIMEOUT_NEWDEVICE     = 3 * time.Second
	METHOD_NAME_UPDATEDEVICE     = "/api/v1.0/user/devices/"
	METHOD_TIMEOUT_UPDATEDEVICE  = 3 * time.Second
	METHOD_NAME_LINKDEVICE       = "/api/v1.0/user/devices/code/"
	METHOD_TIMEOUT_LINKDEVICE    = 3 * time.Second
	METHOD_NAME_DEVICESESSION    = "/api/v1.0/session/device/"
	METHOD_TIMEOUT_DEVICESESSION = 3 * time.Second
)

// post /api/v1.0/user/devices/link/
func CreateDevice(errors binding.Errors, viewdevice models.ViewLongDevice, request *http.Request, r render.Render,
	devicerepository services.DeviceRepository, sessionrepository services.SessionRepository, requestrepository services.RequestRepository) {
	if helpers.CheckFrequence(METHOD_NAME_NEWDEVICE, METHOD_TIMEOUT_NEWDEVICE, request, r, requestrepository,
		config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	found, err := devicerepository.Exists(viewdevice.Serial)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	if found {
		log.Error("Device %v is already linked to user", viewdevice.Serial)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}
	token, err := sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}
	code, err := rand.Int(rand.Reader, big.NewInt(int64(DEVICE_CODE_MAX)))
	if err != nil {
		log.Error("Device code generation error %v", err)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	dtodevice := new(models.DtoDevice)
	dtodevice.User_ID = 0
	dtodevice.OS = viewdevice.OS
	dtodevice.App = viewdevice.App
	dtodevice.Serial = viewdevice.Serial
	dtodevice.Token = token
	dtodevice.Code = fmt.Sprintf("%04d", code)
	hash := sha512.Sum512([]byte(dtodevice.Code + dtodevice.Token))
	dtodevice.Hash = hex.EncodeToString(hash[:])
	dtodevice.Created = time.Now()
	dtodevice.Valid_Till = dtodevice.Created.Add(DEVICE_TTL)
	dtodevice.Active = true

	err = devicerepository.Create(dtodevice)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiDevice(dtodevice.Token, dtodevice.Code))
}

// post /api/v1.0/user/devices/
func UpdateDevice(errors binding.Errors, viewdevice models.ViewHashDevice, request *http.Request, r render.Render,
	devicerepository services.DeviceRepository, requestrepository services.RequestRepository) {
	if helpers.CheckFrequence(METHOD_NAME_UPDATEDEVICE, METHOD_TIMEOUT_UPDATEDEVICE, request, r, requestrepository,
		config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	dtodevice, err := devicerepository.FindByHash(viewdevice.Hash)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	currenttime := time.Now()
	if currenttime.Sub(dtodevice.Valid_Till) > 0 {
		log.Error("Hash has been expired for device %v", dtodevice.Serial)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	if currenttime.Add(DEVICE_TTL_ADJUSTMENT).Sub(dtodevice.Valid_Till) > 0 {
		dtodevice.Valid_Till = currenttime.Add(DEVICE_TTL_ADJUSTMENT)

		err = devicerepository.Update(dtodevice)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
			return
		}
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[config.Configuration.Server.DefaultLanguage].Messages.OK})
}

// post /api/v1.0/user/devices/code/
func LinkDevice(errors binding.Errors, viewdevice models.ViewCodeDevice, request *http.Request, r render.Render,
	devicerepository services.DeviceRepository, requestrepository services.RequestRepository, session *models.DtoSession) {
	if helpers.CheckFrequence(METHOD_NAME_LINKDEVICE, METHOD_TIMEOUT_LINKDEVICE, request, r, requestrepository,
		session.Language) != nil {
		return
	}
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	dtodevice, err := devicerepository.FindByCode(viewdevice.Code)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	if time.Now().Sub(dtodevice.Valid_Till) > 0 {
		log.Error("Code has been expired for device %v", dtodevice.Serial)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	dtodevice.Code = ""
	dtodevice.Hash = ""
	dtodevice.User_ID = session.UserID

	err = devicerepository.Update(dtodevice)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// post /api/v1.0/session/device/
func CreateSessionDevice(errors binding.Errors, viewdevice models.ViewTokenDevice, request *http.Request, r render.Render,
	devicerepository services.DeviceRepository, requestrepository services.RequestRepository, userrepository services.UserRepository,
	sessionrepository services.SessionRepository) {
	if helpers.CheckFrequence(METHOD_NAME_DEVICESESSION, METHOD_TIMEOUT_DEVICESESSION, request, r, requestrepository,
		config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	encryptedtoken, err := hex.DecodeString(viewdevice.Token)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_TOKEN_HASH_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Token_Hash_Wrong})
		return
	}
	cost, err := bcrypt.Cost(encryptedtoken)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_TOKEN_HASH_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Token_Hash_Wrong})
		return
	}
	if cost < DEVICE_TOKEN_COST {
		log.Error("Wrong bcrypt complexity for token %v", viewdevice.Token)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_BCRYPT_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Bcrypt_Wrong})
		return
	}
	dtodevices, err := devicerepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	var dtodevice *models.DtoDevice
	for i, device := range *dtodevices {
		if bcrypt.CompareHashAndPassword(encryptedtoken, []byte(device.Token)) == nil {
			dtodevice = &(*dtodevices)[i]
			break
		}
	}
	if dtodevice == nil {
		log.Error("Can't find device for token %v", viewdevice.Token)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	user, err := userrepository.Get(dtodevice.User_ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DEVICE_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Device_Wrong})
		return
	}
	if !user.Active || !user.Confirmed {
		log.Error("User is not active or confirmed %v", user.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_USER_BLOCKED,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.User_Blocked})
		return
	}

	token, err := sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	dtosession := models.NewDtoSession(token, user.ID, user.Roles, time.Now(), config.Configuration.Server.DefaultLanguage)
	err = sessionrepository.Create(dtosession, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	user.LastLogin = dtosession.LastActivity
	err = userrepository.Update(user, true, false)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiSession(dtosession.LastActivity.Add(config.Configuration.Server.SessionTimeout), dtosession.AccessToken))
}
