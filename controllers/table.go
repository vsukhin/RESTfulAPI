package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"types"

	"application/config"
	"application/helpers"
	"application/models"
	"application/server/middlewares"
	"application/services"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

const (
	TIME_ONE_MONTH = 720
)

// get /api/v1.0/tables/types/
func GetTableTypes(w http.ResponseWriter, r render.Render, tabletyperepository services.TableTypeRepository, session *models.DtoSession) {
	tabletypes, err := tabletyperepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(tabletypes, len(*tabletypes), w, r)
}

// options	/api/v1.0/tables/
func GetMetaUnitTables(request *http.Request, r render.Render, customertablerepository services.CustomerTableRepository,
	session *models.DtoSession) {
	var err error
	query := ""

	var filters *[]models.FilterExp
	filters, err = helpers.GetFilterArray(new(models.TableSearch), nil, request, r, session.Language)
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

	customertable, err := customertablerepository.GetMeta(session.UserID, query, middlewares.IsAdmin(session.Roles))
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, customertable)
}

// get /api/v1.0/tables/
func GetUnitTables(w http.ResponseWriter, request *http.Request, r render.Render, customertablerepository services.CustomerTableRepository,
	session *models.DtoSession) {
	var err error
	query := ""

	var filters *[]models.FilterExp
	filters, err = helpers.GetFilterArray(new(models.TableSearch), nil, request, r, session.Language)
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
	sorts, err = helpers.GetOrderArray(new(models.TableSearch), request, r, session.Language)
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

	customertables, err := customertablerepository.GetByUser(session.UserID, query, middlewares.IsAdmin(session.Roles))
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(customertables, len(*customertables), w, r)
}

// get /api/v1.0/tables/:tid/
func GetTable(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository, session *models.DtoSession) {
	dtocustomertable, err := helpers.CheckTable(r, params, customertablerepository, session.Language)
	if err != nil {
		return
	}

	customertablemeta, err := customertablerepository.GetEx(dtocustomertable.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, customertablemeta)
}

// post /api/v1.0/tables/
func CreateTable(errors binding.Errors, viewcustomertable models.ViewShortCustomerTable, r render.Render,
	userrepository services.UserRepository, customertablerepository services.CustomerTableRepository,
	tabletyperepository services.TableTypeRepository, unitrepository services.UnitRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}

	unitid, typeid, author, err := helpers.CheckCustomerTableParameters(r, viewcustomertable.UnitID, models.TABLE_TYPE_DEFAULT, session.UserID, session.Language,
		userrepository, unitrepository, tabletyperepository)
	if err != nil {
		return
	}

	dtocustomertable := new(models.DtoCustomerTable)
	dtocustomertable.Name = viewcustomertable.Name
	dtocustomertable.Created = time.Now()
	dtocustomertable.TypeID = typeid
	dtocustomertable.UnitID = unitid
	dtocustomertable.Active = true
	dtocustomertable.Permanent = true
	dtocustomertable.Signature = author

	err = customertablerepository.Create(dtocustomertable)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongCustomerTable(dtocustomertable.ID, dtocustomertable.Name, dtocustomertable.TypeID, dtocustomertable.UnitID))
}

// put /api/v1.0/tables/:tid/
func UpdateTable(errors binding.Errors, viewcustomertable models.ViewLongCustomerTable, r render.Render, params martini.Params,
	userrepository services.UserRepository, customertablerepository services.CustomerTableRepository, unitrepository services.UnitRepository,
	tabletyperepository services.TableTypeRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	dtocustomertable, err := helpers.CheckTable(r, params, customertablerepository, session.Language)
	if err != nil {
		return
	}

	unitid, typeid, _, err := helpers.CheckCustomerTableParameters(r, viewcustomertable.UnitID, viewcustomertable.Type, session.UserID, session.Language,
		userrepository, unitrepository, tabletyperepository)
	if err != nil {
		return
	}

	if dtocustomertable.TypeID != typeid {
		if viewcustomertable.Type != models.TABLE_TYPE_DEFAULT {
			log.Error("Can change table type to %v", viewcustomertable.Type)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}
	}

	dtocustomertable.Name = viewcustomertable.Name
	dtocustomertable.TypeID = typeid
	dtocustomertable.UnitID = unitid

	err = customertablerepository.Update(dtocustomertable)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortCustomerTable(dtocustomertable.Name, viewcustomertable.Type, dtocustomertable.UnitID))
}

// delete /api/v1.0/tables/:tid/
func DeleteTable(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository, session *models.DtoSession) {
	dtocustomertable, err := helpers.CheckTable(r, params, customertablerepository, session.Language)
	if err != nil {
		return
	}

	err = customertablerepository.Deactivate(dtocustomertable)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// options /api/v1.0/tables/:tid/data/
func GetTableMetaData(request *http.Request, r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	tablecolumnrepository services.TableColumnRepository, session *models.DtoSession) {
	dtocustomertable, err := helpers.CheckTable(r, params, customertablerepository, session.Language)
	if err != nil {
		return
	}

	filters, err := helpers.GetFilterArray(tablecolumnrepository, dtocustomertable.ID, request, r, session.Language)
	if err != nil {
		return
	}

	query := ""

	tablecolumns, err := tablecolumnrepository.GetByTable(dtocustomertable.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	for _, filter := range *filters {
		var exps []string
		for _, field := range filter.Fields {
			found := false
			for _, tablecolumn := range *tablecolumns {
				value, err := strconv.ParseInt(field, 0, 64)
				if err != nil {
					log.Error("Can't convert to number %v filter column with value %v", err, field)
					r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
					return
				}
				if tablecolumn.ID == value {
					found = true
					exps = append(exps, fmt.Sprintf(" field%v ", tablecolumn.FieldNum)+filter.Op+" "+filter.Value)
					break
				}
			}
			if !found {
				log.Error("Filter column %v doesn't belong table %v", field, dtocustomertable.ID)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
		}
		if len(exps) != 0 {
			query += " (" + strings.Join(exps, " or ") + ")" + " and"
		}
	}

	customertablemeta, err := customertablerepository.GetFullMeta(dtocustomertable, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, customertablemeta)
}

// put /api/v1.0/tables/:tid/price/
func UpdatePriceTable(errors binding.Errors, viewpriceproperties models.ViewApiPriceProperties, r render.Render, params martini.Params,
	customertablerepository services.CustomerTableRepository, pricepropertiesrepository services.PricePropertiesRepository,
	facilityrepository services.FacilityRepository, tablecolumnrepository services.TableColumnRepository,
	columntyperepository services.ColumnTypeRepository, tablerowrepository services.TableRowRepository,
	recognizeproductrepository services.RecognizeProductRepository, verifyproductrepository services.VerifyProductRepository,
	headerproductrepository services.HeaderProductRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	dtocustomertable, err := helpers.CheckTable(r, params, customertablerepository, session.Language)
	if err != nil {
		return
	}

	found, err := pricepropertiesrepository.Exists(dtocustomertable.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	published := false
	created := time.Now()
	if found {
		priceproperties, err := pricepropertiesrepository.Get(dtocustomertable.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}
		if priceproperties.Published {
			if priceproperties.Facility_ID != viewpriceproperties.Facility_ID ||
				priceproperties.After_ID != viewpriceproperties.After_ID ||
				((viewpriceproperties.Begin.IsZero() &&
					!(priceproperties.Begin.Year() == 1 && priceproperties.Begin.Month() == 1 && priceproperties.Begin.Day() == 1)) ||
					(!viewpriceproperties.Begin.IsZero() &&
						(priceproperties.Begin.Year() == 1 && priceproperties.Begin.Month() == 1 && priceproperties.Begin.Day() == 1)) ||
					(!viewpriceproperties.Begin.IsZero() &&
						!(priceproperties.Begin.Year() == 1 && priceproperties.Begin.Month() == 1 && priceproperties.Begin.Day() == 1) &&
						viewpriceproperties.Begin.Sub(priceproperties.Begin) != 0)) ||
				(!viewpriceproperties.End.IsZero() && viewpriceproperties.End.Sub(time.Now()).Hours() < TIME_ONE_MONTH) ||
				priceproperties.Published != viewpriceproperties.Published {
				log.Error("Can't change fields for published price list %v", dtocustomertable.ID)
				r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_CHANGES_DENIED,
					Message: config.Localization[session.Language].Errors.Api.Data_Changes_Denied})
				return
			}
		}
		published = priceproperties.Published
		created = priceproperties.Created
	}

	if !published {
		facility, err := facilityrepository.Get(viewpriceproperties.Facility_ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}
		if !facility.Active {
			log.Error("Facility is not active %v", facility.ID)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}

		if viewpriceproperties.After_ID != 0 {
			customertable, err := helpers.IsTableAvailable(r, customertablerepository, viewpriceproperties.After_ID, session.Language)
			if err != nil {
				return
			}
			if customertable.TypeID != models.TABLE_TYPE_PRICE {
				log.Error("After price is not price type %v", customertable.ID)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
					Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
				return
			}

			priceproperties, err := pricepropertiesrepository.Get(customertable.ID)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
					Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
				return
			}
			if priceproperties.Facility_ID != viewpriceproperties.Facility_ID {
				log.Error("Service is not the same for after price %v", customertable.ID)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
					Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
				return
			}
			if !priceproperties.Published {
				log.Error("After price is not published %v", customertable.ID)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
					Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
				return
			}
		}

		if !viewpriceproperties.Begin.IsZero() && viewpriceproperties.Begin.Sub(time.Now()) < 0 {
			log.Error("Begin date is in the past %v", viewpriceproperties.Begin)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}

		if !viewpriceproperties.End.IsZero() && viewpriceproperties.End.Sub(time.Now()) < 0 {
			log.Error("End date is in the past %v", viewpriceproperties.End)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}

		if !viewpriceproperties.Begin.IsZero() && !viewpriceproperties.End.IsZero() &&
			viewpriceproperties.Begin.Sub(viewpriceproperties.End) > 0 {
			log.Error("Begin date can't be bigger than end date %v", viewpriceproperties.Begin, viewpriceproperties.End)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}

		pricecolumns, err := helpers.CheckPriceColumns(dtocustomertable.ID, facility.Alias, r, tablecolumnrepository, columntyperepository,
			session.Language)
		if err != nil {
			return
		}
		if viewpriceproperties.Published && len(*pricecolumns) != 0 {
			log.Error("Can't find required price list columns for table %v", dtocustomertable.ID)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}
		if len(*pricecolumns) != 0 {
			err = tablecolumnrepository.CreateAll(pricecolumns)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
		}

		if viewpriceproperties.Published &&
			(facility.Alias == models.SERVICE_TYPE_RECOGNIZE || facility.Alias == models.SERVICE_TYPE_VERIFY || facility.Alias == models.SERVICE_TYPE_HEADER) {
			products, err := helpers.CheckProductColumns(dtocustomertable.ID, facility.Alias, r, tablecolumnrepository, tablerowrepository,
				recognizeproductrepository, verifyproductrepository, headerproductrepository, session.Language)
			if err != nil {
				return
			}
			uniqueproducts := make(map[string]int)
			for _, product := range *products {
				if product.Product_ID == 0 {
					uniqueproducts[product.Name] = 0
				}
			}
			if facility.Alias == models.SERVICE_TYPE_RECOGNIZE {
				recognizeproducts := new([]models.DtoRecognizeProduct)
				for productname, _ := range uniqueproducts {
					recognizeproduct := new(models.DtoRecognizeProduct)
					recognizeproduct.Name = productname
					recognizeproduct.Created = time.Now()
					recognizeproduct.Active = true
					*recognizeproducts = append(*recognizeproducts, *recognizeproduct)
				}
				if len(*recognizeproducts) != 0 {
					err = recognizeproductrepository.CreateAll(recognizeproducts)
					if err != nil {
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
							Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
						return
					}
				}
				for _, recognizeproduct := range *recognizeproducts {
					uniqueproducts[recognizeproduct.Name] = recognizeproduct.ID
				}
			}
			if facility.Alias == models.SERVICE_TYPE_VERIFY {
				verifyproducts := new([]models.DtoVerifyProduct)
				for productname, _ := range uniqueproducts {
					verifyproduct := new(models.DtoVerifyProduct)
					verifyproduct.Name = productname
					verifyproduct.Created = time.Now()
					verifyproduct.Active = true
					*verifyproducts = append(*verifyproducts, *verifyproduct)
				}
				if len(*verifyproducts) != 0 {
					err = verifyproductrepository.CreateAll(verifyproducts)
					if err != nil {
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
							Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
						return
					}
				}
				for _, verifyproduct := range *verifyproducts {
					uniqueproducts[verifyproduct.Name] = verifyproduct.ID
				}
			}
			if facility.Alias == models.SERVICE_TYPE_HEADER {
				headerproducts := new([]models.DtoHeaderProduct)
				for productname, _ := range uniqueproducts {
					headerproduct := new(models.DtoHeaderProduct)
					headerproduct.Name = productname
					headerproduct.Created = time.Now()
					headerproduct.Active = true
					*headerproducts = append(*headerproducts, *headerproduct)
				}
				if len(*headerproducts) != 0 {
					err = headerproductrepository.CreateAll(headerproducts)
					if err != nil {
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
							Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
						return
					}
				}
				for _, headerproduct := range *headerproducts {
					uniqueproducts[headerproduct.Name] = headerproduct.ID
				}
			}
			for i, product := range *products {
				if product.Product_ID == 0 {
					productid, ok := uniqueproducts[product.Name]
					if !ok {
						log.Error("Can't find product name %v for price list table %v", product.Name, dtocustomertable.ID)
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
							Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
						return
					}
					(*products)[i].Product_ID = productid
				}
			}
			err = helpers.UpdateProductColumn(products, dtocustomertable.ID, r, tablecolumnrepository, columntyperepository, tablerowrepository,
				session.Language)
			if err != nil {
				return
			}
		}
	}

	dtopriceproperties := models.NewDtoPriceProperties(dtocustomertable.ID, viewpriceproperties.Facility_ID,
		viewpriceproperties.After_ID, viewpriceproperties.Begin, viewpriceproperties.End, created, viewpriceproperties.Published)
	if !found {
		err = pricepropertiesrepository.Create(dtopriceproperties, true)
	} else {
		err = pricepropertiesrepository.Update(dtopriceproperties, true)
	}
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, viewpriceproperties)
}

// get /api/v1.0/tables/services/sms/
func GetSMSTables(w http.ResponseWriter, r render.Render, smstablerepository services.SMSTableRepository, session *models.DtoSession) {
	smstables, err := smstablerepository.GetAll(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(smstables, len(*smstables), w, r)
}

// get /api/v1.0/tables/services/hlr/
func GetHLRTables(w http.ResponseWriter, r render.Render, hlrtablerepository services.HLRTableRepository, session *models.DtoSession) {
	hlrtables, err := hlrtablerepository.GetAll(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(hlrtables, len(*hlrtables), w, r)
}

// get /api/v1.0/tables/services/verification/
func GetVerifyTables(w http.ResponseWriter, r render.Render, verifytablerepository services.VerifyTableRepository, session *models.DtoSession) {
	verifytables, err := verifytablerepository.GetAll(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(verifytables, len(*verifytables), w, r)
}
