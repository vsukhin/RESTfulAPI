package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type ReportSettingsRepository interface {
	GetByReport(report_id int64) (reportsettings *models.DtoReportSettings, err error)
	Create(dtoreportsettings *models.DtoReportSettings, trans *gorp.Transaction) (err error)
	DeleteByReport(report_id int64, trans *gorp.Transaction) (err error)
}

type ReportSettingsService struct {
	*Repository
}

func NewReportSettingsService(repository *Repository) *ReportSettingsService {
	repository.DbContext.AddTableWithName(models.DtoReportSettings{}, repository.Table).SetKeys(false, "report_id")
	return &ReportSettingsService{Repository: repository}
}

func (reportsettingsservice *ReportSettingsService) GetByReport(report_id int64) (reportsettings *models.DtoReportSettings, err error) {
	reportsettings = new(models.DtoReportSettings)
	err = reportsettingsservice.DbContext.SelectOne(reportsettings, "select * from "+reportsettingsservice.Table+
		" where report_id = ?", report_id)
	if err != nil {
		log.Error("Error during getting report settings object from database %v with value %v", err, report_id)
		return nil, err
	}

	return reportsettings, nil
}

func (reportsettingsservice *ReportSettingsService) Create(dtoreportsettings *models.DtoReportSettings, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoreportsettings)
	} else {
		err = reportsettingsservice.DbContext.Insert(dtoreportsettings)
	}
	if err != nil {
		log.Error("Error during creating report settings object in database %v", err)
		return err
	}

	return nil
}

func (reportsettingsservice *ReportSettingsService) DeleteByReport(report_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+reportsettingsservice.Table+" where report_id = ?", report_id)
	} else {
		_, err = reportsettingsservice.DbContext.Exec("delete from "+reportsettingsservice.Table+" where report_id = ?", report_id)
	}
	if err != nil {
		log.Error("Error during deleting report settings objects for report object in database %v with value %v", err, report_id)
		return err
	}

	return nil
}
