package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/martini-contrib/render"
	"net"
	"net/http"
	"time"
	"types"
)

const (
	PARAM_NAME_USER_ID  = "userid"
	PARAMETER_NAME_CODE = "code"
)

func CheckUser(userid int64, language string, r render.Render, userrepository services.UserRepository) (dtouser *models.DtoUser, err error) {
	dtouser, err = userrepository.Get(userid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtouser.Active || !dtouser.Confirmed {
		log.Error("User is not active or confirmed %v", dtouser.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_USER_BLOCKED,
			Message: config.Localization[language].Errors.Api.User_Blocked})
		return nil, errors.New("User not active")
	}

	return dtouser, nil
}

func CheckUserRoles(roles []models.UserRole, language string, r render.Render,
	grouprepository services.GroupRepository) (err error) {
	groups, err := grouprepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}
	for _, role := range roles {
		found := false
		for _, group := range *groups {
			if group.ID == int(role) {
				found = true
				break
			}
		}
		if !found {
			log.Error("Role is unknown %v", role)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return errors.New("Role unknown")
		}
	}

	return nil
}

func SendPasswordRegistration(language string,
	email *models.DtoEmail,
	user *models.DtoUser,
	request *http.Request,
	r render.Render,
	emailrepository services.EmailRepository,
	templaterepository services.TemplateRepository) (err error) {
	subject := config.Localization[email.Language].Messages.RegistrationSubject
	tpl := services.TEMPLATE_PASSWORD_REGISTRATION
	return sendPassword(language, email, user, request, r, emailrepository, templaterepository, tpl, subject)
}

func SendPasswordRecovery(language string,
	email *models.DtoEmail,
	user *models.DtoUser,
	request *http.Request,
	r render.Render,
	emailrepository services.EmailRepository,
	templaterepository services.TemplateRepository) (err error) {
	subject := config.Localization[email.Language].Messages.PasswordSubject
	tpl := services.TEMPLATE_PASSWORD_RECOVERY
	return sendPassword(language, email, user, request, r, emailrepository, templaterepository, tpl, subject)
}

func sendPassword(language string, email *models.DtoEmail, user *models.DtoUser, request *http.Request, r render.Render,
	emailrepository services.EmailRepository, templaterepository services.TemplateRepository, tpl string, subject string) (err error) {
	host := request.Header.Get(REQUEST_HEADER_X_FORWARDED_FOR)
	if host == "" {
		host, _, err = net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			log.Error("Can't detect ip address %v from %v", err, request.RemoteAddr)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return err
		}
	}
	buf, err := templaterepository.GenerateText(models.NewDtoCodeTemplate(
		models.NewDtoTemplate(email.Email, email.Language, request.Host, time.Now(), host), user.Code), tpl, services.TEMPLATE_DIRECTORY_EMAILS,
		services.TEMPLATE_LAYOUT)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}

	err = emailrepository.SendHTML(email.Email, subject, buf.String(), "", "")
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return err
	}

	return nil
}

func SendConfirmation(language string, email *models.DtoEmail, request *http.Request, r render.Render,
	emailrepository services.EmailRepository, templaterepository services.TemplateRepository) (err error) {
	host := request.Header.Get(REQUEST_HEADER_X_FORWARDED_FOR)
	if host == "" {
		host, _, err = net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			log.Error("Can't detect ip address %v from %v", err, request.RemoteAddr)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return err
		}
	}
	subject := config.Localization[email.Language].Messages.ConfirmationSubject
	buf, err := templaterepository.GenerateText(models.NewDtoTemplate(
		email.Email, email.Language, request.Host, time.Now(), host), services.TEMPLATE_CONFIRMATION, services.TEMPLATE_DIRECTORY_EMAILS,
		services.TEMPLATE_LAYOUT)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return err
	}

	err = emailrepository.SendHTML(email.Email, subject, buf.String(), "", "")
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return err
	}

	return nil
}
