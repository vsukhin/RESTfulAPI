package services

import (
	"application/models"
)

type OrderStatusRepository interface {
	Get(order_id int64, status_id models.OrderStatus) (order *models.DtoOrderStatus, err error)
	GetByOrder(order_id int64) (orderstatuses *[]models.DtoOrderStatus, err error)
	Create(orderstatus *models.DtoOrderStatus) (err error)
	Update(orderstatus *models.DtoOrderStatus) (err error)
	Save(orderstatus *models.DtoOrderStatus) (err error)
}

type OrderStatusService struct {
	*Repository
}

func NewOrderStatusService(repository *Repository) *OrderStatusService {
	repository.DbContext.AddTableWithName(models.DtoOrderStatus{}, repository.Table).SetKeys(false, "order_id", "status_id")
	return &OrderStatusService{
		repository,
	}
}

func (orderstatusservice *OrderStatusService) Get(order_id int64, status_id models.OrderStatus) (orderstatus *models.DtoOrderStatus, err error) {
	orderstatus = new(models.DtoOrderStatus)
	err = orderstatusservice.DbContext.SelectOne(orderstatus, "select * from "+orderstatusservice.Table+
		" where order_id = ? and status_id = ?", order_id, status_id)
	if err != nil {
		log.Error("Error during getting order status object from database %v with value %v, %v", err, order_id, status_id)
		return nil, err
	}

	return orderstatus, nil
}

func (orderstatusservice *OrderStatusService) GetByOrder(order_id int64) (orderstatuses *[]models.DtoOrderStatus, err error) {
	orderstatuses = new([]models.DtoOrderStatus)
	_, err = orderstatusservice.DbContext.Select(orderstatuses, "select * from "+orderstatusservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all order status object from database %v with value %v", err, order_id)
		return nil, err
	}

	return orderstatuses, nil
}

func (orderstatusservice *OrderStatusService) Create(orderstatus *models.DtoOrderStatus) (err error) {
	err = orderstatusservice.DbContext.Insert(orderstatus)
	if err != nil {
		log.Error("Error during creating order status object in database %v", err)
		return err
	}

	return nil
}

func (orderstatusservice *OrderStatusService) Update(orderstatus *models.DtoOrderStatus) (err error) {
	_, err = orderstatusservice.DbContext.Update(orderstatus)
	if err != nil {
		log.Error("Error during updating order status object in database %v with value %v, %v", err, orderstatus.Order_ID, orderstatus.Status_ID)
		return err
	}

	return nil
}

func (orderstatusservice *OrderStatusService) Save(orderstatus *models.DtoOrderStatus) (err error) {
	count, err := orderstatusservice.DbContext.SelectInt("select count(*) from "+orderstatusservice.Table+
		" where order_id = ? and status_id = ?", orderstatus.Order_ID, orderstatus.Status_ID)
	if err != nil {
		log.Error("Error during saving order status object in database %v with value %v, %v", err, orderstatus.Order_ID, orderstatus.Status_ID)
		return err
	}
	if count == 0 {
		err = orderstatusservice.Create(orderstatus)
	} else {
		err = orderstatusservice.Update(orderstatus)
	}

	return err
}
