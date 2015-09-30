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
	"strings"
	"time"
	"types"
)

// get /api/v1.0/classification/documents/
func GetDocumentTypes(w http.ResponseWriter, r render.Render, documenttyperepository services.DocumentTypeRepository, session *models.DtoSession) {
	documenttypes, err := documenttyperepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(documenttypes, len(*documenttypes), w, r)
}

// options /api/v1.0/units/documents/
func GetMetaDocuments(request *http.Request, r render.Render, documentrepository services.DocumentRepository, session *models.DtoSession) {
	var err error
	query := ""

	var filters *[]models.FilterExp
	filters, err = helpers.GetFilterArray(new(models.DocumentSearch), nil, request, r, session.Language)
	if err != nil {
		return
	}

	if len(*filters) != 0 {
		var masks []string
		for _, filter := range *filters {
			var exps []string
			for _, field := range filter.Fields {
				exps = append(exps, field+" "+filter.Op+" "+filter.Value)
			}
			masks = append(masks, "("+strings.Join(exps, " or ")+")")
		}
		query += " and "
		query += strings.Join(masks, " and ")
	}

	document, err := documentrepository.GetMeta(session.UserID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, document)
}

// get /api/v1.0/units/documents/
func GetDocuments(w http.ResponseWriter, request *http.Request, r render.Render, documentrepository services.DocumentRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.DocumentSearch), nil, request, r, session.Language)
	if err != nil {
		return
	}
	if len(*filters) != 0 {
		var masks []string
		for _, filter := range *filters {
			var exps []string
			for _, field := range filter.Fields {
				exps = append(exps, field+" "+filter.Op+" "+filter.Value)
			}
			masks = append(masks, "("+strings.Join(exps, " or ")+")")
		}
		query += " and "
		query += strings.Join(masks, " and ")
	}

	var sorts *[]models.OrderExp
	sorts, err = helpers.GetOrderArray(new(models.DocumentSearch), request, r, session.Language)
	if err != nil {
		return
	}
	if len(*sorts) != 0 {
		var orders []string
		for _, sort := range *sorts {
			if strings.ToLower(sort.Field) == "lock" {
				sort.Field = "'" + sort.Field + "'"
			}
			orders = append(orders, " "+sort.Field+" "+sort.Order)
		}
		query += " order by"
		query += strings.Join(orders, ",")
	}

	var limit string
	limit, err = helpers.GetLimitQuery(request, r, session.Language)
	if err != nil {
		return
	}
	query += limit

	documents, err := documentrepository.GetByUser(session.UserID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(documents, len(*documents), w, r)
}

// post /api/v1.0/unit/documents/
func CreateDocument(errors binding.Errors, viewlongdocument models.ViewLongDocument, r render.Render,
	documentrepository services.DocumentRepository, companyrepository services.CompanyRepository, unitrepository services.UnitRepository,
	documenttyperepository services.DocumentTypeRepository, filerepository services.FileRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	_, err := helpers.CheckDocumentType(viewlongdocument.Document_Type_ID, r, documenttyperepository, session.Language)
	if err != nil {
		return
	}
	if viewlongdocument.Company_ID != 0 {
		_, err := helpers.CheckCompanyAvailability(viewlongdocument.Company_ID, session.UserID, r, companyrepository, session.Language)
		if err != nil {
			return
		}
	}
	_, err = filerepository.Get(viewlongdocument.File_ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	unit, err := unitrepository.FindByUser(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	dtodocument := new(models.DtoDocument)
	dtodocument.Document_Type_ID = viewlongdocument.Document_Type_ID
	dtodocument.Unit_ID = unit.ID
	dtodocument.Company_ID = viewlongdocument.Company_ID
	dtodocument.Name = viewlongdocument.Name
	dtodocument.Locked = false
	dtodocument.Pending = false
	dtodocument.File_ID = viewlongdocument.File_ID
	dtodocument.Created = time.Now()
	dtodocument.Updated = time.Now()
	dtodocument.Active = true

	err = documentrepository.Create(dtodocument)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongDocument(dtodocument.ID, dtodocument.Document_Type_ID, dtodocument.Unit_ID, dtodocument.Company_ID,
		dtodocument.Name, dtodocument.Created, dtodocument.Updated, dtodocument.Locked, dtodocument.Pending, dtodocument.File_ID))
}

// post /api/v1.0/unit/documents/matching/
func CreateMatching(errors binding.Errors, viewshortdocument models.ViewShortDocument, r render.Render, emailrepository services.EmailRepository,
	documentrepository services.DocumentRepository, companyrepository services.CompanyRepository, unitrepository services.UnitRepository,
	documenttyperepository services.DocumentTypeRepository, filerepository services.FileRepository,
	templaterepository services.TemplateRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}

	subject := config.Localization[session.Language].Messages.MatchingSubject
	dtodocument, err := helpers.CreateDocumentByType(models.DOCUMENT_TYPE_MATCHING, viewshortdocument.Company_ID,
		fmt.Sprintf(subject, viewshortdocument.Begin_Date, viewshortdocument.End_Date), 0, true, true, r,
		documentrepository, companyrepository, unitrepository, documenttyperepository, filerepository, session)
	if err != nil {
		return
	}

	buf, err := templaterepository.GenerateText(models.NewDtoHTMLTemplate("", session.Language),
		services.TEMPLATE_MATCHING, services.TEMPLATE_DIRECTORY_EMAILS, "")
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	err = emailrepository.SendHTML(config.Configuration.Mail.Receiver, dtodocument.Name, buf.String(), "", config.Configuration.Mail.Sender)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortDocument(dtodocument.ID))
}

// post /api/v1.0/unit/documents/charter/
func CreateCharter(errors binding.Errors, viewmiddledocument models.ViewMiddleDocument, r render.Render,
	documentrepository services.DocumentRepository, companyrepository services.CompanyRepository, unitrepository services.UnitRepository,
	documenttyperepository services.DocumentTypeRepository, filerepository services.FileRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}

	dtodocument, err := helpers.CreateDocumentByType(models.DOCUMENT_TYPE_CHARTER, viewmiddledocument.Company_ID,
		viewmiddledocument.Name, viewmiddledocument.File_ID, false, false, r, documentrepository, companyrepository,
		unitrepository, documenttyperepository, filerepository, session)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortDocument(dtodocument.ID))
}

// post /api/v1.0/unit/documents/extractincorporation/
func CreateExtractIncorporation(errors binding.Errors, viewmiddledocument models.ViewMiddleDocument, r render.Render,
	documentrepository services.DocumentRepository, companyrepository services.CompanyRepository, unitrepository services.UnitRepository,
	documenttyperepository services.DocumentTypeRepository, filerepository services.FileRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}

	dtodocument, err := helpers.CreateDocumentByType(models.DOCUMENT_TYPE_EXTRACTINCORPORATION, viewmiddledocument.Company_ID,
		viewmiddledocument.Name, viewmiddledocument.File_ID, false, false, r, documentrepository, companyrepository,
		unitrepository, documenttyperepository, filerepository, session)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortDocument(dtodocument.ID))
}

//get /api/v1.0/units/documents/:docid/
func GetDocument(r render.Render, params martini.Params, documentrepository services.DocumentRepository, session *models.DtoSession) {
	dtodocument, err := helpers.CheckDocument(r, params, documentrepository, session.Language)
	if err != nil {
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongDocument(dtodocument.ID, dtodocument.Document_Type_ID, dtodocument.Unit_ID, dtodocument.Company_ID,
		dtodocument.Name, dtodocument.Created, dtodocument.Updated, dtodocument.Locked, dtodocument.Pending, dtodocument.File_ID))
}

// delete /api/v1.0/units/documents/:docid/
func DeleteDocument(r render.Render, params martini.Params, documentrepository services.DocumentRepository,
	session *models.DtoSession) {
	dtodocument, err := helpers.CheckDocument(r, params, documentrepository, session.Language)
	if err != nil {
		return
	}
	if dtodocument.Locked {
		log.Error("Document is locked %v", dtodocument.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	err = documentrepository.Deactivate(dtodocument)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
