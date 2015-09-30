package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"time"
	"types"
)

const (
	PARAM_NAME_TABLE_ID = "tid"
)

func CheckCustomerTableParameters(r render.Render, unitparam int64, typeparam int, userid int64,
	language string, userrepository services.UserRepository, unitrepository services.UnitRepository,
	tabletyperepository services.TableTypeRepository) (unitid int64, typeid int, author string, err error) {
	user, err := userrepository.Get(userid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return 0, 0, "", err
	}
	author = user.Name + " " + user.MiddleName + " " + user.Surname
	if unitparam == 0 {
		unitid = user.UnitID
	} else {
		unit, err := unitrepository.Get(unitparam)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return 0, 0, "", err
		}
		unitid = unit.ID
	}

	dtotype, err := tabletyperepository.Get(typeparam)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return 0, 0, "", err
	}
	typeid = dtotype.ID

	return unitid, typeid, author, nil
}

func IsTableActive(r render.Render, customertablerepository services.CustomerTableRepository, tableid int64,
	language string) (dtocustomertable *models.DtoCustomerTable, err error) {
	dtocustomertable, err = customertablerepository.Get(tableid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Table is not found")
	}

	if !dtocustomertable.Active {
		log.Error("Customer table is not active %v", tableid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Table is not active")
	}

	return dtocustomertable, nil
}

func IsTableAvailable(r render.Render, customertablerepository services.CustomerTableRepository, tableid int64,
	language string) (dtocustomertable *models.DtoCustomerTable, err error) {
	dtocustomertable, err = IsTableActive(r, customertablerepository, tableid, language)
	if err != nil {
		return nil, err
	}
	if !dtocustomertable.Permanent {
		log.Error("Customer table is not permanent %v", tableid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Table is not permanent")
	}

	return dtocustomertable, nil
}

func CheckTable(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	language string) (dtocustomertable *models.DtoCustomerTable, err error) {
	var tableid int64
	tableid, err = CheckParameterInt(r, params[PARAM_NAME_TABLE_ID], language)
	if err != nil {
		return nil, err
	}
	dtocustomertable, err = IsTableAvailable(r, customertablerepository, tableid, language)
	if err != nil {
		return nil, err
	}

	return dtocustomertable, nil
}

func IsTableAccessible(table_id int64, user_id int64, r render.Render, customertablerepository services.CustomerTableRepository,
	language string) (err error) {
	allowed, err := customertablerepository.CheckUserAccess(user_id, table_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}
	if !allowed {
		log.Error("Table %v is not accessible for user  %v", table_id, user_id)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return errors.New("Not accessible table")
	}

	return nil
}

func CheckPriceColumns(customer_table_id int64, alias string, r render.Render, tablecolumnrepository services.TableColumnRepository,
	columntyperepository services.ColumnTypeRepository, language string) (pricecolumns *[]models.DtoTableColumn, err error) {
	pricecolumns = new([]models.DtoTableColumn)
	dtotablecolumns, err := tablecolumnrepository.GetByTable(customer_table_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	var columnmobileoperator = false
	var columnprice = false
	var columnproduct = false
	for _, tablecolumn := range *dtotablecolumns {
		if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_MOBILEOPERATOR {
			columnmobileoperator = true
		}
		if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_PRICE {
			columnprice = true
		}
		if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_ID {
			columnproduct = true
		}
	}
	fieldnum, err := FindFreeColumn(customer_table_id, 0, r, tablecolumnrepository, language)
	if err != nil {
		return nil, err
	}
	var position int64 = 0
	if len(*dtotablecolumns) != 0 {
		position, err = tablecolumnrepository.GetDefaultPosition(customer_table_id)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
		position++
	}
	for _, column_type_id := range []int{models.COLUMN_TYPE_PRICELIST_MOBILEOPERATOR, models.COLUMN_TYPE_PRICELIST_ID, models.COLUMN_TYPE_PRICELIST_PRICE} {
		dtotablecolumn := new(models.DtoTableColumn)
		dtotablecolumn.Created = time.Now()
		dtotablecolumn.Active = true
		dtotablecolumn.Edition = 0
		dtotablecolumn.Customer_Table_ID = customer_table_id
		dtotablecolumn.Prebuilt = true
		if alias == models.SERVICE_TYPE_SMS || alias == models.SERVICE_TYPE_HLR {
			if column_type_id == models.COLUMN_TYPE_PRICELIST_MOBILEOPERATOR && !columnmobileoperator {
				dtotablecolumn.Name = "Default mobile operator"
				err = IsColumnTypeActive(r, columntyperepository, column_type_id, language)
				if err != nil {
					return nil, err
				}
				dtotablecolumn.Column_Type_ID = column_type_id
				if fieldnum > models.MAX_COLUMN_NUMBER {
					log.Error("Can't find free column for table %v", customer_table_id)
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
					return nil, errors.New("No free column")
				}
				dtotablecolumn.FieldNum = fieldnum
				fieldnum, err = FindFreeColumn(customer_table_id, fieldnum, r, tablecolumnrepository, language)
				if err != nil {
					return nil, err
				}
				dtotablecolumn.Position = position
				position++
				*pricecolumns = append(*pricecolumns, *dtotablecolumn)
			}
		}
		if alias == models.SERVICE_TYPE_RECOGNIZE || alias == models.SERVICE_TYPE_VERIFY {
			if column_type_id == models.COLUMN_TYPE_PRICELIST_ID && !columnproduct {
				dtotablecolumn.Name = "Default product id"
				err = IsColumnTypeActive(r, columntyperepository, column_type_id, language)
				if err != nil {
					return nil, err
				}
				dtotablecolumn.Column_Type_ID = column_type_id
				if fieldnum > models.MAX_COLUMN_NUMBER {
					log.Error("Can't find free column for table %v", customer_table_id)
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
					return nil, errors.New("No free column")
				}
				dtotablecolumn.FieldNum = fieldnum
				fieldnum, err = FindFreeColumn(customer_table_id, fieldnum, r, tablecolumnrepository, language)
				if err != nil {
					return nil, err
				}
				dtotablecolumn.Position = position
				position++
				*pricecolumns = append(*pricecolumns, *dtotablecolumn)
			}
		}
		if column_type_id == models.COLUMN_TYPE_PRICELIST_PRICE && !columnprice {
			dtotablecolumn.Name = "Default price"
			err = IsColumnTypeActive(r, columntyperepository, column_type_id, language)
			if err != nil {
				return nil, err
			}
			dtotablecolumn.Column_Type_ID = column_type_id
			if fieldnum > models.MAX_COLUMN_NUMBER {
				log.Error("Can't find free column for table %v", customer_table_id)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
					Message: config.Localization[language].Errors.Api.Object_NotExist})
				return nil, errors.New("No free column")
			}
			dtotablecolumn.FieldNum = fieldnum
			fieldnum, err = FindFreeColumn(customer_table_id, fieldnum, r, tablecolumnrepository, language)
			if err != nil {
				return nil, err
			}
			dtotablecolumn.Position = position
			position++
			*pricecolumns = append(*pricecolumns, *dtotablecolumn)
		}
	}

	return pricecolumns, nil
}

func CheckProductColumns(customer_table_id int64, alias string, r render.Render, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository, recognizeproductrepository services.RecognizeProductRepository,
	verifyproductrepository services.VerifyProductRepository, headerproductrepository services.HeaderProductRepository,
	language string) (products *[]models.ApiProduct, err error) {
	products = new([]models.ApiProduct)
	if alias != models.SERVICE_TYPE_RECOGNIZE && alias != models.SERVICE_TYPE_VERIFY && alias != models.SERVICE_TYPE_HEADER {
		return products, nil
	}
	dtotablecolumns, err := tablecolumnrepository.GetByTable(customer_table_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	var columnproduct *models.DtoTableColumn
	var columnname *models.DtoTableColumn
	for i, tablecolumn := range *dtotablecolumns {
		if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_ID {
			if columnproduct != nil {
				log.Error("Can't have multiple product id column in price list %v for service %v",
					customer_table_id, alias)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
					Message: config.Localization[language].Errors.Api.Object_NotExist})
				return nil, errors.New("Multiple product id column")
			}
			columnproduct = &(*dtotablecolumns)[i]
		}
		if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_NAME {
			if columnname != nil {
				log.Error("Can't have multiple name column in price list %v for service %v",
					customer_table_id, alias)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
					Message: config.Localization[language].Errors.Api.Object_NotExist})
				return nil, errors.New("Multiple name column")
			}
			columnname = &(*dtotablecolumns)[i]
		}
	}
	if columnproduct == nil {
		log.Error("Can't find product id column in price list %v for service %v", customer_table_id, alias)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Missed product id column")
	}
	pricecolumns := &[]models.DtoTableColumn{*columnproduct}
	if columnname != nil {
		*pricecolumns = append(*pricecolumns, *columnname)
	}
	apitablerows, err := tablerowrepository.GetAll("", "", customer_table_id, pricecolumns)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	for _, apitablerow := range *apitablerows {
		apiproduct := new(models.ApiProduct)
		apiproduct.Table_Row_ID = apitablerow.ID
		apiproduct.Table_Column_ID = columnproduct.ID
		for _, apitablecell := range apitablerow.Cells {
			if apitablecell.Table_Column_ID == columnproduct.ID {
				value, err := CheckColumnProduct(&apitablecell, r, language)
				if err != nil {
					return nil, err
				}
				apiproduct.Product_ID = value
			}
			if columnname != nil {
				if apitablecell.Table_Column_ID == columnname.ID {
					apiproduct.Name, err = CheckColumnName(&apitablecell, r, language)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		if apiproduct.Name != "" {
			if alias == models.SERVICE_TYPE_RECOGNIZE {
				found, err := recognizeproductrepository.Exists(apiproduct.Name)
				if err != nil {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
					return nil, err
				}
				if !found {
					apiproduct.Product_ID = 0
				} else {
					recognizeproduct, err := recognizeproductrepository.FindByName(apiproduct.Name)
					if err != nil {
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
							Message: config.Localization[language].Errors.Api.Object_NotExist})
						return nil, err
					}
					if apiproduct.Product_ID == recognizeproduct.ID {
						continue
					}
					apiproduct.Product_ID = recognizeproduct.ID
				}
			}
			if alias == models.SERVICE_TYPE_VERIFY {
				found, err := verifyproductrepository.Exists(apiproduct.Name)
				if err != nil {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
					return nil, err
				}
				if !found {
					apiproduct.Product_ID = 0
				} else {
					verifyproduct, err := verifyproductrepository.FindByName(apiproduct.Name)
					if err != nil {
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
							Message: config.Localization[language].Errors.Api.Object_NotExist})
						return nil, err
					}
					if apiproduct.Product_ID == verifyproduct.ID {
						continue
					}
					apiproduct.Product_ID = verifyproduct.ID
				}
			}
			if alias == models.SERVICE_TYPE_HEADER {
				found, err := headerproductrepository.Exists(apiproduct.Name)
				if err != nil {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
					return nil, err
				}
				if !found {
					apiproduct.Product_ID = 0
				} else {
					headerproduct, err := headerproductrepository.FindByName(apiproduct.Name)
					if err != nil {
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
							Message: config.Localization[language].Errors.Api.Object_NotExist})
						return nil, err
					}
					if apiproduct.Product_ID == headerproduct.ID {
						continue
					}
					apiproduct.Product_ID = headerproduct.ID
				}
			}
		} else {
			if alias == models.SERVICE_TYPE_RECOGNIZE {
				_, err = CheckRecognizeProduct(apiproduct.Product_ID, r, recognizeproductrepository, language)
				if err != nil {
					return nil, err
				}
			}
			if alias == models.SERVICE_TYPE_VERIFY {
				_, err = CheckVerifyProduct(apiproduct.Product_ID, r, verifyproductrepository, language)
				if err != nil {
					return nil, err
				}
			}
			if alias == models.SERVICE_TYPE_HEADER {
				_, err = CheckHeaderProduct(apiproduct.Product_ID, r, headerproductrepository, language)
				if err != nil {
					return nil, err
				}
			}
			continue
		}
		*products = append(*products, *apiproduct)
	}

	return products, nil
}
