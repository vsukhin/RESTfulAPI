package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
	"time"
	"types"
)

const (
	PARAM_NAME_ORDER_ID = "oid"
)

func CheckOrder(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	language string) (dtoorder *models.DtoOrder, err error) {
	orderid, err := CheckParameterInt(r, params[PARAM_NAME_ORDER_ID], language)
	if err != nil {
		return
	}
	dtoorder, err = orderrepository.Get(orderid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return dtoorder, nil
}

func CheckOrderValidity(order_id int64, r render.Render, orderrepository services.OrderRepository,
	language string) (dtoorder *models.DtoOrder, err error) {
	dtoorder, err = orderrepository.Get(order_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return dtoorder, nil
}

func CheckOrderAccess(order_id int64, user_id int64, r render.Render, orderrepository services.OrderRepository,
	language string) (err error) {
	allowed, err := orderrepository.CheckUserAccess(user_id, order_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}
	if !allowed {
		log.Error("Order %v is not accessible for user %v", order_id, user_id)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return errors.New("Order not accessible")
	}

	return nil
}

func CheckOrderEditability(dtoorder *models.DtoOrder, r render.Render, orderstatusrepository services.OrderStatusRepository,
	language string) (err error) {
	dtoorderstatuses, err := orderstatusrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}
	for _, orderstatus := range *dtoorderstatuses {
		if (orderstatus.Status_ID == models.ORDER_STATUS_ARCHIVE && orderstatus.Value == true) ||
			(orderstatus.Status_ID == models.ORDER_STATUS_DEL && orderstatus.Value == true) {
			log.Error("Can't update order %v", dtoorder.ID)
			r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_CHANGES_DENIED,
				Message: config.Localization[language].Errors.Api.Data_Changes_Denied})
			return errors.New("Not editable order")
		}
	}

	return nil
}

func CheckOrderFacilityEditability(dtoorder *models.DtoOrder, r render.Render, orderstatusrepository services.OrderStatusRepository,
	language string) (err error) {
	dtoorderstatuses, err := orderstatusrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}
	for _, orderstatus := range *dtoorderstatuses {
		if (orderstatus.Status_ID == models.ORDER_STATUS_COMPLETED && orderstatus.Value == true) ||
			(orderstatus.Status_ID == models.ORDER_STATUS_ARCHIVE && orderstatus.Value == true) ||
			(orderstatus.Status_ID == models.ORDER_STATUS_DEL && orderstatus.Value == true) {
			log.Error("Can't update order %v", dtoorder.ID)
			r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_CHANGES_DENIED,
				Message: config.Localization[language].Errors.Api.Data_Changes_Denied})
			return errors.New("Not editable order")
		}
	}
	if dtoorder.Step < 2 || dtoorder.Step > 3 {
		log.Error("Can't update order %v for step %v", dtoorder.ID, dtoorder.Step)
		r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_CHANGES_DENIED,
			Message: config.Localization[language].Errors.Api.Data_Changes_Denied})
		return errors.New("Not editable order step")
	}

	return nil
}

func UpdateLongOrder(dtoorder *models.DtoOrder, vieworder *models.ViewLongOrder, r render.Render, params martini.Params,
	orderrepository services.OrderRepository, unitrepository services.UnitRepository, facilityrepository services.FacilityRepository,
	language string) (apiorder *models.ApiLongOrder, err error) {
	if vieworder.Step > models.MAX_STEP_NUMBER {
		log.Error("Order step number is too big %v", vieworder.Step)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Wrong step")
	}
	var unitid int64 = 0
	var facilityid int64 = 0
	if vieworder.Facility_ID != 0 {
		err = CheckFacility(vieworder.Facility_ID, r, facilityrepository, language)
		if err != nil {
			return nil, err
		}
		facilityid = vieworder.Facility_ID
	}
	if vieworder.Supplier_ID != 0 {
		err = CheckUnitValidity(vieworder.Supplier_ID, language, r, unitrepository)
		if err != nil {
			return nil, err
		}
		unitid = vieworder.Supplier_ID
	}
	dtoorder.Supplier_ID = unitid
	dtoorder.Facility_ID = facilityid
	dtoorder.Name = vieworder.Name
	dtoorder.Step = vieworder.Step
	dtoorder.Proposed_Price = vieworder.Proposed_Price
	dtoorder.Charged_Fee = vieworder.Charged_Fee
	dtoorder.Execution_Forecast = vieworder.Execution_Forecast
	if vieworder.IsPaid && (dtoorder.Begin_Date.IsZero() ||
		(dtoorder.Begin_Date.Year() == 1 && dtoorder.Begin_Date.Month() == 1 && dtoorder.Begin_Date.Day() == 1)) {
		dtoorder.Begin_Date = time.Now()
	}
	if vieworder.IsExecuted && (dtoorder.End_Date.IsZero() ||
		(dtoorder.End_Date.Year() == 1 && dtoorder.End_Date.Month() == 1 && dtoorder.End_Date.Day() == 1)) {
		dtoorder.End_Date = time.Now()
	}

	dtoorderstatuses := vieworder.ToOrderStatuses(dtoorder.ID)

	err = orderrepository.Update(dtoorder, dtoorderstatuses, nil, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	return models.NewApiLongOrderFromDto(dtoorder, dtoorderstatuses), nil
}

func UpdateFullOrder(dtoorder *models.DtoOrder, vieworder *models.ViewFullOrder, r render.Render, params martini.Params,
	orderrepository services.OrderRepository, unitrepository services.UnitRepository, facilityrepository services.FacilityRepository,
	userrepository services.UserRepository, projectrepository services.ProjectRepository,
	language string) (apiorder *models.ApiFullOrder, err error) {
	var dtouser *models.DtoUser
	if vieworder.Creator_ID != 0 {
		dtouser, err = CheckUser(vieworder.Creator_ID, language, r, userrepository)
		if err != nil {
			return nil, err
		}
	}
	err = CheckUnitValidity(vieworder.Unit_ID, language, r, unitrepository)
	if err != nil {
		return nil, err
	}
	if vieworder.Creator_ID != 0 {
		if dtouser.UnitID != vieworder.Unit_ID {
			log.Error("User %v doesn't belong to unit %v", vieworder.Creator_ID, vieworder.Unit_ID)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, errors.New("User doesn't match unit")
		}
	}
	dtoproject, err := projectrepository.Get(dtoorder.Project_ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if dtoproject.Unit_ID != vieworder.Unit_ID {
		log.Error("Order project %v doesn't belong to unit %v", dtoproject.ID, vieworder.Unit_ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Order project doesn't match unit")
	}
	dtoorder.Unit_ID = vieworder.Unit_ID
	dtoorder.Creator_ID = vieworder.Creator_ID

	apilongorder, err := UpdateLongOrder(dtoorder, &vieworder.ViewLongOrder, r, params, orderrepository, unitrepository, facilityrepository, language)
	if err != nil {
		return nil, err
	}

	return models.NewApiFullOrder(dtoorder.Creator_ID, dtoorder.Unit_ID, dtoorder.Created, *apilongorder), nil
}

func GetOrderTables(order_id int64, r render.Render, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, language string) (resulttables *[]models.ApiResultTable,
	worktables *[]models.ApiWorkTable, err error) {
	resulttables, err = resulttablerepository.GetByOrder(order_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, err
	}

	worktables, err = worktablerepository.GetByOrder(order_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, err
	}

	return resulttables, worktables, nil
}

func GetSMSOrderData(dtoorder *models.DtoOrder, dtosmsfacility *models.DtoSMSFacility, r render.Render,
	mobileoperatoroperationrepository services.MobileOperatorOperationRepository,
	smsperiodrepository services.SMSPeriodRepository, smseventrepository services.SMSEventRepository,
	language string) (operators *[]models.ViewApiMobileOperatorOperation,
	periods *[]models.ViewApiSMSPeriod, events *[]models.ViewApiSMSEvent, deliveryType string, err error) {
	operators, err = mobileoperatoroperationrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, "", err
	}

	periods, err = smsperiodrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, "", err
	}

	events, err = smseventrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, "", err
	}

	deliveryType = ""
	switch dtosmsfacility.DeliveryType {
	case models.TYPE_DELIVERY_ONCE:
		deliveryType = models.TYPE_DELIVERY_ONCE_VALUE
	case models.TYPE_DELIVERY_SCHEDULED:
		deliveryType = models.TYPE_DELIVERY_SCHEDULED_VALUE
	case models.TYPE_DELIVERY_EVENTTRIGGERED:
		deliveryType = models.TYPE_DELIVERY_EVENTTRIGGERED_VALUE
	}

	return operators, periods, events, deliveryType, nil
}

func GetSMSOrder(dtoorder *models.DtoOrder, r render.Render, facilityrepository services.FacilityRepository,
	smsfacilityrepository services.SMSFacilityRepository, mobileoperatoroperationrepository services.MobileOperatorOperationRepository,
	smsperiodrepository services.SMSPeriodRepository, smseventrepository services.SMSEventRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository,
	language string) (apismsfacility *models.ApiSMSFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_SMS, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}

	found, err := smsfacilityrepository.Exists(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !found {
		return &models.ApiSMSFacility{}, nil
	}

	dtosmsfacility, err := smsfacilityrepository.Get(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	operators, periods, events, deliveryType, err := GetSMSOrderData(dtoorder, dtosmsfacility, r,
		mobileoperatoroperationrepository, smsperiodrepository, smseventrepository, language)
	if err != nil {
		return nil, err
	}

	resulttables, worktables, err := GetOrderTables(dtoorder.ID, r, resulttablerepository, worktablerepository, language)
	if err != nil {
		return nil, err
	}

	return models.NewApiSMSFacility(dtosmsfacility.EstimatedNumbersShipments, dtosmsfacility.EstimatedMessageInCyrillic,
		dtosmsfacility.EstimatedNumberCharacters, dtosmsfacility.EstimatedNumberSmsInMessage, *operators,
		deliveryType, dtosmsfacility.DeliveryTime, *periods, *events, dtosmsfacility.DeliveryTimeStart, dtosmsfacility.DeliveryTimeEnd,
		dtosmsfacility.DeliveryBaseTime, dtosmsfacility.DeliveryDataId, dtosmsfacility.DeliveryDataDelete, dtosmsfacility.MessageFromId,
		dtosmsfacility.MessageFromInColumnId, dtosmsfacility.MessageToInColumnId, dtosmsfacility.MessageBody, dtosmsfacility.MessageBodyInColumnId,
		dtosmsfacility.TimeCorrection, dtosmsfacility.Cost, dtosmsfacility.CostFactual, *resulttables, *worktables), nil
}

func UpdateSMSOrder(dtoorder *models.DtoOrder, viewsmsfacility models.ViewSMSFacility, r render.Render,
	facilityrepository services.FacilityRepository, smsfacilityrepository services.SMSFacilityRepository,
	orderstatusrepository services.OrderStatusRepository, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository,
	smssenderrepository services.SMSSenderRepository, mobileoperatorrepository services.MobileOperatorRepository,
	periodrepository services.PeriodRepository, eventrepository services.EventRepository,
	smsperiodrepository services.SMSPeriodRepository, smseventrepository services.SMSEventRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository,
	mobileoperatoroperationrepository services.MobileOperatorOperationRepository, checkaccess bool, userid int64,
	language string) (apismsfacility *models.ApiSMSFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_SMS, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}
	err = CheckOrderFacilityEditability(dtoorder, r, orderstatusrepository, language)
	if err != nil {
		return nil, err
	}

	found, err := smsfacilityrepository.Exists(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	var dtosmsfacility *models.DtoSMSFacility
	if !found {
		dtosmsfacility = new(models.DtoSMSFacility)
		dtosmsfacility.Order_ID = dtoorder.ID
	} else {
		dtosmsfacility, err = smsfacilityrepository.Get(dtoorder.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
	}

	if dtoorder.Step == 2 || dtoorder.Step == 3 {
		dtosmsfacility.EstimatedNumbersShipments = viewsmsfacility.EstimatedNumbersShipments
		dtosmsfacility.EstimatedMessageInCyrillic = viewsmsfacility.EstimatedMessageInCyrillic
		dtosmsfacility.EstimatedNumberCharacters = viewsmsfacility.EstimatedNumberCharacters
		dtosmsfacility.EstimatedNumberSmsInMessage = viewsmsfacility.EstimatedNumberSmsInMessage

		var totalpercent byte = 0
		var totalmessages int = 0
		for _, apioperator := range viewsmsfacility.EstimatedOperators {
			totalpercent += apioperator.Percent
			totalmessages += apioperator.Count
			_, err = CheckMobileOperator(apioperator.MobileOperator_ID, r, mobileoperatorrepository, language)
			if err != nil {
				return nil, err
			}
			dtosmsfacility.EstimatedOperators = append(dtosmsfacility.EstimatedOperators, *models.NewDtoMobileOperatorOperation(dtoorder.ID,
				apioperator.MobileOperator_ID, apioperator.Percent, apioperator.Count))
		}
		if totalpercent != 100 {
			log.Error("Total percent sum is not equal 100, %v", totalpercent)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return nil, errors.New("Wrong percent")
		}
		if totalmessages != viewsmsfacility.EstimatedNumbersShipments {
			log.Error("Total count sum is not equal estimated number of shipments, %v", totalmessages)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return nil, errors.New("Wrong count")
		}
	}

	if dtoorder.Step == 2 {
		smsperiods, err := smsperiodrepository.GetAll(dtoorder.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
		dtosmsfacility.Periods = *smsperiods
		smsevents, err := smseventrepository.GetAll(dtoorder.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
		dtosmsfacility.Events = *smsevents
	}
	if dtoorder.Step == 3 {
		var deliveryType models.DeliveryType
		switch strings.ToLower(viewsmsfacility.DeliveryType) {
		case models.TYPE_DELIVERY_ONCE_VALUE:
			deliveryType = models.TYPE_DELIVERY_ONCE
		case models.TYPE_DELIVERY_SCHEDULED_VALUE:
			deliveryType = models.TYPE_DELIVERY_SCHEDULED
		case models.TYPE_DELIVERY_EVENTTRIGGERED_VALUE:
			deliveryType = models.TYPE_DELIVERY_EVENTTRIGGERED
		default:
			log.Error("Unknown delivery type %v", viewsmsfacility.DeliveryType)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return nil, errors.New("Wrong delivery type")
		}
		dtosmsfacility.DeliveryType = deliveryType
		dtosmsfacility.DeliveryTime = viewsmsfacility.DeliveryTime

		if dtosmsfacility.DeliveryTime {
			for _, apiperiod := range viewsmsfacility.Periods {
				dtoperiod, err := CheckPeriod(apiperiod.Period_ID, r, periodrepository, language)
				if err != nil {
					return nil, err
				}
				dtosmsfacility.Periods = append(dtosmsfacility.Periods, *models.NewDtoSMSPeriod(dtoorder.ID, dtoperiod.ID))
			}

			for _, apievent := range viewsmsfacility.Events {
				dtoevent, err := CheckEvent(apievent.Event_ID, r, eventrepository, language)
				if err != nil {
					return nil, err
				}
				dtosmsfacility.Events = append(dtosmsfacility.Events, *models.NewDtoSMSEvent(dtoorder.ID, dtoevent.ID))
			}

			if !viewsmsfacility.DeliveryTimeStart.IsZero() && viewsmsfacility.DeliveryTimeStart.Sub(time.Now()) < 0 {
				log.Error("Time start is in the past %v", viewsmsfacility.DeliveryTimeStart)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return nil, errors.New("Wrong time start")
			}
			if !viewsmsfacility.DeliveryTimeEnd.IsZero() && viewsmsfacility.DeliveryTimeEnd.Sub(time.Now()) < 0 {
				log.Error("Time end is in the past %v", viewsmsfacility.DeliveryTimeEnd)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return nil, errors.New("Wrong time end")
			}
			if !viewsmsfacility.DeliveryTimeStart.IsZero() && !viewsmsfacility.DeliveryTimeEnd.IsZero() &&
				viewsmsfacility.DeliveryTimeStart.Sub(viewsmsfacility.DeliveryTimeEnd) > 0 {
				log.Error("Time start can't be bigger than time end %v", viewsmsfacility.DeliveryTimeStart, viewsmsfacility.DeliveryTimeEnd)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return nil, errors.New("Wrong time start and end")
			}
			dtosmsfacility.DeliveryTimeStart = viewsmsfacility.DeliveryTimeStart
			dtosmsfacility.DeliveryTimeEnd = viewsmsfacility.DeliveryTimeEnd
			dtosmsfacility.DeliveryBaseTime = viewsmsfacility.DeliveryBaseTime
		}

		_, err = IsTableAvailable(r, customertablerepository, viewsmsfacility.DeliveryDataId, language)
		if err != nil {
			return nil, err
		}
		if checkaccess {
			err = IsTableAccessible(viewsmsfacility.DeliveryDataId, userid, r, customertablerepository, language)
			if err != nil {
				return nil, err
			}
		}
		dtosmsfacility.DeliveryDataId = viewsmsfacility.DeliveryDataId
		dtosmsfacility.DeliveryDataDelete = viewsmsfacility.DeliveryDataDelete

		if viewsmsfacility.MessageFromInColumnId == 0 {
			_, err = IsSMSSenderActive(viewsmsfacility.MessageFromId, r, smssenderrepository, language)
			if err != nil {
				return nil, err
			}
			if checkaccess {
				err = IsSMSSenderAccessible(viewsmsfacility.MessageFromId, userid, r, smssenderrepository, language)
				if err != nil {
					return nil, err
				}
			}
		}
		dtosmsfacility.MessageFromId = viewsmsfacility.MessageFromId

		if viewsmsfacility.MessageFromInColumnId != 0 {
			_, err = CheckColumnValidity(viewsmsfacility.DeliveryDataId, viewsmsfacility.MessageFromInColumnId, r, columntyperepository,
				tablecolumnrepository, language)
			if err != nil {
				return nil, err
			}
		}
		dtosmsfacility.MessageFromInColumnId = viewsmsfacility.MessageFromInColumnId

		if viewsmsfacility.MessageToInColumnId != 0 {
			_, err = CheckColumnValidity(viewsmsfacility.DeliveryDataId, viewsmsfacility.MessageToInColumnId, r, columntyperepository,
				tablecolumnrepository, language)
			if err != nil {
				return nil, err
			}
		}
		dtosmsfacility.MessageToInColumnId = viewsmsfacility.MessageToInColumnId
		dtosmsfacility.MessageBody = viewsmsfacility.MessageBody

		if viewsmsfacility.MessageBodyInColumnId != 0 {
			_, err = CheckColumnValidity(viewsmsfacility.DeliveryDataId, viewsmsfacility.MessageBodyInColumnId, r, columntyperepository,
				tablecolumnrepository, language)
			if err != nil {
				return nil, err
			}
		}
		dtosmsfacility.MessageBodyInColumnId = viewsmsfacility.MessageBodyInColumnId
		dtosmsfacility.TimeCorrection = viewsmsfacility.TimeCorrection
	}

	resulttables, worktables, err := GetOrderTables(dtoorder.ID, r, resulttablerepository, worktablerepository, language)
	if err != nil {
		return nil, err
	}

	for _, resulttable := range *resulttables {
		dtosmsfacility.ResultTables = append(dtosmsfacility.ResultTables, *models.NewDtoResultTable(dtoorder.ID, resulttable.Customer_Table_ID))
	}
	for _, worktable := range *worktables {
		dtosmsfacility.WorkTables = append(dtosmsfacility.WorkTables, *models.NewDtoWorkTable(dtoorder.ID, worktable.Customer_Table_ID))
	}

	if !found {
		err = smsfacilityrepository.Create(dtosmsfacility, false, true)
	} else {
		err = smsfacilityrepository.Update(dtosmsfacility, false, true)
	}
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	operators, periods, events, deliveryType, err := GetSMSOrderData(dtoorder, dtosmsfacility, r,
		mobileoperatoroperationrepository, smsperiodrepository, smseventrepository, language)
	if err != nil {
		return nil, err
	}

	return models.NewApiSMSFacility(dtosmsfacility.EstimatedNumbersShipments, dtosmsfacility.EstimatedMessageInCyrillic,
		dtosmsfacility.EstimatedNumberCharacters, dtosmsfacility.EstimatedNumberSmsInMessage, *operators, deliveryType, dtosmsfacility.DeliveryTime,
		*periods, *events, dtosmsfacility.DeliveryTimeStart, dtosmsfacility.DeliveryTimeEnd, dtosmsfacility.DeliveryBaseTime, dtosmsfacility.DeliveryDataId,
		dtosmsfacility.DeliveryDataDelete, dtosmsfacility.MessageFromId, dtosmsfacility.MessageFromInColumnId, dtosmsfacility.MessageToInColumnId,
		dtosmsfacility.MessageBody, dtosmsfacility.MessageBodyInColumnId, dtosmsfacility.TimeCorrection, dtosmsfacility.Cost, dtosmsfacility.CostFactual,
		*resulttables, *worktables), nil
}

func GetHLROrder(dtoorder *models.DtoOrder, r render.Render, facilityrepository services.FacilityRepository,
	hlrfacilityrepository services.HLRFacilityRepository, mobileoperatoroperationrepository services.MobileOperatorOperationRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository,
	language string) (apihlrfacility *models.ApiHLRFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_HLR, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}

	found, err := hlrfacilityrepository.Exists(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !found {
		return &models.ApiHLRFacility{}, nil
	}

	dtohlrfacility, err := hlrfacilityrepository.Get(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	operators, err := mobileoperatoroperationrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	resulttables, worktables, err := GetOrderTables(dtoorder.ID, r, resulttablerepository, worktablerepository, language)
	if err != nil {
		return nil, err
	}

	return models.NewApiHLRFacility(dtohlrfacility.EstimatedNumbersShipments, *operators, dtohlrfacility.DeliveryDataId,
		dtohlrfacility.DeliveryDataDelete, dtohlrfacility.MessageToInColumnId, dtohlrfacility.Cost, dtohlrfacility.CostFactual,
		*resulttables, *worktables), nil
}

func UpdateHLROrder(dtoorder *models.DtoOrder, viewhlrfacility models.ViewHLRFacility, r render.Render,
	facilityrepository services.FacilityRepository, hlrfacilityrepository services.HLRFacilityRepository,
	orderstatusrepository services.OrderStatusRepository, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository,
	mobileoperatorrepository services.MobileOperatorRepository, mobileoperatoroperationrepository services.MobileOperatorOperationRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository, checkaccess bool, userid int64,
	language string) (apihlrfacility *models.ApiHLRFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_HLR, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}
	err = CheckOrderFacilityEditability(dtoorder, r, orderstatusrepository, language)
	if err != nil {
		return nil, err
	}

	found, err := hlrfacilityrepository.Exists(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	var dtohlrfacility *models.DtoHLRFacility
	if !found {
		dtohlrfacility = new(models.DtoHLRFacility)
		dtohlrfacility.Order_ID = dtoorder.ID
	} else {
		dtohlrfacility, err = hlrfacilityrepository.Get(dtoorder.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
	}

	if dtoorder.Step == 2 || dtoorder.Step == 3 {
		dtohlrfacility.EstimatedNumbersShipments = viewhlrfacility.EstimatedNumbersShipments

		var totalpercent byte = 0
		var totalrequests int = 0
		for _, apioperator := range viewhlrfacility.EstimatedOperators {
			totalpercent += apioperator.Percent
			totalrequests += apioperator.Count
			_, err = CheckMobileOperator(apioperator.MobileOperator_ID, r, mobileoperatorrepository, language)
			if err != nil {
				return nil, err
			}
			dtohlrfacility.EstimatedOperators = append(dtohlrfacility.EstimatedOperators, *models.NewDtoMobileOperatorOperation(dtoorder.ID,
				apioperator.MobileOperator_ID, apioperator.Percent, apioperator.Count))
		}
		if totalpercent != 100 {
			log.Error("Total percent sum is not equal 100, %v", totalpercent)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return nil, errors.New("Wrong percent")
		}
		if totalrequests != viewhlrfacility.EstimatedNumbersShipments {
			log.Error("Total count sum is not equal estimated number of shipments, %v", totalrequests)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return nil, errors.New("Wrong count")
		}
	}

	if dtoorder.Step == 2 {
	}
	if dtoorder.Step == 3 {
		_, err = IsTableAvailable(r, customertablerepository, viewhlrfacility.DeliveryDataId, language)
		if err != nil {
			return nil, err
		}
		if checkaccess {
			err = IsTableAccessible(viewhlrfacility.DeliveryDataId, userid, r, customertablerepository, language)
			if err != nil {
				return nil, err
			}
		}
		dtohlrfacility.DeliveryDataId = viewhlrfacility.DeliveryDataId
		dtohlrfacility.DeliveryDataDelete = viewhlrfacility.DeliveryDataDelete

		if viewhlrfacility.MessageToInColumnId != 0 {
			_, err = CheckColumnValidity(viewhlrfacility.DeliveryDataId, viewhlrfacility.MessageToInColumnId, r, columntyperepository,
				tablecolumnrepository, language)
			if err != nil {
				return nil, err
			}
		}
		dtohlrfacility.MessageToInColumnId = viewhlrfacility.MessageToInColumnId
	}

	resulttables, worktables, err := GetOrderTables(dtoorder.ID, r, resulttablerepository, worktablerepository, language)
	if err != nil {
		return nil, err
	}
	for _, resulttable := range *resulttables {
		dtohlrfacility.ResultTables = append(dtohlrfacility.ResultTables, *models.NewDtoResultTable(dtoorder.ID, resulttable.Customer_Table_ID))
	}
	for _, worktable := range *worktables {
		dtohlrfacility.WorkTables = append(dtohlrfacility.WorkTables, *models.NewDtoWorkTable(dtoorder.ID, worktable.Customer_Table_ID))
	}

	if !found {
		err = hlrfacilityrepository.Create(dtohlrfacility, false, true)
	} else {
		err = hlrfacilityrepository.Update(dtohlrfacility, false, true)
	}
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	operators, err := mobileoperatoroperationrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return models.NewApiHLRFacility(dtohlrfacility.EstimatedNumbersShipments, *operators, dtohlrfacility.DeliveryDataId,
		dtohlrfacility.DeliveryDataDelete, dtohlrfacility.MessageToInColumnId, dtohlrfacility.Cost, dtohlrfacility.CostFactual,
		*resulttables, *worktables), nil
}

func GetRecognizeData(dtoorder *models.DtoOrder, r render.Render, inputfieldrepository services.InputFieldRepository,
	inputproductrepository services.InputProductRepository, inputfilerepository services.InputFileRepository,
	language string) (inputfields *[]models.ViewApiInputField, inputproducts *[]models.ViewApiInputProduct,
	inputfiles *[]models.ApiInputFile, err error) {
	inputfields, err = inputfieldrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, err
	}

	inputproducts, err = inputproductrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, err
	}

	inputfiles, err = inputfilerepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, err
	}

	return inputfields, inputproducts, inputfiles, nil
}

func GetRecognizeOrder(dtoorder *models.DtoOrder, r render.Render, facilityrepository services.FacilityRepository,
	recognizefacilityrepository services.RecognizeFacilityRepository, inputfieldrepository services.InputFieldRepository,
	inputproductrepository services.InputProductRepository, inputfilerepository services.InputFileRepository,
	supplierrequestrepository services.SupplierRequestRepository, inputftprepository services.InputFtpRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository,
	language string) (apirecognizefacility *models.ApiRecognizeFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_RECOGNIZE, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}

	found, err := recognizefacilityrepository.Exists(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !found {
		return &models.ApiRecognizeFacility{}, nil
	}

	dtorecognizefacility, err := recognizefacilityrepository.Get(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	inputfields, inputproducts, inputfiles, err := GetRecognizeData(dtoorder, r, inputfieldrepository, inputproductrepository, inputfilerepository, language)
	if err != nil {
		return nil, err
	}

	supplierrequests, err := supplierrequestrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	inputftp, err := CheckInputFtp(dtoorder.ID, r, inputftprepository, language)
	if err != nil {
		return nil, err
	}

	resulttables, worktables, err := GetOrderTables(dtoorder.ID, r, resulttablerepository, worktablerepository, language)
	if err != nil {
		return nil, err
	}

	return models.NewApiRecognizeFacility(dtorecognizefacility.EstimatedNumbersForm, dtorecognizefacility.EstimatedCalculationOnFields,
		*inputfields, *inputproducts, dtorecognizefacility.RequiredFields, dtorecognizefacility.LoadDefectiveForms,
		dtorecognizefacility.CommentsForSupplier, *inputfiles, dtorecognizefacility.RequestsSend, dtorecognizefacility.RequestsCancel, *supplierrequests,
		dtorecognizefacility.Cost, dtorecognizefacility.CostFactual, *models.NewApiInputFtp(inputftp.Ready, inputftp.Customer_Table_ID,
			inputftp.Host, inputftp.Port, inputftp.Path, inputftp.Login, inputftp.Password), *resulttables, *worktables), nil
}

func UpdateRecognizeOrder(dtoorder *models.DtoOrder, viewrecognizefacility models.ViewRecognizeFacility, r render.Render,
	facilityrepository services.FacilityRepository, recognizefacilityrepository services.RecognizeFacilityRepository,
	orderstatusrepository services.OrderStatusRepository, columntyperepository services.ColumnTypeRepository,
	recognizeproductrepository services.RecognizeProductRepository, filerepository services.FileRepository,
	inputfieldrepository services.InputFieldRepository, inputproductrepository services.InputProductRepository,
	inputfilerepository services.InputFileRepository, supplierrequestrepository services.SupplierRequestRepository,
	inputftprepository services.InputFtpRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, language string) (apirecognizefacility *models.ApiRecognizeFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_RECOGNIZE, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}
	err = CheckOrderFacilityEditability(dtoorder, r, orderstatusrepository, language)
	if err != nil {
		return nil, err
	}

	found, err := recognizefacilityrepository.Exists(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	var dtorecognizefacility *models.DtoRecognizeFacility
	if !found {
		dtorecognizefacility = new(models.DtoRecognizeFacility)
		dtorecognizefacility.Order_ID = dtoorder.ID
	} else {
		dtorecognizefacility, err = recognizefacilityrepository.Get(dtoorder.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
	}

	if dtoorder.Step == 2 || dtoorder.Step == 3 {
		dtorecognizefacility.EstimatedNumbersForm = viewrecognizefacility.EstimatedNumbersForm
		dtorecognizefacility.EstimatedCalculationOnFields = viewrecognizefacility.EstimatedCalculationOnFields

		for _, inputfield := range viewrecognizefacility.EstimatedFields {
			dtorecognizeproduct, err := CheckRecognizeProduct(inputfield.Product_ID, r, recognizeproductrepository, language)
			if err != nil {
				return nil, err
			}
			if dtorecognizeproduct.Increase {
				log.Error("Can't use a discount position as a field %v", inputfield.Product_ID)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
					Message: config.Localization[language].Errors.Api.Object_NotExist})
				return nil, errors.New("Not a field")
			}

			dtorecognizefacility.EstimatedFields = append(dtorecognizefacility.EstimatedFields, *models.NewDtoInputField(dtoorder.ID,
				inputfield.Product_ID, inputfield.Count))
		}
	}

	if dtoorder.Step == 2 {
		dtoinputproducts, err := inputproductrepository.GetAll(dtoorder.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
		dtorecognizefacility.InputProducts = *dtoinputproducts
		dtoinputfiles, err := inputfilerepository.GetAll(dtoorder.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
		dtorecognizefacility.EstimatedFormFiles = *dtoinputfiles
	}
	if dtoorder.Step == 3 {
		for _, inputproduct := range viewrecognizefacility.InputProducts {
			dtorecognizeproduct, err := CheckRecognizeProduct(inputproduct.Product_ID, r, recognizeproductrepository, language)
			if err != nil {
				return nil, err
			}
			if !dtorecognizeproduct.Increase {
				log.Error("Can't use a field position as a discount %v", inputproduct.Product_ID)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
					Message: config.Localization[language].Errors.Api.Object_NotExist})
				return nil, errors.New("Not a discount")
			}

			dtorecognizefacility.InputProducts = append(dtorecognizefacility.InputProducts, *models.NewDtoInputProduct(dtoorder.ID,
				inputproduct.Product_ID))
		}

		dtorecognizefacility.RequiredFields = viewrecognizefacility.RequiredFields
		dtorecognizefacility.LoadDefectiveForms = viewrecognizefacility.LoadDefectiveForms
		dtorecognizefacility.CommentsForSupplier = viewrecognizefacility.CommentsForSupplier

		for _, inputfile := range viewrecognizefacility.EstimatedFormFiles {
			_, err = filerepository.Get(inputfile.File_ID)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
					Message: config.Localization[language].Errors.Api.Object_NotExist})
				return nil, err
			}
			dtorecognizefacility.EstimatedFormFiles = append(dtorecognizefacility.EstimatedFormFiles, *models.NewDtoInputFile(dtoorder.ID,
				inputfile.File_ID))
		}
	}

	supplierrequests, err := supplierrequestrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	for _, supplierrequest := range *supplierrequests {
		dtorecognizefacility.SupplierRequests = append(dtorecognizefacility.SupplierRequests,
			*models.NewDtoSupplierRequest(dtoorder.ID, supplierrequest.Supplier_ID, supplierrequest.RequestDate, supplierrequest.Responded,
				supplierrequest.RespondedDate, supplierrequest.EstimatedCost, supplierrequest.MyChoice))
	}

	inputftp, err := CheckInputFtp(dtoorder.ID, r, inputftprepository, language)
	if err != nil {
		return nil, err
	}
	dtorecognizefacility.Ftp = *inputftp

	resulttables, worktables, err := GetOrderTables(dtoorder.ID, r, resulttablerepository, worktablerepository, language)
	if err != nil {
		return nil, err
	}
	for _, resulttable := range *resulttables {
		dtorecognizefacility.ResultTables = append(dtorecognizefacility.ResultTables, *models.NewDtoResultTable(dtoorder.ID, resulttable.Customer_Table_ID))
	}
	for _, worktable := range *worktables {
		dtorecognizefacility.WorkTables = append(dtorecognizefacility.WorkTables, *models.NewDtoWorkTable(dtoorder.ID, worktable.Customer_Table_ID))
	}

	if !found {
		err = recognizefacilityrepository.Create(dtorecognizefacility, false, true)
	} else {
		err = recognizefacilityrepository.Update(dtorecognizefacility, false, true)
	}
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	inputfields, inputproducts, inputfiles, err := GetRecognizeData(dtoorder, r, inputfieldrepository, inputproductrepository, inputfilerepository, language)
	if err != nil {
		return nil, err
	}

	return models.NewApiRecognizeFacility(dtorecognizefacility.EstimatedNumbersForm, dtorecognizefacility.EstimatedCalculationOnFields,
		*inputfields, *inputproducts, dtorecognizefacility.RequiredFields, dtorecognizefacility.LoadDefectiveForms,
		dtorecognizefacility.CommentsForSupplier, *inputfiles, dtorecognizefacility.RequestsSend, dtorecognizefacility.RequestsCancel,
		*supplierrequests, dtorecognizefacility.Cost, dtorecognizefacility.CostFactual,
		*models.NewApiInputFtp(inputftp.Ready, inputftp.Customer_Table_ID, inputftp.Host, inputftp.Port, inputftp.Path, inputftp.Login, inputftp.Password),
		*resulttables, *worktables), nil
}

func GetVerifyOrder(dtoorder *models.DtoOrder, r render.Render, facilityrepository services.FacilityRepository,
	verifyfacilityrepository services.VerifyFacilityRepository, dataproductrepository services.DataProductRepository,
	datacolumnrepository services.DataColumnRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, language string) (apiverifyfacility *models.ApiVerifyFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_VERIFY, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}

	found, err := verifyfacilityrepository.Exists(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !found {
		return &models.ApiVerifyFacility{}, nil
	}

	dtoverifyfacility, err := verifyfacilityrepository.Get(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	dataproducts, err := dataproductrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	datacolumns, err := datacolumnrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	resulttables, worktables, err := GetOrderTables(dtoorder.ID, r, resulttablerepository, worktablerepository, language)
	if err != nil {
		return nil, err
	}

	return models.NewApiVerifyFacility(dtoverifyfacility.EstimatedNumbersRecords, *dataproducts, dtoverifyfacility.TablesDataId,
		dtoverifyfacility.TablesDataDelete, *datacolumns, dtoverifyfacility.Cost, dtoverifyfacility.CostFactual,
		*resulttables, *worktables), nil
}

func UpdateVerifyOrder(dtoorder *models.DtoOrder, viewverifyfacility models.ViewVerifyFacility, r render.Render,
	facilityrepository services.FacilityRepository, verifyfacilityrepository services.VerifyFacilityRepository,
	orderstatusrepository services.OrderStatusRepository, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, verifyproductrepository services.VerifyProductRepository,
	tablecolumnrepository services.TableColumnRepository, dataproductrepository services.DataProductRepository,
	datacolumnrepository services.DataColumnRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, checkaccess bool, userid int64,
	language string) (apiverifyfacility *models.ApiVerifyFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_VERIFY, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}
	err = CheckOrderFacilityEditability(dtoorder, r, orderstatusrepository, language)
	if err != nil {
		return nil, err
	}

	found, err := verifyfacilityrepository.Exists(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	var dtoverifyfacility *models.DtoVerifyFacility
	if !found {
		dtoverifyfacility = new(models.DtoVerifyFacility)
		dtoverifyfacility.Order_ID = dtoorder.ID
	} else {
		dtoverifyfacility, err = verifyfacilityrepository.Get(dtoorder.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
	}

	if dtoorder.Step == 2 || dtoorder.Step == 3 {
		dtoverifyfacility.EstimatedNumbersRecords = viewverifyfacility.EstimatedNumbersRecords
		for _, dataproduct := range viewverifyfacility.DataProducts {
			_, err = CheckVerifyProduct(dataproduct.Product_ID, r, verifyproductrepository, language)
			if err != nil {
				return nil, err
			}
			dtoverifyfacility.DataProducts = append(dtoverifyfacility.DataProducts, *models.NewDtoDataProduct(dtoorder.ID,
				dataproduct.Product_ID))
		}
	}

	if dtoorder.Step == 2 {
		dtodatacolumns, err := datacolumnrepository.GetAll(dtoorder.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
		dtoverifyfacility.DataColumns = *dtodatacolumns
	}
	if dtoorder.Step == 3 {
		_, err = IsTableAvailable(r, customertablerepository, viewverifyfacility.TablesDataId, language)
		if err != nil {
			return nil, err
		}
		if checkaccess {
			err = IsTableAccessible(viewverifyfacility.TablesDataId, userid, r, customertablerepository, language)
			if err != nil {
				return nil, err
			}
		}
		dtoverifyfacility.TablesDataId = viewverifyfacility.TablesDataId
		dtoverifyfacility.TablesDataDelete = viewverifyfacility.TablesDataDelete

		for _, viewcolumn := range viewverifyfacility.DataColumns {
			_, err = CheckColumnValidity(viewverifyfacility.TablesDataId, viewcolumn.Table_Column_ID, r, columntyperepository,
				tablecolumnrepository, language)
			if err != nil {
				return nil, err
			}
			dtoverifyfacility.DataColumns = append(dtoverifyfacility.DataColumns, *models.NewDtoDataColumn(dtoorder.ID,
				viewcolumn.Table_Column_ID))
		}
	}

	resulttables, worktables, err := GetOrderTables(dtoorder.ID, r, resulttablerepository, worktablerepository, language)
	if err != nil {
		return nil, err
	}
	for _, resulttable := range *resulttables {
		dtoverifyfacility.ResultTables = append(dtoverifyfacility.ResultTables, *models.NewDtoResultTable(dtoorder.ID, resulttable.Customer_Table_ID))
	}
	for _, worktable := range *worktables {
		dtoverifyfacility.WorkTables = append(dtoverifyfacility.WorkTables, *models.NewDtoWorkTable(dtoorder.ID, worktable.Customer_Table_ID))
	}

	if !found {
		err = verifyfacilityrepository.Create(dtoverifyfacility, false, true)
	} else {
		err = verifyfacilityrepository.Update(dtoverifyfacility, false, true)
	}
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	dataproducts, err := dataproductrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	datacolumns, err := datacolumnrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return models.NewApiVerifyFacility(dtoverifyfacility.EstimatedNumbersRecords, *dataproducts, dtoverifyfacility.TablesDataId,
		dtoverifyfacility.TablesDataDelete, *datacolumns, dtoverifyfacility.Cost, dtoverifyfacility.CostFactual,
		*resulttables, *worktables), nil
}

func GetHeaderOrder(dtoorder *models.DtoOrder, r render.Render, facilityrepository services.FacilityRepository,
	headerfacilityrepository services.HeaderFacilityRepository, language string) (apiheaderfacility *models.ApiHeaderFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_HEADER, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}

	found, err := headerfacilityrepository.Exists(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !found {
		return &models.ApiHeaderFacility{}, nil
	}

	dtoheaderfacility, err := headerfacilityrepository.Get(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return models.NewApiHeaderFacility(dtoheaderfacility.CreateBegin.Format(models.FORMAT_DATE), dtoheaderfacility.CreateEnd.Format(models.FORMAT_DATE),
		dtoheaderfacility.Name, dtoheaderfacility.Begin.Format(models.FORMAT_DATE), dtoheaderfacility.End.Format(models.FORMAT_DATE),
		dtoheaderfacility.AutoRenew, dtoheaderfacility.Cost, dtoheaderfacility.CostFactual), nil
}

func UpdateHeaderOrder(dtoorder *models.DtoOrder, viewheaderfacility models.ViewHeaderFacility, r render.Render,
	facilityrepository services.FacilityRepository, headerfacilityrepository services.HeaderFacilityRepository,
	orderstatusrepository services.OrderStatusRepository, language string) (apiheaderfacility *models.ApiHeaderFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_HEADER, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}
	err = CheckOrderFacilityEditability(dtoorder, r, orderstatusrepository, language)
	if err != nil {
		return nil, err
	}

	found, err := headerfacilityrepository.Exists(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	var dtoheaderfacility *models.DtoHeaderFacility
	if !found {
		dtoheaderfacility = new(models.DtoHeaderFacility)
		dtoheaderfacility.Order_ID = dtoorder.ID
	} else {
		dtoheaderfacility, err = headerfacilityrepository.Get(dtoorder.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
	}

	if dtoorder.Step == 2 || dtoorder.Step == 3 {

	}

	if dtoorder.Step == 2 {
		if viewheaderfacility.CreateBegin != "" {
			dtoheaderfacility.CreateBegin, err = models.ParseDate(viewheaderfacility.CreateBegin)
			if err != nil {
				log.Error("Can't parse create begin date %v", viewheaderfacility.CreateBegin)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
					Message: config.Localization[language].Errors.Api.Object_NotExist})
				return nil, err
			}
		} else {
			dtoheaderfacility.CreateBegin = time.Time{}
		}
		if !dtoheaderfacility.CreateBegin.IsZero() && dtoheaderfacility.CreateBegin.Sub(time.Now()) < 0 {
			log.Error("Create begin date is in the past %v", dtoheaderfacility.CreateBegin)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return nil, errors.New("Wrong begin date")
		}
		if viewheaderfacility.CreateEnd != "" {
			dtoheaderfacility.CreateEnd, err = models.ParseDate(viewheaderfacility.CreateEnd)
			if err != nil {
				log.Error("Can't parse create end date %v", viewheaderfacility.CreateEnd)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
					Message: config.Localization[language].Errors.Api.Object_NotExist})
				return nil, err
			}
		} else {
			dtoheaderfacility.CreateEnd = time.Time{}
		}
		if !dtoheaderfacility.CreateEnd.IsZero() && dtoheaderfacility.CreateEnd.Sub(time.Now()) < 0 {
			log.Error("Create end date is in the past %v", dtoheaderfacility.CreateEnd)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return nil, errors.New("Wrong end date")
		}
		if !dtoheaderfacility.CreateBegin.IsZero() && !dtoheaderfacility.CreateEnd.IsZero() &&
			dtoheaderfacility.CreateBegin.Sub(dtoheaderfacility.CreateEnd) > 0 {
			log.Error("Create begin date can't be bigger than create end date %v", dtoheaderfacility.CreateBegin, dtoheaderfacility.CreateEnd)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return nil, errors.New("Wrong dates")
		}
		dtoheaderfacility.Name = viewheaderfacility.Name
	}
	if dtoorder.Step == 3 {
		dtoheaderfacility.AutoRenew = viewheaderfacility.Renew
	}

	if !found {
		err = headerfacilityrepository.Create(dtoheaderfacility)
	} else {
		err = headerfacilityrepository.Update(dtoheaderfacility)
	}
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	return models.NewApiHeaderFacility(dtoheaderfacility.CreateBegin.Format(models.FORMAT_DATE), dtoheaderfacility.CreateEnd.Format(models.FORMAT_DATE),
		dtoheaderfacility.Name, dtoheaderfacility.Begin.Format(models.FORMAT_DATE), dtoheaderfacility.End.Format(models.FORMAT_DATE),
		dtoheaderfacility.AutoRenew, dtoheaderfacility.Cost, dtoheaderfacility.CostFactual), nil
}
