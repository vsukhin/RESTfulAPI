package services

import (
	"application/models"
	"fmt"
	"github.com/coopernurse/gorp"
)

type GroupService struct {
	*Repository
}

func NewGroupService(repository *Repository) *GroupService {
	repository.DbContext.AddTableWithName(models.DtoGroup{}, repository.Table).SetKeys(false, "id")
	return &GroupService{
		repository,
	}
}

func (groupservice *GroupService) GetByUser(userid int64) (groups *[]models.UserRole, err error) {
	groups = new([]models.UserRole)
	_, err = groupservice.DbContext.Select(groups,
		"select id from "+groupservice.Table+" g inner join usergroups u on g.id = u.group_id where g.active = 1 and u.user_id = ?", userid)
	if err != nil {
		log.Error("Error during getting group objects for user object from database %v with value %v", err, userid)
		return nil, err
	}

	return groups, nil
}

func (groupservice *GroupService) GetDefault() (groups *[]models.UserRole, err error) {
	groups = new([]models.UserRole)
	_, err = groupservice.DbContext.Select(groups,
		"select id from "+groupservice.Table+" where active = 1 and `default` = 1")
	if err != nil {
		log.Error("Error during getting default group object from database %v", err)
		return nil, err
	}

	return groups, nil
}

func (groupservice *GroupService) GetByUserExt(userid int64) (groups *[]models.ApiGroup, err error) {
	groups = new([]models.ApiGroup)
	_, err = groupservice.DbContext.Select(groups,
		"select id, name from "+groupservice.Table+" g inner join usergroups u on g.id = u.group_id where g.active = 1 and u.user_id = ?", userid)
	if err != nil {
		log.Error("Error during getting group objects for user object from database %v with value %v", err, userid)
		return nil, err
	}

	return groups, nil
}

func (groupservice *GroupService) GetBySession(token string) (groups *[]models.UserRole, err error) {
	groups = new([]models.UserRole)
	_, err = groupservice.DbContext.Select(groups,
		"select id from "+groupservice.Table+" g inner join sessiongroups s on g.id = s.group_id where g.active = 1 and s.session_token = ?", token)
	if err != nil {
		log.Error("Error during getting group objects for session object from database %v with value %v", err, token)
		return nil, err
	}

	return groups, nil
}

func (groupservice *GroupService) GetBySessionExt(token string) (groups *[]models.ApiGroup, err error) {
	groups = new([]models.ApiGroup)
	_, err = groupservice.DbContext.Select(groups,
		"select id, name from "+groupservice.Table+" g inner join sessiongroups s on g.id = s.group_id where g.active = 1 and s.session_token = ?", token)
	if err != nil {
		log.Error("Error during getting group objects for session object from database %v with value %v", err, token)
		return nil, err
	}

	return groups, nil
}

func (groupservice *GroupService) GetAll() (groups *[]models.ApiGroup, err error) {
	groups = new([]models.ApiGroup)
	_, err = groupservice.DbContext.Select(groups, "select id, name from "+groupservice.Table+" where active = 1")
	if err != nil {
		log.Error("Error during getting all group objects from database %v", err)
		return nil, err
	}

	return groups, nil
}

func (groupservice *GroupService) SetByUser(userid int64, groups *[]models.UserRole, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = groupservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during setting group objects for user object in database %v", err)
			return err
		}
	}

	_, err = groupservice.DbContext.Exec("delete from usergroups where user_id = ?", userid)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during setting group objects for user object in database %v with value %v", err, userid)
		return err
	}

	if len(*groups) > 0 {
		statement := ""
		for _, value := range *groups {
			if statement != "" {
				statement += " union"
			}
			statement += fmt.Sprintf(" select %v, %v", userid, value)
		}
		_, err = groupservice.DbContext.Exec("insert into usergroups (user_id, group_id)" + statement)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during setting group objects for user object in database %v with value %v", err, userid)
			return err
		}
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during setting group objects for user object in database %v", err)
			return err
		}
	}

	return nil
}

func (groupservice *GroupService) SetBySession(token string, groups *[]models.UserRole, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = groupservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during setting group objects for session object in database %v", err)
			return err
		}
	}

	_, err = groupservice.DbContext.Exec("delete from sessiongroups where session_token = ?", token)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during setting group objects for session object in database %v with value %v", err, token)
		return err
	}

	if len(*groups) > 0 {
		statement := ""
		for _, value := range *groups {
			if statement != "" {
				statement += " union"
			}
			statement += fmt.Sprintf(" select '%v', %v", token, value) // поскольку токен сессии хранится в base64, то escape не требуется
		}
		_, err = groupservice.DbContext.Exec("insert into sessiongroups (session_token, group_id)" + statement)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during setting group objects for session object in database %v with value %v", err, token)
			return err
		}
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during setting group objects for session object in database %v", err)
			return err
		}
	}

	return nil
}

func (groupservice *GroupService) DeleteByUser(userid int64) (err error) {
	_, err = groupservice.DbContext.Exec("delete from sessiongroups where session_token in (select token from sessions where user_id = ?)", userid)
	if err != nil {
		log.Error("Error during deleting group objects for session object in database %v with value %v", err, userid)
		return err
	}

	return nil
}
