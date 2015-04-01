package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"time"
	"types"
)

// get /api/v1.0/files/:key/
func GetFile(w http.ResponseWriter, r render.Render, params martini.Params, filerepository services.FileRepository, session *models.DtoSession) {
	fileid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_KEY], session.Language)
	if err != nil {
		return
	}

	file, err := filerepository.Get(fileid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)

	w.Write(file.FileData)
}

// post /api/v1.0/files/
func UploadFile(data models.ViewFile, r render.Render, filerepository services.FileRepository, session *models.DtoSession) {
	file := new(models.DtoFile)
	file.Created = time.Now()
	file.Name = data.FileData.Filename
	file.Path = "/" + fmt.Sprintf("%04d/%02d/%02d/", file.Created.Year(), file.Created.Month(), file.Created.Day())
	file.Permanent = false
	file.Export_Ready = true
	file.Export_Percentage = 100
	file.Export_Object_ID = 0

	err := filerepository.Create(file, &data)
	if err != nil {
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.ApiFile{ID: file.ID})
}

// delete /api/v1.0/files/:key/
func DeleteFile(r render.Render, params martini.Params, filerepository services.FileRepository, session *models.DtoSession) {
	fileid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_KEY], session.Language)
	if err != nil {
		return
	}

	file, err := filerepository.Get(fileid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	if file.Permanent {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	err = filerepository.Delete(file)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// get /api/v1.0/images/:type/
func GetImage(r render.Render, params martini.Params, filerepository services.FileRepository, session *models.DtoSession) {
	filetype := params[helpers.PARAM_NAME_TYPE]
	if filetype == "" || len(filetype) > helpers.PARAM_LENGTH_MAX {
		log.Error("Wrong parameter length %v", filetype)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	file, err := filerepository.FindByType(filetype)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.ApiImage{ID: file.ID})
}
