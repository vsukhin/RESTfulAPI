package services

import (
	"application/models"
)

type RequestRepository interface {
	Exists(ip_address string, method string) (found bool, err error)
	Get(ip_address string, method string) (request *models.DtoRequest, err error)
	Create(request *models.DtoRequest) (err error)
	Update(request *models.DtoRequest) (err error)
	Save(request *models.DtoRequest) (err error)
}

type RequestService struct {
	*Repository
}

func NewRequestService(repository *Repository) *RequestService {
	repository.DbContext.AddTableWithName(models.DtoRequest{}, repository.Table).SetKeys(false, "ip_address", "method")
	return &RequestService{
		repository,
	}
}

func (requestservice *RequestService) Exists(ip_address string, method string) (found bool, err error) {
	var count int64
	count, err = requestservice.DbContext.SelectInt("select count(*) from "+requestservice.Table+
		" where ip_address = ? and method = ?", ip_address, method)
	if err != nil {
		log.Error("Error during getting request object from database %v with value %v, %v", err, ip_address, method)
		return false, err
	}

	return count != 0, nil
}

func (requestservice *RequestService) Get(ip_address string, method string) (request *models.DtoRequest, err error) {
	request = new(models.DtoRequest)
	err = requestservice.DbContext.SelectOne(request, "select * from "+requestservice.Table+
		" where ip_address = ? and method = ?", ip_address, method)
	if err != nil {
		log.Error("Error during getting request object from database %v with value %v, %v", err, ip_address, method)
		return nil, err
	}

	return request, nil
}

func (requestservice *RequestService) Create(request *models.DtoRequest) (err error) {
	err = requestservice.DbContext.Insert(request)
	if err != nil {
		log.Error("Error during creating request object in database %v", err)
		return err
	}

	return nil
}

func (requestservice *RequestService) Update(request *models.DtoRequest) (err error) {
	_, err = requestservice.DbContext.Update(request)
	if err != nil {
		log.Error("Error during updating request object in database %v with value %v, %v", err, request.IP_Address, request.Method)
		return err
	}

	return nil
}

func (requestservice *RequestService) Save(request *models.DtoRequest) (err error) {
	count, err := requestservice.DbContext.SelectInt("select count(*) from "+requestservice.Table+
		" where ip_address = ? and method = ?", request.IP_Address, request.Method)
	if err != nil {
		log.Error("Error during saving request object in database %v with value %v, %v", err, request.IP_Address, request.Method)
		return err
	}
	if count == 0 {
		err = requestservice.Create(request)
	} else {
		err = requestservice.Update(request)
	}

	return err
}
