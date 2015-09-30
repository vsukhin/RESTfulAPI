package services

import (
	"application/models"
)

type TariffPlanRepository interface {
	Get(id int) (tariffplan *models.DtoTariffPlan, err error)
	GetAll() (tariffplans *[]models.ApiTariffPlan, err error)
}

type TariffPlanService struct {
	*Repository
}

func NewTariffPlanService(repository *Repository) *TariffPlanService {
	repository.DbContext.AddTableWithName(models.DtoTariffPlan{}, repository.Table).SetKeys(false, "id")
	return &TariffPlanService{Repository: repository}
}

func (tariffplanservice *TariffPlanService) Get(id int) (tariffplan *models.DtoTariffPlan, err error) {
	tariffplan = new(models.DtoTariffPlan)
	err = tariffplanservice.DbContext.SelectOne(tariffplan, "select * from "+tariffplanservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting tariff plan object from database %v with value %v", err, id)
		return nil, err
	}

	return tariffplan, nil
}

func (tariffplanservice *TariffPlanService) GetAll() (tariffplans *[]models.ApiTariffPlan, err error) {
	tariffplans = new([]models.ApiTariffPlan)
	_, err = tariffplanservice.DbContext.Select(tariffplans,
		"select id, name, position, public from "+tariffplanservice.Table+" where active = 1 and public = 1 order by position asc")
	if err != nil {
		log.Error("Error during getting all tariff plan object from database %v", err)
		return nil, err
	}

	return tariffplans, nil
}
