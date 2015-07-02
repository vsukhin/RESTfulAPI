package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type ReportComplexStatusRepository interface {
	Get(report_id int64, complexstatus_id int) (reportcomplexstatus *models.DtoReportComplexStatus, err error)
	GetByReport(report_id int64) (reportcomplexstatuses *[]models.ViewApiReportComplexStatus, err error)
	Create(dtoreportcomplexstatus *models.DtoReportComplexStatus, trans *gorp.Transaction) (err error)
	DeleteByReport(report_id int64, trans *gorp.Transaction) (err error)
}

type ReportComplexStatusService struct {
	*Repository
}

func NewReportComplexStatusService(repository *Repository) *ReportComplexStatusService {
	repository.DbContext.AddTableWithName(models.DtoReportComplexStatus{}, repository.Table).SetKeys(false, "report_id", "complex_status_id")
	return &ReportComplexStatusService{Repository: repository}
}

func (reportcomplexstatusservice *ReportComplexStatusService) Get(report_id int64,
	complexstatus_id int) (reportcomplexstatus *models.DtoReportComplexStatus, err error) {
	reportcomplexstatus = new(models.DtoReportComplexStatus)
	err = reportcomplexstatusservice.DbContext.SelectOne(reportcomplexstatus, "select * from "+reportcomplexstatusservice.Table+
		" where report_id = ? and complex_status_id = ?", report_id, complexstatus_id)
	if err != nil {
		log.Error("Error during getting report complex status object from database %v with value %v, %v", err, report_id, complexstatus_id)
		return nil, err
	}

	return reportcomplexstatus, nil
}

func (reportcomplexstatusservice *ReportComplexStatusService) GetByReport(
	report_id int64) (reportcomplexstatuses *[]models.ViewApiReportComplexStatus, err error) {
	reportcomplexstatuses = new([]models.ViewApiReportComplexStatus)
	_, err = reportcomplexstatusservice.DbContext.Select(reportcomplexstatuses,
		"select complex_status_id from "+reportcomplexstatusservice.Table+" where report_id = ?", report_id)
	if err != nil {
		log.Error("Error during getting all report complex status object from database %v with value %v", err, report_id)
		return nil, err
	}

	return reportcomplexstatuses, nil
}

func (reportcomplexstatusservice *ReportComplexStatusService) Create(dtoreportcomplexstatus *models.DtoReportComplexStatus,
	trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoreportcomplexstatus)
	} else {
		err = reportcomplexstatusservice.DbContext.Insert(dtoreportcomplexstatus)
	}
	if err != nil {
		log.Error("Error during creating report complex status object in database %v", err)
		return err
	}

	return nil
}

func (reportcomplexstatusservice *ReportComplexStatusService) DeleteByReport(report_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+reportcomplexstatusservice.Table+" where report_id = ?", report_id)
	} else {
		_, err = reportcomplexstatusservice.DbContext.Exec("delete from "+reportcomplexstatusservice.Table+" where report_id = ?", report_id)
	}
	if err != nil {
		log.Error("Error during deleting report complex status objects for report object in database %v with value %v", err, report_id)
		return err
	}

	return nil
}
