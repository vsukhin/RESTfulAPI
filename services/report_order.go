package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type ReportOrderRepository interface {
	Get(report_id int64, order_id int64) (reportorder *models.DtoReportOrder, err error)
	GetByReport(report_id int64) (reportorders *[]models.ViewApiReportOrder, err error)
	Create(dtoreportorder *models.DtoReportOrder, trans *gorp.Transaction) (err error)
	DeleteByReport(report_id int64, trans *gorp.Transaction) (err error)
}

type ReportOrderService struct {
	*Repository
}

func NewReportOrderService(repository *Repository) *ReportOrderService {
	repository.DbContext.AddTableWithName(models.DtoReportOrder{}, repository.Table).SetKeys(false, "report_id", "order_id")
	return &ReportOrderService{Repository: repository}
}

func (reportorderservice *ReportOrderService) Get(report_id int64, order_id int64) (reportorder *models.DtoReportOrder, err error) {
	reportorder = new(models.DtoReportOrder)
	err = reportorderservice.DbContext.SelectOne(reportorder, "select * from "+reportorderservice.Table+
		" where report_id = ? and order_id = ?", report_id, order_id)
	if err != nil {
		log.Error("Error during getting report order object from database %v with value %v, %v", err, report_id, order_id)
		return nil, err
	}

	return reportorder, nil
}

func (reportorderservice *ReportOrderService) GetByReport(report_id int64) (reportorders *[]models.ViewApiReportOrder, err error) {
	reportorders = new([]models.ViewApiReportOrder)
	_, err = reportorderservice.DbContext.Select(reportorders,
		"select order_id from "+reportorderservice.Table+" where report_id = ?", report_id)
	if err != nil {
		log.Error("Error during getting all report order object from database %v with value %v", err, report_id)
		return nil, err
	}

	return reportorders, nil
}

func (reportorderservice *ReportOrderService) Create(dtoreportorder *models.DtoReportOrder, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoreportorder)
	} else {
		err = reportorderservice.DbContext.Insert(dtoreportorder)
	}
	if err != nil {
		log.Error("Error during creating report order object in database %v", err)
		return err
	}

	return nil
}

func (reportorderservice *ReportOrderService) DeleteByReport(report_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+reportorderservice.Table+" where report_id = ?", report_id)
	} else {
		_, err = reportorderservice.DbContext.Exec("delete from "+reportorderservice.Table+" where report_id = ?", report_id)
	}
	if err != nil {
		log.Error("Error during deleting report order objects for report object in database %v with value %v", err, report_id)
		return err
	}

	return nil
}
