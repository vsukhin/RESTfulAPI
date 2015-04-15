package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
	"time"
	"types"
)

// options /api/v1.0/services/
// options /api/v1.0/suppliers/services/
func GetFacilities(r render.Render, facilityrepository services.FacilityRepository, session *models.DtoSession) {
	facilities, err := facilityrepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, facilities)
}

// get /api/v1.0/suppliers/services/
func GetSupplierFacilities(r render.Render, facilityrepository services.FacilityRepository, session *models.DtoSession) {
	facilities, err := facilityrepository.GetByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, facilities)
}

// put /api/v1.0/suppliers/services/
func UpdateSupplierFacilities(errors binding.Errors, viewfacilities models.ViewFacilities, r render.Render,
	facilityrepository services.FacilityRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}

	facilities := new([]int64)
	for _, viewfacility := range viewfacilities {
		dtofacility, err := facilityrepository.Get(viewfacility.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}
		if !dtofacility.Active {
			log.Error("Service is not active %v", dtofacility.ID)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}

		*facilities = append(*facilities, dtofacility.ID)
	}

	err := facilityrepository.SetByUser(session.UserID, facilities, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	apifacilities, err := facilityrepository.GetByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, apifacilities)
}

// options /api/v1.0/suppliers/orders/
func GetMetaOrders(r render.Render, orderrepository services.OrderRepository, session *models.DtoSession) {
	order, err := orderrepository.GetMeta(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, order)
}

// get /api/v1.0/suppliers/orders/
func GetOrders(request *http.Request, r render.Render, orderrepository services.OrderRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.OrderSearch), nil, request, r, session.Language)
	if err != nil {
		return
	}
	if len(*filters) != 0 {
		var masks []string
		for _, filter := range *filters {
			var exps []string
			for _, field := range filter.Fields {
				exps = append(exps, field+" "+filter.Op+" "+filter.Value)
			}
			masks = append(masks, "("+strings.Join(exps, " or ")+")")
		}
		query += " and "
		query += strings.Join(masks, " and ")
	}

	var sorts *[]models.OrderExp
	sorts, err = helpers.GetOrderArray(new(models.OrderSearch), request, r, session.Language)
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

	orders, err := orderrepository.GetByUser(session.UserID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, orders)
}

// get /api/v1.0/suppliers/orders/:oid/
func GetOrder(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	orderstatusrepository services.OrderStatusRepository, session *models.DtoSession) {
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}

	dtoorderstatuses, err := orderstatusrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongOrderFromDto(dtoorder, dtoorderstatuses))
}

// put /api/v1.0/suppliers/orders/:oid/
func UpdateOrder(errors binding.Errors, vieworder models.ViewLongOrder, r render.Render, params martini.Params,
	orderrepository services.OrderRepository, facilityrepository services.FacilityRepository, unitrepository services.UnitRepository,
	session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}

	apiorder, err := helpers.UpdateOrder(dtoorder, &vieworder, r, params, orderrepository, unitrepository, facilityrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apiorder)
}

// delete /api/v1.0/suppliers/orders/:oid/
func DeleteOrder(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	orderstatusrepository services.OrderStatusRepository, session *models.DtoSession) {
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}

	orderstatus := models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_CANCEL, true, "", time.Now())
	err = orderstatusrepository.Save(orderstatus)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// get /api/v1.0/suppliers/orders/:oid/services/:sid/
func GetOrderInfo(r render.Render, session *models.DtoSession) {

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// put /api/v1.0/suppliers/orders/:oid/services/:sid/
func UpdateOrderInfo(r render.Render, session *models.DtoSession) {

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
