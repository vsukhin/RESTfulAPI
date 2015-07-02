package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type ReportProjectRepository interface {
	Get(report_id int64, project_id int64) (reportproject *models.DtoReportProject, err error)
	GetByReport(report_id int64) (reportprojects *[]models.ViewApiReportProject, err error)
	Create(dtoreportproject *models.DtoReportProject, trans *gorp.Transaction) (err error)
	DeleteByReport(report_id int64, trans *gorp.Transaction) (err error)
}

type ReportProjectService struct {
	*Repository
}

func NewReportProjectService(repository *Repository) *ReportProjectService {
	repository.DbContext.AddTableWithName(models.DtoReportProject{}, repository.Table).SetKeys(false, "report_id", "project_id")
	return &ReportProjectService{Repository: repository}
}

func (reportprojectservice *ReportProjectService) Get(report_id int64, project_id int64) (reportproject *models.DtoReportProject, err error) {
	reportproject = new(models.DtoReportProject)
	err = reportprojectservice.DbContext.SelectOne(reportproject, "select * from "+reportprojectservice.Table+
		" where report_id = ? and project_id = ?", report_id, project_id)
	if err != nil {
		log.Error("Error during getting report project object from database %v with value %v, %v", err, report_id, project_id)
		return nil, err
	}

	return reportproject, nil
}

func (reportprojectservice *ReportProjectService) GetByReport(report_id int64) (reportprojects *[]models.ViewApiReportProject, err error) {
	reportprojects = new([]models.ViewApiReportProject)
	_, err = reportprojectservice.DbContext.Select(reportprojects,
		"select project_id from "+reportprojectservice.Table+" where report_id = ?", report_id)
	if err != nil {
		log.Error("Error during getting all report project object from database %v with value %v", err, report_id)
		return nil, err
	}

	return reportprojects, nil
}

func (reportprojectservice *ReportProjectService) Create(dtoreportproject *models.DtoReportProject, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoreportproject)
	} else {
		err = reportprojectservice.DbContext.Insert(dtoreportproject)
	}
	if err != nil {
		log.Error("Error during creating report project object in database %v", err)
		return err
	}

	return nil
}

func (reportprojectservice *ReportProjectService) DeleteByReport(report_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+reportprojectservice.Table+" where report_id = ?", report_id)
	} else {
		_, err = reportprojectservice.DbContext.Exec("delete from "+reportprojectservice.Table+" where report_id = ?", report_id)
	}
	if err != nil {
		log.Error("Error during deleting report project objects for report object in database %v with value %v", err, report_id)
		return err
	}

	return nil
}
