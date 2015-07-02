package services

import (
	"application/models"
)

type PriceRepository interface {
	GetSupplierPrices(alias string) (supplierprices *[]models.ApiSupplierPrice, err error)
}

type PriceService struct {
	*Repository
}

func NewPriceService(repository *Repository) *PriceService {
	return &PriceService{Repository: repository}
}

func (priceservice *PriceService) GetSupplierPrices(alias string) (supplierprices *[]models.ApiSupplierPrice, err error) {
	supplierprices = new([]models.ApiSupplierPrice)
	_, err = priceservice.DbContext.Select(supplierprices, "select supplier_id, customer_table_id from price_properties p inner join supplier_services s on "+
		" p.service_id = s.service_id where supplier_id in (select id from units where active = 1)"+
		" and p.service_id in (select id from services where active = 1 and alias = ?)"+
		" and published = 1 and customer_table_id in (select id from customer_tables where active = 1 and permanent = 1 and type_id = ? and unit_id = supplier_id)"+
		" and ((date(end) != '0001-01-01' and now() <= end) or (date(end) = '0001-01-01'))"+
		" and ((date(begin) != '0001-01-01' and now() >= begin) or (date(begin) = '0001-01-01' and after_id = 0)"+
		" or (date(begin) = '0001-01-01' and after_id != 0 and"+
		" (select date(end) from price_properties where customer_table_id = p.after_id) != '0001-01-01' and"+
		" now() > (select end from price_properties where customer_table_id = p.after_id)))", alias, models.TABLE_TYPE_PRICE)
	if err != nil {
		log.Error("Error during getting all supplier price object from database %v", err)
		return nil, err
	}

	return supplierprices, nil
}
