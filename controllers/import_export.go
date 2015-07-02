package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
	"types"
)

const (
	URL_EXPORT_DATA = "/api/v1.0/tables/export"
)

// post /api/v1.0/tables/import/
func ImportDataFromFile(errors binding.Errors, viewimporttable models.ViewImportTable, r render.Render, userrepository services.UserRepository,
	filerepository services.FileRepository, unitrepository services.UnitRepository, tabletyperepository services.TableTypeRepository,
	customertablerepository services.CustomerTableRepository, importsteprepository services.ImportStepRepository,
	columntyperepository services.ColumnTypeRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewimporttable, errors, r, session.Language) != nil {
		return
	}

	fileid, err := strconv.ParseInt(viewimporttable.File_ID, 0, 64)
	if err != nil {
		log.Error("Can't convert to number %v with value %v", err, viewimporttable.File_ID)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	file, err := filerepository.GetBriefly(fileid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	unitid, typeid, err := helpers.CheckCustomerTableParameters(r, 0, models.TABLE_TYPE_DEFAULT, session.UserID, session.Language,
		userrepository, unitrepository, tabletyperepository)
	if err != nil {
		return
	}

	dtocustomertable := new(models.DtoCustomerTable)
	dtocustomertable.Name = file.Name
	dtocustomertable.Created = time.Now()
	dtocustomertable.TypeID = typeid
	dtocustomertable.UnitID = unitid
	dtocustomertable.Active = true
	dtocustomertable.Permanent = false
	dtocustomertable.Import_Ready = false
	dtocustomertable.Import_Percentage = 25

	err = customertablerepository.Create(dtocustomertable)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	// 1
	dtoimportstep := models.NewDtoImportStep(dtocustomertable.ID, 1, true, 100, time.Now(), time.Now())
	err = importsteprepository.Save(dtoimportstep)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	go helpers.ImportData(viewimporttable, file, dtocustomertable, customertablerepository, importsteprepository, columntyperepository)

	r.JSON(http.StatusOK, models.NewApiImportTable(dtocustomertable.ID))
}

// get /api/v1.0/tables/import/:tmpid/columns/
func GetImportDataColumns(w http.ResponseWriter, r render.Render, params martini.Params, customertablerepository services.CustomerTableRepository,
	tablecolumnrepository services.TableColumnRepository, session *models.DtoSession) {
	tableid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_TEMPORABLE_TABLE_ID], session.Language)
	if err != nil {
		return
	}
	dtocustomertable, err := helpers.IsTableActive(r, customertablerepository, tableid, session.Language)
	if err != nil {
		return
	}
	if dtocustomertable.Permanent {
		log.Error("Can't inquire permanent table %v", tableid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	tablecolumns, err := tablecolumnrepository.GetByTable(tableid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	importcolumns := new([]models.ApiImportColumn)
	for _, tablecolumn := range *tablecolumns {
		*importcolumns = append(*importcolumns, *models.NewApiImportColumn(tablecolumn.ID, tablecolumn.Name, tablecolumn.Position))
	}

	helpers.RenderJSONArray(importcolumns, len(*importcolumns), w, r)
}

// put /api/v1.0/tables/import/:tmpid/columns/
func UpdateImportDataColumns(errors binding.Errors, viewimportcolumns models.ViewImportColumns, r render.Render, params martini.Params,
	customertablerepository services.CustomerTableRepository, tablecolumnrepository services.TableColumnRepository,
	columntyperepository services.ColumnTypeRepository, tablerowrepository services.TableRowRepository,
	importsteprepository services.ImportStepRepository, session *models.DtoSession) {
	if helpers.CheckValidation(&viewimportcolumns, errors, r, session.Language) != nil {
		return
	}
	tableid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_TEMPORABLE_TABLE_ID], session.Language)
	if err != nil {
		return
	}
	dtocustomertable, err := helpers.IsTableActive(r, customertablerepository, tableid, session.Language)
	if err != nil {
		return
	}
	if dtocustomertable.Permanent {
		log.Error("Can't update permanent table %v", tableid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	_, err = helpers.CheckColumnSet(viewimportcolumns, tableid, r, tablecolumnrepository, session.Language)
	if err != nil {
		return
	}

	dtoimportstep := models.NewDtoImportStep(tableid, 4, false, 0, time.Now(), time.Now())
	err = importsteprepository.Save(dtoimportstep)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	dtotablecolumns := new([]models.DtoTableColumn)
	for _, viewimportcolumn := range viewimportcolumns {
		dtotablecolumn, err := tablecolumnrepository.Get(viewimportcolumn.ID)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return
		}
		if dtotablecolumn.Prebuilt {
			log.Error("Can't update prebuilt column %v", dtotablecolumn.ID)
			r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_DATA_CHANGES_DENIED,
				Message: config.Localization[session.Language].Errors.Api.Data_Changes_Denied})
			return
		}

		dtotablecolumn.Name = viewimportcolumn.Name
		dtotablecolumn.Position = viewimportcolumn.Position
		if helpers.IsColumnTypeActive(r, columntyperepository, viewimportcolumn.TypeID, session.Language) != nil {
			return
		}

		dtotablecolumn.Column_Type_ID = viewimportcolumn.TypeID
		if viewimportcolumn.Use {
			dtotablecolumn.Active = true
		} else {
			dtotablecolumn.Active = false
		}
		*dtotablecolumns = append(*dtotablecolumns, *dtotablecolumn)
	}

	dtocustomertable.Permanent = true
	dtocustomertable.Import_Ready = true
	dtocustomertable.Import_Percentage = 100
	err = customertablerepository.UpdateImportStructure(dtocustomertable, dtotablecolumns, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	// 4
	dtoimportstep.Ready = true
	dtoimportstep.Percentage = 100
	dtoimportstep.Completed = time.Now()
	err = importsteprepository.Save(dtoimportstep)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	go helpers.CheckTableCells(dtocustomertable, tablecolumnrepository, columntyperepository, tablerowrepository, importsteprepository)

	r.JSON(http.StatusOK, models.NewApiLongCustomerTable(dtocustomertable.ID, dtocustomertable.Name, dtocustomertable.TypeID, dtocustomertable.UnitID))
}

// options /api/v1.0/tables/import/:tmpid/
func GetImportDataStatus(r render.Render, params martini.Params, filerepository services.FileRepository,
	customertablerepository services.CustomerTableRepository, importsteprepository services.ImportStepRepository,
	session *models.DtoSession) {
	tableid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_TEMPORABLE_TABLE_ID], session.Language)
	if err != nil {
		return
	}
	dtocustomertable, err := helpers.IsTableActive(r, customertablerepository, tableid, session.Language)
	if err != nil {
		return
	}

	importsteps, err := importsteprepository.GetByTable(tableid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiImportStatus(dtocustomertable.Import_Ready, dtocustomertable.Import_Percentage, *importsteps,
		dtocustomertable.Import_Columns, dtocustomertable.Import_Rows, dtocustomertable.Import_WrongRows))
}

// options /api/v1.0/tables/export/
func GetExportDataMeta(request *http.Request, r render.Render, dataformatrepository services.DataFormatRepository,
	virtualdirrepository services.VirtualDirRepository, sessionrepository services.SessionRepository, session *models.DtoSession) {
	dataformats, err := dataformatrepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	token, err := sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	virtualdir := models.NewDtoVirtualDir(token, time.Now())
	err = virtualdirrepository.Create(virtualdir)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiMetaExportTable(*dataformats, "http://"+path.Join(request.Host, URL_EXPORT_DATA, virtualdir.Token)))
}

// get /api/v1.0/tables/:tid/export/
func ExportDataToFile(request *http.Request, r render.Render, params martini.Params, filerepository services.FileRepository,
	customertablerepository services.CustomerTableRepository, dataformatrepository services.DataFormatRepository,
	tablecolumnrepository services.TableColumnRepository, session *models.DtoSession) {
	dtocustomertable, err := helpers.CheckTable(r, params, customertablerepository, session.Language)
	if err != nil {
		return
	}
	tablecolumns, err := tablecolumnrepository.GetByTable(dtocustomertable.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if len(*tablecolumns) == 0 {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		log.Error("Can't find any data in table %v", dtocustomertable.ID)
		return
	}

	format, err := url.QueryUnescape(request.URL.Query().Get(helpers.PARAM_QUERY_FORMAT))
	if err != nil {
		log.Error("Can't unescape %v url data", err)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	formatid, err := strconv.Atoi(format)
	if err != nil {
		log.Error("Can't convert data format to number %v with value %v", err, format)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	dataformat, err := dataformatrepository.Get(formatid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	rowtype, err := url.QueryUnescape(request.URL.Query().Get(helpers.PARAM_QUERY_TYPE))
	if err != nil {
		log.Error("Can't unescape %v url data", err)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	if rowtype != models.EXPORT_DATA_ALL && rowtype != models.EXPORT_DATA_VALID && rowtype != models.EXPORT_DATA_INVALID {
		log.Error("Error during looking up for existing export type %v", rowtype)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	file := new(models.DtoFile)
	file.Created = time.Now()
	file.Name = dtocustomertable.Name
	file.Path = "/" + fmt.Sprintf("%04d/%02d/%02d/", file.Created.Year(), file.Created.Month(), file.Created.Day())
	file.Permanent = false
	file.Export_Ready = false
	file.Export_Percentage = 0
	file.Export_Object_ID = dtocustomertable.ID

	err = filerepository.Create(file, nil)
	if err != nil {
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	viewexporttable := new(models.ViewExportTable)
	viewexporttable.Data_Format_ID = dataformat.ID
	viewexporttable.Type = rowtype
	go helpers.ExportData(viewexporttable, file, dtocustomertable, tablecolumns, filerepository, customertablerepository)

	r.JSON(http.StatusOK, models.ApiFile{ID: file.ID})
}

// options /api/v1.0/tables/:tid/export/:fid/
func GetExportDataStatus(r render.Render, params martini.Params, filerepository services.FileRepository,
	customertablerepository services.CustomerTableRepository, session *models.DtoSession) {
	dtocustomertable, err := helpers.CheckTable(r, params, customertablerepository, session.Language)
	if err != nil {
		return
	}

	fileid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_FILE_ID], session.Language)
	if err != nil {
		return
	}

	file, err := filerepository.Get(fileid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	if dtocustomertable.ID != file.Export_Object_ID {
		log.Error("Linked file object %v and exported table %v don't match", file.Export_Object_ID, dtocustomertable.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiExportStatus(file.Export_Ready, file.Export_Percentage,
		fmt.Sprintf("%v", file.Created.Add(config.Configuration.FileTimeout))))
}

// get /api/v1.0/tables/export/:token/:fid/
func GetExportData(w http.ResponseWriter, r render.Render, params martini.Params, virtualdirrepository services.VirtualDirRepository,
	filerepository services.FileRepository) {
	if params[helpers.PARAM_NAME_TOKEN] == "" || len(params[helpers.PARAM_NAME_TOKEN]) > helpers.PARAM_LENGTH_MAX {
		log.Error("Parameter is too long or too short %v", params[helpers.PARAM_NAME_TOKEN])
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	virtualdir, err := virtualdirrepository.Get(params[helpers.PARAM_NAME_TOKEN])
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	if time.Now().Sub(virtualdir.Created) > config.Configuration.FileTimeout {
		log.Error("File token has been expired %v with value %v", virtualdir.Created, virtualdir.Token)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	fileid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_FILE_ID], config.Configuration.Server.DefaultLanguage)
	if err != nil {
		return
	}

	file, err := filerepository.Get(fileid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	if file.Export_Object_ID == 0 {
		log.Error("File doesn't contain exported object %v", fileid)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	if !file.Export_Ready {
		log.Error("File is not yet ready %v", fileid)
		r.JSON(http.StatusNoContent, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)

	w.Write(file.FileData)
}
