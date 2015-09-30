package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
	"strings"
	"types"
)

func CheckColumnName(apitablecell *models.ApiLongTableCell, r render.Render, language string) (name string, err error) {
	name = ""
	if !apitablecell.Valid {
		log.Error("Name column value is not valid")
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return "", errors.New("Not valid name")
	}
	name = apitablecell.Value
	if name == "" || len([]rune(name)) > PARAM_LENGTH_MAX {
		log.Error("Wrong name length %v", len([]rune(name)))
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return "", errors.New("Name length is wrong")
	}

	return name, nil
}

func CheckColumnMobileOperatorInternal(apitablecell *models.ApiLongTableCell) (mobile_operator_id int, err error) {
	mobile_operator_id = 0
	if !apitablecell.Valid {
		log.Error("Mobile operator column value is not valid")
		return 0, errors.New("Not valid mobile operator")
	}
	mobile_operator_id, err = strconv.Atoi(apitablecell.Value)
	if err != nil {
		log.Error("Can't convert to number mobile operator %v with value %v", err, apitablecell.Value)
		return 0, errors.New("Mobile operator value is wrong")
	}

	return mobile_operator_id, nil
}

func CheckColumnMobileOperator(apitablecell *models.ApiLongTableCell, r render.Render, language string) (mobile_operator_id int, err error) {
	mobile_operator_id, err = CheckColumnMobileOperatorInternal(apitablecell)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return 0, err
	}

	return mobile_operator_id, nil
}

func CheckColumnProductInternal(apitablecell *models.ApiLongTableCell) (product_id int, err error) {
	product_id = 0
	if !apitablecell.Valid {
		log.Error("Product value is not valid")
		return 0, errors.New("Not valid verify")
	}
	product_id, err = strconv.Atoi(apitablecell.Value)
	if err != nil {
		log.Error("Can't convert to number product %v with value %v", err, apitablecell.Value)
		return 0, errors.New("Product value is wrong")
	}

	return product_id, nil
}

func CheckColumnProduct(apitablecell *models.ApiLongTableCell, r render.Render, language string) (product_id int, err error) {
	product_id, err = CheckColumnProductInternal(apitablecell)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return 0, err
	}

	return product_id, nil
}

func CheckColumnDiscountInternal(apitablecell *models.ApiLongTableCell) (discount float64, err error) {
	discount = 0
	if !apitablecell.Valid {
		log.Error("Discount column value is not valid")
		return 0, errors.New("Not valid discount")
	}
	discount, err = strconv.ParseFloat(strings.Replace(apitablecell.Value, ",", ".", -1), 64)
	if err != nil {
		log.Error("Can't convert to number discount %v with value %v", err, apitablecell.Value)
		return 0, errors.New("Discount value is wrong")
	}
	if discount < 0 {
		log.Error("Discount value can't be negative %v", discount)
		return 0, errors.New("Discount value is negative")
	}

	return discount, nil
}

func CheckColumnDiscount(apitablecell *models.ApiLongTableCell, r render.Render, language string) (discount float64, err error) {
	discount, err = CheckColumnDiscountInternal(apitablecell)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return 0, err
	}

	return discount, nil
}

func CheckColumnRangeInternal(apitablecell *models.ApiLongTableCell) (apirange *models.ApiRange, err error) {
	apirange = new(models.ApiRange)
	if !apitablecell.Valid {
		log.Error("Range column value is not valid")
		return nil, errors.New("Not valid range")
	}
	values := strings.Split(apitablecell.Value, "-")
	if len(values) != 2 {
		log.Error("Wrong range format %v", apitablecell.Value)
		return nil, errors.New("Wrong range format")
	}
	var begin float64 = 0
	if values[0] != "" {
		begin, err = strconv.ParseFloat(strings.Replace(values[0], ",", ".", -1), 64)
		if err != nil {
			log.Error("Can't convert to number range begin %v with value %v", err, values[0])
			return nil, errors.New("Range begin value is wrong")
		}
	}
	var end float64 = 0
	if values[1] != "" {
		end, err = strconv.ParseFloat(strings.Replace(values[1], ",", ".", -1), 64)
		if err != nil {
			log.Error("Can't convert to number range end %v with value %v", err, values[1])
			return nil, errors.New("Range end value is wrong")
		}
	}
	if begin < 0 {
		log.Error("Range begin can't be negative %v", begin)
		return nil, errors.New("Range begin value is negative")
	}
	if end < 0 {
		log.Error("Range end can't be negative %v", end)
		return nil, errors.New("Range end value is negative")
	}
	if begin != 0 && end != 0 {
		if begin > end {
			log.Error("Range begin %v can't be bigger than range end %v", begin, end)
			return nil, errors.New("Wrong range begin and end")
		}
	}
	apirange.Begin = int(begin)
	apirange.End = int(end)

	return apirange, nil
}

func CheckColumnRange(apitablecell *models.ApiLongTableCell, r render.Render, language string) (apirange *models.ApiRange, err error) {
	apirange, err = CheckColumnRangeInternal(apitablecell)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	return apirange, nil
}

func CheckColumnPriceInternal(apitablecell *models.ApiLongTableCell) (price float64, err error) {
	price = 0
	if !apitablecell.Valid {
		log.Error("Price column value is not valid")
		return 0, errors.New("Not valid price")
	}
	price, err = strconv.ParseFloat(strings.Replace(apitablecell.Value, ",", ".", -1), 64)
	if err != nil {
		log.Error("Can't convert to number price %v with value %v", err, apitablecell.Value)
		return 0, errors.New("Price value is wrong")
	}
	if price < 0 {
		log.Error("Price value can't be negative %v", price)
		return 0, errors.New("Price value is negative")
	}

	return price, nil
}

func CheckColumnPrice(apitablecell *models.ApiLongTableCell, r render.Render, language string) (price float64, err error) {
	price, err = CheckColumnPriceInternal(apitablecell)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return 0, err
	}

	return price, nil
}

func GetSMSHLRPriceColumns(supplierprice *models.ApiSupplierPrice, tablecolumnrepository services.TableColumnRepository) (columnmobileoperator,
	columnrange, columnprice *models.DtoTableColumn, err error) {
	columnmobileoperator = nil
	columnrange = nil
	columnprice = nil
	dtotablecolumns, err := tablecolumnrepository.GetByTable(supplierprice.Customer_Table_ID)
	if err != nil {
		return nil, nil, nil, err
	}
	for i, tablecolumn := range *dtotablecolumns {
		if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_MOBILEOPERATOR {
			if columnmobileoperator != nil {
				log.Error("Can't have multiple mobile operator column in price list %v", supplierprice.Customer_Table_ID)
				return nil, nil, nil, errors.New("Multiple mobile operator column")
			}
			columnmobileoperator = &(*dtotablecolumns)[i]
		}
		if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_RANGE {
			if columnrange != nil {
				log.Error("Can't have multiple range column in price list %v", supplierprice.Customer_Table_ID)
				return nil, nil, nil, errors.New("Multiple range column")
			}
			columnrange = &(*dtotablecolumns)[i]
		}
		if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_PRICE {
			if columnprice != nil {
				log.Error("Can't have multiple price column in price list %v", supplierprice.Customer_Table_ID)
				return nil, nil, nil, errors.New("Multiple price column")
			}
			columnprice = &(*dtotablecolumns)[i]
		}
	}
	if columnmobileoperator == nil {
		log.Error("Can't find mobile operator column in price list %v", supplierprice.Customer_Table_ID)
		return nil, nil, nil, errors.New("Missed mobile opertor column")
	}
	if columnprice == nil {
		log.Error("Can't find price column in price list %v", supplierprice.Customer_Table_ID)
		return nil, nil, nil, errors.New("Missed price column")
	}

	return columnmobileoperator, columnrange, columnprice, nil
}

func GetSMSHLRPriceRows(columnmobileoperator, columnrange, columnprice *models.DtoTableColumn, supplierprice *models.ApiSupplierPrice,
	tablerowrepository services.TableRowRepository, mobileoperatorrepository services.MobileOperatorRepository,
	smshlrprices *[]models.ApiSMSHLRPrice) (err error) {
	pricecolumns := &[]models.DtoTableColumn{*columnmobileoperator, *columnprice}
	if columnrange != nil {
		*pricecolumns = append(*pricecolumns, *columnrange)
	}
	apitablerows, err := tablerowrepository.GetAll("", "", supplierprice.Customer_Table_ID, pricecolumns)
	if err != nil {
		return err
	}
	for _, apitablerow := range *apitablerows {
		apismshlrprice := new(models.ApiSMSHLRPrice)
		apismshlrprice.Supplier_ID = supplierprice.Supplier_ID
		for _, apitablecell := range apitablerow.Cells {
			if apitablecell.Table_Column_ID == columnmobileoperator.ID {
				value, err := CheckColumnMobileOperatorInternal(&apitablecell)
				if err != nil {
					return err
				}
				dtomobileoperator, err := CheckMobileOperatorInternal(value, mobileoperatorrepository)
				if err != nil {
					return err
				}
				apismshlrprice.Mobile_Operator_ID = dtomobileoperator.ID
			}
			if columnrange != nil {
				if apitablecell.Table_Column_ID == columnrange.ID {
					apirange, err := CheckColumnRangeInternal(&apitablecell)
					if err != nil {
						return err
					}
					apismshlrprice.AmountRange = *apirange
				}
			}
			if apitablecell.Table_Column_ID == columnprice.ID {
				apismshlrprice.Price, err = CheckColumnPriceInternal(&apitablecell)
				if err != nil {
					return err
				}
			}
		}
		*smshlrprices = append(*smshlrprices, *apismshlrprice)
	}

	return nil
}

func GetSMSHLRPrices(alias string, supplier_id int64, r render.Render, pricerepository services.PriceRepository,
	tablecolumnrepository services.TableColumnRepository, tablerowrepository services.TableRowRepository,
	mobileoperatorrepository services.MobileOperatorRepository, language string, internal bool) (smshlrprices *[]models.ApiSMSHLRPrice, err error) {
	smshlrprices = new([]models.ApiSMSHLRPrice)
	supplierprices, err := pricerepository.GetSupplierPrices(alias)
	if err != nil {
		if !internal {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
		}
		return nil, err
	}
	for _, supplierprice := range *supplierprices {
		if supplier_id == 0 || supplier_id == supplierprice.Supplier_ID {
			columnmobileoperator, columnrange, columnprice, err := GetSMSHLRPriceColumns(&supplierprice, tablecolumnrepository)
			if err != nil {
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, err
			}
			err = GetSMSHLRPriceRows(columnmobileoperator, columnrange, columnprice, &supplierprice, tablerowrepository, mobileoperatorrepository, smshlrprices)
			if err != nil {
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[language].Errors.Api.Data_Wrong})
				}
				return nil, err
			}
		}
	}

	return smshlrprices, nil
}

func GetRecognizePrices(alias string, supplier_id int64, r render.Render, pricerepository services.PriceRepository, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository, recognizeproductrepository services.RecognizeProductRepository,
	language string, internal bool) (recognizeprices *[]models.ApiRecognizePrice, err error) {
	recognizeprices = new([]models.ApiRecognizePrice)
	supplierprices, err := pricerepository.GetSupplierPrices(alias)
	if err != nil {
		if !internal {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
		}
		return nil, err
	}
	for _, supplierprice := range *supplierprices {
		if supplier_id == 0 || supplier_id == supplierprice.Supplier_ID {
			dtotablecolumns, err := tablecolumnrepository.GetByTable(supplierprice.Customer_Table_ID)
			if err != nil {
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, err
			}
			var columnproduct *models.DtoTableColumn
			var columnprice *models.DtoTableColumn
			var columndiscount *models.DtoTableColumn
			for i, tablecolumn := range *dtotablecolumns {
				if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_ID {
					if columnproduct != nil {
						log.Error("Can't have multiple product id column in price list %v for supplier %v",
							supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
						if !internal {
							r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
								Message: config.Localization[language].Errors.Api.Object_NotExist})
						}
						return nil, errors.New("Multiple product id column")
					}
					columnproduct = &(*dtotablecolumns)[i]
				}
				if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_PRICE {
					if columnprice != nil {
						log.Error("Can't have multiple price column in price list %v for supplier %v",
							supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
						if !internal {
							r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
								Message: config.Localization[language].Errors.Api.Object_NotExist})
						}
						return nil, errors.New("Multiple price column")
					}
					columnprice = &(*dtotablecolumns)[i]
				}
				if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_DISCOUNT {
					if columndiscount != nil {
						log.Error("Can't have multiple discount column in price list %v for supplier %v",
							supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
						if !internal {
							r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
								Message: config.Localization[language].Errors.Api.Object_NotExist})
						}
						return nil, errors.New("Multiple discount column")
					}
					columndiscount = &(*dtotablecolumns)[i]
				}
			}
			if columnproduct == nil {
				log.Error("Can't find product id column in price list %v for supplier %v", supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, errors.New("Missed product id column")
			}
			if columnprice == nil {
				log.Error("Can't find price column in price list %v for supplier %v", supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, errors.New("Missed price column")
			}
			pricecolumns := &[]models.DtoTableColumn{*columnproduct, *columnprice}
			if columndiscount != nil {
				*pricecolumns = append(*pricecolumns, *columndiscount)
			}
			apitablerows, err := tablerowrepository.GetAll("", "", supplierprice.Customer_Table_ID, pricecolumns)
			if err != nil {
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, err
			}
			for _, apitablerow := range *apitablerows {
				apirecognizeprice := new(models.ApiRecognizePrice)
				apirecognizeprice.Supplier_ID = supplierprice.Supplier_ID
				for _, apitablecell := range apitablerow.Cells {
					if apitablecell.Table_Column_ID == columnproduct.ID {
						var value int
						if !internal {
							value, err = CheckColumnProduct(&apitablecell, r, language)
						} else {
							value, err = CheckColumnProductInternal(&apitablecell)
						}
						if err != nil {
							return nil, err
						}
						var dtorecognizeproduct *models.DtoRecognizeProduct
						if !internal {
							dtorecognizeproduct, err = CheckRecognizeProduct(value, r, recognizeproductrepository, language)
						} else {
							dtorecognizeproduct, err = CheckRecognizeProductInternal(value, recognizeproductrepository)
						}
						if err != nil {
							return nil, err
						}
						apirecognizeprice.Product_ID = dtorecognizeproduct.ID
						apirecognizeprice.Increase = dtorecognizeproduct.Increase
					}
					if apitablecell.Table_Column_ID == columnprice.ID {
						if !internal {
							apirecognizeprice.Price, err = CheckColumnPrice(&apitablecell, r, language)
						} else {
							apirecognizeprice.Price, err = CheckColumnPriceInternal(&apitablecell)
						}
						if err != nil {
							return nil, err
						}
					}
					if columndiscount != nil {
						if apitablecell.Table_Column_ID == columndiscount.ID {
							if !internal {
								apirecognizeprice.PriceIncrease, err = CheckColumnDiscount(&apitablecell, r, language)
							} else {
								apirecognizeprice.PriceIncrease, err = CheckColumnDiscountInternal(&apitablecell)
							}
							if err != nil {
								return nil, err
							}
						}
					}
				}
				*recognizeprices = append(*recognizeprices, *apirecognizeprice)
			}
		}
	}

	return recognizeprices, nil
}

func GetVerifyPrices(alias string, supplier_id int64, r render.Render, pricerepository services.PriceRepository,
	tablecolumnrepository services.TableColumnRepository, tablerowrepository services.TableRowRepository,
	verifyproductrepository services.VerifyProductRepository, language string, internal bool) (verifyprices *[]models.ApiVerifyPrice, err error) {
	verifyprices = new([]models.ApiVerifyPrice)
	supplierprices, err := pricerepository.GetSupplierPrices(alias)
	if err != nil {
		if !internal {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
		}
		return nil, err
	}
	for _, supplierprice := range *supplierprices {
		if supplier_id == 0 || supplier_id == supplierprice.Supplier_ID {
			dtotablecolumns, err := tablecolumnrepository.GetByTable(supplierprice.Customer_Table_ID)
			if err != nil {
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, err
			}
			var columnproduct *models.DtoTableColumn
			var columnprice *models.DtoTableColumn
			for i, tablecolumn := range *dtotablecolumns {
				if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_ID {
					if columnproduct != nil {
						log.Error("Can't have multiple product id column in price list %v for supplier %v",
							supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
						if !internal {
							r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
								Message: config.Localization[language].Errors.Api.Object_NotExist})
						}
						return nil, errors.New("Multiple product id column")
					}
					columnproduct = &(*dtotablecolumns)[i]
				}
				if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_PRICE {
					if columnprice != nil {
						log.Error("Can't have multiple price column in price list %v for supplier %v",
							supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
						if !internal {
							r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
								Message: config.Localization[language].Errors.Api.Object_NotExist})
						}
						return nil, errors.New("Multiple price column")
					}
					columnprice = &(*dtotablecolumns)[i]
				}
			}
			if columnproduct == nil {
				log.Error("Can't find product id column in price list %v for supplier %v", supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, errors.New("Missed product id column")
			}
			if columnprice == nil {
				log.Error("Can't find price column in price list %v for supplier %v", supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, errors.New("Missed price column")
			}
			pricecolumns := &[]models.DtoTableColumn{*columnproduct, *columnprice}
			apitablerows, err := tablerowrepository.GetAll("", "", supplierprice.Customer_Table_ID, pricecolumns)
			if err != nil {
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, err
			}
			for _, apitablerow := range *apitablerows {
				apiverifyprice := new(models.ApiVerifyPrice)
				apiverifyprice.Supplier_ID = supplierprice.Supplier_ID
				for _, apitablecell := range apitablerow.Cells {
					if apitablecell.Table_Column_ID == columnproduct.ID {
						var value int
						if !internal {
							value, err = CheckColumnProduct(&apitablecell, r, language)
						} else {
							value, err = CheckColumnProductInternal(&apitablecell)
						}
						if err != nil {
							return nil, err
						}
						var dtoverifyproduct *models.DtoVerifyProduct
						if !internal {
							dtoverifyproduct, err = CheckVerifyProduct(value, r, verifyproductrepository, language)
						} else {
							dtoverifyproduct, err = CheckVerifyProductInternal(value, verifyproductrepository)
						}
						if err != nil {
							return nil, err
						}
						apiverifyprice.Product_ID = dtoverifyproduct.ID
					}
					if apitablecell.Table_Column_ID == columnprice.ID {
						if !internal {
							apiverifyprice.Price, err = CheckColumnPrice(&apitablecell, r, language)
						} else {
							apiverifyprice.Price, err = CheckColumnPriceInternal(&apitablecell)
						}
						if err != nil {
							return nil, err
						}
					}
				}
				*verifyprices = append(*verifyprices, *apiverifyprice)
			}
		}
	}

	return verifyprices, nil
}

func GetHeaderPrices(alias string, supplier_id int64, r render.Render, pricerepository services.PriceRepository,
	tablecolumnrepository services.TableColumnRepository, tablerowrepository services.TableRowRepository,
	headerproductrepository services.HeaderProductRepository, language string, internal bool) (headerprices *[]models.ApiHeaderPrice, err error) {
	headerprices = new([]models.ApiHeaderPrice)
	supplierprices, err := pricerepository.GetSupplierPrices(alias)
	if err != nil {
		if !internal {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
		}
		return nil, err
	}
	for _, supplierprice := range *supplierprices {
		if supplier_id == 0 || supplier_id == supplierprice.Supplier_ID {
			dtotablecolumns, err := tablecolumnrepository.GetByTable(supplierprice.Customer_Table_ID)
			if err != nil {
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, err
			}
			var columnproduct *models.DtoTableColumn
			var columnrange *models.DtoTableColumn
			var columnprice *models.DtoTableColumn
			var columndiscount *models.DtoTableColumn
			for i, tablecolumn := range *dtotablecolumns {
				if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_ID {
					if columnproduct != nil {
						log.Error("Can't have multiple product id column in price list %v for supplier %v",
							supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
						if !internal {
							r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
								Message: config.Localization[language].Errors.Api.Object_NotExist})
						}
						return nil, errors.New("Multiple product id column")
					}
					columnproduct = &(*dtotablecolumns)[i]
				}
				if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_RANGE {
					if columnrange != nil {
						log.Error("Can't have multiple range column in price list %v for supplier %v",
							supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
						return nil, errors.New("Multiple range column")
					}
					columnrange = &(*dtotablecolumns)[i]
				}
				if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_PRICE {
					if columnprice != nil {
						log.Error("Can't have multiple price column in price list %v for supplier %v",
							supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
						if !internal {
							r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
								Message: config.Localization[language].Errors.Api.Object_NotExist})
						}
						return nil, errors.New("Multiple price column")
					}
					columnprice = &(*dtotablecolumns)[i]
				}

				if tablecolumn.Column_Type_ID == models.COLUMN_TYPE_PRICELIST_DISCOUNT {
					if columndiscount != nil {
						log.Error("Can't have multiple discount column in price list %v for supplier %v",
							supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
						if !internal {
							r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
								Message: config.Localization[language].Errors.Api.Object_NotExist})
						}
						return nil, errors.New("Multiple discount column")
					}
					columndiscount = &(*dtotablecolumns)[i]
				}
			}
			if columnproduct == nil {
				log.Error("Can't find product id column in price list %v for supplier %v", supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, errors.New("Missed product id column")
			}
			if columnprice == nil {
				log.Error("Can't find price column in price list %v for supplier %v", supplierprice.Customer_Table_ID, supplierprice.Supplier_ID)
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, errors.New("Missed price column")
			}
			pricecolumns := &[]models.DtoTableColumn{*columnproduct, *columnprice}
			if columnrange != nil {
				*pricecolumns = append(*pricecolumns, *columnrange)
			}
			if columndiscount != nil {
				*pricecolumns = append(*pricecolumns, *columndiscount)
			}
			apitablerows, err := tablerowrepository.GetAll("", "", supplierprice.Customer_Table_ID, pricecolumns)
			if err != nil {
				if !internal {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[language].Errors.Api.Object_NotExist})
				}
				return nil, err
			}
			for _, apitablerow := range *apitablerows {
				apiheaderprice := new(models.ApiHeaderPrice)
				apiheaderprice.Supplier_ID = supplierprice.Supplier_ID
				for _, apitablecell := range apitablerow.Cells {
					if apitablecell.Table_Column_ID == columnproduct.ID {
						var value int
						if !internal {
							value, err = CheckColumnProduct(&apitablecell, r, language)
						} else {
							value, err = CheckColumnProductInternal(&apitablecell)
						}
						if err != nil {
							return nil, err
						}
						var dtoheaderproduct *models.DtoHeaderProduct
						if !internal {
							dtoheaderproduct, err = CheckHeaderProduct(value, r, headerproductrepository, language)
						} else {
							dtoheaderproduct, err = CheckHeaderProductInternal(value, headerproductrepository)
						}
						if err != nil {
							return nil, err
						}
						apiheaderprice.Product_ID = dtoheaderproduct.ID
						apiheaderprice.Increase = dtoheaderproduct.Increase
						apiheaderprice.FeeOnce = dtoheaderproduct.FeeOnce
						apiheaderprice.FeeMonthly = dtoheaderproduct.FeeMonthly
					}
					if columnrange != nil {
						if apitablecell.Table_Column_ID == columnrange.ID {
							apirange, err := CheckColumnRangeInternal(&apitablecell)
							if err != nil {
								return nil, err
							}
							apiheaderprice.AmountRange = *apirange
						}
					}
					if apitablecell.Table_Column_ID == columnprice.ID {
						if !internal {
							apiheaderprice.Price, err = CheckColumnPrice(&apitablecell, r, language)
						} else {
							apiheaderprice.Price, err = CheckColumnPriceInternal(&apitablecell)
						}
						if err != nil {
							return nil, err
						}
					}
					if columndiscount != nil {
						if apitablecell.Table_Column_ID == columndiscount.ID {
							if !internal {
								apiheaderprice.PriceIncrease, err = CheckColumnDiscount(&apitablecell, r, language)
							} else {
								apiheaderprice.PriceIncrease, err = CheckColumnDiscountInternal(&apitablecell)
							}
							if err != nil {
								return nil, err
							}
						}
					}
				}
				*headerprices = append(*headerprices, *apiheaderprice)
			}
		}
	}

	return headerprices, nil
}
