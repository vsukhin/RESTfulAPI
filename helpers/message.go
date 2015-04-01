package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"time"
	"types"
)

const (
	PARAM_NAME_MESSAGE = "mid"
)

func CheckMessage(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	messagerepository services.MessageRepository, language string) (dtomessage *models.DtoMessage, err error) {
	message_id, err := CheckParameterInt(r, params[PARAM_NAME_MESSAGE], language)
	if err != nil {
		return nil, err
	}
	dtoorder, err := CheckOrder(r, params, orderrepository, language)
	if err != nil {
		return nil, err
	}

	dtomessage, err = messagerepository.Get(message_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	if dtomessage.Order_ID != dtoorder.ID {
		log.Error("Message %v doesn't belong to order %v", dtomessage.ID, dtoorder.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, errors.New("Message doesn't belong to order")
	}

	return dtomessage, nil
}

func CheckChangeableMessage(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	messagerepository services.MessageRepository, user_id int64, language string, hastimeout bool) (dtomessage *models.DtoMessage, err error) {
	dtomessage, err = CheckMessage(r, params, orderrepository, messagerepository, language)
	if err != nil {
		return nil, err
	}
	if dtomessage.User_ID != user_id {
		log.Error("Data can't be changed by other user %v for message %v", user_id, dtomessage.ID)
		r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Changes_Denied})
		return nil, errors.New("Wrong user")
	}

	last, err := messagerepository.IsLastForUser(user_id, dtomessage.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}
	if !last {
		log.Error("Data can't be changed when message is not last for user %v for message %v", user_id, dtomessage.ID)
		r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Changes_Denied})
		return nil, errors.New("Message is not last")
	}

	read, err := messagerepository.IsViewed(dtomessage.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}
	if read {
		old := time.Now().Sub(dtomessage.Created) > config.Configuration.Server.MessageTimeout
		if !hastimeout || (hastimeout && old) {
			log.Error("Data can't be changed when message is viewed %v", dtomessage.ID)
			r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Changes_Denied})
			return nil, errors.New("Message is viewed")
		}
	}

	return dtomessage, nil
}
