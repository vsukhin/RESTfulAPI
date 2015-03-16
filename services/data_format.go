package services

import (
	"application/models"
)

type DataFormatRepository interface {
	Get(id int) (dataformat *models.DtoDataFormat, err error)
	GetAll() (dataformats *[]models.ApiDataFormat, err error)
	FindByName(name string) (id int64, err error)
}

type DataFormatService struct {
	*Repository
}

func NewDataFormatService(repository *Repository) *DataFormatService {
	repository.DbContext.AddTableWithName(models.DtoDataFormat{}, repository.Table).SetKeys(false, "id")
	return &DataFormatService{
		repository,
	}
}

func (dataformatservice *DataFormatService) Get(id int) (dataformat *models.DtoDataFormat, err error) {
	dataformat = new(models.DtoDataFormat)
	err = dataformatservice.DbContext.SelectOne(dataformat, "select * from "+dataformatservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting data format object from database %v with value %v", err, id)
		return nil, err
	}

	return dataformat, nil
}

func (dataformatservice *DataFormatService) GetAll() (dataformats *[]models.ApiDataFormat, err error) {
	dataformats = new([]models.ApiDataFormat)
	_, err = dataformatservice.DbContext.Select(dataformats, "select * from "+dataformatservice.Table)
	if err != nil {
		log.Error("Error during getting all data format object from database %v", err)
		return nil, err
	}

	return dataformats, nil
}

func (dataformatservice *DataFormatService) FindByName(name string) (id int64, err error) {
	err = dataformatservice.DbContext.SelectOne(&id, "select id from "+dataformatservice.Table+" where name = ?", name)
	if err != nil {
		log.Error("Error during finding data format object from database %v with value %v", err, name)
		return 0, err
	}

	return id, nil
}
