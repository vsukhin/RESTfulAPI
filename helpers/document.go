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
	PARAM_NAME_DOCUMENT_ID = "docid"
)

func CheckDocument(r render.Render, params martini.Params, documentrepository services.DocumentRepository,
	language string) (dtodocument *models.DtoDocument, err error) {
	document_id, err := CheckParameterInt(r, params[PARAM_NAME_DOCUMENT_ID], language)
	if err != nil {
		return nil, err
	}

	dtodocument, err = documentrepository.Get(document_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	if !dtodocument.Active {
		log.Error("Document is not active %v", dtodocument.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Not active document")
	}

	return dtodocument, nil
}

func CreateDocumentByType(document_type_id int, company_id int64, name string, file_id int64, locked bool, pending bool, r render.Render,
	documentrepository services.DocumentRepository, companyrepository services.CompanyRepository, unitrepository services.UnitRepository,
	documenttyperepository services.DocumentTypeRepository, filerepository services.FileRepository,
	session *models.DtoSession) (dtodocument *models.DtoDocument, err error) {
	dtodocumenttype, err := CheckDocumentType(document_type_id, r, documenttyperepository, session.Language)
	if err != nil {
		return nil, err
	}
	dtocompany, err := CheckCompanyAvailability(company_id, session.UserID, r, companyrepository, session.Language)
	if err != nil {
		return nil, err
	}
	if dtocompany.Locked {
		log.Error("Company is locked %v", dtocompany.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return nil, errors.New("Locked company")
	}

	if file_id != 0 {
		_, err = filerepository.Get(file_id)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
			return nil, err
		}
	}

	unit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return nil, err
	}

	dtodocument = new(models.DtoDocument)
	dtodocument.Document_Type_ID = document_type_id
	dtodocument.Unit_ID = unit.ID
	dtodocument.Company_ID = company_id
	if name != "" {
		dtodocument.Name = name
	} else {
		dtodocument.Name = dtodocumenttype.Name
	}
	dtodocument.Locked = locked
	dtodocument.Pending = pending
	dtodocument.File_ID = file_id
	dtodocument.Created = time.Now()
	dtodocument.Updated = time.Now()
	dtodocument.Active = true

	err = documentrepository.Create(dtodocument)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return nil, err
	}

	return dtodocument, nil
}
