package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

func CheckDocumentType(document_type_id int, r render.Render, documenttyperepository services.DocumentTypeRepository,
	language string) (dtodocumenttype *models.DtoDocumentType, err error) {
	dtodocumenttype, err = documenttyperepository.Get(document_type_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtodocumenttype.Active {
		log.Error("Document type is not active %v", dtodocumenttype.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Document type not active")
	}

	return dtodocumenttype, nil
}
