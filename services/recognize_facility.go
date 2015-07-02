package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type RecognizeFacilityRepository interface {
	Exists(orderid int64) (found bool, err error)
	Get(order_id int64) (recognizefacility *models.DtoRecognizeFacility, err error)
	SetArrays(recognizefacility *models.DtoRecognizeFacility, trans *gorp.Transaction) (err error)
	Create(recognizefacility *models.DtoRecognizeFacility, inTrans bool) (err error)
	Update(recognizefacility *models.DtoRecognizeFacility, inTrans bool) (err error)
	Save(recognizefacility *models.DtoRecognizeFacility, inTrans bool) (err error)
}

type RecognizeFacilityService struct {
	InputFieldRepository      InputFieldRepository
	InputFileRepository       InputFileRepository
	SupplierRequestRepository SupplierRequestRepository
	InputFtpRepository        InputFtpRepository
	ResultTableRepository     ResultTableRepository
	WorkTableRepository       WorkTableRepository
	InputProductRepository    InputProductRepository
	*Repository
}

func NewRecognizeFacilityService(repository *Repository) *RecognizeFacilityService {
	repository.DbContext.AddTableWithName(models.DtoRecognizeFacility{}, repository.Table).SetKeys(false, "order_id")
	return &RecognizeFacilityService{Repository: repository}
}

func (recognizefacilityservice *RecognizeFacilityService) Exists(orderid int64) (found bool, err error) {
	var count int64
	count, err = recognizefacilityservice.DbContext.SelectInt("select count(*) from "+recognizefacilityservice.Table+
		" where order_id = ?", orderid)
	if err != nil {
		log.Error("Error during checking recognize facility object in database %v with value %v", err, orderid)
		return false, err
	}

	return count != 0, nil
}

func (recognizefacilityservice *RecognizeFacilityService) Get(orderid int64) (recognizefacility *models.DtoRecognizeFacility, err error) {
	recognizefacility = new(models.DtoRecognizeFacility)
	err = recognizefacilityservice.DbContext.SelectOne(recognizefacility, "select * from "+recognizefacilityservice.Table+" where order_id = ?", orderid)
	if err != nil {
		log.Error("Error during getting recognize facility object from database %v with value %v", err, orderid)
		return nil, err
	}

	return recognizefacility, nil
}

func (recognizefacilityservice *RecognizeFacilityService) SetArrays(recognizefacility *models.DtoRecognizeFacility, trans *gorp.Transaction) (err error) {
	err = recognizefacilityservice.InputFieldRepository.DeleteByOrder(recognizefacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
		return err
	}
	for _, dtoinputfield := range recognizefacility.EstimatedFields {
		err = recognizefacilityservice.InputFieldRepository.Create(&dtoinputfield, trans)
		if err != nil {
			log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
			return err
		}
	}
	err = recognizefacilityservice.InputProductRepository.DeleteByOrder(recognizefacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
		return err
	}
	for _, dtoinputproduct := range recognizefacility.InputProducts {
		err = recognizefacilityservice.InputProductRepository.Create(&dtoinputproduct, trans)
		if err != nil {
			log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
			return err
		}
	}
	err = recognizefacilityservice.InputFileRepository.DeleteByOrder(recognizefacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
		return err
	}
	for _, dtoinputfile := range recognizefacility.EstimatedFromFiles {
		err = recognizefacilityservice.InputFileRepository.Create(&dtoinputfile, trans)
		if err != nil {
			log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
			return err
		}
	}
	err = recognizefacilityservice.SupplierRequestRepository.DeleteByOrder(recognizefacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
		return err
	}
	for _, dtosupplierrequest := range recognizefacility.SupplierRequests {
		err = recognizefacilityservice.SupplierRequestRepository.Create(&dtosupplierrequest, trans)
		if err != nil {
			log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
			return err
		}
	}
	err = recognizefacilityservice.InputFtpRepository.Delete(recognizefacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
		return err
	}
	err = recognizefacilityservice.InputFtpRepository.Create(&recognizefacility.Ftp, trans)
	if err != nil {
		log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
		return err
	}
	err = recognizefacilityservice.ResultTableRepository.DeleteByOrder(recognizefacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
		return err
	}
	for _, dtoresulttable := range recognizefacility.ResultTables {
		err = recognizefacilityservice.ResultTableRepository.Create(&dtoresulttable, trans)
		if err != nil {
			log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
			return err
		}
	}
	err = recognizefacilityservice.WorkTableRepository.DeleteByOrder(recognizefacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
		return err
	}
	for _, dtoworktable := range recognizefacility.WorkTables {
		err = recognizefacilityservice.WorkTableRepository.Create(&dtoworktable, trans)
		if err != nil {
			log.Error("Error during setting recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
			return err
		}
	}

	return nil
}

func (recognizefacilityservice *RecognizeFacilityService) Create(recognizefacility *models.DtoRecognizeFacility, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = recognizefacilityservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating recognize facility object in database %v", err)
			return err
		}
	}

	if inTrans {
		err = trans.Insert(recognizefacility)
	} else {
		err = recognizefacilityservice.DbContext.Insert(recognizefacility)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating recognize facility object in database %v", err)
		return err
	}

	err = recognizefacilityservice.SetArrays(recognizefacility, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during creating recognize facility object in database %v", err)
			return err
		}
	}

	return nil
}

func (recognizefacilityservice *RecognizeFacilityService) Update(recognizefacility *models.DtoRecognizeFacility, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = recognizefacilityservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating recognize facility object in database %v", err)
			return err
		}
	}

	if inTrans {
		_, err = trans.Update(recognizefacility)
	} else {
		_, err = recognizefacilityservice.DbContext.Update(recognizefacility)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
		return err
	}

	err = recognizefacilityservice.SetArrays(recognizefacility, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during updating recognize facility object in database %v", err)
			return err
		}
	}

	return nil
}

func (recognizefacilityservice *RecognizeFacilityService) Save(recognizefacility *models.DtoRecognizeFacility, inTrans bool) (err error) {
	count, err := recognizefacilityservice.DbContext.SelectInt("select count(*) from "+recognizefacilityservice.Table+
		" where order_id = ?", recognizefacility.Order_ID)
	if err != nil {
		log.Error("Error during saving recognize facility object in database %v with value %v", err, recognizefacility.Order_ID)
		return err
	}
	if count == 0 {
		err = recognizefacilityservice.Create(recognizefacility, inTrans)
	} else {
		err = recognizefacilityservice.Update(recognizefacility, inTrans)
	}

	return err
}
