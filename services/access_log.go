package services

import (
	"application/models"
)

type AccessLogRepository interface {
	Get(id int64) (dtoaccesslog *models.DtoAccessLog, err error)
	Create(dtoaccesslog *models.DtoAccessLog) (err error)
}

type AccessLogService struct {
	*Repository
}

func NewAccessLogService(repository *Repository) *AccessLogService {
	repository.DbContext.AddTableWithName(models.DtoAccessLog{}, repository.Table).SetKeys(true, "id")
	return &AccessLogService{
		repository,
	}
}

func (accesslogservice *AccessLogService) Get(id int64) (dtoaccesslog *models.DtoAccessLog, err error) {
	dtoaccesslog = new(models.DtoAccessLog)
	err = accesslogservice.DbContext.SelectOne(dtoaccesslog, "select * from "+accesslogservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting access log object from database %v with value %v", err, id)
		return nil, err
	}

	return dtoaccesslog, nil
}

func (accesslogservice *AccessLogService) Create(dtoaccesslog *models.DtoAccessLog) (err error) {
	err = accesslogservice.DbContext.Insert(dtoaccesslog)
	if err != nil {
		log.Error("Error during creating access log object in database %v", err)
		return err
	}

	return nil
}
