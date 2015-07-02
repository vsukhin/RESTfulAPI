package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type ReportFacilityRepository interface {
	Get(report_id int64, facility_id int64) (reportfacility *models.DtoReportFacility, err error)
	GetByReport(report_id int64) (reportfacilities *[]models.ViewApiReportFacility, err error)
	Create(dtoreportfacility *models.DtoReportFacility, trans *gorp.Transaction) (err error)
	DeleteByReport(report_id int64, trans *gorp.Transaction) (err error)
}

type ReportFacilityService struct {
	*Repository
}

func NewReportFacilityService(repository *Repository) *ReportFacilityService {
	repository.DbContext.AddTableWithName(models.DtoReportFacility{}, repository.Table).SetKeys(false, "report_id", "service_id")
	return &ReportFacilityService{Repository: repository}
}

func (reportfacilityservice *ReportFacilityService) Get(report_id int64, facility_id int64) (reportfacility *models.DtoReportFacility, err error) {
	reportfacility = new(models.DtoReportFacility)
	err = reportfacilityservice.DbContext.SelectOne(reportfacility, "select * from "+reportfacilityservice.Table+
		" where report_id = ? and service_id = ?", report_id, facility_id)
	if err != nil {
		log.Error("Error during getting report facility object from database %v with value %v, %v", err, report_id, facility_id)
		return nil, err
	}

	return reportfacility, nil
}

func (reportfacilityservice *ReportFacilityService) GetByReport(report_id int64) (reportfacilities *[]models.ViewApiReportFacility, err error) {
	reportfacilities = new([]models.ViewApiReportFacility)
	_, err = reportfacilityservice.DbContext.Select(reportfacilities,
		"select service_id from "+reportfacilityservice.Table+" where report_id = ?", report_id)
	if err != nil {
		log.Error("Error during getting all report facility object from database %v with value %v", err, report_id)
		return nil, err
	}

	return reportfacilities, nil
}

func (reportfacilityservice *ReportFacilityService) Create(dtoreportfacility *models.DtoReportFacility, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoreportfacility)
	} else {
		err = reportfacilityservice.DbContext.Insert(dtoreportfacility)
	}
	if err != nil {
		log.Error("Error during creating report facility object in database %v", err)
		return err
	}

	return nil
}

func (reportfacilityservice *ReportFacilityService) DeleteByReport(report_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+reportfacilityservice.Table+" where report_id = ?", report_id)
	} else {
		_, err = reportfacilityservice.DbContext.Exec("delete from "+reportfacilityservice.Table+" where report_id = ?", report_id)
	}
	if err != nil {
		log.Error("Error during deleting report facility objects for report object in database %v with value %v", err, report_id)
		return err
	}

	return nil
}
