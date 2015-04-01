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
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, facilities)
}

// get /api/v1.0/suppliers/services/
func GetSupplierFacilities(r render.Render, facilityrepository services.FacilityRepository, session *models.DtoSession) {
	facilities, err := facilityrepository.GetByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
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
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}
		if !dtofacility.Active {
			log.Error("Service is not active %v", dtofacility.ID)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
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
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, apifacilities)
}

// options /api/v1.0/suppliers/orders/
func GetMetaOrders(r render.Render, orderrepository services.OrderRepository, session *models.DtoSession) {
	order, err := orderrepository.GetMeta(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
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

	orders, err := orderrepository.GetAll(session.UserID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
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

	orderstatuses, err := orderstatusrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	apiorder := new(models.ApiLongOrder)
	apiorder.ID = dtoorder.ID
	apiorder.Name = dtoorder.Name
	apiorder.Step = dtoorder.Step
	apiorder.Facility_ID = dtoorder.Facility_ID
	apiorder.Supplier_ID = dtoorder.Supplier_ID
	apiorder.Proposed_Price = dtoorder.Proposed_Price
	apiorder.Charged_Fee = dtoorder.Charged_Fee
	for _, orderstatus := range *orderstatuses {
		switch orderstatus.Status_ID {
		case models.ORDER_STATUS_COMPLETED:
			apiorder.IsAssembled = orderstatus.Value
		case models.ORDER_STATUS_MODERATOR_CONFIRMED:
			apiorder.IsConfirmed = orderstatus.Value
		case models.ORDER_STATUS_NEW:
			apiorder.IsNew = orderstatus.Value
		case models.ORDER_STATUS_OPEN:
			apiorder.IsOpen = orderstatus.Value
		case models.ORDER_STATUS_CANCEL:
			apiorder.IsCancelled = orderstatus.Value
			apiorder.Reason = orderstatus.Comments
		case models.ORDER_STATUS_SUPPLIER_COST_NEW:
			apiorder.IsNewCost = orderstatus.Value
		case models.ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED:
			apiorder.IsNewCostConfirmed = orderstatus.Value
		case models.ORDER_STATUS_PAID:
			apiorder.IsPaid = orderstatus.Value
		case models.ORDER_STATUS_MODERATOR_BEGIN:
			apiorder.IsStarted = orderstatus.Value
		case models.ORDER_STATUS_SUPPLIER_CLOSE:
			apiorder.IsExecuted = orderstatus.Value
		case models.ORDER_STATUS_MODERATOR_DOCUMENTS_GOTTEN:
			apiorder.IsDocumented = orderstatus.Value
		case models.ORDER_STATUS_MODERATOR_CLOSE:
			apiorder.IsClosed = orderstatus.Value
		case models.ORDER_STATUS_ARCHIVE:
			apiorder.IsArchived = orderstatus.Value
		case models.ORDER_STATUS_DEL:
			apiorder.IsDeleted = orderstatus.Value
		}
	}

	r.JSON(http.StatusOK, apiorder)
}

// put /api/v1.0/suppliers/orders/:oid/
func UpdateOrder(errors binding.Errors, vieworder models.ViewOrder, r render.Render, params martini.Params,
	orderrepository services.OrderRepository, orderstatusrepository services.OrderStatusRepository,
	facilityrepository services.FacilityRepository, unitrepository services.UnitRepository,
	session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}
	if vieworder.Step > models.MAX_STEP_NUMBER {
		log.Error("Order step number is too big %v", vieworder.Step)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	dtofacility, err := facilityrepository.Get(vieworder.Facility_ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	if !dtofacility.Active {
		log.Error("Service is not active %v", dtofacility.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	dtounit, err := unitrepository.Get(vieworder.Supplier_ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	dtoorder.Supplier_ID = dtounit.ID
	dtoorder.Facility_ID = dtofacility.ID
	dtoorder.Name = vieworder.Name
	dtoorder.Step = vieworder.Step
	dtoorder.Proposed_Price = vieworder.Proposed_Price
	dtoorder.Charged_Fee = vieworder.Charged_Fee

	orderstatuses := []models.DtoOrderStatus{
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_COMPLETED, Value: vieworder.IsAssembled, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_MODERATOR_CONFIRMED, Value: vieworder.IsConfirmed, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_NEW, Value: vieworder.IsNew, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_OPEN, Value: vieworder.IsOpen, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_CANCEL, Value: vieworder.IsCancelled, Comments: vieworder.Reason, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_SUPPLIER_COST_NEW, Value: vieworder.IsNewCost, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED, Value: vieworder.IsNewCostConfirmed, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_PAID, Value: vieworder.IsPaid, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_MODERATOR_BEGIN, Value: vieworder.IsStarted, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_SUPPLIER_CLOSE, Value: vieworder.IsExecuted, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_MODERATOR_DOCUMENTS_GOTTEN, Value: vieworder.IsDocumented, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_MODERATOR_CLOSE, Value: vieworder.IsClosed, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_ARCHIVE, Value: vieworder.IsArchived, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_DEL, Value: vieworder.IsDeleted, Created: time.Now()},
	}

	err = orderrepository.Update(dtoorder, &orderstatuses, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongOrder(dtoorder.ID, dtoorder.Name, dtoorder.Step, vieworder.IsAssembled, vieworder.IsConfirmed,
		dtoorder.Facility_ID, dtoorder.Supplier_ID, vieworder.IsNew, vieworder.IsOpen, vieworder.IsCancelled, vieworder.Reason,
		vieworder.Proposed_Price, vieworder.IsNewCost, vieworder.IsNewCostConfirmed, vieworder.IsPaid, vieworder.IsStarted,
		vieworder.Charged_Fee, vieworder.IsExecuted, vieworder.IsDocumented, vieworder.IsClosed, vieworder.IsArchived, vieworder.IsDeleted))
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
