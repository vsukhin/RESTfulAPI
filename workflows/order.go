package workflows

import (
	"application/models"
	"application/services"
	"time"
)

type OrderWorkflow struct {
	OrderRepository    services.OrderRepository
	FacilityRepository services.FacilityRepository
	HeaderWorkflow     Executor
	SMSWorkflow        Executor
	HLRWorkflow        Executor
	VerifyWorkflow     Executor
}

func NewOrderWorkflow(orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	headerworkflow, smsworkflow, hlrworkflow, verifyworkflow Executor) *OrderWorkflow {
	return &OrderWorkflow{
		OrderRepository:    orderrepository,
		FacilityRepository: facilityrepository,
		HeaderWorkflow:     headerworkflow,
		SMSWorkflow:        smsworkflow,
		HLRWorkflow:        hlrworkflow,
		VerifyWorkflow:     verifyworkflow,
	}
}

func (orderworkflow *OrderWorkflow) Execute() {
	for {
		orders, err := orderworkflow.OrderRepository.Get4Processing()
		if err == nil {
			for _, order := range *orders {
				dtofacility, err := orderworkflow.FacilityRepository.Get(order.Facility_ID)
				if err == nil {
					switch dtofacility.Alias {
					case models.SERVICE_TYPE_HEADER:
						go orderworkflow.HeaderWorkflow.ExecuteOrder(order.ID)
					case models.SERVICE_TYPE_SMS:
						go orderworkflow.SMSWorkflow.ExecuteOrder(order.ID)
					case models.SERVICE_TYPE_HLR:
						go orderworkflow.HLRWorkflow.ExecuteOrder(order.ID)
					case models.SERVICE_TYPE_VERIFY:
						go orderworkflow.VerifyWorkflow.ExecuteOrder(order.ID)
					}
				}
			}
		}
		time.Sleep(time.Minute)
	}
}
