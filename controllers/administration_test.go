package controllers

import (
	"application/models"
	"errors"
	"net/http"
	"testing"
	"types"
)

type TestGroupRepository struct {
	Groups *[]models.ApiGroup
	Err    error
}

func (testGroupRepository *TestGroupRepository) GetByUser(userid int64) (groups *[]models.UserRole, err error) {
	return nil, nil
}

func (testGroupRepository *TestGroupRepository) GetDefault() (groups *[]models.UserRole, err error) {
	return nil, nil
}

func (testGroupRepository *TestGroupRepository) GetByUserExt(userid int64) (groups *[]models.ApiGroup, err error) {
	return nil, nil
}

func (testGroupRepository *TestGroupRepository) GetBySession(token string) (groups *[]models.UserRole, err error) {
	return nil, nil
}

func (testGroupRepository *TestGroupRepository) GetBySessionExt(token string) (groups *[]models.ApiGroup, err error) {
	return nil, nil
}

func (testGroupRepository *TestGroupRepository) GetAll() (groups *[]models.ApiGroup, err error) {
	return testGroupRepository.Groups, testGroupRepository.Err
}

func (testGroupRepository *TestGroupRepository) SetByUser(userid int64, groups *[]models.UserRole, inTrans bool) (err error) {
	return nil
}

func (testGroupRepository *TestGroupRepository) SetBySession(token string, groups *[]models.UserRole, inTrans bool) (err error) {
	return nil
}

func (testGroupRepository *TestGroupRepository) DeleteByUser(userid int64) (err error) {
	return nil
}

func TestGetGroupsInfoError(t *testing.T) {
	var session = &(models.DtoSession{Language: "eng"})
	var r = new(Renderer)
	var grouprepository = new(TestGroupRepository)
	grouprepository.Groups = nil
	grouprepository.Err = errors.New("Groups error")

	GetGroupsInfo(r, grouprepository, session)
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Get groups info wrong http status and error code")
	}
}

func TestGetGroupsInfoOk(t *testing.T) {
	var session = &(models.DtoSession{Language: "eng"})
	var r = new(Renderer)
	var grouprepository = new(TestGroupRepository)
	grouprepository.Groups = &([]models.ApiGroup{{1, "Developer"}, {2, "Administrator"}, {3, "Supplier"}, {4, "Customer"}})
	grouprepository.Err = nil

	GetGroupsInfo(r, grouprepository, session)
	if r.StatusValue != http.StatusOK {
		t.Error("Get groups info wrong http status")
	}
}
