package helpers

import (
	"application/models"
	"bytes"
	"errors"
	"github.com/coopernurse/gorp"
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
	User *models.DtoUser
	Err  error
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
	return testUserRepository.User, testUserRepository.Err
}

func (testUserRepository *TestUserRepository) GetAll(filter string) (users *[]models.ApiUserShort, err error) {
	return nil, nil
}

func (testUserRepository *TestUserRepository) GetByUnit(unitid int64) (users *[]models.ApiUserTiny, err error) {
	return nil, nil
}

func (testUserRepository *TestUserRepository) GetMeta() (usermeta *models.ApiUserMeta, err error) {
	return nil, nil
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
	return nil
}

type TestUnitRepository struct {
	Unit *models.DtoUnit
	Err  error
}

func (testUnitRepository *TestUnitRepository) FindByUser(userid int64) (unit *models.DtoUnit, err error) {
	return nil, nil
}

func (testUnitRepository *TestUnitRepository) Get(unitid int64) (unit *models.DtoUnit, err error) {
	return testUnitRepository.Unit, testUnitRepository.Err
}

func (testUnitRepository *TestUnitRepository) GetMeta() (unit *models.ApiShortMetaUnit, err error) {
	return nil, nil
}

func (testUnitRepository *TestUnitRepository) GetAll(filter string) (units *[]models.ApiShortUnit, err error) {
	return nil, nil
}

func (testUnitRepository *TestUnitRepository) Create(unit *models.DtoUnit) (err error) {
	return nil
}

func (testUnitRepository *TestUnitRepository) Update(unit *models.DtoUnit) (err error) {
	return nil
}

func (testUnitRepository *TestUnitRepository) Deactivate(*models.DtoUnit) (err error) {
	return nil
}

type TestEmailRepository struct {
	Found     bool
	Email     *models.DtoEmail
	ExistsErr error
	GetErr    error
	SendErr   error
}

func (testEmailRepository *TestEmailRepository) SendEmail(email string, subject string, body string) (err error) {
	return testEmailRepository.SendErr
}

func (testEmailRepository *TestEmailRepository) Exists(email string) (found bool, err error) {
	return testEmailRepository.Found, testEmailRepository.ExistsErr
}

func (testEmailRepository *TestEmailRepository) FindByCode(code string) (email *models.DtoEmail, err error) {
	return nil, nil
}

func (testEmailRepository *TestEmailRepository) Get(email string) (dtoemail *models.DtoEmail, err error) {
	return testEmailRepository.Email, testEmailRepository.GetErr
}

func (testEmailRepository *TestEmailRepository) GetByUser(userid int64) (emails *[]models.DtoEmail, err error) {
	return nil, nil
}

func (testEmailRepository *TestEmailRepository) Create(email *models.DtoEmail) (err error) {
	return nil
}

func (testEmailRepository *TestEmailRepository) Update(email *models.DtoEmail) (err error) {
	return nil
}

func (testEmailRepository *TestEmailRepository) Delete(email string) (err error) {
	return nil
}

func (testEmailRepository *TestEmailRepository) DeleteByUser(userid int64) (err error) {
	return nil
}

type TestTemplateRepository struct {
	Buf *bytes.Buffer
	Err error
}

func (testTemplateRepository *TestTemplateRepository) GenerateText(
	dtotemplate *models.DtoTemplate, name string, layout string) (buf *bytes.Buffer, err error) {
	return testTemplateRepository.Buf, testTemplateRepository.Err
}

func (testTemplateRepository *TestTemplateRepository) GenerateHTML(name string, w http.ResponseWriter, object interface{}) (err error) {
	return nil
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

func TestCheckUserRolesMatching(t *testing.T) {
	var roles = []models.UserRole{models.USER_ROLE_DEVELOPER, models.USER_ROLE_ADMINISTRATOR}
	var language = "eng"
	var r = new(Renderer)
	var grouprepository = new(TestGroupRepository)
	grouprepository.Groups = &([]models.ApiGroup{{1, "Developer"}, {2, "Administrator"}})
	grouprepository.Err = nil

	err := CheckUserRoles(roles, language, r, grouprepository)
	if err != nil {
		t.Error("Check user roles should not return error")
	}
}

func TestCheckUserRolesNonMatching(t *testing.T) {
	var roles = []models.UserRole{models.USER_ROLE_DEVELOPER, models.USER_ROLE_ADMINISTRATOR}
	var language = "eng"
	var r = new(Renderer)
	var grouprepository = new(TestGroupRepository)
	grouprepository.Groups = &([]models.ApiGroup{{3, "Supplier"}, {4, "Customer"}})
	grouprepository.Err = nil
	var testlogger = new(TestLogger)
	InitLogger(testlogger)

	err := CheckUserRoles(roles, language, r, grouprepository)
	if err == nil {
		t.Error("Check user roles should return error")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Check user roles wrong http status and error code")
	}
}

func TestCheckUserNotFound(t *testing.T) {
	var userid int64 = 1
	var language = "eng"
	var r = new(Renderer)
	var userrepository = new(TestUserRepository)
	userrepository.User = nil
	userrepository.Err = errors.New("User not found")

	dtouser, err := CheckUser(userid, language, r, userrepository)
	if dtouser != nil {
		t.Error("Check user should not return data")
	}
	if err == nil {
		t.Error("Check user should return error")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Check user wrong http status and error code")
	}
}

func TestCheckUserNotActive(t *testing.T) {
	var userid int64 = 1
	var language = "eng"
	var r = new(Renderer)
	var userrepository = new(TestUserRepository)
	userrepository.User = &(models.DtoUser{Active: false, Confirmed: true})
	userrepository.Err = nil
	var testlogger = new(TestLogger)
	InitLogger(testlogger)

	dtouser, err := CheckUser(userid, language, r, userrepository)
	if dtouser != nil {
		t.Error("Check user should not return data")
	}
	if err == nil {
		t.Error("Check user should return error")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_USER_BLOCKED {
		t.Error("Check user wrong http status and error code")
	}
}

func TestCheckUserNotConfirmed(t *testing.T) {
	var userid int64 = 1
	var language = "eng"
	var r = new(Renderer)
	var userrepository = new(TestUserRepository)
	userrepository.User = &(models.DtoUser{Active: true, Confirmed: false})
	userrepository.Err = nil
	var testlogger = new(TestLogger)
	InitLogger(testlogger)

	dtouser, err := CheckUser(userid, language, r, userrepository)
	if dtouser != nil {
		t.Error("Check user should not return data")
	}
	if err == nil {
		t.Error("Check user should return error")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_USER_BLOCKED {
		t.Error("Check user wrong http status and error code")
	}
}

func TestCheckUserFound(t *testing.T) {
	var userid int64 = 1
	var language = "eng"
	var r = new(Renderer)
	var userrepository = new(TestUserRepository)
	userrepository.User = &(models.DtoUser{ID: userid, Active: true, Confirmed: true})
	userrepository.Err = nil

	dtouser, err := CheckUser(userid, language, r, userrepository)
	if dtouser == nil {
		t.Error("Check user should return data")
	}
	if err != nil {
		t.Error("Check user should not return error")
	}
	if dtouser != nil {
		if dtouser.ID != userid {
			t.Error("Check user wrong return data")
		}
	}
}

func TestCheckUnitNotFound(t *testing.T) {
	var unitid int64 = 1
	var language = "eng"
	var r = new(Renderer)
	var unitrepository = new(TestUnitRepository)
	unitrepository.Unit = nil
	unitrepository.Err = errors.New("Unit not found")

	err := CheckUnitValidity(unitid, language, r, unitrepository)
	if err == nil {
		t.Error("Check unit should return error")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Check unit wrong http status and error code")
	}
}

func TestCheckUnitFound(t *testing.T) {
	var unitid int64 = 1
	var language = "eng"
	var r = new(Renderer)
	var unitrepository = new(TestUnitRepository)
	unitrepository.Unit = &(models.DtoUnit{})
	unitrepository.Err = nil

	err := CheckUnitValidity(unitid, language, r, unitrepository)
	if err != nil {
		t.Error("Check unit should not return error")
	}
}

func TestCheckPrimaryEmailZero(t *testing.T) {
	var language = "eng"
	var r = new(Renderer)
	var user = new(models.ViewApiUserFull)
	var testlogger = new(TestLogger)
	InitLogger(testlogger)

	err := CheckPrimaryEmail(user, language, r)
	if err == nil {
		t.Error("Check primary email should return error")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Check primary email wrong http status and error code")
	}
}

func TestCheckPrimaryEmailMany(t *testing.T) {
	var language = "eng"
	var r = new(Renderer)
	var user = new(models.ViewApiUserFull)
	user.Emails = []models.ViewApiEmail{{Primary: true}, {Primary: true}}
	var testlogger = new(TestLogger)
	InitLogger(testlogger)

	err := CheckPrimaryEmail(user, language, r)
	if err == nil {
		t.Error("Check primary email should return error")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Check primary email wrong http status and error code")
	}
}

func TestCheckPrimaryEmailUserNotConfirmed(t *testing.T) {
	var language = "eng"
	var r = new(Renderer)
	var user = new(models.ViewApiUserFull)
	user.Confirmed = false
	user.Emails = []models.ViewApiEmail{{Primary: true, Confirmed: true}}
	var testlogger = new(TestLogger)
	InitLogger(testlogger)

	err := CheckPrimaryEmail(user, language, r)
	if err == nil {
		t.Error("Check primary email should return error")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Check primary email wrong http status and error code")
	}
}

func TestCheckPrimaryEmailEmailNotConfirmed(t *testing.T) {
	var language = "eng"
	var r = new(Renderer)
	var user = new(models.ViewApiUserFull)
	user.Confirmed = true
	user.Emails = []models.ViewApiEmail{{Primary: true, Confirmed: false}}
	var testlogger = new(TestLogger)
	InitLogger(testlogger)

	err := CheckPrimaryEmail(user, language, r)
	if err == nil {
		t.Error("Check primary email should return error")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Check primary email wrong http status and error code")
	}
}

func TestCheckPrimaryEmailOk(t *testing.T) {
	var language = "eng"
	var r = new(Renderer)
	var user = new(models.ViewApiUserFull)
	user.Confirmed = true
	user.Emails = []models.ViewApiEmail{{Primary: true, Confirmed: true}}

	err := CheckPrimaryEmail(user, language, r)
	if err != nil {
		t.Error("Check primary email should not return error")
	}
}

func TestCheckEmailAvailabilityExistsError(t *testing.T) {
	var value = "email"
	var language = "eng"
	var r = new(Renderer)
	var emailrepository = new(TestEmailRepository)
	emailrepository.Found = false
	emailrepository.ExistsErr = errors.New("Email error")
	var testlogger = new(TestLogger)
	InitLogger(testlogger)

	_, err := CheckEmailAvailability(value, language, r, emailrepository)
	if err == nil {
		t.Error("Check email availability should return error")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Check email availability wrong http status and error code")
	}
}

func TestCheckEmailAvailabilityNotExists(t *testing.T) {
	var value = "email"
	var language = "eng"
	var r = new(Renderer)
	var emailrepository = new(TestEmailRepository)
	emailrepository.Found = false
	emailrepository.ExistsErr = nil

	emailExists, err := CheckEmailAvailability(value, language, r, emailrepository)
	if err != nil {
		t.Error("Check email availability should not return error")
	}
	if emailExists {
		t.Error("Check email availability should return not exists status")
	}
}

func TestCheckEmailAvailabilityGetError(t *testing.T) {
	var value = "email"
	var language = "eng"
	var r = new(Renderer)
	var emailrepository = new(TestEmailRepository)
	emailrepository.Found = true
	emailrepository.ExistsErr = nil
	emailrepository.Email = nil
	emailrepository.GetErr = errors.New("Email error")

	emailExists, err := CheckEmailAvailability(value, language, r, emailrepository)
	if err == nil {
		t.Error("Check email availability should return error")
	}
	if !emailExists {
		t.Error("Check email availability should return exists status")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Check email availability wrong http status and error code")
	}
}

func TestCheckEmailAvailabilityExistsConfirmed(t *testing.T) {
	var value = "email"
	var language = "eng"
	var r = new(Renderer)
	var emailrepository = new(TestEmailRepository)
	emailrepository.Found = true
	emailrepository.ExistsErr = nil
	emailrepository.Email = &(models.DtoEmail{Confirmed: true})
	emailrepository.GetErr = nil

	emailExists, err := CheckEmailAvailability(value, language, r, emailrepository)
	if err == nil {
		t.Error("Check email availability should return error")
	}
	if !emailExists {
		t.Error("Check email availability should return exists status")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_EMAIL_INUSE {
		t.Error("Check email availability wrong http status and error code")
	}
}

func TestCheckEmailAvailabilityGet(t *testing.T) {
	var value = "email"
	var language = "eng"
	var r = new(Renderer)
	var emailrepository = new(TestEmailRepository)
	emailrepository.Found = true
	emailrepository.ExistsErr = nil
	emailrepository.Email = &(models.DtoEmail{Confirmed: false})
	emailrepository.GetErr = nil

	emailExists, err := CheckEmailAvailability(value, language, r, emailrepository)
	if err != nil {
		t.Error("Check email availability should not return error")
	}
	if !emailExists {
		t.Error("Check email availability should return exists status")
	}
}

func TestSendConfirmationsZero(t *testing.T) {
	var dtouser = &(models.DtoUser{})
	dtouser.Emails = &([]models.DtoEmail{{Confirmed: true}})
	var session = &(models.DtoSession{Language: "eng"})
	var request = &(http.Request{Host: "http://host.com"})
	var r = new(Renderer)
	var emailrepository = new(TestEmailRepository)
	var templaterepository = new(TestTemplateRepository)

	err := SendConfirmations(dtouser, session, request, r, emailrepository, templaterepository)
	if err != nil {
		t.Error("Send confirmations should not return error")
	}
}

func TestSendConfirmationTemplateErr(t *testing.T) {
	var dtouser = &(models.DtoUser{})
	dtouser.Emails = &([]models.DtoEmail{{Confirmed: false, Primary: true, Language: "eng"}})
	var session = &(models.DtoSession{Language: "eng"})
	var request = &(http.Request{Host: "http://host.com"})
	var r = new(Renderer)
	var emailrepository = new(TestEmailRepository)
	var templaterepository = new(TestTemplateRepository)
	templaterepository.Err = errors.New("Template error")

	err := SendConfirmations(dtouser, session, request, r, emailrepository, templaterepository)
	if err == nil {
		t.Error("Send confirmations should return error")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Send confirmations wrong http status and error code")
	}
}

func TestSendConfirmationEmailErr(t *testing.T) {
	var dtouser = &(models.DtoUser{})
	dtouser.Emails = &([]models.DtoEmail{{Confirmed: false, Primary: false, Language: "rus"}})
	var session = &(models.DtoSession{Language: "eng"})
	var request = &(http.Request{Host: "http://host.com"})
	var r = new(Renderer)
	var emailrepository = new(TestEmailRepository)
	emailrepository.SendErr = errors.New("Email error")
	var templaterepository = new(TestTemplateRepository)
	templaterepository.Buf = new(bytes.Buffer)
	templaterepository.Err = nil

	err := SendConfirmations(dtouser, session, request, r, emailrepository, templaterepository)
	if err == nil {
		t.Error("Send confirmations should return error")
	}
	if r.StatusValue != http.StatusNotFound && r.ErrorValue.Code != types.TYPE_ERROR_DATA_WRONG {
		t.Error("Send confirmations wrong http status and error code")
	}
}

func TestSendConfirmationsOk(t *testing.T) {
	var dtouser = &(models.DtoUser{})
	dtouser.Emails = &([]models.DtoEmail{{Confirmed: false, Primary: true, Language: "rus"}})
	var session = &(models.DtoSession{Language: "eng"})
	var request = &(http.Request{Host: "http://host.com"})
	var r = new(Renderer)
	var emailrepository = new(TestEmailRepository)
	emailrepository.SendErr = nil
	var templaterepository = new(TestTemplateRepository)
	templaterepository.Buf = new(bytes.Buffer)
	templaterepository.Err = nil

	err := SendConfirmations(dtouser, session, request, r, emailrepository, templaterepository)
	if err != nil {
		t.Error("Send confirmations should not return error")
	}
}
