package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type OrderRepository interface {
	CheckUserAccess(user_id int64, id int64) (allowed bool, err error)
	CheckSupplierAccess(user_id int64, id int64) (allowed bool, err error)
	Get(id int64) (order *models.DtoOrder, err error)
	GetAll(user_id int64, filter string) (orders *[]models.ApiShortOrder, err error)
	GetMeta(user_id int64) (order *models.ApiMetaOrder, err error)
	Create(order *models.DtoOrder) (err error)
	Update(order *models.DtoOrder, orderstatuses *[]models.DtoOrderStatus, inTrans bool) (err error)
}

type OrderService struct {
	OrderStatusRepository OrderStatusRepository
	*Repository
}

func NewOrderService(repository *Repository) *OrderService {
	repository.DbContext.AddTableWithName(models.DtoOrder{}, repository.Table).SetKeys(true, "id")
	return &OrderService{Repository: repository}
}

func (orderservice *OrderService) CheckUserAccess(user_id int64, id int64) (allowed bool, err error) {
	count, err := orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" where id = ? and (unit_id = (select unit_id from users where id = ?) or supplier_id = "+
		"(select unit_id from users where id = ? and active = 1 and confirmed = 1))",
		id, user_id, user_id)
	if err != nil {
		log.Error("Error during checking order object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (orderservice *OrderService) CheckSupplierAccess(user_id int64, id int64) (allowed bool, err error) {
	count, err := orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" where id = ? and supplier_id = (select unit_id from users where id = ? and active = 1 and confirmed = 1)", id, user_id)
	if err != nil {
		log.Error("Error during checking order object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (orderservice *OrderService) Get(id int64) (order *models.DtoOrder, err error) {
	order = new(models.DtoOrder)
	err = orderservice.DbContext.SelectOne(order, "select * from "+orderservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting order object from database %v with value %v", err, id)
		return nil, err
	}

	return order, nil
}

func (orderservice *OrderService) GetAll(user_id int64, filter string) (orders *[]models.ApiShortOrder, err error) {
	orders = new([]models.ApiShortOrder)
	_, err = orderservice.DbContext.Select(orders, "select o.id, o.name, o.service_id as type, o.supplier_id as supplierId,"+
		"coalesce(c.value,0) as completed, coalesce(n.value, 0) as new, coalesce(p.value, 0) as open from "+orderservice.Table+" o"+
		" left join order_statuses c on o.id = c.order_id and c.status_id = 1"+
		" left join order_statuses n on o.id = n.order_id and n.status_id = 2"+
		" left join order_statuses p on o.id = p.order_id and p.status_id = 3"+
		" where o.supplier_id = (select unit_id from users where id = ? and active = 1 and confirmed = 1)"+filter, user_id)
	if err != nil {
		log.Error("Error during getting all order object from database %v with value %v", err, user_id)
		return nil, err
	}

	return orders, nil
}

func (orderservice *OrderService) GetMeta(user_id int64) (order *models.ApiMetaOrder, err error) {
	order = new(models.ApiMetaOrder)
	order.Total, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" where supplier_id = (select unit_id from users where id = ? and active = 1 and confirmed = 1)", user_id)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, user_id)
		return nil, err
	}
	order.NumOfNew, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where o.supplier_id = (select unit_id from users where id = ? and active = 1 and confirmed = 1)"+
		" and s.status_id = ? and s.value = 1", user_id, models.ORDER_STATUS_NEW)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, user_id)
		return nil, err
	}
	order.NumOfOpen, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where o.supplier_id = (select unit_id from users where id = ? and active = 1 and confirmed = 1)"+
		" and s.status_id = ? and s.value = 1", user_id, models.ORDER_STATUS_OPEN)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, user_id)
		return nil, err
	}
	order.NumOfClosed, err = orderservice.DbContext.SelectInt("select count(distinct o.id) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where o.supplier_id = (select unit_id from users where id = ? and active = 1 and confirmed = 1)"+
		" and s.status_id in (?, ?, ?) and s.value = 1", user_id, models.ORDER_STATUS_CANCEL,
		models.ORDER_STATUS_SUPPLIER_CLOSE, models.ORDER_STATUS_MODERATOR_CLOSE)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, user_id)
		return nil, err
	}
	order.NumOfArchived, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where o.supplier_id = (select unit_id from users where id = ? and active = 1 and confirmed = 1)"+
		" and s.status_id = ? and s.value = 1", user_id, models.ORDER_STATUS_ARCHIVE)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, user_id)
		return nil, err
	}
	order.NumOfAlert, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o where id in (select order_id from messages m where m.user_id in (select id from users where unit_id = o.unit_id)"+
		" and m.id not in (select message_id from user_messages um where um.user_id in "+
		"(select id from users where unit_id = o.supplier_id))) and o.supplier_id = "+
		"(select unit_id from users where id = ? and active = 1 and confirmed = 1)", user_id)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, user_id)
		return nil, err
	}

	return order, nil
}

func (orderservice *OrderService) Create(order *models.DtoOrder) (err error) {
	err = orderservice.DbContext.Insert(order)
	if err != nil {
		log.Error("Error during creating order object in database %v", err)
		return err
	}

	return nil
}

func (orderservice *OrderService) Update(order *models.DtoOrder, orderstatuses *[]models.DtoOrderStatus, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = orderservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating order object in database %v", err)
			return err
		}
	}

	for _, orderstatus := range *orderstatuses {
		err = orderservice.OrderStatusRepository.Save(&orderstatus)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during updating order object in database %v with value %v", err, order.ID)
			return err
		}
	}

	_, err = orderservice.DbContext.Update(order)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating order object in database %v with value %v", err, order.ID)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during updating order object in database %v", err)
			return err
		}
	}

	return nil
}
