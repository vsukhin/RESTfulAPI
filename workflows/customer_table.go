package workflows

import (
	"application/config"
	"application/services"
	"time"
)

type CustomerTableWorkflow struct {
	CustomerTableRepository services.CustomerTableRepository
}

func NewCustomerTableWorkflow(customertablerepository services.CustomerTableRepository) *CustomerTableWorkflow {
	return &CustomerTableWorkflow{
		CustomerTableRepository: customertablerepository,
	}
}

func (customertableworkflow *CustomerTableWorkflow) ClearExpired() {
	for {
		tables, err := customertableworkflow.CustomerTableRepository.GetExpired(config.Configuration.TableTimeout)
		if err == nil {
			for _, table := range *tables {
				err = customertableworkflow.CustomerTableRepository.Deactivate(&table)
			}
		}
		time.Sleep(time.Minute)
	}
}
