package services

import (
	"application/models"
)

type DataEncodingRepository interface {
	Get(id int) (dataencoding *models.DtoDataEncoding, err error)
	GetAll() (dataencodings *[]models.ApiDataEncoding, err error)
	FindByName(name string) (id int64, err error)
}

type DataEncodingService struct {
	*Repository
}

func NewDataEncodingService(repository *Repository) *DataEncodingService {
	repository.DbContext.AddTableWithName(models.DtoDataEncoding{}, repository.Table).SetKeys(false, "id")
	return &DataEncodingService{
		repository,
	}
}

func (dataencodingservice *DataEncodingService) Get(id int) (dataencoding *models.DtoDataEncoding, err error) {
	dataencoding = new(models.DtoDataEncoding)
	err = dataencodingservice.DbContext.SelectOne(dataencoding, "select * from "+dataencodingservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting data encoding object from database %v with value %v", err, id)
		return nil, err
	}

	return dataencoding, nil
}

func (dataencodingservice *DataEncodingService) GetAll() (dataencodings *[]models.ApiDataEncoding, err error) {
	dataencodings = new([]models.ApiDataEncoding)
	_, err = dataencodingservice.DbContext.Select(dataencodings, "select * from "+dataencodingservice.Table)
	if err != nil {
		log.Error("Error during getting all data encoding object from database %v", err)
		return nil, err
	}

	return dataencodings, nil
}

func (dataencodingservice *DataEncodingService) FindByName(name string) (id int64, err error) {
	err = dataencodingservice.DbContext.SelectOne(&id, "select id from "+dataencodingservice.Table+" where name = ?", name)
	if err != nil {
		log.Error("Error during finding data encoding object from database %v with value %v", err, name)
		return 0, err
	}

	return id, nil
}
