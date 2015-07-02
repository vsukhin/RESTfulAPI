package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type ReportRepository interface {
	CheckCustomerAccess(user_id int64, id int64) (allowed bool, err error)
	Get(id int64) (report *models.DtoReport, err error)
	GetMeta(user_id int64) (report *models.ApiMetaReport, err error)
	SetArrays(report *models.DtoReport, trans *gorp.Transaction) (err error)
	Create(report *models.DtoReport, inTrans bool) (err error)
}

type ReportService struct {
	UserRepository                UserRepository
	ReportPeriodRepository        ReportPeriodRepository
	ReportProjectRepository       ReportProjectRepository
	ReportOrderRepository         ReportOrderRepository
	ReportFacilityRepository      ReportFacilityRepository
	ReportComplexStatusRepository ReportComplexStatusRepository
	ReportSupplierRepository      ReportSupplierRepository
	ReportSettingsRepository      ReportSettingsRepository
	*Repository
}

func NewReportService(repository *Repository) *ReportService {
	repository.DbContext.AddTableWithName(models.DtoReport{}, repository.Table).SetKeys(true, "id")
	return &ReportService{Repository: repository}
}

func (reportservice *ReportService) CheckCustomerAccess(user_id int64, id int64) (allowed bool, err error) {
	count, err := reportservice.DbContext.SelectInt("select count(*) from "+reportservice.Table+
		" where id = ? and unit_id = (select unit_id from users where id = ?)", id, user_id)
	if err != nil {
		log.Error("Error during checking report object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (reportservice *ReportService) Get(id int64) (report *models.DtoReport, err error) {
	report = new(models.DtoReport)
	err = reportservice.DbContext.SelectOne(report, "select * from "+reportservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting report object from database %v with value %v", err, id)
		return nil, err
	}

	return report, nil
}

func (reportservice *ReportService) GetMeta(user_id int64) (report *models.ApiMetaReport, err error) {
	report = new(models.ApiMetaReport)
	report.Access, err = reportservice.UserRepository.CheckReportAccess(user_id)
	if err != nil {
		log.Error("Error during getting report meta object from database %v", err)
		return nil, err
	}
	reports := new([]models.DtoReport)
	_, err = reportservice.DbContext.Select(reports, "select * from "+reportservice.Table+" where user_id = ? and active = 1 order by created desc limit 1",
		user_id)
	if err != nil {
		log.Error("Error during getting report object from database %v with value %v", err, user_id)
		return nil, err
	}
	if len(*reports) != 0 {
		report.ID = (*reports)[0].ID
	}

	return report, nil
}

func (reportservice *ReportService) SetArrays(report *models.DtoReport, trans *gorp.Transaction) (err error) {
	err = reportservice.ReportPeriodRepository.DeleteByReport(report.ID, trans)
	if err != nil {
		log.Error("Error during setting report object in database %v with value %v", err, report.ID)
		return err
	}
	for _, dtoreportperiod := range report.Periods {
		dtoreportperiod.Report_ID = report.ID
		err = reportservice.ReportPeriodRepository.Create(&dtoreportperiod, trans)
		if err != nil {
			log.Error("Error during setting report object in database %v with value %v", err, report.ID)
			return err
		}
	}
	err = reportservice.ReportProjectRepository.DeleteByReport(report.ID, trans)
	if err != nil {
		log.Error("Error during setting report object in database %v with value %v", err, report.ID)
		return err
	}
	for _, dtoreportproject := range report.Projects {
		dtoreportproject.Report_ID = report.ID
		err = reportservice.ReportProjectRepository.Create(&dtoreportproject, trans)
		if err != nil {
			log.Error("Error during setting report object in database %v with value %v", err, report.ID)
			return err
		}
	}
	err = reportservice.ReportOrderRepository.DeleteByReport(report.ID, trans)
	if err != nil {
		log.Error("Error during setting report object in database %v with value %v", err, report.ID)
		return err
	}
	for _, dtoreportorder := range report.Orders {
		dtoreportorder.Report_ID = report.ID
		err = reportservice.ReportOrderRepository.Create(&dtoreportorder, trans)
		if err != nil {
			log.Error("Error during setting report object in database %v with value %v", err, report.ID)
			return err
		}
	}
	err = reportservice.ReportFacilityRepository.DeleteByReport(report.ID, trans)
	if err != nil {
		log.Error("Error during setting report object in database %v with value %v", err, report.ID)
		return err
	}
	for _, dtoreportfacility := range report.Facilities {
		dtoreportfacility.Report_ID = report.ID
		err = reportservice.ReportFacilityRepository.Create(&dtoreportfacility, trans)
		if err != nil {
			log.Error("Error during setting report object in database %v with value %v", err, report.ID)
			return err
		}
	}
	err = reportservice.ReportComplexStatusRepository.DeleteByReport(report.ID, trans)
	if err != nil {
		log.Error("Error during setting report object in database %v with value %v", err, report.ID)
		return err
	}
	for _, dtoreportcomplexstatus := range report.ComplexStatuses {
		dtoreportcomplexstatus.Report_ID = report.ID
		err = reportservice.ReportComplexStatusRepository.Create(&dtoreportcomplexstatus, trans)
		if err != nil {
			log.Error("Error during setting report object in database %v with value %v", err, report.ID)
			return err
		}
	}
	err = reportservice.ReportSupplierRepository.DeleteByReport(report.ID, trans)
	if err != nil {
		log.Error("Error during setting report object in database %v with value %v", err, report.ID)
		return err
	}
	for _, dtoreportsupplier := range report.Suppliers {
		dtoreportsupplier.Report_ID = report.ID
		err = reportservice.ReportSupplierRepository.Create(&dtoreportsupplier, trans)
		if err != nil {
			log.Error("Error during setting report object in database %v with value %v", err, report.ID)
			return err
		}
	}
	err = reportservice.ReportSettingsRepository.DeleteByReport(report.ID, trans)
	if err != nil {
		log.Error("Error during setting report object in database %v with value %v", err, report.ID)
		return err
	}
	report.Settings.Report_ID = report.ID
	err = reportservice.ReportSettingsRepository.Create(&report.Settings, trans)
	if err != nil {
		log.Error("Error during setting report object in database %v with value %v", err, report.ID)
		return err
	}

	return nil
}

func (reportservice *ReportService) Create(report *models.DtoReport, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = reportservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating report object in database %v", err)
			return err
		}
	}

	if inTrans {
		err = trans.Insert(report)
	} else {
		err = reportservice.DbContext.Insert(report)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating report object in database %v", err)
		return err
	}

	err = reportservice.SetArrays(report, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during creating report object in database %v", err)
			return err
		}
	}

	return nil
}
