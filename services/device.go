package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type DeviceRepository interface {
	Exists(serial string) (found bool, err error)
	FindByHash(hash string) (dtodevice *models.DtoDevice, err error)
	FindByCode(code string) (dtodevice *models.DtoDevice, err error)
	FindByToken(token string) (dtodevice *models.DtoDevice, err error)
	Get(id int64) (dtodevice *models.DtoDevice, err error)
	GetAll() (dtodevices *[]models.DtoDevice, err error)
	Create(dtodevice *models.DtoDevice) (err error)
	Update(dtodevice *models.DtoDevice) (err error)
	DeleteByUser(userid int64, trans *gorp.Transaction) (err error)
}

type DeviceService struct {
	*Repository
}

func NewDeviceService(repository *Repository) *DeviceService {
	repository.DbContext.AddTableWithName(models.DtoDevice{}, repository.Table).SetKeys(true, "id")
	return &DeviceService{
		repository,
	}
}

func (deviceservice *DeviceService) Exists(serial string) (found bool, err error) {
	var count int64
	count, err = deviceservice.DbContext.SelectInt("select count(*) from "+deviceservice.Table+" where serial = ? and user_id != 0 and active = 1", serial)
	if err != nil {
		log.Error("Error during getting device object from database %v with value %v", err, serial)
		return false, err
	}

	return count != 0, nil
}

func (deviceservice *DeviceService) FindByHash(hash string) (dtodevice *models.DtoDevice, err error) {
	dtodevice = new(models.DtoDevice)
	err = deviceservice.DbContext.SelectOne(dtodevice, "select * from "+deviceservice.Table+" where hash = ? and active = 1", hash)
	if err != nil {
		log.Error("Error during finding device object in database %v with value %v", err, hash)
		return nil, err
	}

	return dtodevice, nil
}

func (deviceservice *DeviceService) FindByCode(code string) (dtodevice *models.DtoDevice, err error) {
	dtodevice = new(models.DtoDevice)
	err = deviceservice.DbContext.SelectOne(dtodevice, "select * from "+deviceservice.Table+" where code = ? and active = 1", code)
	if err != nil {
		log.Error("Error during finding device object in database %v with value %v", err, code)
		return nil, err
	}

	return dtodevice, nil
}

func (deviceservice *DeviceService) FindByToken(token string) (dtodevice *models.DtoDevice, err error) {
	dtodevice = new(models.DtoDevice)
	err = deviceservice.DbContext.SelectOne(dtodevice, "select * from "+deviceservice.Table+" where token = ? and active = 1", token)
	if err != nil {
		log.Error("Error during finding device object in database %v with value %v", err, token)
		return nil, err
	}

	return dtodevice, nil
}

func (deviceservice *DeviceService) Get(id int64) (dtodevice *models.DtoDevice, err error) {
	dtodevice = new(models.DtoDevice)
	err = deviceservice.DbContext.SelectOne(dtodevice, "select * from "+deviceservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting device object from database %v with value %v", err, id)
		return nil, err
	}

	return dtodevice, nil
}

func (deviceservice *DeviceService) GetAll() (dtodevices *[]models.DtoDevice, err error) {
	dtodevices = new([]models.DtoDevice)
	_, err = deviceservice.DbContext.Select(dtodevices, "select * from "+deviceservice.Table+" where user_id != 0 and active = 1")
	if err != nil {
		log.Error("Error during getting all device object in database %v", err)
		return nil, err
	}

	return dtodevices, nil
}

func (deviceservice *DeviceService) Create(dtodevice *models.DtoDevice) (err error) {
	err = deviceservice.DbContext.Insert(dtodevice)
	if err != nil {
		log.Error("Error during creating device object in database %v", err)
		return err
	}

	return nil
}

func (deviceservice *DeviceService) Update(dtodevice *models.DtoDevice) (err error) {
	_, err = deviceservice.DbContext.Update(dtodevice)
	if err != nil {
		log.Error("Error during updating device object in database %v with value %v", err, dtodevice.ID)
		return err
	}

	return nil
}

func (deviceservice *DeviceService) DeleteByUser(userid int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+deviceservice.Table+" where user_id = ?", userid)
	} else {
		_, err = deviceservice.DbContext.Exec("delete from "+deviceservice.Table+" where user_id = ?", userid)
	}
	if err != nil {
		log.Error("Error during deleting device object in database %v with value %v", err, userid)
		return err
	}

	return nil
}
