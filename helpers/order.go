package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
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

func CheckFacility(facilityid int64, r render.Render, facilityrepository services.FacilityRepository,
	language string) (err error) {
	dtofacility, err := facilityrepository.Get(facilityid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}
	if !dtofacility.Active {
		log.Error("Service is not active %v", dtofacility.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return errors.New("Service not active")
	}

	return nil
}

func UpdateOrder(dtoorder *models.DtoOrder, vieworder *models.ViewLongOrder, r render.Render, params martini.Params,
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

	apilongorder, err := UpdateOrder(dtoorder, &vieworder.ViewLongOrder, r, params, orderrepository, unitrepository, facilityrepository, language)
	if err != nil {
		return nil, err
	}

	return models.NewApiFullOrder(dtoorder.Creator_ID, dtoorder.Unit_ID, dtoorder.Created, *apilongorder), nil
}
