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
	PARAM_NAME_TABLE_ID = "tid"
)

func CheckCustomerTableParameters(r render.Render, unitparam int64, typeparam string, userid int64,
	language string, userservice *services.UserService, unitservice *services.UnitService,
	tabletypeservice *services.TableTypeService) (unitid int64, typeid int64, err error) {
	if unitparam == 0 {
		var user *models.DtoUser
		user, err = userservice.Get(userid)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return 0, 0, err
		}
		unitid = user.UnitID
	} else {
		var unit *models.DtoUnit
		unit, err = unitservice.Get(unitparam)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return 0, 0, err
		}
		unitid = unit.ID
	}

	typeid, err = tabletypeservice.FindByName(typeparam)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return 0, 0, err
	}

	return unitid, typeid, nil
}

func IsTableActive(r render.Render, customertableservice *services.CustomerTableService, tableid int64,
	language string) (dtocustomertable *models.DtoCustomerTable, err error) {
	dtocustomertable, err = customertableservice.Get(tableid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, errors.New("Table is not found")
	}

	if !dtocustomertable.Active {
		log.Error("Customer table is not active %v", tableid)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Table is not active")
	}

	return dtocustomertable, nil
}

func IsTableAvailable(r render.Render, customertableservice *services.CustomerTableService, tableid int64,
	language string) (dtocustomertable *models.DtoCustomerTable, err error) {
	dtocustomertable, err = IsTableActive(r, customertableservice, tableid, language)
	if err != nil {
		return nil, err
	}
	if !dtocustomertable.Permanent {
		log.Error("Customer table is not permanent %v", tableid)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Table is not permanent")
	}

	return dtocustomertable, nil
}

func CheckTable(r render.Render, params martini.Params, customertableservice *services.CustomerTableService,
	language string) (dtocustomertable *models.DtoCustomerTable, err error) {
	var tableid int64
	tableid, err = CheckParameterInt(r, params[PARAM_NAME_TABLE_ID], language)
	if err != nil {
		return nil, err
	}
	dtocustomertable, err = IsTableAvailable(r, customertableservice, tableid, language)
	if err != nil {
		return nil, err
	}

	return dtocustomertable, nil
}