package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/nfnt/resize"
	"github.com/saintfish/chardet"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"types"
)

// get /api/v1.0/files/:key/
// get /api/v1.0/files/:key/:modeId/
// get /api/v1.0/files/:key/:modeId/:size/
func GetFile(w http.ResponseWriter, r render.Render, params martini.Params, filerepository services.FileRepository, session *models.DtoSession) {
	fileid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_KEY], session.Language)
	if err != nil {
		return
	}
	modeid := params[helpers.PARAM_NAME_MODEID]
	if len(modeid) > helpers.PARAM_LENGTH_MAX {
		log.Error("Parameter is too long", modeid)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	if modeid == "" {
		modeid = helpers.PARAM_MODEID_ATTACH
	}
	size := params[helpers.PARAM_NAME_SIZE]
	if len(size) > helpers.PARAM_LENGTH_MAX {
		log.Error("Parameter is too long", size)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	resizewidth := 0
	resizeheight := 0
	if size != "" {
		sizes := strings.Split(size, ":")
		if len(sizes) != helpers.PARAM_SIZE_NUMBER {
			log.Error("Wrong number of size parameter elements %v", len(sizes))
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
			return
		}

		sizeparams := []int{helpers.PARAM_SIZE_WIDTH, helpers.PARAM_SIZE_HEIGHT}
		for _, sizeparam := range sizeparams {
			paramvalue, err := strconv.Atoi(sizes[sizeparam])
			if err != nil {
				log.Error("Wrong parameter value %v, %v", sizeparam, sizes[sizeparam])
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
			if paramvalue < 0 {
				log.Error("Wrong parameter value %v, %v", sizeparam, sizes[sizeparam])
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
			if sizeparam == helpers.PARAM_SIZE_WIDTH {
				resizewidth = paramvalue
			}
			if sizeparam == helpers.PARAM_SIZE_HEIGHT {
				resizeheight = paramvalue
			}
		}
	}

	file, err := filerepository.Get(fileid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	var width, height int
	data := file.FileData
	imagetype, image, imageconfig := helpers.DetectImage(data)
	if imagetype != helpers.IMAGE_TYPE_UNKNOWN {
		width = imageconfig.Width
		height = imageconfig.Height
	}

	convertimagetype := helpers.IMAGE_TYPE_UNKNOWN
	switch strings.ToLower(modeid) {
	case helpers.PARAM_MODEID_ATTACH:
		w.Header().Set("Content-Type", helpers.CONTENT_TYPE_DEFAULT)
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)
		w.Write(data)
	case helpers.PARAM_MODEID_DATAURLGIF:
		fallthrough
	case helpers.PARAM_MODEID_DATAURLJPEG:
		fallthrough
	case helpers.PARAM_MODEID_DATAURLPNG:
		fallthrough
	case helpers.PARAM_MODEID_DATAURL:
		if resizewidth != 0 || resizeheight != 0 {
			if imagetype != helpers.IMAGE_TYPE_UNKNOWN {
				image = resize.Resize(uint(resizewidth), uint(resizeheight), image, resize.Lanczos3)
				data, err = helpers.ConvertImage(image, imagetype)
				if err != nil {
					log.Error("Can't convert image to format %v", imagetype)
					r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
					return
				}
			} else {
				log.Error("File data is not known image")
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
		}
		if strings.ToLower(modeid) == helpers.PARAM_MODEID_DATAURLGIF {
			convertimagetype = helpers.IMAGE_TYPE_GIF
		}
		if strings.ToLower(modeid) == helpers.PARAM_MODEID_DATAURLJPEG {
			convertimagetype = helpers.IMAGE_TYPE_JPEG
		}
		if strings.ToLower(modeid) == helpers.PARAM_MODEID_DATAURLPNG {
			convertimagetype = helpers.IMAGE_TYPE_PNG
		}
		if convertimagetype != helpers.IMAGE_TYPE_UNKNOWN {
			if imagetype != helpers.IMAGE_TYPE_UNKNOWN {
				data, err = helpers.ConvertImage(image, convertimagetype)
				if err != nil {
					log.Error("Can't convert image to format %v", convertimagetype)
					r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
					return
				}
			} else {
				log.Error("File data is not known image")
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return
			}
		}
		w.Write([]byte(base64.StdEncoding.EncodeToString(data)))
	case helpers.PARAM_MODEID_HEAD:
	case helpers.PARAM_MODEID_OBJECT:
		charset := helpers.CONTENT_ENCODING_DEFAULT
		if imagetype == helpers.IMAGE_TYPE_UNKNOWN {
			detector := chardet.NewTextDetector()
			result, err := detector.DetectBest(file.FileData)
			if err == nil {
				charset = strings.ToLower(result.Charset)
			}
		}
		contenttype := http.DetectContentType(file.FileData)
		ext := filepath.Ext(file.Name)
		if ext == "" {
			ext = helpers.MimeTypes[contenttype]
		}
		hash := sha512.Sum512(file.FileData)
		if strings.ToLower(modeid) == helpers.PARAM_MODEID_HEAD {
			r.JSON(http.StatusOK, models.NewApiFileHead(file.ID, file.Name, ext,
				imagetype != helpers.IMAGE_TYPE_UNKNOWN, width, height,
				len(file.FileData), file.Created.Format(time.UnixDate), contenttype, charset, hex.EncodeToString(hash[:])))
		}
		if strings.ToLower(modeid) == helpers.PARAM_MODEID_OBJECT {
			r.JSON(http.StatusOK, models.NewApiFileObject(*models.NewApiFileHead(file.ID, file.Name, ext,
				imagetype != helpers.IMAGE_TYPE_UNKNOWN, width, height,
				len(file.FileData), file.Created.Format(time.UnixDate), contenttype, charset, hex.EncodeToString(hash[:])),
				base64.StdEncoding.EncodeToString(file.FileData)))
		}
	default:
		log.Error("Uknown file mode")
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
	}
}

// post /api/v1.0/files/
func UploadFile(data models.ViewFile, r render.Render, filerepository services.FileRepository, session *models.DtoSession) {
	if data.FileData == nil {
		log.Error("Empty data file field")
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	file := new(models.DtoFile)
	file.Created = time.Now()
	file.Name = data.FileData.Filename
	file.Path = "/" + fmt.Sprintf("%04d/%02d/%02d/", file.Created.Year(), file.Created.Month(), file.Created.Day())
	file.Permanent = false
	file.Export_Ready = true
	file.Export_Percentage = 100
	file.Export_Object_ID = 0
	file.Export_Error = false
	file.Export_ErrorDescription = ""

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
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
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
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
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
