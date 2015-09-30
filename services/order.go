package services

import (
	"application/models"
	"fmt"
	"github.com/coopernurse/gorp"
)

const (
	ORDER_NAME_TEMPLATE = "Заказ №%v"
)

type OrderRepository interface {
	CheckUserAccess(user_id int64, id int64) (allowed bool, err error)
	CheckUnitAccess(unit_id int64, id int64) (allowed bool, err error)
	CheckSupplierAccess(user_id int64, id int64) (allowed bool, err error)
	IsConfirmed(id int64) (confirmed bool, err error)
	Get(id int64) (order *models.DtoOrder, err error)
	GetByUser(user_id int64, filter string) (orders *[]models.ApiShortOrder, err error)
	GetByProject(project_id int64, filter string) (orders *[]models.ApiMiddleOrder, err error)
	GetByUnit(unit_id int64) (orders *[]models.ApiBriefOrder, err error)
	GetAll(filter string) (orders *[]models.ApiListOrder, err error)
	Get4Processing() (orders *[]models.ApiTinyOrder, err error)
	GetFinance(unit_id int64, filter string) (orders *[]models.ApiFinanceOrder, err error)
	GetResult(unit_id int64, filter string) (order *models.ApiResultOrder, err error)
	GetMeta(user_id int64) (order *models.ApiMetaOrder, err error)
	GetMetaByProject(project_id int64) (order *models.ApiMetaOrderByProject, err error)
	GetFullMeta() (order *models.ApiFullMetaOrder, err error)
	Create(order *models.DtoOrder, orderstatuses *[]models.DtoOrderStatus, trans *gorp.Transaction, inTrans bool) (err error)
	Update(order *models.DtoOrder, orderstatuses *[]models.DtoOrderStatus, trans *gorp.Transaction, inTrans bool) (err error)
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
		" where id = ? and (unit_id = (select unit_id from users where id = ?)"+
		" or (supplier_id != 0 and supplier_id = (select unit_id from users where id = ?))"+
		" or (supplier_id = 0 and (select unit_id from users where id = ?) in (select supplier_id from supplier_requests where order_id = id)))",
		id, user_id, user_id, user_id)
	if err != nil {
		log.Error("Error during checking order object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (orderservice *OrderService) CheckUnitAccess(unit_id int64, id int64) (allowed bool, err error) {
	count, err := orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" where id = ? and (unit_id = ? or (supplier_id != 0 and supplier_id = ?)"+
		" or (supplier_id = 0 and ? in (select supplier_id from supplier_requests where order_id = id)))",
		id, unit_id, unit_id, unit_id)
	if err != nil {
		log.Error("Error during checking order object from database %v with value %v, %v", err, unit_id, id)
		return false, err
	}

	return count != 0, nil
}

func (orderservice *OrderService) CheckSupplierAccess(user_id int64, id int64) (allowed bool, err error) {
	count, err := orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" where id = ? and ((supplier_id != 0 and supplier_id = (select unit_id from users where id = ?))"+
		" or (supplier_id = 0 and (select unit_id from users where id = ?) in (select supplier_id from supplier_requests where order_id = id)))",
		id, user_id, user_id)
	if err != nil {
		log.Error("Error during checking order object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (orderservice *OrderService) IsConfirmed(id int64) (confirmed bool, err error) {
	count, err := orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" where id = ? and id in (select order_id from order_statuses where status_id = "+
		fmt.Sprintf("%v", models.ORDER_STATUS_MODERATOR_CONFIRMED)+" and value = 1)", id)
	if err != nil {
		log.Error("Error during checking order object from database %v with value %v", err, id)
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

func (orderservice *OrderService) GetByUser(user_id int64, filter string) (orders *[]models.ApiShortOrder, err error) {
	orders = new([]models.ApiShortOrder)
	_, err = orderservice.DbContext.Select(orders, "select o.id, o.name, o.service_id as type, o.supplier_id as supplierId,"+
		"coalesce(c.value, 0) as completed, coalesce(n.value, 0) as new, coalesce(p.value, 0) as open from "+orderservice.Table+" o"+
		" left join order_statuses c on o.id = c.order_id and c.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_COMPLETED)+
		" left join order_statuses n on o.id = n.order_id and n.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_NEW)+
		" left join order_statuses p on o.id = p.order_id and p.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_OPEN)+
		" where (o.id not in (select order_id from order_statuses where status_id in ("+
		fmt.Sprintf("%v, %v", models.ORDER_STATUS_ARCHIVE, models.ORDER_STATUS_DEL)+") and value = 1)) and"+
		" ((o.supplier_id != 0 and o.supplier_id = (select unit_id from users where id = ?))"+
		" or (o.supplier_id = 0 and (select unit_id from users where id = ?) in (select supplier_id from supplier_requests where order_id = o.id)))"+filter,
		user_id, user_id)
	if err != nil {
		log.Error("Error during getting all order object from database %v with value %v", err, user_id)
		return nil, err
	}

	return orders, nil
}

func (orderservice *OrderService) GetByProject(project_id int64, filter string) (orders *[]models.ApiMiddleOrder, err error) {
	orders = new([]models.ApiMiddleOrder)
	_, err = orderservice.DbContext.Select(orders, "select o.id, o.name, o.step, o.service_id as type, o.supplier_id as supplierId,"+
		" o.execution_forecast as supplierForecastWorkDays, o.proposed_price as supplierCost, o.charged_fee as supplierFactualCost,"+
		" coalesce(c.value, 0) as completed, coalesce(m.value, 0) as moderatorConfirmed, coalesce(n.value, 0) as new, coalesce(v.value, 0) as open,"+
		" coalesce(a.value, 0) as cancel, coalesce(a.comments, '') as cancelDescription, coalesce(s.value, 0) as supplierCostNew,"+
		" coalesce(u.value, 0) as customerNewCostConfirmed, coalesce(p.value, 0) as paid, coalesce(d.value, 0) as moderatorBegin,"+
		" coalesce(l.value, 0) as supplierFactualCostNew, coalesce(i.value, 0) as supplierClose, coalesce(e.value, 0) as moderatorDocumentsGotten,"+
		" coalesce(r.value, 0) as moderatorClose, coalesce(h.value, 0) as archive from "+orderservice.Table+" o"+
		" left join order_statuses c on o.id = c.order_id and c.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_COMPLETED)+
		" left join order_statuses m on o.id = m.order_id and m.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_MODERATOR_CONFIRMED)+
		" left join order_statuses n on o.id = n.order_id and n.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_NEW)+
		" left join order_statuses v on o.id = v.order_id and v.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_OPEN)+
		" left join order_statuses a on o.id = a.order_id and a.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_CANCEL)+
		" left join order_statuses s on o.id = s.order_id and s.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_SUPPLIER_COST_NEW)+
		" left join order_statuses u on o.id = u.order_id and u.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED)+
		" left join order_statuses p on o.id = p.order_id and p.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_PAID)+
		" left join order_statuses d on o.id = d.order_id and d.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_MODERATOR_BEGIN)+
		" left join order_statuses l on o.id = l.order_id and l.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_SUPPLIER_FACTUAL_COST_NEW)+
		" left join order_statuses i on o.id = i.order_id and i.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_SUPPLIER_CLOSE)+
		" left join order_statuses e on o.id = e.order_id and e.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_MODERATOR_DOCUMENTS_GOTTEN)+
		" left join order_statuses r on o.id = r.order_id and r.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_MODERATOR_CLOSE)+
		" left join order_statuses h on o.id = h.order_id and h.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_ARCHIVE)+
		" where (o.id not in (select order_id from order_statuses where status_id in ("+
		fmt.Sprintf("%v", models.ORDER_STATUS_DEL)+") and value = 1)) and"+
		" (o.project_id = ?)"+filter, project_id)
	if err != nil {
		log.Error("Error during getting all order object from database %v with value %v", err, project_id)
		return nil, err
	}

	return orders, nil
}

func (orderservice *OrderService) GetByUnit(unit_id int64) (orders *[]models.ApiBriefOrder, err error) {
	orders = new([]models.ApiBriefOrder)
	_, err = orderservice.DbContext.Select(orders, "select o.id, o.name, coalesce(p.value, 0) as paid,"+
		" coalesce(a.value, 0) as archive, coalesce(d.value, 0) as del from "+orderservice.Table+" o"+
		" left join order_statuses p on o.id = p.order_id and p.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_PAID)+
		" left join order_statuses a on o.id = a.order_id and a.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_ARCHIVE)+
		" left join order_statuses d on o.id = d.order_id and d.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_DEL)+
		" where o.unit_id = ? or (o.supplier_id != 0 and o.supplier_id = ?)"+
		" or (o.supplier_id = 0 and ? in (select supplier_id from supplier_requests where order_id = o.id))", unit_id, unit_id, unit_id)
	if err != nil {
		log.Error("Error during getting all order object from database %v with value %v", err, unit_id)
		return nil, err
	}

	return orders, nil
}

func (orderservice *OrderService) GetAll(filter string) (orders *[]models.ApiListOrder, err error) {
	orders = new([]models.ApiListOrder)
	_, err = orderservice.DbContext.Select(orders, "select o.id, o.name, o.step, o.service_id as type, o.supplier_id as supplierId,"+
		" o.unit_id as unitId, o.user_id as customerId, o.charged_fee as cost, coalesce(c.value, 0) as completed,"+
		" coalesce(n.value, 0) as new, coalesce(p.value, 0) as open, coalesce(a.value, 0) as cancel,"+
		" coalesce(i.value, 0) as paid, coalesce(r.value, 0) as archive, coalesce(e.value, 0) as del from "+orderservice.Table+" o"+
		" left join order_statuses c on o.id = c.order_id and c.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_COMPLETED)+
		" left join order_statuses n on o.id = n.order_id and n.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_NEW)+
		" left join order_statuses p on o.id = p.order_id and p.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_OPEN)+
		" left join order_statuses a on o.id = a.order_id and a.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_CANCEL)+
		" left join order_statuses i on o.id = i.order_id and i.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_PAID)+
		" left join order_statuses r on o.id = r.order_id and r.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_ARCHIVE)+
		" left join order_statuses e on o.id = e.order_id and e.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_DEL)+
		filter)
	if err != nil {
		log.Error("Error during getting all order object from database %v", err)
		return nil, err
	}

	return orders, nil
}

func (orderservice *OrderService) Get4Processing() (orders *[]models.ApiTinyOrder, err error) {
	orders = new([]models.ApiTinyOrder)
	_, err = orderservice.DbContext.Select(orders, "select id, service_id from "+orderservice.Table+
		" where (id not in (select order_id from order_statuses where status_id in ("+
		fmt.Sprintf("%v, %v, %v, %v, %v", models.ORDER_STATUS_OPEN, models.ORDER_STATUS_CANCEL, models.ORDER_STATUS_MODERATOR_CLOSE,
			models.ORDER_STATUS_ARCHIVE, models.ORDER_STATUS_DEL)+") and value = 1))"+
		" and (id in (select order_id from order_statuses where status_id in ("+
		fmt.Sprintf("%v, %v", models.ORDER_STATUS_COMPLETED, models.ORDER_STATUS_MODERATOR_CONFIRMED)+") and value = 1))")
	if err != nil {
		log.Error("Error during getting all order object from database %v ", err)
		return nil, err
	}

	return orders, nil
}

func (orderservice *OrderService) GetFinance(unit_id int64, filter string) (orders *[]models.ApiFinanceOrder, err error) {
	orders = new([]models.ApiFinanceOrder)
	_, err = orderservice.DbContext.Select(orders, "select o.id as orderId, o.project_id as projectId, o.begin_date as beginDate,"+
		" o.charged_fee as cost, o.service_id as type, o.supplier_id as supplierId, coalesce(a.file_id, 0) as actFileId,"+
		" coalesce(i.file_id, 0) as invoiceFileId, coalesce(s.complex_status_id, 0) as statusId from "+orderservice.Table+" o"+
		" left join documents a on o.act_id = a.id left join documents i on o.eInvoice_id = i.id left join order_complex_statuses s on o.id = s.order_id"+
		" where (o.unit_id = ?)"+filter, unit_id)
	if err != nil {
		log.Error("Error during getting unit order object from database %v with value %v", err, unit_id)
		return nil, err
	}

	return orders, nil
}

func (orderservice *OrderService) GetResult(unit_id int64, filter string) (order *models.ApiResultOrder, err error) {
	order = new(models.ApiResultOrder)
	order.Total, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o left join documents a on o.act_id = a.id left join documents i on o.eInvoice_id = i.id"+
		" left join order_complex_statuses s on o.id = s.order_id where (o.unit_id = ?)"+filter, unit_id)
	if err != nil {
		log.Error("Error during getting unit order object from database %v with value %v", err, unit_id)
		return nil, err
	}
	order.Processing, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses t on o.id = t.order_id left join documents a on o.act_id = a.id"+
		" left join documents i on o.eInvoice_id = i.id left join order_complex_statuses s on o.id = s.order_id"+
		" where (t.status_id = ? and t.value = 1) and o.id not in (select order_id from order_statuses where status_id = ? and value = 1)"+
		" and (o.unit_id = ?)"+filter,
		models.ORDER_STATUS_MODERATOR_BEGIN, models.ORDER_STATUS_MODERATOR_DOCUMENTS_GOTTEN, unit_id)
	if err != nil {
		log.Error("Error during getting unit order object from database %v with value %v", err, unit_id)
		return nil, err
	}
	order.Charged_Fee, err = orderservice.DbContext.SelectFloat("select coalesce(sum(o.charged_fee), 0) from "+orderservice.Table+
		" o left join documents a on o.act_id = a.id left join documents i on o.eInvoice_id = i.id"+
		" left join order_complex_statuses s on o.id = s.order_id where (o.unit_id = ?)"+filter, unit_id)
	if err != nil {
		log.Error("Error during getting unit order object from database %v with value %v", err, unit_id)
		return nil, err
	}
	order.Charged_Fee = models.Round(order.Charged_Fee, 0.5, 2)

	return order, nil
}

func (orderservice *OrderService) GetMeta(user_id int64) (order *models.ApiMetaOrder, err error) {
	order = new(models.ApiMetaOrder)
	order.Total, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" where (supplier_id != 0 and supplier_id = (select unit_id from users where id = ?))"+
		" or (supplier_id = 0 and (select unit_id from users where id = ?) in (select supplier_id from supplier_requests where order_id = id))",
		user_id, user_id)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, user_id)
		return nil, err
	}
	order.NumOfNew, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where ((o.supplier_id != 0 and o.supplier_id = (select unit_id from users where id = ?))"+
		" or (o.supplier_id = 0 and (select unit_id from users where id = ?) in (select supplier_id from supplier_requests where order_id = o.id)))"+
		" and s.status_id = ? and s.value = 1", user_id, user_id, models.ORDER_STATUS_NEW)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, user_id)
		return nil, err
	}
	order.NumOfOpen, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where ((o.supplier_id != 0 and o.supplier_id = (select unit_id from users where id = ?))"+
		" or (o.supplier_id = 0 and (select unit_id from users where id = ?) in (select supplier_id from supplier_requests where order_id = o.id)))"+
		" and s.status_id = ? and s.value = 1", user_id, user_id, models.ORDER_STATUS_OPEN)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, user_id)
		return nil, err
	}
	order.NumOfClosed, err = orderservice.DbContext.SelectInt("select count(distinct o.id) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where ((o.supplier_id != 0 and o.supplier_id = (select unit_id from users where id = ?))"+
		" or (o.supplier_id = 0 and (select unit_id from users where id = ?) in (select supplier_id from supplier_requests where order_id = o.id)))"+
		" and s.status_id in (?, ?, ?) and s.value = 1", user_id, user_id, models.ORDER_STATUS_CANCEL,
		models.ORDER_STATUS_SUPPLIER_CLOSE, models.ORDER_STATUS_MODERATOR_CLOSE)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, user_id)
		return nil, err
	}
	order.NumOfArchived, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where ((o.supplier_id != 0 and o.supplier_id = (select unit_id from users where id = ?))"+
		" or (o.supplier_id = 0 and (select unit_id from users where id = ?) in (select supplier_id from supplier_requests where order_id = o.id)))"+
		" and s.status_id = ? and s.value = 1", user_id, user_id, models.ORDER_STATUS_ARCHIVE)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, user_id)
		return nil, err
	}
	order.NumOfAlert, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" where id in (select order_id from messages where receiver_id = (select unit_id from users where id = ?)"+
		" and id not in (select message_id from user_messages where user_id in (select id from users where unit_id = receiver_id)))"+
		" and ((supplier_id != 0 and supplier_id = (select unit_id from users where id = ?))"+
		" or (supplier_id = 0 and (select unit_id from users where id = ?) in (select supplier_id from supplier_requests where order_id = id)))",
		user_id, user_id, user_id)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, user_id)
		return nil, err
	}

	return order, nil
}

func (orderservice *OrderService) GetMetaByProject(project_id int64) (order *models.ApiMetaOrderByProject, err error) {
	order = new(models.ApiMetaOrderByProject)
	order.Total, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" where project_id = ?", project_id)
	if err != nil {
		log.Error("Error during getting meta order object from database %v with value %v", err, project_id)
		return nil, err
	}
	order.NumOfAlert, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o where id in (select order_id from messages where receiver_id = unit_id"+
		" and id not in (select message_id from user_messages where user_id in"+
		" (select id from users where unit_id = o.unit_id))) and project_id = ?", project_id)
	if err != nil {
		log.Error("Error during getting meta project object from database %v with value %v", err, project_id)
		return nil, err
	}

	return order, nil
}

func (orderservice *OrderService) GetFullMeta() (order *models.ApiFullMetaOrder, err error) {
	order = new(models.ApiFullMetaOrder)
	order.Total, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" where id not in (select order_id from order_statuses where status_id = ? and value = 1) and"+
		" id not in (select order_id from order_statuses where status_id = ? and value = 1)",
		models.ORDER_STATUS_ARCHIVE, models.ORDER_STATUS_DEL)
	if err != nil {
		log.Error("Error during getting meta order object from database %v", err)
		return nil, err
	}
	order.NumOfCompleted, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where s.status_id = ? and s.value = 1", models.ORDER_STATUS_COMPLETED)
	order.NumOfNew, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where s.status_id = ? and s.value = 1", models.ORDER_STATUS_NEW)
	order.NumOfOpen, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where s.status_id = ? and s.value = 1", models.ORDER_STATUS_OPEN)
	order.NumOfClosed, err = orderservice.DbContext.SelectInt("select count(distinct o.id) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where s.status_id in (?, ?, ?) and s.value = 1", models.ORDER_STATUS_CANCEL,
		models.ORDER_STATUS_SUPPLIER_CLOSE, models.ORDER_STATUS_MODERATOR_CLOSE)
	order.NumOfNotPaid, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" where id not in (select order_id from order_statuses where status_id = ? and value = 1)",
		models.ORDER_STATUS_PAID)
	order.NumOfOnTheGo, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id where s.status_id = ? and s.value = 1"+
		" and o.id not in (select order_id from order_statuses where status_id = ? and value = 1)"+
		" and o.id not in (select order_id from order_statuses where status_id = ? and value = 1)"+
		" and o.id not in (select order_id from order_statuses where status_id = ? and value = 1)",
		models.ORDER_STATUS_MODERATOR_BEGIN, models.ORDER_STATUS_MODERATOR_DOCUMENTS_GOTTEN, models.ORDER_STATUS_ARCHIVE, models.ORDER_STATUS_DEL)
	order.NumOfNoDocuments, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" where id not in (select order_id from order_statuses where status_id = ? and value = 1)"+
		" and id not in (select order_id from order_statuses where status_id = ? and value = 1)"+
		" and id not in (select order_id from order_statuses where status_id = ? and value = 1)",
		models.ORDER_STATUS_MODERATOR_DOCUMENTS_GOTTEN, models.ORDER_STATUS_ARCHIVE, models.ORDER_STATUS_DEL)
	order.NumOfArchived, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where s.status_id = ? and s.value = 1", models.ORDER_STATUS_ARCHIVE)
	order.NumOfAlert, err = orderservice.DbContext.SelectInt("select count(*) from " + orderservice.Table +
		" where id in (select order_id from messages where id not in (select message_id from user_messages where user_id in" +
		" (select id from users where unit_id = receiver_id)))")
	order.NumOfDeleted, err = orderservice.DbContext.SelectInt("select count(*) from "+orderservice.Table+
		" o inner join order_statuses s on o.id = s.order_id "+
		" where s.status_id = ? and s.value = 1", models.ORDER_STATUS_DEL)

	return order, nil
}

func (orderservice *OrderService) Create(order *models.DtoOrder, orderstatuses *[]models.DtoOrderStatus,
	trans *gorp.Transaction, inTrans bool) (err error) {
	if inTrans {
		trans, err = orderservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating order object in database %v", err)
			return err
		}
	}

	if trans != nil {
		err = trans.Insert(order)
	} else {
		err = orderservice.DbContext.Insert(order)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating order object in database %v", err)
		return err
	}

	if order.Name == "" {
		order.Name = fmt.Sprintf(ORDER_NAME_TEMPLATE, order.ID)
		if trans != nil {
			_, err = trans.Update(order)
		} else {
			_, err = orderservice.DbContext.Update(order)
		}
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during creating order object in database %v", err)
			return err
		}
	}

	for _, orderstatus := range *orderstatuses {
		orderstatus.Order_ID = order.ID

		err = orderservice.OrderStatusRepository.Save(&orderstatus, trans)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			return err
		}
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during creating order object in database %v", err)
			return err
		}
	}

	return nil
}

func (orderservice *OrderService) Update(order *models.DtoOrder, orderstatuses *[]models.DtoOrderStatus,
	trans *gorp.Transaction, inTrans bool) (err error) {
	if inTrans {
		trans, err = orderservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating order object in database %v", err)
			return err
		}
	}

	if order.Name == "" {
		order.Name = fmt.Sprintf(ORDER_NAME_TEMPLATE, order.ID)
	}
	if trans != nil {
		_, err = trans.Update(order)
	} else {
		_, err = orderservice.DbContext.Update(order)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating order object in database %v with value %v", err, order.ID)
		return err
	}

	for _, orderstatus := range *orderstatuses {
		err = orderservice.OrderStatusRepository.Save(&orderstatus, trans)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			return err
		}
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
