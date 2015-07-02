package services

import (
	"application/models"
	"fmt"
	"github.com/coopernurse/gorp"
)

type ReportPeriodRepository interface {
	Get(id int64) (reportperiod *models.DtoReportPeriod, err error)
	GetByReport(report_id int64) (reportperiods *[]models.ViewApiReportPeriod, err error)
	Create(dtoreportperiod *models.DtoReportPeriod, trans *gorp.Transaction) (err error)
	DeleteByReport(report_id int64, trans *gorp.Transaction) (err error)
}

type ReportPeriodService struct {
	*Repository
}

func NewReportPeriodService(repository *Repository) *ReportPeriodService {
	repository.DbContext.AddTableWithName(models.DtoReportPeriod{}, repository.Table).SetKeys(true, "id")
	return &ReportPeriodService{Repository: repository}
}

func (reportperiodservice *ReportPeriodService) Get(id int64) (reportperiod *models.DtoReportPeriod, err error) {
	reportperiod = new(models.DtoReportPeriod)
	err = reportperiodservice.DbContext.SelectOne(reportperiod, "select * from "+reportperiodservice.Table+
		" where id = ?", id)
	if err != nil {
		log.Error("Error during getting report period object from database %v with value %v", err, id)
		return nil, err
	}

	return reportperiod, nil
}

func (reportperiodservice *ReportPeriodService) GetByReport(report_id int64) (reportperiods *[]models.ViewApiReportPeriod, err error) {
	reportperiods = new([]models.ViewApiReportPeriod)
	periods := new([]models.DtoReportPeriod)
	_, err = reportperiodservice.DbContext.Select(periods,
		"select * from "+reportperiodservice.Table+" where report_id = ?", report_id)
	if err != nil {
		log.Error("Error during getting all report period object from database %v with value %v", err, report_id)
		return nil, err
	}
	for _, period := range *periods {
		*reportperiods = append(*reportperiods, *models.NewViewApiReportPeriod(fmt.Sprintf("%v", period.Begin), fmt.Sprintf("%v", period.End)))
	}

	return reportperiods, nil
}

func (reportperiodservice *ReportPeriodService) Create(dtoreportperiod *models.DtoReportPeriod, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoreportperiod)
	} else {
		err = reportperiodservice.DbContext.Insert(dtoreportperiod)
	}
	if err != nil {
		log.Error("Error during creating report period object in database %v", err)
		return err
	}

	return nil
}

func (reportperiodservice *ReportPeriodService) DeleteByReport(report_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+reportperiodservice.Table+" where report_id = ?", report_id)
	} else {
		_, err = reportperiodservice.DbContext.Exec("delete from "+reportperiodservice.Table+" where report_id = ?", report_id)
	}
	if err != nil {
		log.Error("Error during deleting report period objects for report object in database %v with value %v", err, report_id)
		return err
	}

	return nil
}
