package services

import (
	"application/models"
)

type PeriodRepository interface {
	Get(id int) (period *models.DtoPeriod, err error)
	GetAll() (periods *[]models.ApiPeriod, err error)
}

type PeriodService struct {
	*Repository
}

func NewPeriodService(repository *Repository) *PeriodService {
	repository.DbContext.AddTableWithName(models.DtoPeriod{}, repository.Table).SetKeys(true, "id")
	return &PeriodService{Repository: repository}
}

func (periodservice *PeriodService) Get(id int) (period *models.DtoPeriod, err error) {
	period = new(models.DtoPeriod)
	err = periodservice.DbContext.SelectOne(period, "select * from "+periodservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting period object from database %v with value %v", err, id)
		return nil, err
	}

	return period, nil
}

func (periodservice *PeriodService) GetAll() (periods *[]models.ApiPeriod, err error) {
	periods = new([]models.ApiPeriod)
	_, err = periodservice.DbContext.Select(periods, "select id, name, position from "+periodservice.Table+" where active = 1 order by position asc")
	if err != nil {
		log.Error("Error during getting all period object from database %v", err)
		return nil, err
	}

	return periods, nil
}
