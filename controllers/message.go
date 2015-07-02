package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
	"time"
	"types"
)

// options /api/v1.0/messages/orders/:oid/
func GetMetaMessages(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	messagerepository services.MessageRepository, session *models.DtoSession) {
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}

	apimessage, err := messagerepository.GetMetaByOrder(dtoorder.ID, session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, apimessage)
}

// get /api/v1.0/messages/orders/:oid/
func GetMessages(w http.ResponseWriter, request *http.Request, r render.Render, params martini.Params, orderrepository services.OrderRepository,
	messagerepository services.MessageRepository, session *models.DtoSession) {
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}

	query := ""
	var filters *[]models.FilterExp
	filters, err = helpers.GetFilterArray(new(models.MessageSearch), nil, request, r, session.Language)
	if err != nil {
		return
	}
	if len(*filters) != 0 {
		var masks []string
		for _, filter := range *filters {
			var exps []string
			for _, field := range filter.Fields {
				switch field {
				case "isMine":
					field = fmt.Sprintf("(m.user_id = %v)", session.UserID)
				case "new":
					field = fmt.Sprintf("(coalesce((select u.user_id from user_messages u where u.message_id = m.id and u.user_id = %v), 0) <> %v)",
						session.UserID, session.UserID)
				}
				exps = append(exps, field+" "+filter.Op+" "+filter.Value)
			}
			masks = append(masks, "("+strings.Join(exps, " or ")+")")
		}
		query += " and "
		query += strings.Join(masks, " and ")
	}

	var sorts *[]models.OrderExp
	sorts, err = helpers.GetOrderArray(new(models.MessageSearch), request, r, session.Language)
	if err != nil {
		return
	}
	if len(*sorts) != 0 {
		var orders []string
		for _, sort := range *sorts {
			orders = append(orders, " "+sort.Field+" "+sort.Order)
		}
		query += " order by"
		query += strings.Join(orders, ",")
	}

	var limit string
	limit, err = helpers.GetLimitQuery(request, r, session.Language)
	if err != nil {
		return
	}
	query += limit

	messages, err := messagerepository.GetByOrder(dtoorder.ID, session.UserID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(messages, len(*messages), w, r)
}

// post /api/v1.0/messages/order/:oid/
func CreateMessage(errors binding.Errors, viewmessage models.ViewLongMessage, r render.Render, params martini.Params,
	orderrepository services.OrderRepository, messagerepository services.MessageRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewmessage, errors, r, session.Language) != nil {
		return
	}
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}
	allowed, err := orderrepository.CheckUnitAccess(viewmessage.Receiver_ID, dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if !allowed {
		log.Error("Receiver %v is not allowed for order %v", viewmessage.Receiver_ID, dtoorder.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	dtomessage := new(models.DtoMessage)
	dtomessage.User_ID = session.UserID
	dtomessage.Order_ID = dtoorder.ID
	dtomessage.Created = time.Now()
	dtomessage.Content = viewmessage.Content
	dtomessage.Receiver_ID = viewmessage.Receiver_ID
	err = messagerepository.Create(dtomessage, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortMessage(dtomessage.ID))
}

// get /api/v1.0/messages/orders/:oid/message/:mid/
func GetMessage(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	messagerepository services.MessageRepository, session *models.DtoSession) {
	dtomessage, err := helpers.CheckMessage(r, params, orderrepository, messagerepository, session.Language)
	if err != nil {
		return
	}

	read, err := messagerepository.IsReadByUser(session.UserID, dtomessage.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongMessage(dtomessage.ID, dtomessage.Created, !read,
		dtomessage.User_ID == session.UserID, dtomessage.User_ID, dtomessage.Receiver_ID, dtomessage.Content))
}

// patch /api/v1.0/messages/orders/:oid/message/:mid/
func MarkMessage(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	messagerepository services.MessageRepository, session *models.DtoSession) {
	dtomessage, err := helpers.CheckMessage(r, params, orderrepository, messagerepository, session.Language)
	if err != nil {
		return
	}

	err = messagerepository.SetReadByUser(session.UserID, dtomessage.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// patch /api/v1.0/messages/orders/:oid/messages/
func MarkMessages(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	messagerepository services.MessageRepository, session *models.DtoSession) {
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}

	err = messagerepository.SetReadByUserForOrder(session.UserID, dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// put /api/v1.0/messages/orders/:oid/message/:mid/
func UpdateMessage(errors binding.Errors, viewmessage models.ViewShortMessage, r render.Render, params martini.Params,
	orderrepository services.OrderRepository, messagerepository services.MessageRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewmessage, errors, r, session.Language) != nil {
		return
	}
	dtomessage, err := helpers.CheckChangeableMessage(r, params, orderrepository, messagerepository, session.UserID, session.Language, true)
	if err != nil {
		return
	}

	dtomessage.Content = viewmessage.Content
	err = messagerepository.Update(dtomessage)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	read, err := messagerepository.IsReadByUser(session.UserID, dtomessage.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongMessage(dtomessage.ID, dtomessage.Created, !read,
		dtomessage.User_ID == session.UserID, dtomessage.User_ID, dtomessage.Receiver_ID, dtomessage.Content))
}

// delete /api/v1.0/messages/orders/:oid/message/:mid/
func DeleteMessage(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	messagerepository services.MessageRepository, session *models.DtoSession) {
	dtomessage, err := helpers.CheckChangeableMessage(r, params, orderrepository, messagerepository, session.UserID, session.Language, false)
	if err != nil {
		return
	}

	err = messagerepository.Delete(dtomessage, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
