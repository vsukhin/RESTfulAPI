package administration

import (
	"application/helpers"
	"application/models"
	"errors"
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
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

type TestUserRepository struct {
	User   *models.DtoUser
	Meta   *models.ApiUserMeta
	GetErr error
	DelErr error
}

func (testUserRepository *TestUserRepository) GetUserArrays(user *models.DtoUser) (*models.DtoUser, error) {
	return nil, nil
}

func (testUserRepository *TestUserRepository) FindByLogin(login string) (user *models.DtoUser, err error) {
	return nil, nil
}

func (testUserRepository *TestUserRepository) FindByCode(code string) (user *models.DtoUser, err error) {
	return nil, nil
}

func (testUserRepository *TestUserRepository) Get(userid int64) (user *models.DtoUser, err error) {
	return testUserRepository.User, testUserRepository.GetErr
}

func (testUserRepository *TestUserRepository) GetAll(filter string) (users *[]models.ApiUserShort, err error) {
	return nil, nil
}

func (testUserRepository *TestUserRepository) GetByUnit(unitid int64) (users *[]models.ApiUserTiny, err error) {
	return nil, nil
}

func (testUserRepository *TestUserRepository) GetMeta() (usermeta *models.ApiUserMeta, err error) {
	return testUserRepository.Meta, testUserRepository.GetErr
}

func (testUserRepository *TestUserRepository) InitUnit(trans *gorp.Transaction, inTrans bool) (unitid int64, err error) {
	return 0, nil
}

func (testUserRepository *TestUserRepository) Create(user *models.DtoUser, inTrans bool) (err error) {
	return nil
}

func (testUserRepository *TestUserRepository) Update(user *models.DtoUser, briefly bool, inTrans bool) (err error) {
	return nil
}

func (testUserRepository *TestUserRepository) Delete(userid int64, inTrans bool) (err error) {
	return testUserRepository.DelErr
}

type TestLogger struct {
}

func (testLogger *TestLogger) Info(query string, args ...interface{}) {
}

func (testLogger *TestLogger) Warning(query string, args ...interface{}) {
}

func (testLogger *TestLogger) Error(query string, args ...interface{}) {
}

func (testLogger *TestLogger) Fatalf(query string, args ...interface{}) {
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

func TestGetUserFullInfoParameterError(t *testing.T) {
	var session = &(models.DtoSession{Language: "eng"})
	var r = new(Renderer)
	var params = martini.Params{helpers.PARAM_NAME_USER_ID: "None"}
	var userrepository = new(TestUserRepository)
	var testlogger = new(TestLogger)
	helpers.InitLogger(testlogger)

	GetUserFullInfo(r, params, userrepository, session)
	if r.StatusValue != http.StatusBadRequest && r.ErrorValue.Code != types.TYPE_ERROR_OBJECT_NOTEXIST {
		t.Error("Get user full info wrong http status and error code")
	}
}

func TestGetUserFullInfoUserError(t *testing.T) {
	var session = &(models.DtoSession{Language: "eng"})
	var r = new(Renderer)
	var params = martini.Params{helpers.PARAM_NAME_USER_ID: "1"}
	var userrepository = new(TestUserRepository)
	userrepository.GetErr = errors.New("User error")
	userrepository.User = nil

	GetUserFullInfo(r, params, userrepository, session)
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_OBJECT_NOTEXIST {
		t.Error("Get user full info wrong http status and error code")
	}
}

func TestGetUserFullInfoUserOk(t *testing.T) {
	var session = &(models.DtoSession{Language: "eng"})
	var r = new(Renderer)
	var params = martini.Params{helpers.PARAM_NAME_USER_ID: "1"}
	var userrepository = new(TestUserRepository)
	userrepository.GetErr = nil
	userrepository.User = &(models.DtoUser{Emails: new([]models.DtoEmail)})

	GetUserFullInfo(r, params, userrepository, session)
	if r.StatusValue != http.StatusOK {
		t.Error("Get user full info wrong http status and error code")
	}
}

func TestGetUserMetaDataError(t *testing.T) {
	var session = &(models.DtoSession{Language: "eng"})
	var r = new(Renderer)
	var userrepository = new(TestUserRepository)
	userrepository.GetErr = errors.New("User error")
	userrepository.Meta = nil

	GetUserMetaData(r, userrepository, session)
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Get user full info wrong http status and error code")
	}
}

func TestGetUserMetaDataOk(t *testing.T) {
	var session = &(models.DtoSession{Language: "eng"})
	var r = new(Renderer)
	var userrepository = new(TestUserRepository)
	userrepository.GetErr = nil
	userrepository.Meta = new(models.ApiUserMeta)

	GetUserMetaData(r, userrepository, session)
	if r.StatusValue != http.StatusOK {
		t.Error("Get user full info wrong http status and error code")
	}
}

func TestDeleteUserParameterError(t *testing.T) {
	var session = &(models.DtoSession{Language: "eng"})
	var r = new(Renderer)
	var params = martini.Params{helpers.PARAM_NAME_USER_ID: "None"}
	var userrepository = new(TestUserRepository)
	var testlogger = new(TestLogger)
	helpers.InitLogger(testlogger)

	DeleteUser(r, params, userrepository, session)
	if r.StatusValue != http.StatusBadRequest && r.ErrorValue.Code != types.TYPE_ERROR_OBJECT_NOTEXIST {
		t.Error("Delete user wrong http status and error code")
	}
}

func TestDeleteUserGetError(t *testing.T) {
	var session = &(models.DtoSession{Language: "eng"})
	var r = new(Renderer)
	var params = martini.Params{helpers.PARAM_NAME_USER_ID: "1"}
	var userrepository = new(TestUserRepository)
	userrepository.GetErr = errors.New("User error")
	userrepository.User = nil

	DeleteUser(r, params, userrepository, session)
	if r.StatusValue == http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_OBJECT_NOTEXIST {
		t.Error("Delete user wrong http status and error code")
	}
}

func TestDeleteUserDelError(t *testing.T) {
	var session = &(models.DtoSession{Language: "eng"})
	var r = new(Renderer)
	var params = martini.Params{helpers.PARAM_NAME_USER_ID: "1"}
	var userrepository = new(TestUserRepository)
	userrepository.GetErr = nil
	userrepository.User = &(models.DtoUser{Emails: new([]models.DtoEmail)})
	userrepository.DelErr = errors.New("User error")
	DeleteUser(r, params, userrepository, session)

	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_OBJECT_NOTEXIST {
		t.Error("Delete user wrong http status and error code")
	}
}

func TestDeleteUserOk(t *testing.T) {
	var session = &(models.DtoSession{Language: "eng"})
	var r = new(Renderer)
	var params = martini.Params{helpers.PARAM_NAME_USER_ID: "1"}
	var userrepository = new(TestUserRepository)
	userrepository.GetErr = nil
	userrepository.DelErr = nil
	userrepository.User = &(models.DtoUser{Emails: new([]models.DtoEmail)})

	DeleteUser(r, params, userrepository, session)
	if r.StatusValue != http.StatusOK {
		t.Error("Delete user wrong http status and error code")
	}
}
