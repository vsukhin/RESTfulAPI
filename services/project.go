package services

import (
	"application/models"
	"fmt"
)

type ProjectRepository interface {
	CheckCustomerAccess(user_id int64, id int64) (allowed bool, err error)
	HasNotCompletedOrder(id int64) (has bool, err error)
	HasNotPaidOrder(id int64) (has bool, err error)
	Get(id int64) (project *models.DtoProject, err error)
	GetMeta(user_id int64) (project *models.ApiMetaProject, err error)
	GetByUser(userid int64, filter string) (projects *[]models.ApiMiddleProject, err error)
	GetByUserWithStatus(userid int64, active bool, filter string) (projects *[]models.ApiShortProject, err error)
	GetByUnit(unitid int64) (projects *[]models.ApiShortProject, err error)
	Create(project *models.DtoProject) (err error)
	Update(project *models.DtoProject) (err error)
	Deactivate(project *models.DtoProject) (err error)
}

type ProjectService struct {
	*Repository
}

func NewProjectService(repository *Repository) *ProjectService {
	repository.DbContext.AddTableWithName(models.DtoProject{}, repository.Table).SetKeys(true, "id")
	return &ProjectService{Repository: repository}
}

func (projectservice *ProjectService) CheckCustomerAccess(user_id int64, id int64) (allowed bool, err error) {
	count, err := projectservice.DbContext.SelectInt("select count(*) from "+projectservice.Table+
		" where id = ? and unit_id = (select unit_id from users where id = ?)", id, user_id)
	if err != nil {
		log.Error("Error during checking project object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (projectservice *ProjectService) HasNotCompletedOrder(id int64) (has bool, err error) {
	count, err := projectservice.DbContext.SelectInt("select count(*) from "+projectservice.Table+
		" where id = ? and id in (select project_id from orders where id not in"+
		" (select s.order_id from order_statuses s inner join order_statuses d on s.order_id = d.order_id"+
		" where s.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_SUPPLIER_CLOSE)+
		" and s.value = 1 and d.status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_MODERATOR_DOCUMENTS_GOTTEN)+" and d.value = 1"+
		" union (select order_id from order_statuses where status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_CANCEL)+" and value = 1)))", id)
	if err != nil {
		log.Error("Error during checking project object from database %v with value %v", err, id)
		return false, err
	}

	return count != 0, nil
}

func (projectservice *ProjectService) HasNotPaidOrder(id int64) (has bool, err error) {
	count, err := projectservice.DbContext.SelectInt("select count(*) from "+projectservice.Table+
		" where id = ? and id in (select project_id from orders where id not in"+
		" (select order_id from order_statuses where status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_PAID)+" and value = 1"+
		" union (select order_id from order_statuses where status_id = "+fmt.Sprintf("%v", models.ORDER_STATUS_CANCEL)+" and value = 1)))", id)
	if err != nil {
		log.Error("Error during checking project object from database %v with value %v", err, id)
		return false, err
	}

	return count != 0, nil
}

func (projectservice *ProjectService) Get(id int64) (project *models.DtoProject, err error) {
	project = new(models.DtoProject)
	err = projectservice.DbContext.SelectOne(project, "select * from "+projectservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting project object from database %v with value %v", err, id)
		return nil, err
	}

	return project, nil
}

func (projectservice *ProjectService) GetMeta(user_id int64) (project *models.ApiMetaProject, err error) {
	project = new(models.ApiMetaProject)
	project.Total, err = projectservice.DbContext.SelectInt("select count(*) from "+projectservice.Table+
		" where unit_id = (select unit_id from users where id = ?)", user_id)
	if err != nil {
		log.Error("Error during getting meta project object from database %v with value %v", err, user_id)
		return nil, err
	}
	project.NumOfArchive, err = projectservice.DbContext.SelectInt("select count(*) from "+projectservice.Table+
		" where active = 0 and unit_id = (select unit_id from users where id = ?)", user_id)
	if err != nil {
		log.Error("Error during getting meta project object from database %v with value %v", err, user_id)
		return nil, err
	}

	return project, nil
}

func (projectservice *ProjectService) GetByUser(userid int64, filter string) (projects *[]models.ApiMiddleProject, err error) {
	projects = new([]models.ApiMiddleProject)
	_, err = projectservice.DbContext.Select(projects, "select id, name, not active as archive from "+projectservice.Table+
		" where unit_id = (select unit_id from users where id = ?)"+filter, userid)
	if err != nil {
		log.Error("Error during getting unit project object from database %v with value %v", err, userid)
		return nil, err
	}

	return projects, nil
}

func (projectservice *ProjectService) GetByUserWithStatus(userid int64, active bool, filter string) (projects *[]models.ApiShortProject, err error) {
	projects = new([]models.ApiShortProject)
	constraint := " and "
	if active {
		constraint += " active = 1"
	} else {
		constraint += " active = 0"
	}
	_, err = projectservice.DbContext.Select(projects, "select id, name from "+projectservice.Table+
		" where unit_id = (select unit_id from users where id = ?)"+constraint+filter, userid)
	if err != nil {
		log.Error("Error during getting unit project object from database %v with value %v", err, userid)
		return nil, err
	}

	return projects, nil
}

func (projectservice *ProjectService) GetByUnit(unitid int64) (projects *[]models.ApiShortProject, err error) {
	projects = new([]models.ApiShortProject)
	_, err = projectservice.DbContext.Select(projects, "select id, name from "+projectservice.Table+" where unit_id = ?", unitid)
	if err != nil {
		log.Error("Error during getting unit project object from database %v with value %v", err, unitid)
		return nil, err
	}

	return projects, nil
}

func (projectservice *ProjectService) Create(project *models.DtoProject) (err error) {
	err = projectservice.DbContext.Insert(project)
	if err != nil {
		log.Error("Error during creating project object in database %v", err)
		return err
	}

	return nil
}

func (projectservice *ProjectService) Update(project *models.DtoProject) (err error) {
	_, err = projectservice.DbContext.Update(project)
	if err != nil {
		log.Error("Error during updating project object in database %v with value %v", err, project.ID)
		return err
	}

	return nil
}

func (projectservice *ProjectService) Deactivate(project *models.DtoProject) (err error) {
	_, err = projectservice.DbContext.Exec("update "+projectservice.Table+" set active = 0 where id = ?", project.ID)
	if err != nil {
		log.Error("Error during deactivating project object in database %v with value %v", err, project.ID)
		return err
	}

	return nil
}
