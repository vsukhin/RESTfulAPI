package services

import (
	"application/models"
	"fmt"
	"github.com/coopernurse/gorp"
)

type SupplierFacilityRepository interface {
	Get(supplier_id int64, service_id int64) (supplierfacility *models.DtoSupplierFacility, err error)
	GetByAlias(alias string) (supplierfacilities *[]models.ApiShortSupplierFacility, err error)
	GetByUnit(unit_id int64, filter string) (supplierfacilities *[]models.ApiLongSupplierFacility, err error)
	SetArrayByUser(user_id int64, facilities *[]int64, inTrans bool) (err error)
}

type SupplierFacilityService struct {
	*Repository
}

func NewSupplierFacilityService(repository *Repository) *SupplierFacilityService {
	repository.DbContext.AddTableWithName(models.DtoSupplierFacility{}, repository.Table).SetKeys(false, "supplier_id", "service_id")
	return &SupplierFacilityService{Repository: repository}
}

func (supplierfacilityservice *SupplierFacilityService) Get(supplier_id int64, service_id int64) (supplierfacility *models.DtoSupplierFacility, err error) {
	supplierfacility = new(models.DtoSupplierFacility)
	err = supplierfacilityservice.DbContext.SelectOne(supplierfacility, "select * from "+supplierfacilityservice.Table+
		" where supplier_id = ? and service_id = ?", supplier_id, service_id)
	if err != nil {
		log.Error("Error during getting supplier facility object from database %v with value %v, %v", err, supplier_id, service_id)
		return nil, err
	}

	return supplierfacility, nil
}

func (supplierfacilityservice *SupplierFacilityService) GetByAlias(alias string) (
	supplierfacilities *[]models.ApiShortSupplierFacility, err error) {
	supplierfacilities = new([]models.ApiShortSupplierFacility)
	_, err = supplierfacilityservice.DbContext.Select(supplierfacilities,
		"select f.supplier_id as id, u.name, f.position, f.rating, f.throughput from "+
			supplierfacilityservice.Table+" f inner join services s on f.service_id = s.id inner join units u on f.supplier_id = u.id"+
			" where s.active = 1 and u.active = 1 and s.alias = ? and f.service_id in (select service_id from price_properties p"+
			" where published = 1 and customer_table_id in (select id from customer_tables"+
			" where active = 1 and permanent = 1 and type_id = ? and unit_id = f.supplier_id)"+
			" and ((date(end) != '0001-01-01' and now() <= end) or (date(end) = '0001-01-01'))"+
			" and ((date(begin) != '0001-01-01' and now() >= begin) or (date(begin) = '0001-01-01' and after_id = 0)"+
			" or (date(begin) = '0001-01-01' and after_id != 0 and"+
			" (select date(end) from price_properties where customer_table_id = p.after_id) != '0001-01-01' and"+
			" now() > (select end from price_properties where customer_table_id = p.after_id)))) order by f.position asc", alias, models.TABLE_TYPE_PRICE)
	if err != nil {
		log.Error("Error during getting all supplier facility object from database %v", err)
		return nil, err
	}

	return supplierfacilities, nil
}

func (supplierfacilityservice *SupplierFacilityService) GetByUnit(unit_id int64, filter string) (
	supplierfacilities *[]models.ApiLongSupplierFacility, err error) {
	supplierfacilities = new([]models.ApiLongSupplierFacility)
	_, err = supplierfacilityservice.DbContext.Select(supplierfacilities,
		"select f.supplier_id as id, u.name, f.service_id as serviceId, f.position, f.rating, f.throughput from "+
			supplierfacilityservice.Table+" f inner join services s on f.service_id = s.id inner join units u on f.supplier_id = u.id"+
			" where s.active = 1 and u.active = 1 and f.supplier_id in (select supplier_id from orders where unit_id = ?)"+filter, unit_id)
	if err != nil {
		log.Error("Error during getting all supplier facility object from database %v", err)
		return nil, err
	}

	return supplierfacilities, nil
}

func (supplierfacilityservice *SupplierFacilityService) SetArrayByUser(user_id int64, facilities *[]int64, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = supplierfacilityservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during setting supplier facility objects for user object in database %v", err)
			return err
		}
	}

	ids := ""
	if len(*facilities) != 0 {
		ids += " and service_id not in ("
		for i, value := range *facilities {
			if i != 0 {
				ids += ","
			}
			ids += fmt.Sprintf("%v", value)
		}
		ids += ")"
	}

	if inTrans {
		_, err = trans.Exec("delete from "+supplierfacilityservice.Table+
			" where supplier_id = (select unit_id from users where id = ?)"+ids, user_id)
	} else {
		_, err = supplierfacilityservice.DbContext.Exec("delete from "+supplierfacilityservice.Table+" where supplier_id = "+
			"(select unit_id from users where id = ?)"+ids, user_id)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during setting supplier facility objects for user object in database %v with value %v", err, user_id)
		return err
	}

	if len(*facilities) > 0 {
		statement := ""
		for _, value := range *facilities {
			if statement != "" {
				statement += " union"
			}
			statement += fmt.Sprintf(" select s.* from (select unit_id, %v as s, 0 as p, 0 as r, 0 as t from users where id = %v) as s where not exists "+
				"(select * from "+supplierfacilityservice.Table+" where service_id = %v and supplier_id = (select unit_id from users where id = %v))",
				value, user_id, value, user_id)
		}
		if inTrans {
			_, err = trans.Exec("insert into " + supplierfacilityservice.Table +
				" (supplier_id, service_id, position, rating, throughput)" + statement)
		} else {
			_, err = supplierfacilityservice.DbContext.Exec("insert into " + supplierfacilityservice.Table +
				" (supplier_id, service_id, position, rating, throughput)" + statement)
		}
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during setting supplier facility objects for user object in database %v with value %v", err, user_id)
			return err
		}
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during setting supplier facility objects for user object in database %v", err)
			return err
		}
	}

	return nil
}
