package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type ReportSupplierRepository interface {
	Get(report_id int64, supplier_id int64) (reportsupplier *models.DtoReportSupplier, err error)
	GetByReport(report_id int64) (reportsuppliers *[]models.ViewApiReportSupplier, err error)
	Create(dtoreportsupplier *models.DtoReportSupplier, trans *gorp.Transaction) (err error)
	DeleteByReport(report_id int64, trans *gorp.Transaction) (err error)
}

type ReportSupplierService struct {
	*Repository
}

func NewReportSupplierService(repository *Repository) *ReportSupplierService {
	repository.DbContext.AddTableWithName(models.DtoReportSupplier{}, repository.Table).SetKeys(false, "report_id", "supplier_id")
	return &ReportSupplierService{Repository: repository}
}

func (reportsupplierservice *ReportSupplierService) Get(report_id int64, supplier_id int64) (reportsupplier *models.DtoReportSupplier, err error) {
	reportsupplier = new(models.DtoReportSupplier)
	err = reportsupplierservice.DbContext.SelectOne(reportsupplier, "select * from "+reportsupplierservice.Table+
		" where report_id = ? and supplier_id = ?", report_id, supplier_id)
	if err != nil {
		log.Error("Error during getting report supplier object from database %v with value %v, %v", err, report_id, supplier_id)
		return nil, err
	}

	return reportsupplier, nil
}

func (reportsupplierservice *ReportSupplierService) GetByReport(report_id int64) (reportsuppliers *[]models.ViewApiReportSupplier, err error) {
	reportsuppliers = new([]models.ViewApiReportSupplier)
	_, err = reportsupplierservice.DbContext.Select(reportsuppliers,
		"select supplier_id from "+reportsupplierservice.Table+" where report_id = ?", report_id)
	if err != nil {
		log.Error("Error during getting all report supplier object from database %v with value %v", err, report_id)
		return nil, err
	}

	return reportsuppliers, nil
}

func (reportsupplierservice *ReportSupplierService) Create(dtoreportsupplier *models.DtoReportSupplier, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoreportsupplier)
	} else {
		err = reportsupplierservice.DbContext.Insert(dtoreportsupplier)
	}
	if err != nil {
		log.Error("Error during creating report supplier object in database %v", err)
		return err
	}

	return nil
}

func (reportsupplierservice *ReportSupplierService) DeleteByReport(report_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+reportsupplierservice.Table+" where report_id = ?", report_id)
	} else {
		_, err = reportsupplierservice.DbContext.Exec("delete from "+reportsupplierservice.Table+" where report_id = ?", report_id)
	}
	if err != nil {
		log.Error("Error during deleting report supplier objects for report object in database %v with value %v", err, report_id)
		return err
	}

	return nil
}
