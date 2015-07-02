package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"sort"
	"strconv"
	"time"
	"types"
)

const (
	PARAM_NAME_COLUMN_ID = "cid"
)

func CheckColumnValidity(tableid int64, columnid int64, r render.Render, columntyperepository services.ColumnTypeRepository,
	tablecolumnrepository services.TableColumnRepository, language string) (dtotablecolumn *models.DtoTableColumn, err error) {
	dtotablecolumn, err = tablecolumnrepository.Get(columnid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtotablecolumn.Active {
		log.Error("Table column is not active %v", columnid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Column is not active")
	}

	err = IsColumnTypeActive(r, columntyperepository, dtotablecolumn.Column_Type_ID, language)
	if err != nil {
		return nil, err
	}

	if dtotablecolumn.Customer_Table_ID != tableid {
		log.Error("Column %v doesn't belong table %v", columnid, tableid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Not matched table and column")
	}

	return dtotablecolumn, nil
}

func CheckTableColumn(r render.Render, params martini.Params, columntyperepository services.ColumnTypeRepository,
	customertablerepository services.CustomerTableRepository, tablecolumnrepository services.TableColumnRepository,
	language string) (dtotablecolumn *models.DtoTableColumn, err error) {
	var tableid int64
	var columnid int64

	tableid, err = CheckParameterInt(r, params[PARAM_NAME_TABLE_ID], language)
	if err != nil {
		return nil, err
	}
	columnid, err = CheckParameterInt(r, params[PARAM_NAME_COLUMN_ID], language)
	if err != nil {
		return nil, err
	}
	_, err = IsTableAvailable(r, customertablerepository, tableid, language)
	if err != nil {
		return nil, err
	}

	return CheckColumnValidity(tableid, columnid, r, columntyperepository, tablecolumnrepository, language)
}

func IsColumnTypeActive(r render.Render, columntyperepository services.ColumnTypeRepository, typeid int, language string) (err error) {
	dtocolumntype, err := columntyperepository.Get(typeid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return errors.New("Column is not found")
	}
	if !dtocolumntype.Active {
		log.Error("Column type is not active %v", typeid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return errors.New("Colimn is not active")
	}

	return nil
}

func CheckColumnSet(ids models.IDs, tableid int64, r render.Render, tablecolumnrepository services.TableColumnRepository,
	language string) (tablecolumns *[]models.DtoTableColumn, err error) {
	tablecolumns, err = tablecolumnrepository.GetByTable(tableid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	for _, tablecolumn := range *tablecolumns {
		found := false
		for _, id := range ids.GetIDs() {
			if tablecolumn.ID == id {
				found = true
				break
			}
		}
		if !found {
			log.Error("Can't find column %v", tablecolumn.ID)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, errors.New("Column not found")
		}
	}

	for _, id := range ids.GetIDs() {
		found := false
		for _, tablecolumn := range *tablecolumns {
			if id == tablecolumn.ID {
				found = true
				break
			}
		}
		if !found {
			log.Error("Can't find column %v for table", id)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, errors.New("Column not found")
		}
	}

	return tablecolumns, nil
}

func FindFreeColumnInternal(tableid int64, afterfield byte, tablecolumnrepository services.TableColumnRepository) (fieldnum byte, err error) {
	tablecolumns, err := tablecolumnrepository.GetByTable(tableid)
	if err != nil {
		return 0, err
	}

	allcolumns := make(map[int]bool)
	for i := 0; i < models.MAX_COLUMN_NUMBER; i++ {
		allcolumns[i] = true
	}
	for _, tablecolumn := range *tablecolumns {
		allcolumns[int(tablecolumn.FieldNum)-1] = false
	}

	var keys []int
	for i, avialable := range allcolumns {
		if avialable {
			keys = append(keys, i)
		}
	}
	sort.Ints(keys)

	found := false
	for _, number := range keys {
		if afterfield != 0 && byte(number)+1 <= afterfield {
			continue
		}
		fieldnum = byte(number) + 1
		found = true
		break
	}
	if !found {
		log.Error("Can't find free column for table %v", tableid)
		return 0, errors.New("No free column")
	}

	return fieldnum, nil
}

func FindFreeColumn(tableid int64, afterfield byte, r render.Render, tablecolumnrepository services.TableColumnRepository,
	language string) (fieldnum byte, err error) {
	fieldnum, err = FindFreeColumnInternal(tableid, afterfield, tablecolumnrepository)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return 0, err
	}

	return fieldnum, nil
}

func UpdateProductColumn(products *[]models.ApiProduct, tableid int64, r render.Render,
	tablecolumnrepository services.TableColumnRepository, columntyperepository services.ColumnTypeRepository,
	tablerowrepository services.TableRowRepository, language string) (err error) {
	tablecolumns, err := tablecolumnrepository.GetByTable(tableid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}

	for _, product := range *products {
		dtotablecolumn, err := CheckColumnValidity(tableid, product.Table_Column_ID, r, columntyperepository, tablecolumnrepository, language)
		if err != nil {
			return err
		}
		oldtablerow, err := CheckRowValidity(tableid, product.Table_Row_ID, r, tablerowrepository, language)
		if err != nil {
			return err
		}

		tablecells, err := oldtablerow.TableRowToDtoTableCells(tablecolumns)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return err
		}

		newtablerow := new(models.DtoTableRow)
		*newtablerow = *oldtablerow

		newtablerow.Created = time.Now()
		newtablerow.Edition += 1

		oldtablerow.ID = 0
		oldtablerow.Active = false
		oldtablerow.Original_ID = newtablerow.ID

		dtotablecell := new(models.DtoTableCell)
		dtotablecell.Table_Column_ID = product.Table_Column_ID

		for i, _ := range *tablecells {
			if (*tablecells)[i].Table_Column_ID == dtotablecell.Table_Column_ID {
				dtotablecell.Value = strconv.Itoa(product.Product_ID)
				dtotablecell.Checked = true
				dtotablecell.Valid, err = ValidateCell(dtotablecell.Value, dtotablecolumn, r, columntyperepository, language)
				if err != nil {
					return err
				}
				(*tablecells)[i] = *dtotablecell
				break
			}
		}
		err = newtablerow.TableCellsToTableRow(tablecells, tablecolumns)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return err
		}

		err = tablerowrepository.Update(newtablerow, oldtablerow, false, true)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return err
		}
	}

	return nil
}
