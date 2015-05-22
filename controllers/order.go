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
	"time"
	"types"
)

// options /api/v1.0/projects/:prid/orders/
func GetMetaProjectOrders(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	orderrepository services.OrderRepository, session *models.DtoSession) {
	dtoproject, err := helpers.CheckProject(r, params, projectrepository, session.Language)
	if err != nil {
		return
	}

	order, err := orderrepository.GetMetaByProject(dtoproject.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, order)
}

// get /api/v1.0/projects/:prid/orders/
func GetProjectOrders(w http.ResponseWriter, r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	orderrepository services.OrderRepository, session *models.DtoSession) {
	dtoproject, err := helpers.CheckProject(r, params, projectrepository, session.Language)
	if err != nil {
		return
	}

	orders, err := orderrepository.GetByProject(dtoproject.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(orders, len(*orders), w, r)
}

// post /api/v1.0/projects/:prid/orders/
func CreateProjectOrder(errors binding.Errors, vieworder models.ViewShortOrder, r render.Render, params martini.Params,
	projectrepository services.ProjectRepository, orderrepository services.OrderRepository, unitrepository services.UnitRepository,
	session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	dtoproject, err := helpers.CheckProject(r, params, projectrepository, session.Language)
	if err != nil {
		return
	}
	unit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	dtoorder := new(models.DtoOrder)
	dtoorder.Project_ID = dtoproject.ID
	dtoorder.Creator_ID = session.UserID
	dtoorder.Unit_ID = unit.ID
	dtoorder.Name = vieworder.Name
	dtoorder.Step = 0
	dtoorder.Created = time.Now()

	order := new(models.ViewLongOrder)
	dtoorderstatuses := order.ToOrderStatuses(dtoorder.ID)
	err = orderrepository.Create(dtoorder, dtoorderstatuses, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongOrderFromDto(dtoorder, dtoorderstatuses))
}

// get /api/v1.0/projects/:prid/orders/:oid/
func GetProjectOrder(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	orderrepository services.OrderRepository, orderstatusrepository services.OrderStatusRepository, session *models.DtoSession) {
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
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

// patch /api/v1.0/projects/:prid/orders/:oid/
func UpdateProjectOrder(errors binding.Errors, vieworder models.ViewMiddleOrder, r render.Render, params martini.Params,
	projectrepository services.ProjectRepository, orderrepository services.OrderRepository, unitrepository services.UnitRepository,
	facilityrepository services.FacilityRepository, orderstatusrepository services.OrderStatusRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
	if err != nil {
		return
	}
	if vieworder.Step > models.MAX_STEP_NUMBER {
		log.Error("Order step number is too big %v", vieworder.Step)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	var unitid int64 = 0
	var facilityid int64 = 0
	if vieworder.Facility_ID != 0 {
		err = helpers.CheckFacility(vieworder.Facility_ID, r, facilityrepository, session.Language)
		if err != nil {
			return
		}
		facilityid = vieworder.Facility_ID
	}
	if vieworder.Supplier_ID != 0 {
		err = helpers.CheckUnitValidity(vieworder.Supplier_ID, session.Language, r, unitrepository)
		if err != nil {
			return
		}
		unitid = vieworder.Supplier_ID
	}
	dtoorder.Supplier_ID = unitid
	dtoorder.Facility_ID = facilityid
	dtoorder.Name = vieworder.Name
	dtoorder.Step = vieworder.Step
	dtoorderstatuses := &[]models.DtoOrderStatus{
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_COMPLETED, Value: vieworder.IsAssembled, Created: time.Now()},
		{Order_ID: dtoorder.ID, Status_ID: models.ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED, Value: vieworder.IsNewCostConfirmed, Created: time.Now()},
	}

	err = orderrepository.Update(dtoorder, dtoorderstatuses, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	dtoorderstatuses, err = orderstatusrepository.GetByOrder(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongOrderFromDto(dtoorder, dtoorderstatuses))
}

// delete /api/v1.0/projects/:prid/orders/:oid/
func DeleteProjectOrder(r render.Render, params martini.Params, projectrepository services.ProjectRepository,
	orderrepository services.OrderRepository, orderstatusrepository services.OrderStatusRepository,
	session *models.DtoSession) {
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
	if err != nil {
		return
	}
	confirmed, err := orderrepository.IsConfirmed(dtoorder.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if confirmed {
		log.Error("Can't delete confirmed order %v", dtoorder.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	orderstatus := models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_DEL, true, "", time.Now())
	err = orderstatusrepository.Save(orderstatus, nil)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// get /api/v1.0/projects/:prid/orders/:oid/service/sms/
func GetProjectSMSOrder(r render.Render, params martini.Params, projectrepository services.ProjectRepository, orderrepository services.OrderRepository,
	facilityrepository services.FacilityRepository, smsfacilityrepository services.SMSFacilityRepository,
	mobileoperatoroperationrepository services.MobileOperatorOperationRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
	if err != nil {
		return
	}

	apismsfacility, err := helpers.GetSMSOrder(dtoorder, r, facilityrepository, smsfacilityrepository, mobileoperatoroperationrepository,
		resulttablerepository, worktablerepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apismsfacility)
}

// put /api/v1.0/projects/:prid/orders/:oid/service/sms/
func UpdateProjectSMSOrder(errors binding.Errors, viewsmsfacility models.ViewSMSFacility, r render.Render, params martini.Params,
	projectrepository services.ProjectRepository, orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	smsfacilityrepository services.SMSFacilityRepository, orderstatusrepository services.OrderStatusRepository,
	customertablerepository services.CustomerTableRepository, columntyperepository services.ColumnTypeRepository,
	tablecolumnrepository services.TableColumnRepository, smssenderrepository services.SMSSenderRepository,
	mobileoperatorrepository services.MobileOperatorRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
	if err != nil {
		return
	}

	apismsfacility, err := helpers.UpdateSMSOrder(dtoorder, viewsmsfacility, r, facilityrepository, smsfacilityrepository,
		orderstatusrepository, customertablerepository, columntyperepository, tablecolumnrepository, smssenderrepository,
		mobileoperatorrepository, resulttablerepository, worktablerepository, true, session.UserID, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apismsfacility)
}

// get /api/v1.0/projects/:prid/orders/:oid/service/hlr/
func GetProjectHLROrder(r render.Render, params martini.Params, projectrepository services.ProjectRepository, orderrepository services.OrderRepository,
	facilityrepository services.FacilityRepository, hlrfacilityrepository services.HLRFacilityRepository,
	mobileoperatoroperationrepository services.MobileOperatorOperationRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
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

// put /api/v1.0/projects/:prid/orders/:oid/service/hlr/
func UpdateProjectHLROrder(errors binding.Errors, viewhlrfacility models.ViewHLRFacility, r render.Render, params martini.Params,
	projectrepository services.ProjectRepository, orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	hlrfacilityrepository services.HLRFacilityRepository, orderstatusrepository services.OrderStatusRepository,
	customertablerepository services.CustomerTableRepository, columntyperepository services.ColumnTypeRepository,
	tablecolumnrepository services.TableColumnRepository, mobileoperatorrepository services.MobileOperatorRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
	if err != nil {
		return
	}

	apihlrfacility, err := helpers.UpdateHLROrder(dtoorder, viewhlrfacility, r, facilityrepository, hlrfacilityrepository,
		orderstatusrepository, customertablerepository, columntyperepository, tablecolumnrepository, mobileoperatorrepository,
		resulttablerepository, worktablerepository, true, session.UserID, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apihlrfacility)
}

// get /api/v1.0/projects/:prid/orders/:oid/service/recognize/
func GetProjectRecognizeOrder(r render.Render, params martini.Params, projectrepository services.ProjectRepository, orderrepository services.OrderRepository,
	facilityrepository services.FacilityRepository, recognizefacilityrepository services.RecognizeFacilityRepository,
	inputfieldrepository services.InputFieldRepository, inputfilerepository services.InputFileRepository,
	supplierrequestrepository services.SupplierRequestRepository, inputftprepository services.InputFtpRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
	if err != nil {
		return
	}

	apirecognizefacility, err := helpers.GetRecognizeOrder(dtoorder, r, facilityrepository, recognizefacilityrepository, inputfieldrepository,
		inputfilerepository, supplierrequestrepository, inputftprepository, resulttablerepository, worktablerepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apirecognizefacility)
}

// put /api/v1.0/projects/:prid/orders/:oid/service/recognize/
func UpdateProjectRecognizeOrder(errors binding.Errors, viewrecognizefacility models.ViewRecognizeFacility, r render.Render, params martini.Params,
	projectrepository services.ProjectRepository, orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	recognizefacilityrepository services.RecognizeFacilityRepository, orderstatusrepository services.OrderStatusRepository,
	columntyperepository services.ColumnTypeRepository, filerepository services.FileRepository, inputfilerepository services.InputFileRepository,
	supplierrequestrepository services.SupplierRequestRepository, inputftprepository services.InputFtpRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
	if err != nil {
		return
	}

	apirecognizefacility, err := helpers.UpdateRecognizeOrder(dtoorder, viewrecognizefacility, r, facilityrepository, recognizefacilityrepository,
		orderstatusrepository, columntyperepository, filerepository, inputfilerepository, supplierrequestrepository, inputftprepository,
		resulttablerepository, worktablerepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apirecognizefacility)
}

// get /api/v1.0/projects/:prid/orders/:oid/service/verification/
func GetProjectVerifyOrder(r render.Render, params martini.Params, projectrepository services.ProjectRepository, orderrepository services.OrderRepository,
	facilityrepository services.FacilityRepository, verifyfacilityrepository services.VerifyFacilityRepository,
	datacolumnrepository services.DataColumnRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
	if err != nil {
		return
	}

	apiverifyfacility, err := helpers.GetVerifyOrder(dtoorder, r, facilityrepository, verifyfacilityrepository, datacolumnrepository,
		resulttablerepository, worktablerepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apiverifyfacility)
}

// put /api/v1.0/projects/:prid/orders/:oid/service/verification/
func UpdateProjectVerifyOrder(errors binding.Errors, viewverifyfacility models.ViewVerifyFacility, r render.Render, params martini.Params,
	projectrepository services.ProjectRepository, orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	verifyfacilityrepository services.VerifyFacilityRepository, orderstatusrepository services.OrderStatusRepository,
	customertablerepository services.CustomerTableRepository, columntyperepository services.ColumnTypeRepository,
	tablecolumnrepository services.TableColumnRepository, datacolumnrepository services.DataColumnRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	_, dtoorder, err := helpers.CheckProjectOrder(r, params, projectrepository, orderrepository, session.Language)
	if err != nil {
		return
	}

	apiverifyfacility, err := helpers.UpdateVerifyOrder(dtoorder, viewverifyfacility, r, facilityrepository, verifyfacilityrepository,
		orderstatusrepository, customertablerepository, columntyperepository, tablecolumnrepository, datacolumnrepository,
		resulttablerepository, worktablerepository, true, session.UserID, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, apiverifyfacility)
}
