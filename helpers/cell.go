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

func ValidateCell(value string, dtotablecolumn *models.DtoTableColumn, r render.Render,
	columntyperepository services.ColumnTypeRepository, language string) (valid bool, err error) {
	valid = true

	dtocolumntype, err := columntyperepository.Get(dtotablecolumn.Column_Type_ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return false, err
	}

	valid, _, err = columntyperepository.Validate(dtocolumntype, nil, value)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return false, err
	}

	return valid, nil
}

func CheckTableCell(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository,
	language string) (dtotablecell *models.DtoTableCell, dtotablecolumn *models.DtoTableColumn, dtotablerow *models.DtoTableRow, err error) {
	dtotablecell = new(models.DtoTableCell)
	dtotablerow, err = CheckTableRow(r, params, customertablerepository, tablerowrepository, language)
	if err != nil {
		return nil, nil, nil, err
	}
	dtotablecolumn, err = CheckTableColumn(r, params, columntyperepository, customertablerepository, tablecolumnrepository, language)
	if err != nil {
		return nil, nil, nil, err
	}

	tablecolumns, err := tablecolumnrepository.GetByTable(dtotablerow.Customer_Table_ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, err
	}

	tablecells, err := dtotablerow.TableRowToDtoTableCells(tablecolumns)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, err
	}

	found := false
	for _, tablecell := range *tablecells {
		if tablecell.Table_Column_ID == dtotablecolumn.ID {
			found = true
			*dtotablecell = tablecell
			break
		}
	}
	if !found {
		log.Error("Can't find cell for table row %v column %v", dtotablerow.ID, dtotablecolumn.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, errors.New("Cell not found")
	}

	return dtotablecell, dtotablecolumn, dtotablerow, nil
}

func SaveTableCell(value string, r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository, language string) (dtotablecell *models.DtoTableCell, err error) {
	dtotablecell, dtotablecolumn, oldtablerow, err := CheckTableCell(r, params, customertablerepository, columntyperepository,
		tablecolumnrepository, tablerowrepository, language)
	if err != nil {
		return nil, err
	}

	tablecolumns, err := tablecolumnrepository.GetByTable(oldtablerow.Customer_Table_ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	tablecells, err := oldtablerow.TableRowToDtoTableCells(tablecolumns)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	newtablerow := new(models.DtoTableRow)
	*newtablerow = *oldtablerow

	newtablerow.Created = time.Now()
	newtablerow.Edition += 1

	oldtablerow.ID = 0
	oldtablerow.Active = false
	oldtablerow.Original_ID = newtablerow.ID

	for i, _ := range *tablecells {
		if (*tablecells)[i].Table_Column_ID == dtotablecell.Table_Column_ID {
			dtotablecell.Value = value
			dtotablecell.Checked = true
			dtotablecell.Valid, err = ValidateCell(dtotablecell.Value, dtotablecolumn, r, columntyperepository, language)
			if err != nil {
				return nil, err
			}
			(*tablecells)[i] = *dtotablecell
			break
		}
	}
	err = newtablerow.TableCellsToTableRow(tablecells, tablecolumns)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	err = tablerowrepository.Update(newtablerow, oldtablerow, false, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	return dtotablecell, nil
}
