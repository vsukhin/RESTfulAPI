package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type PricePropertiesService struct {
	*Repository
}

func NewPricePropertiesService(repository *Repository) *PricePropertiesService {
	repository.DbContext.AddTableWithName(models.DtoPriceProperties{}, repository.Table).SetKeys(false, "customer_table_id")
	return &PricePropertiesService{
		repository,
	}
}

func (pricepropertiesservice *PricePropertiesService) Get(tableid int64) (priceproperties *models.DtoPriceProperties, err error) {
	priceproperties = new(models.DtoPriceProperties)
	err = pricepropertiesservice.DbContext.SelectOne(priceproperties, "select * from "+pricepropertiesservice.Table+" where customer_table_id = ?", tableid)
	if err != nil {
		log.Error("Error during getting price properties object from database %v with value %v", err, tableid)
		return nil, err
	}

	return priceproperties, nil
}

func (pricepropertiesservice *PricePropertiesService) Exists(tableid int64) (found bool, err error) {
	var count int64
	count, err = pricepropertiesservice.DbContext.SelectInt("select count(*) from "+pricepropertiesservice.Table+" where customer_table_id = ?", tableid)
	if err != nil {
		log.Error("Error during getting price properties object from database %v with value %v", err, tableid)
		return false, err
	}

	return count != 0, nil
}

func (pricepropertiesservice *PricePropertiesService) Create(priceproperties *models.DtoPriceProperties, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = pricepropertiesservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating price properties object in database %v", err)
			return err
		}
	}

	err = pricepropertiesservice.DbContext.Insert(priceproperties)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating price properties object in database %v", err)
		return err
	}

	_, err = pricepropertiesservice.DbContext.Exec(
		"update customer_tables set type_id =(select id from table_types where name = ?) where id = ?",
		models.TABLE_TYPE_PRICE, priceproperties.Customer_Table_ID)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating price properties object in database %v with value %v", err, priceproperties.Customer_Table_ID)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during creating price properties object in database %v", err)
			return err
		}
	}

	return nil
}
