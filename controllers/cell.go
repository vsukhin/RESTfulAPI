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
	"types"
)

// options /api/v1.0/tables/:tid/cell/:rid/:cid/
func GetTableMetaCell(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository, session *models.DtoSession) {
	dtotablecell, dtotablecolumn, _, err := helpers.CheckTableCell(r, params, customertablerepository, columntyperepository,
		tablecolumnrepository, tablerowrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiMetaTableCell(dtotablecell.Checked, dtotablecell.Valid, dtotablecolumn.Column_Type_ID))
}

// get /api/v1.0/tables/:tid/cell/:rid/:cid/
func GetTableCell(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository, session *models.DtoSession) {
	dtotablecell, _, _, err := helpers.CheckTableCell(r, params, customertablerepository, columntyperepository, tablecolumnrepository,
		tablerowrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortTableCell(dtotablecell.Value, dtotablecell.Valid))
}

// put /api/v1.0/tables/:tid/cell/:rid/:cid/
func UpdateTableCell(errors binding.Errors, viewtablecell models.ViewTableCell,
	r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}

	dtotablecell, err := helpers.SaveTableCell(viewtablecell.Value, r, params, customertablerepository, columntyperepository,
		tablecolumnrepository, tablerowrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortTableCell(dtotablecell.Value, dtotablecell.Valid))
}

// delete /api/v1.0/tables/:tid/cell/:rid/:cid
func DeleteTableCell(r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	columntyperepository services.ColumnTypeRepository, tablecolumnrepository services.TableColumnRepository,
	tablerowrepository services.TableRowRepository, session *models.DtoSession) {
	_, err := helpers.SaveTableCell("", r, params, customertablerepository, columntyperepository,
		tablecolumnrepository, tablerowrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
