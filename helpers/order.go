package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"types"
	"unicode"
)

const (
	PARAM_NAME_ORDER_ID = "oid"
	SMS_LENGTH_ASCII    = 153
	SMS_LENGTH_UICODE   = 67
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

func CheckOrderEditability(orderid int64, r render.Render, orderstatusrepository services.OrderStatusRepository,
	language string) (err error) {
	dtoorderstatuses, err := orderstatusrepository.GetByOrder(orderid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}
	for _, orderstatus := range *dtoorderstatuses {
		if orderstatus.Status_ID == models.ORDER_STATUS_COMPLETED && orderstatus.Value == true {
			log.Error("Can't update completed order %v", orderid)
			r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_CHANGES_DENIED,
				Message: config.Localization[language].Errors.Api.Data_Changes_Denied})
			return errors.New("Not editable order")
		}
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

	err = orderrepository.Update(dtoorder, dtoorderstatuses, true)
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

func GetSMSOrder(dtoorder *models.DtoOrder, r render.Render, facilityrepository services.FacilityRepository,
	smsfacilityrepository services.SMSFacilityRepository, mobileoperatoroperationrepository services.MobileOperatorOperationRepository,
	smsperiodrepository services.SMSPeriodRepository, smseventrepository services.SMSEventRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository,
	language string) (apismsfacility *models.ApiSMSFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_SMS, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}

	dtosmsfacility, err := smsfacilityrepository.Get(dtoorder.ID)
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

	periods, err := smsperiodrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	events, err := smseventrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	resulttables, worktables, err := GetOrderTables(dtoorder.ID, r, resulttablerepository, worktablerepository, language)
	if err != nil {
		return nil, err
	}

	deliveryType := ""
	switch dtosmsfacility.DeliveryType {
	case models.TYPE_DELIVERY_ONCE:
		deliveryType = models.TYPE_DELIVERY_ONCE_VALUE
	case models.TYPE_DELIVERY_SCHEDULED:
		deliveryType = models.TYPE_DELIVERY_SCHEDULED_VALUE
	case models.TYPE_DELIVERY_EVENTTRIGGERED:
		deliveryType = models.TYPE_DELIVERY_EVENTTRIGGERED_VALUE
	default:
		log.Error("Unknown delivery type %v", dtosmsfacility.DeliveryType)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, errors.New("Wrong delivery type")
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
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository,
	checkaccess bool, userid int64, language string) (apismsfacility *models.ApiSMSFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_SMS, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}
	err = CheckOrderEditability(dtoorder.ID, r, orderstatusrepository, language)
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

	return models.NewApiSMSFacility(viewsmsfacility.EstimatedNumbersShipments, viewsmsfacility.EstimatedMessageInCyrillic,
		viewsmsfacility.EstimatedNumberCharacters, viewsmsfacility.EstimatedNumberSmsInMessage, viewsmsfacility.EstimatedOperators,
		viewsmsfacility.DeliveryType, viewsmsfacility.DeliveryTime, viewsmsfacility.Periods, viewsmsfacility.Events,
		viewsmsfacility.DeliveryTimeStart, viewsmsfacility.DeliveryTimeEnd, viewsmsfacility.DeliveryBaseTime, viewsmsfacility.DeliveryDataId,
		viewsmsfacility.DeliveryDataDelete, viewsmsfacility.MessageFromId, viewsmsfacility.MessageFromInColumnId, viewsmsfacility.MessageToInColumnId,
		viewsmsfacility.MessageBody, viewsmsfacility.MessageBodyInColumnId, viewsmsfacility.TimeCorrection, dtosmsfacility.Cost, dtosmsfacility.CostFactual,
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
	mobileoperatorrepository services.MobileOperatorRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, checkaccess bool, userid int64,
	language string) (apihlrfacility *models.ApiHLRFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_HLR, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}
	err = CheckOrderEditability(dtoorder.ID, r, orderstatusrepository, language)
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
		err = hlrfacilityrepository.Create(dtohlrfacility, true)
	} else {
		err = hlrfacilityrepository.Update(dtohlrfacility, true)
	}
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	return models.NewApiHLRFacility(viewhlrfacility.EstimatedNumbersShipments, viewhlrfacility.EstimatedOperators, viewhlrfacility.DeliveryDataId,
		viewhlrfacility.DeliveryDataDelete, viewhlrfacility.MessageToInColumnId, dtohlrfacility.Cost, dtohlrfacility.CostFactual,
		*resulttables, *worktables), nil
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

	dtorecognizefacility, err := recognizefacilityrepository.Get(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	inputfields, err := inputfieldrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	inputproducts, err := inputproductrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	inputfiles, err := inputfilerepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
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
		*inputfields, *inputproducts, dtorecognizefacility.PriceIncreaseUrgent, dtorecognizefacility.PriceIncreaseNano,
		dtorecognizefacility.PriceIncreaseBackgroundBlack, dtorecognizefacility.RequiredFields, dtorecognizefacility.LoadDefectiveForms,
		dtorecognizefacility.CommentsForSupplier, *inputfiles, dtorecognizefacility.RequestsSend, dtorecognizefacility.RequestsCancel, *supplierrequests,
		dtorecognizefacility.Cost, dtorecognizefacility.CostFactual, *models.NewApiInputFtp(inputftp.Ready, inputftp.Customer_Table_ID,
			inputftp.Host, inputftp.Port, inputftp.Path, inputftp.Login, inputftp.Password), *resulttables, *worktables), nil
}

func UpdateRecognizeOrder(dtoorder *models.DtoOrder, viewrecognizefacility models.ViewRecognizeFacility, r render.Render,
	facilityrepository services.FacilityRepository, recognizefacilityrepository services.RecognizeFacilityRepository,
	orderstatusrepository services.OrderStatusRepository, columntyperepository services.ColumnTypeRepository,
	recognizeproductrepository services.RecognizeProductRepository, filerepository services.FileRepository,
	inputfilerepository services.InputFileRepository, supplierrequestrepository services.SupplierRequestRepository,
	inputftprepository services.InputFtpRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, language string) (apirecognizefacility *models.ApiRecognizeFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_RECOGNIZE, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}
	err = CheckOrderEditability(dtoorder.ID, r, orderstatusrepository, language)
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

	dtorecognizefacility.EstimatedNumbersForm = viewrecognizefacility.EstimatedNumbersForm
	dtorecognizefacility.EstimatedCalculationOnFields = viewrecognizefacility.EstimatedCalculationOnFields

	for _, inputfield := range viewrecognizefacility.EstimatedFields {
		err = IsColumnTypeActive(r, columntyperepository, inputfield.Column_Type_ID, language)
		if err != nil {
			return nil, err
		}
		dtorecognizefacility.EstimatedFields = append(dtorecognizefacility.EstimatedFields, *models.NewDtoInputField(dtoorder.ID,
			inputfield.Column_Type_ID, inputfield.Count))
	}

	for _, inputproduct := range viewrecognizefacility.InputProducts {
		_, err = CheckRecognizeProduct(inputproduct.Product_ID, r, recognizeproductrepository, language)
		if err != nil {
			return nil, err
		}
		dtorecognizefacility.InputProducts = append(dtorecognizefacility.InputProducts, *models.NewDtoInputProduct(dtoorder.ID,
			inputproduct.Product_ID))
	}

	dtorecognizefacility.PriceIncreaseUrgent = viewrecognizefacility.PriceIncreaseUrgent
	dtorecognizefacility.PriceIncreaseNano = viewrecognizefacility.PriceIncreaseNano
	dtorecognizefacility.PriceIncreaseBackgroundBlack = viewrecognizefacility.PriceIncreaseBackgroundBlack
	dtorecognizefacility.RequiredFields = viewrecognizefacility.RequiredFields
	dtorecognizefacility.LoadDefectiveForms = viewrecognizefacility.LoadDefectiveForms
	dtorecognizefacility.CommentsForSupplier = viewrecognizefacility.CommentsForSupplier

	for _, inputfile := range viewrecognizefacility.EstimatedFromFiles {
		_, err = filerepository.Get(inputfile.File_ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
		dtorecognizefacility.EstimatedFromFiles = append(dtorecognizefacility.EstimatedFromFiles, *models.NewDtoInputFile(dtoorder.ID,
			inputfile.File_ID))
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
		err = recognizefacilityrepository.Create(dtorecognizefacility, true)
	} else {
		err = recognizefacilityrepository.Update(dtorecognizefacility, true)
	}
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	inputfiles, err := inputfilerepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return models.NewApiRecognizeFacility(viewrecognizefacility.EstimatedNumbersForm, viewrecognizefacility.EstimatedCalculationOnFields,
		viewrecognizefacility.EstimatedFields, viewrecognizefacility.InputProducts, viewrecognizefacility.PriceIncreaseUrgent,
		viewrecognizefacility.PriceIncreaseNano, viewrecognizefacility.PriceIncreaseBackgroundBlack, viewrecognizefacility.RequiredFields,
		viewrecognizefacility.LoadDefectiveForms, viewrecognizefacility.CommentsForSupplier, *inputfiles, dtorecognizefacility.RequestsSend,
		dtorecognizefacility.RequestsCancel, *supplierrequests, dtorecognizefacility.Cost, dtorecognizefacility.CostFactual,
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
	tablecolumnrepository services.TableColumnRepository, datacolumnrepository services.DataColumnRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository, checkaccess bool, userid int64,
	language string) (apiverifyfacility *models.ApiVerifyFacility, err error) {
	err = CheckFacilityAlias(dtoorder.Facility_ID, models.SERVICE_TYPE_VERIFY, r, facilityrepository, language)
	if err != nil {
		return nil, err
	}
	err = CheckOrderEditability(dtoorder.ID, r, orderstatusrepository, language)
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

	dtoverifyfacility.EstimatedNumbersRecords = viewverifyfacility.EstimatedNumbersRecords
	for _, dataproduct := range viewverifyfacility.DataProducts {
		_, err = CheckVerifyProduct(dataproduct.Product_ID, r, verifyproductrepository, language)
		if err != nil {
			return nil, err
		}
		dtoverifyfacility.DataProducts = append(dtoverifyfacility.DataProducts, *models.NewDtoDataProduct(dtoorder.ID,
			dataproduct.Product_ID))
	}

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
		err = verifyfacilityrepository.Create(dtoverifyfacility, true)
	} else {
		err = verifyfacilityrepository.Update(dtoverifyfacility, true)
	}
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	datacolumns, err := datacolumnrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return models.NewApiVerifyFacility(viewverifyfacility.EstimatedNumbersRecords, viewverifyfacility.DataProducts, viewverifyfacility.TablesDataId,
		viewverifyfacility.TablesDataDelete, *datacolumns, dtoverifyfacility.Cost, dtoverifyfacility.CostFactual,
		*resulttables, *worktables), nil
}

func CalculateSMSQuantity(sms string) (count int) {
	count = 0
	runes := []rune(sms)
	isASCII := true
	for _, current_rune := range runes {
		if current_rune > unicode.MaxASCII {
			isASCII = false
			break
		}
	}
	length := 0
	if isASCII {
		length = SMS_LENGTH_ASCII
		count = len(runes) / length

	} else {
		length = SMS_LENGTH_UICODE
		count = len(runes) / length
	}
	if math.Mod(float64(len(runes)), float64(length)) != 0 {
		count++
	}

	return count
}

func ExecuteSMSOrder(order_id int64, orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	smsfacilityrepository services.SMSFacilityRepository, orderstatusrepository services.OrderStatusRepository,
	customertablerepository services.CustomerTableRepository, smstablerepository services.SMSTableRepository,
	smssenderrepository services.SMSSenderRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, invoicerepository services.InvoiceRepository,
	companyrepository services.CompanyRepository, operationrepository services.OperationRepository,
	transactiontyperepository services.TransactionTypeRepository, tablecolumnrepository services.TableColumnRepository,
	unitrepository services.UnitRepository, tablerowrepository services.TableRowRepository, pricerepository services.PriceRepository,
	mobileoperatorrepository services.MobileOperatorRepository) {
	dtoorder, err := orderrepository.Get(order_id)
	if err != nil {
		return
	}

	dtofacility, err := facilityrepository.Get(dtoorder.Facility_ID)
	if err != nil {
		return
	}
	if dtofacility.Alias != models.SERVICE_TYPE_SMS {
		log.Error("Order service is not macthed to the service method %v", dtoorder.Facility_ID)
		return
	}

	dtosmsfacility, err := smsfacilityrepository.Get(dtoorder.ID)
	if err != nil {
		return
	}

	dtoorderstatuses, err := orderstatusrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		return
	}
	order := models.NewApiLongOrderFromDto(dtoorder, dtoorderstatuses)

	if order.IsAssembled {
		dtoorderstatus := models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_MODERATOR_CONFIRMED, true, "", time.Now())
		err = orderstatusrepository.Save(dtoorderstatus, nil)
		if err != nil {
			return
		}

		if dtosmsfacility.DeliveryType == models.TYPE_DELIVERY_ONCE {
			schhour, schmin, schsec := dtosmsfacility.DeliveryBaseTime.Local().Clock()
			hour, min, sec := time.Now().Clock()
			if schhour == hour && schmin == min && schsec == sec {
				dtodatatable, err := customertablerepository.Get(dtosmsfacility.DeliveryDataId)
				if err != nil {
					return
				}
				smstable, err := smstablerepository.Get(dtodatatable.ID)
				if err != nil {
					return
				}
				dtoorderstatus = models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_OPEN, true, "", time.Now())
				err = orderstatusrepository.Save(dtoorderstatus, nil)
				if err != nil {
					return
				}

				if dtosmsfacility.MessageFromInColumnId != 0 {
					found := false
					for _, smssender := range smstable.SMSSenders {
						if smssender.ID == dtosmsfacility.MessageFromInColumnId {
							found = true
							break
						}
					}
					if !found {
						log.Error("Can find sms sender column %v in table %v", dtosmsfacility.MessageFromInColumnId, dtosmsfacility.DeliveryDataId)
						return
					}
					dtotablecolumn, err := tablecolumnrepository.Get(dtosmsfacility.MessageFromInColumnId)
					if err != nil {
						return
					}
					found, err = smssenderrepository.Belongs(dtotablecolumn, dtoorder.Unit_ID)
					if err != nil {
						return
					}
					if !found {
						log.Error("SMS sender column %v data doesn't belong to unit %v", dtotablecolumn.ID, dtoorder.Unit_ID)
						return
					}
				} else {
					dtosmsender, err := smssenderrepository.Get(dtosmsfacility.MessageFromInColumnId)
					if err != nil {
						return
					}
					if dtosmsender.Unit_ID != dtoorder.Unit_ID {
						log.Error("SMS sender %v unit doesn't match order %v unit", dtosmsender.ID, dtoorder.ID)
						return
					}
				}

				var columnmobilephone_id int64
				if dtosmsfacility.MessageToInColumnId != 0 {
					found := false
					for _, mobilephone := range smstable.MobilePhones {
						if mobilephone.ID == dtosmsfacility.MessageToInColumnId {
							found = true
							columnmobilephone_id = dtosmsfacility.MessageToInColumnId
							break
						}
					}
					if !found {
						log.Error("Can find mobile phone column %v in table %v", dtosmsfacility.MessageToInColumnId, dtosmsfacility.DeliveryDataId)
						return
					}
				} else {
					if len(smstable.MobilePhones) == 0 {
						log.Error("Can find mobile column in table %v", dtosmsfacility.DeliveryDataId)
						return
					} else {
						columnmobilephone_id = smstable.MobilePhones[0].ID
					}
				}

				if dtosmsfacility.MessageBodyInColumnId != 0 {
					found := false
					for _, message := range smstable.Messages {
						if message.ID == dtosmsfacility.MessageBodyInColumnId {
							found = true
							break
						}
					}
					if !found {
						log.Error("Can find message column %v in table %v", dtosmsfacility.MessageBodyInColumnId, dtosmsfacility.DeliveryDataId)
						return
					}
				}
				dtoorderstatus = models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_SUPPLIER_COST_NEW, true, "", time.Now())
				err = orderstatusrepository.Save(dtoorderstatus, nil)
				if err != nil {
					return
				}

				tablecolumns := []models.DtoTableColumn{}
				var columnmessage, columnmobilephone *models.DtoTableColumn
				if dtosmsfacility.MessageBodyInColumnId != 0 {
					columnmessage, err = tablecolumnrepository.Get(dtosmsfacility.MessageBodyInColumnId)
					if err != nil {
						return
					}
					tablecolumns = append(tablecolumns, *columnmessage)
				}
				columnmobilephone, err = tablecolumnrepository.Get(columnmobilephone_id)
				if err != nil {
					return
				}
				tablecolumns = append(tablecolumns, *columnmobilephone)

				apitablerows, err := tablerowrepository.GetAll("", "", dtosmsfacility.DeliveryDataId, &tablecolumns)
				if err != nil {
					return
				}
				mobileoperatorsmses := make(map[int]int)
				for _, apitablerow := range *apitablerows {
					smscount := 0
					mobileoperator_id := 0
					for _, apitablecell := range apitablerow.Cells {
						if apitablecell.Table_Column_ID == columnmessage.ID {
							if dtosmsfacility.MessageBodyInColumnId != 0 {
								smscount = CalculateSMSQuantity(apitablecell.Value)
							} else {
								smscount = CalculateSMSQuantity(dtosmsfacility.MessageBody)
							}
						}
						if apitablecell.Table_Column_ID == columnmobilephone_id {
							mobileoperator_id, _ = strconv.Atoi(apitablecell.Value)
						}
					}
					mobileoperatorsmses[mobileoperator_id] += smscount
				}

				smshlrprices, err := GetSMSHLRPricesInternal(models.SERVICE_TYPE_SMS, pricerepository, tablecolumnrepository, tablerowrepository,
					mobileoperatorrepository)
				if err != nil {
					return
				}
				dtosmsfacility.Cost = 0
				for _, smshlrprice := range *smshlrprices {
					count := mobileoperatorsmses[smshlrprice.Mobile_Operator_ID]
					if smshlrprice.Supplier_ID == dtoorder.Supplier_ID && count != 0 &&
						count >= smshlrprice.AmountRange.Begin && (count <= smshlrprice.AmountRange.End || smshlrprice.AmountRange.End == 0) {
						dtosmsfacility.Cost += float64(count) * smshlrprice.Price
						delete(mobileoperatorsmses, smshlrprice.Mobile_Operator_ID)
					}
				}
				if len(mobileoperatorsmses) != 0 {
					log.Error("Can't find price list for supplier%v", dtoorder.Supplier_ID)
					return
				}

				dtosmsfacility.CostFactual = dtosmsfacility.Cost
				err = smsfacilityrepository.Update(dtosmsfacility, true, false)
				if err != nil {
					return
				}
				dtoorderstatus = models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED, true, "", time.Now())
				err = orderstatusrepository.Save(dtoorderstatus, nil)
				if err != nil {
					return
				}
				dtoorderstatus = models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_MODERATOR_BEGIN, true, "", time.Now())
				err = orderstatusrepository.Save(dtoorderstatus, nil)
				if err != nil {
					return
				}

				dtocompany, err := companyrepository.GetPrimaryByUnit(dtoorder.Unit_ID)
				if err != nil {
					return
				}

				balance, err := operationrepository.CalculateBalance(dtoorder.Unit_ID)
				if err != nil {
					return
				}
				dtoorderstatus = models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_MODERATOR_BEGIN, true, "", time.Now())
				err = orderstatusrepository.Save(dtoorderstatus, nil)
				if err != nil {
					return
				}

				if balance < dtosmsfacility.Cost {
					log.Error("Not enough money at unit balance %v to pay for order %v execution", dtoorder.Unit_ID, dtoorder.ID)
					return
				}

				dtounit, err := unitrepository.Get(config.Configuration.SystemAccount)
				if err != nil {
					return
				}
				dtotransactiontype, err := transactiontyperepository.Get(models.TRANSACTION_TYPE_SERVICE_FEE_SMS)
				if err != nil {
					return
				}

				dtotransaction := new(models.DtoTransaction)
				dtotransaction.Source_ID = dtoorder.Unit_ID
				dtotransaction.Destination_ID = dtounit.ID
				dtotransaction.Type_ID = dtotransactiontype.ID

				dtoinvoice := new(models.DtoInvoice)
				dtoinvoice.Company_ID = dtocompany.ID
				if dtocompany.VAT != 0 {
					dtoinvoice.VAT = (dtosmsfacility.Cost / (1 + float64(dtocompany.VAT)/100)) * float64(dtocompany.VAT) / 100
				}
				dtoinvoice.Total = dtosmsfacility.Cost
				dtoinvoice.Paid = true
				dtoinvoice.Created = time.Now()
				dtoinvoice.Active = true
				dtoinvoice.InvoiceItems = []models.DtoInvoiceItem{*models.NewDtoInvoiceItem(0, 0, dtotransactiontype.Name, models.INVOICE_ITEM_TYPE_ROUBLE,
					1, dtoinvoice.Total, dtoinvoice.Total)}

				err = invoicerepository.PaidForOrder(dtoorder.ID, dtoinvoice, dtotransaction, true)
				if err != nil {
					return
				}

				dtoworkdatatable, err := customertablerepository.Copy(dtodatatable, true)
				if err != nil {
					return
				}
				dtodatatable.TypeID = models.TABLE_TYPE_READONLY
				err = customertablerepository.Update(dtodatatable)
				if err != nil {
					return
				}
				dtoworkdatatable.TypeID = models.TABLE_TYPE_HIDDEN
				err = customertablerepository.Update(dtoworkdatatable)
				if err != nil {
					return
				}

				dtoworktable := models.NewDtoWorkTable(dtoorder.ID, dtoworkdatatable.ID)
				err = worktablerepository.Create(dtoworktable, nil)
				if err != nil {
					return
				}

				position, err := tablecolumnrepository.GetDefaultPosition(dtoworkdatatable.ID)
				if err != nil {
					return
				}
				position++
				fieldnum, err := FindFreeColumnInternal(dtoworkdatatable.ID, 0, tablecolumnrepository)
				if err != nil {
					return
				}

				dtotablecolumn := new(models.DtoTableColumn)
				dtotablecolumn.Created = time.Now()
				dtotablecolumn.Position = position
				dtotablecolumn.Name = "Status"
				dtotablecolumn.Customer_Table_ID = dtoworkdatatable.ID
				dtotablecolumn.Column_Type_ID = models.COLUMN_TYPE_DEFAULT
				dtotablecolumn.Prebuilt = true
				dtotablecolumn.FieldNum = fieldnum
				dtotablecolumn.Active = true
				dtotablecolumn.Edition = 0

				err = tablecolumnrepository.Create(dtotablecolumn, nil)
				if err != nil {
					return
				}

				tablecolumns = []models.DtoTableColumn{*dtotablecolumn}
				apitablerows, err = tablerowrepository.GetAll("", "", dtoworkdatatable.ID, &tablecolumns)
				if err != nil {
					return
				}

				workdatatablecolumns, err := tablecolumnrepository.GetByTable(dtoworkdatatable.ID)
				if err != nil {
					return
				}

				for _, apitablerow := range *apitablerows {
					tablerow, err := tablerowrepository.Get(apitablerow.ID)
					if err != nil {
						return
					}

					tablecells, err := tablerow.TableRowToDtoTableCells(workdatatablecolumns)
					if err != nil {
						return
					}

					for i, _ := range *tablecells {
						for j, _ := range tablecolumns {
							if (*tablecells)[i].Table_Column_ID == tablecolumns[j].ID {
								(*tablecells)[i].Value = strconv.Itoa(i * j)
								(*tablecells)[i].Checked = true
								(*tablecells)[i].Valid = true
							}
						}
					}

					err = tablerow.TableCellsToTableRow(tablecells, workdatatablecolumns)
					if err != nil {
						return
					}

					err = tablerowrepository.Update(tablerow, nil, true, false)
					if err != nil {
						return
					}
				}

				if dtosmsfacility.DeliveryDataDelete {
					err = customertablerepository.Deactivate(dtodatatable)
					if err != nil {
						return
					}
				} else {
					dtodatatable.TypeID = models.TABLE_TYPE_DEFAULT
					err = customertablerepository.Update(dtodatatable)
					if err != nil {
						return
					}
				}

				worktables, err := worktablerepository.GetByOrder(dtoorder.ID)
				if err != nil {
					return
				}
				for _, worktable := range *worktables {
					dtocustomertable, err := customertablerepository.Get(worktable.Customer_Table_ID)
					if err != nil {
						return
					}
					err = customertablerepository.Deactivate(dtocustomertable)
					if err != nil {
						return
					}
				}

				dtoorderstatus = models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_SUPPLIER_CLOSE, true, "", time.Now())
				err = orderstatusrepository.Save(dtoorderstatus, nil)
				if err != nil {
					return
				}

			}
		}
	}

}
