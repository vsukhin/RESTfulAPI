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

// get /api/v1.0/suppliers/services/
func GetSupplierFacilities(w http.ResponseWriter, r render.Render, facilityrepository services.FacilityRepository, session *models.DtoSession) {
	facilities, err := facilityrepository.GetByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(facilities, len(*facilities), w, r)
}

// put /api/v1.0/suppliers/services/
func UpdateSupplierFacilities(errors binding.Errors, viewfacilities models.ViewFacilities, w http.ResponseWriter, r render.Render,
	facilityrepository services.FacilityRepository, supplierfacilityrepository services.SupplierFacilityRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewfacilities, errors, r, session.Language) != nil {
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

	err := supplierfacilityrepository.SetArrayByUser(session.UserID, facilities, true)
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

	helpers.RenderJSONArray(apifacilities, len(*apifacilities), w, r)
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
func GetOrders(w http.ResponseWriter, request *http.Request, r render.Render, orderrepository services.OrderRepository, session *models.DtoSession) {
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

	helpers.RenderJSONArray(orders, len(*orders), w, r)
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
	if helpers.CheckValidation(&vieworder, errors, r, session.Language) != nil {
		return
	}
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}

	apiorder, err := helpers.UpdateLongOrder(dtoorder, &vieworder, r, params, orderrepository, unitrepository, facilityrepository, session.Language)
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
	err = orderstatusrepository.Save(orderstatus, nil)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// get /api/v1.0/suppliers/orders/:oid:/service/sms/
func GetSMSOrder(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	facilityrepository services.FacilityRepository, smsfacilityrepository services.SMSFacilityRepository,
	mobileoperatoroperationrepository services.MobileOperatorOperationRepository, smsperiodrepository services.SMSPeriodRepository,
	smseventrepository services.SMSEventRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}
	apismsfacility, err := helpers.GetSMSOrder(dtoorder, r, facilityrepository, smsfacilityrepository, mobileoperatoroperationrepository,
		smsperiodrepository, smseventrepository, resulttablerepository, worktablerepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apismsfacility)
}

// put /api/v1.0/suppliers/orders/:oid:/service/sms/
func UpdateSMSOrder(errors binding.Errors, viewsmsfacility models.ViewSMSFacility, r render.Render, params martini.Params,
	orderrepository services.OrderRepository, facilityrepository services.FacilityRepository, smsfacilityrepository services.SMSFacilityRepository,
	orderstatusrepository services.OrderStatusRepository, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository,
	smssenderrepository services.SMSSenderRepository, mobileoperatorrepository services.MobileOperatorRepository,
	periodrepository services.PeriodRepository, eventrepository services.EventRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewsmsfacility, errors, r, session.Language) != nil {
		return
	}
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}
	apismsfacility, err := helpers.UpdateSMSOrder(dtoorder, viewsmsfacility, r, facilityrepository, smsfacilityrepository,
		orderstatusrepository, customertablerepository, columntyperepository, tablecolumnrepository, smssenderrepository,
		mobileoperatorrepository, periodrepository, eventrepository, resulttablerepository, worktablerepository,
		false, session.UserID, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apismsfacility)
}

// get /api/v1.0/suppliers/orders/:oid:/service/hlr/
func GetHLROrder(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	facilityrepository services.FacilityRepository, hlrfacilityrepository services.HLRFacilityRepository,
	mobileoperatoroperationrepository services.MobileOperatorOperationRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository,
	session *models.DtoSession) {
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}
	apihlrfacility, err := helpers.GetHLROrder(dtoorder, r, facilityrepository, hlrfacilityrepository, mobileoperatoroperationrepository,
		resulttablerepository, worktablerepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apihlrfacility)
}

// put /api/v1.0/suppliers/orders/:oid:/service/hlr/
func UpdateHLROrder(errors binding.Errors, viewhlrfacility models.ViewHLRFacility, r render.Render, params martini.Params,
	orderrepository services.OrderRepository, facilityrepository services.FacilityRepository, hlrfacilityrepository services.HLRFacilityRepository,
	orderstatusrepository services.OrderStatusRepository, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository,
	mobileoperatorrepository services.MobileOperatorRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewhlrfacility, errors, r, session.Language) != nil {
		return
	}
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}
	apihlrfacility, err := helpers.UpdateHLROrder(dtoorder, viewhlrfacility, r, facilityrepository, hlrfacilityrepository,
		orderstatusrepository, customertablerepository, columntyperepository, tablecolumnrepository, mobileoperatorrepository,
		resulttablerepository, worktablerepository, false, session.UserID, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apihlrfacility)
}

// get /api/v1.0/suppliers/orders/:oid:/service/recognize/
func GetRecognizeOrder(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	facilityrepository services.FacilityRepository, recognizefacilityrepository services.RecognizeFacilityRepository,
	inputproductrepository services.InputProductRepository, inputfieldrepository services.InputFieldRepository,
	inputfilerepository services.InputFileRepository, supplierrequestrepository services.SupplierRequestRepository,
	inputftprepository services.InputFtpRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}
	apirecognizefacility, err := helpers.GetRecognizeOrder(dtoorder, r, facilityrepository, recognizefacilityrepository, inputfieldrepository,
		inputproductrepository, inputfilerepository, supplierrequestrepository, inputftprepository, resulttablerepository, worktablerepository,
		session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apirecognizefacility)
}

// put /api/v1.0/suppliers/orders/:oid:/service/recognize/
func UpdateRecognizeOrder(errors binding.Errors, viewrecognizefacility models.ViewRecognizeFacility, r render.Render, params martini.Params,
	orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	recognizefacilityrepository services.RecognizeFacilityRepository, orderstatusrepository services.OrderStatusRepository,
	columntyperepository services.ColumnTypeRepository, recognizeproductrepository services.RecognizeProductRepository,
	filerepository services.FileRepository, inputfilerepository services.InputFileRepository, supplierrequestrepository services.SupplierRequestRepository,
	inputftprepository services.InputFtpRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewrecognizefacility, errors, r, session.Language) != nil {
		return
	}
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}
	apirecognizefacility, err := helpers.UpdateRecognizeOrder(dtoorder, viewrecognizefacility, r, facilityrepository, recognizefacilityrepository,
		orderstatusrepository, columntyperepository, recognizeproductrepository, filerepository, inputfilerepository, supplierrequestrepository,
		inputftprepository, resulttablerepository, worktablerepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apirecognizefacility)
}

// get /api/v1.0/suppliers/orders/:oid:/service/verification/
func GetVerifyOrder(r render.Render, params martini.Params, orderrepository services.OrderRepository,
	facilityrepository services.FacilityRepository, verifyfacilityrepository services.VerifyFacilityRepository,
	dataproductrepository services.DataProductRepository, datacolumnrepository services.DataColumnRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}
	apiverifyfacility, err := helpers.GetVerifyOrder(dtoorder, r, facilityrepository, verifyfacilityrepository, dataproductrepository, datacolumnrepository,
		resulttablerepository, worktablerepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apiverifyfacility)
}

// put /api/v1.0/suppliers/orders/:oid:/service/verification/
func UpdateVerifyOrder(errors binding.Errors, viewverifyfacility models.ViewVerifyFacility, r render.Render, params martini.Params,
	orderrepository services.OrderRepository, facilityrepository services.FacilityRepository, verifyfacilityrepository services.VerifyFacilityRepository,
	orderstatusrepository services.OrderStatusRepository, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, verifyproductrepository services.VerifyProductRepository,
	tablecolumnrepository services.TableColumnRepository, datacolumnrepository services.DataColumnRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&verifyfacilityrepository, errors, r, session.Language) != nil {
		return
	}
	dtoorder, err := helpers.CheckOrder(r, params, orderrepository, session.Language)
	if err != nil {
		return
	}
	apiverifyfacility, err := helpers.UpdateVerifyOrder(dtoorder, viewverifyfacility, r, facilityrepository, verifyfacilityrepository,
		orderstatusrepository, customertablerepository, columntyperepository, verifyproductrepository, tablecolumnrepository, datacolumnrepository,
		resulttablerepository, worktablerepository, false, session.UserID, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apiverifyfacility)
}
