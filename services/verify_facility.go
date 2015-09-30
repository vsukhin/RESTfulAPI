package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type VerifyFacilityRepository interface {
	Exists(orderid int64) (found bool, err error)
	Get(order_id int64) (verifyfacility *models.DtoVerifyFacility, err error)
	SetArrays(verifyfacility *models.DtoVerifyFacility, trans *gorp.Transaction) (err error)
	Create(verifyfacility *models.DtoVerifyFacility, briefly bool, inTrans bool) (err error)
	Update(verifyfacility *models.DtoVerifyFacility, briefly bool, inTrans bool) (err error)
	Save(verifyfacility *models.DtoVerifyFacility, briefly bool, inTrans bool) (err error)
}

type VerifyFacilityService struct {
	DataColumnRepository  DataColumnRepository
	ResultTableRepository ResultTableRepository
	WorkTableRepository   WorkTableRepository
	DataProductRepository DataProductRepository
	*Repository
}

func NewVerifyFacilityService(repository *Repository) *VerifyFacilityService {
	repository.DbContext.AddTableWithName(models.DtoVerifyFacility{}, repository.Table).SetKeys(false, "order_id")
	return &VerifyFacilityService{Repository: repository}
}

func (verifyfacilityservice *VerifyFacilityService) Exists(orderid int64) (found bool, err error) {
	var count int64
	count, err = verifyfacilityservice.DbContext.SelectInt("select count(*) from "+verifyfacilityservice.Table+
		" where order_id = ?", orderid)
	if err != nil {
		log.Error("Error during checking verify facility object in database %v with value %v", err, orderid)
		return false, err
	}

	return count != 0, nil
}

func (verifyfacilityservice *VerifyFacilityService) Get(orderid int64) (verifyfacility *models.DtoVerifyFacility, err error) {
	verifyfacility = new(models.DtoVerifyFacility)
	err = verifyfacilityservice.DbContext.SelectOne(verifyfacility, "select * from "+verifyfacilityservice.Table+" where order_id = ?", orderid)
	if err != nil {
		log.Error("Error during getting verify facility object from database %v with value %v", err, orderid)
		return nil, err
	}

	return verifyfacility, nil
}

func (verifyfacilityservice *VerifyFacilityService) SetArrays(verifyfacility *models.DtoVerifyFacility, trans *gorp.Transaction) (err error) {
	err = verifyfacilityservice.DataProductRepository.DeleteByOrder(verifyfacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting verify facility object in database %v with value %v", err, verifyfacility.Order_ID)
		return err
	}
	for _, dtodataproduct := range verifyfacility.DataProducts {
		err = verifyfacilityservice.DataProductRepository.Create(&dtodataproduct, trans)
		if err != nil {
			log.Error("Error during setting verify facility object in database %v with value %v", err, verifyfacility.Order_ID)
			return err
		}
	}
	err = verifyfacilityservice.DataColumnRepository.DeleteByOrder(verifyfacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting verify facility object in database %v with value %v", err, verifyfacility.Order_ID)
		return err
	}
	for _, dtodatacolumn := range verifyfacility.DataColumns {
		err = verifyfacilityservice.DataColumnRepository.Create(&dtodatacolumn, trans)
		if err != nil {
			log.Error("Error during setting verify facility object in database %v with value %v", err, verifyfacility.Order_ID)
			return err
		}
	}
	err = verifyfacilityservice.ResultTableRepository.DeleteByOrder(verifyfacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting verify facility object in database %v with value %v", err, verifyfacility.Order_ID)
		return err
	}
	for _, dtoresulttable := range verifyfacility.ResultTables {
		err = verifyfacilityservice.ResultTableRepository.Create(&dtoresulttable, trans)
		if err != nil {
			log.Error("Error during setting verify facility object in database %v with value %v", err, verifyfacility.Order_ID)
			return err
		}
	}
	err = verifyfacilityservice.WorkTableRepository.DeleteByOrder(verifyfacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting verify facility object in database %v with value %v", err, verifyfacility.Order_ID)
		return err
	}
	for _, dtoworktable := range verifyfacility.WorkTables {
		err = verifyfacilityservice.WorkTableRepository.Create(&dtoworktable, trans)
		if err != nil {
			log.Error("Error during setting verify facility object in database %v with value %v", err, verifyfacility.Order_ID)
			return err
		}
	}

	return nil
}

func (verifyfacilityservice *VerifyFacilityService) Create(verifyfacility *models.DtoVerifyFacility, briefly bool, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = verifyfacilityservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating verify facility object in database %v", err)
			return err
		}
	}

	if inTrans {
		err = trans.Insert(verifyfacility)
	} else {
		err = verifyfacilityservice.DbContext.Insert(verifyfacility)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating verify facility object in database %v", err)
		return err
	}

	if !briefly {
		err = verifyfacilityservice.SetArrays(verifyfacility, trans)
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
			log.Error("Error during creating verify facility object in database %v", err)
			return err
		}
	}

	return nil
}

func (verifyfacilityservice *VerifyFacilityService) Update(verifyfacility *models.DtoVerifyFacility, briefly bool, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = verifyfacilityservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating verify facility object in database %v", err)
			return err
		}
	}

	if inTrans {
		_, err = trans.Update(verifyfacility)
	} else {
		_, err = verifyfacilityservice.DbContext.Update(verifyfacility)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating verify facility object in database %v with value %v", err, verifyfacility.Order_ID)
		return err
	}

	if !briefly {
		err = verifyfacilityservice.SetArrays(verifyfacility, trans)
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
			log.Error("Error during updating verify facility object in database %v", err)
			return err
		}
	}

	return nil
}

func (verifyfacilityservice *VerifyFacilityService) Save(verifyfacility *models.DtoVerifyFacility, briefly bool, inTrans bool) (err error) {
	count, err := verifyfacilityservice.DbContext.SelectInt("select count(*) from "+verifyfacilityservice.Table+
		" where order_id = ?", verifyfacility.Order_ID)
	if err != nil {
		log.Error("Error during saving verify facility object in database %v with value %v", err, verifyfacility.Order_ID)
		return err
	}
	if count == 0 {
		err = verifyfacilityservice.Create(verifyfacility, briefly, inTrans)
	} else {
		err = verifyfacilityservice.Update(verifyfacility, briefly, inTrans)
	}

	return err
}
