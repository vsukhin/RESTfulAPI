package services

import (
	"application/models"
)

type AddressTypeRepository interface {
	Get(id int) (addresstype *models.DtoAddressType, err error)
	GetAll() (addresstypes *[]models.ApiAddressType, err error)
}

type AddressTypeService struct {
	*Repository
}

func NewAddressTypeService(repository *Repository) *AddressTypeService {
	repository.DbContext.AddTableWithName(models.DtoAddressType{}, repository.Table).SetKeys(false, "id")
	return &AddressTypeService{Repository: repository}
}

func (addresstypeservice *AddressTypeService) Get(id int) (addresstype *models.DtoAddressType, err error) {
	addresstype = new(models.DtoAddressType)
	err = addresstypeservice.DbContext.SelectOne(addresstype, "select * from "+addresstypeservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting address type object from database %v with value %v", err, id)
		return nil, err
	}

	return addresstype, nil
}

func (addresstypeservice *AddressTypeService) GetAll() (addresstypes *[]models.ApiAddressType, err error) {
	addresstypes = new([]models.ApiAddressType)
	_, err = addresstypeservice.DbContext.Select(addresstypes,
		"select id, name, required, position from "+addresstypeservice.Table+" where active = 1 order by position asc")
	if err != nil {
		log.Error("Error during getting all address type object from database %v", err)
		return nil, err
	}

	return addresstypes, nil
}
