package services

import (
	"application/models"
)

type DocumentRepository interface {
	CheckUserAccess(user_id int64, id int64) (allowed bool, err error)
	Get(id int64) (document *models.DtoDocument, err error)
	GetMeta(user_id int64, filter string) (document *models.ApiMetaDocument, err error)
	GetByUser(userid int64, filter string) (documents *[]models.ApiLongDocument, err error)
	Create(document *models.DtoDocument) (err error)
	Update(document *models.DtoDocument) (err error)
	Deactivate(document *models.DtoDocument) (err error)
}

type DocumentService struct {
	*Repository
}

func NewDocumentService(repository *Repository) *DocumentService {
	repository.DbContext.AddTableWithName(models.DtoDocument{}, repository.Table).SetKeys(true, "id")
	return &DocumentService{Repository: repository}
}

func (documentservice *DocumentService) CheckUserAccess(user_id int64, id int64) (allowed bool, err error) {
	count, err := documentservice.DbContext.SelectInt("select count(*) from "+documentservice.Table+
		" where id = ? and unit_id = (select unit_id from users where id = ?)", id, user_id)
	if err != nil {
		log.Error("Error during checking document object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (documentservice *DocumentService) Get(id int64) (document *models.DtoDocument, err error) {
	document = new(models.DtoDocument)
	err = documentservice.DbContext.SelectOne(document, "select * from "+documentservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting document object from database %v with value %v", err, id)
		return nil, err
	}

	return document, nil
}

func (documentservice *DocumentService) GetMeta(user_id int64, filter string) (document *models.ApiMetaDocument, err error) {
	document = new(models.ApiMetaDocument)
	document.Total, err = documentservice.DbContext.SelectInt(
		"select count(*) from "+documentservice.Table+" where active = 1 and unit_id = (select unit_id from users where id = ?)"+filter, user_id)
	if err != nil {
		log.Error("Error during getting meta document object from database %v with value %v", err, user_id)
		return nil, err
	}

	return document, nil
}

func (documentservice *DocumentService) GetByUser(userid int64, filter string) (documents *[]models.ApiLongDocument, err error) {
	documents = new([]models.ApiLongDocument)
	_, err = documentservice.DbContext.Select(documents,
		"select id, document_type_id as categoryId, unit_id as unitId, company_id as organisationId, name, created, updated as edited,"+
			" locked as `lock`, pending, file_id as fileId from "+documentservice.Table+
			" where unit_id = (select unit_id from users where id = ?) and active = 1"+filter, userid)
	if err != nil {
		log.Error("Error during getting unit document object from database %v with value %v", err, userid)
		return nil, err
	}

	return documents, nil
}

func (documentservice *DocumentService) Create(document *models.DtoDocument) (err error) {
	err = documentservice.DbContext.Insert(document)
	if err != nil {
		log.Error("Error during creating document object in database %v", err)
		return err
	}

	return nil
}

func (documentservice *DocumentService) Update(document *models.DtoDocument) (err error) {
	_, err = documentservice.DbContext.Update(document)
	if err != nil {
		log.Error("Error during updating document object in database %v with value %v", err, document.ID)
		return err
	}

	return nil
}

func (documentservice *DocumentService) Deactivate(document *models.DtoDocument) (err error) {
	_, err = documentservice.DbContext.Exec("update "+documentservice.Table+" set active = 0 where id = ?", document.ID)
	if err != nil {
		log.Error("Error during deactivating document object in database %v with value %v", err, document.ID)
		return err
	}

	return nil
}
