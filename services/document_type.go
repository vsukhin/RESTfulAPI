package services

import (
	"application/models"
)

type DocumentTypeRepository interface {
	Get(id int) (documenttype *models.DtoDocumentType, err error)
	GetAll() (documenttypes *[]models.ApiDocumentType, err error)
}

type DocumentTypeService struct {
	*Repository
}

func NewDocumentTypeService(repository *Repository) *DocumentTypeService {
	repository.DbContext.AddTableWithName(models.DtoDocumentType{}, repository.Table).SetKeys(false, "id")
	return &DocumentTypeService{Repository: repository}
}

func (documenttypeservice *DocumentTypeService) Get(id int) (documenttype *models.DtoDocumentType, err error) {
	documenttype = new(models.DtoDocumentType)
	err = documenttypeservice.DbContext.SelectOne(documenttype, "select * from "+documenttypeservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting document type object from database %v with value %v", err, id)
		return nil, err
	}

	return documenttype, nil
}

func (documenttypeservice *DocumentTypeService) GetAll() (documenttypes *[]models.ApiDocumentType, err error) {
	documenttypes = new([]models.ApiDocumentType)
	_, err = documenttypeservice.DbContext.Select(documenttypes,
		"select id, position, name, description from "+documenttypeservice.Table+" where active = 1 order by position asc")
	if err != nil {
		log.Error("Error during getting all document type object from database %v", err)
		return nil, err
	}

	return documenttypes, nil
}
