package services

import (
	"application/models"
)

type TableTypeRepository interface {
	Get(id int) (tabletype *models.DtoTableType, err error)
	GetAll() (tabletypes *[]models.ApiTableType, err error)
	FindByName(name string) (id int, err error)
}

type TableTypeService struct {
	*Repository
}

func NewTableTypeService(repository *Repository) *TableTypeService {
	repository.DbContext.AddTableWithName(models.DtoTableType{}, repository.Table).SetKeys(false, "id")
	return &TableTypeService{
		repository,
	}
}

func (tabletypeservice *TableTypeService) Get(id int) (tabletype *models.DtoTableType, err error) {
	tabletype = new(models.DtoTableType)
	err = tabletypeservice.DbContext.SelectOne(tabletype, "select * from "+tabletypeservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting table type object from database %v with value %v", err, id)
		return nil, err
	}

	return tabletype, nil
}

func (tabletypeservice *TableTypeService) GetAll() (tabletypes *[]models.ApiTableType, err error) {
	tabletypes = new([]models.ApiTableType)
	_, err = tabletypeservice.DbContext.Select(tabletypes, "select id, name from "+tabletypeservice.Table)
	if err != nil {
		log.Error("Error during getting all table type object from database %v", err)
		return nil, err
	}

	return tabletypes, nil
}

func (tabletypeservice *TableTypeService) FindByName(name string) (id int, err error) {
	err = tabletypeservice.DbContext.SelectOne(&id, "select id from "+tabletypeservice.Table+" where name = ?", name)
	if err != nil {
		log.Error("Error during finding table type object from database %v with value %v", err, name)
		return 0, err
	}

	return id, nil
}
