package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

const (
	PARAM_NAME_USER_ID = "userid"
)

func CheckUserRoles(roles []models.UserRole, language string, r render.Render,
	grouprepository services.GroupRepository) (err error) {
	groups, err := grouprepository.GetAll()
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
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
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return errors.New("Role unknown")
		}
	}

	return nil
}

func CheckUser(userid int64, language string, r render.Render, userrepository services.UserRepository) (err error) {
	user, err := userrepository.Get(userid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return err
	}
	if !user.Active || !user.Confirmed {
		log.Error("User is not active or confirmed %v", user.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_USER_BLOCKED,
			Message: config.Localization[language].Errors.Api.User_Blocked})
		return errors.New("User not active")
	}

	return nil
}

func CheckUnit(unitid int64, language string, r render.Render, unitrepository services.UnitRepository) (err error) {
	_, err = unitrepository.Get(unitid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return err
	}

	return nil
}

func CheckPrimaryEmail(user *models.ViewApiUserFull, language string, r render.Render) (err error) {
	count := 0
	for _, checkEmail := range user.Emails {
		if checkEmail.Primary {
			if user.Confirmed != checkEmail.Confirmed {
				log.Error("Confirmation statuses for user and email are different")
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return errors.New("Mismatched statuses")
			}
			count++
		}
	}
	if count != 1 {
		log.Error("Only one primary email is allowed")
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return errors.New("Wrong primary emails amount")
	}

	return nil
}

func CheckEmailAvailability(value string, language string, r render.Render,
	emailrepository services.EmailRepository) (emailExists bool, err error) {
	emailExists, err = emailrepository.Exists(value)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return emailExists, err
	}

	if emailExists {
		email, err := emailrepository.Get(value)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return emailExists, err
		}
		if email.Confirmed {
			log.Error("Email exists in database %v", value)
			r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_EMAIL_INUSE,
				Message: config.Localization[language].Errors.Api.Email_InUse})
			return emailExists, errors.New("Email exists")
		}
	}

	return emailExists, nil
}

func SendConfirmations(dtouser *models.DtoUser, session *models.DtoSession, request *http.Request, r render.Render,
	emailrepository services.EmailRepository, templaterepository services.TemplateRepository) (err error) {
	for _, confEmail := range *dtouser.Emails {
		if !confEmail.Confirmed {
			subject := ""
			if confEmail.Primary {
				subject = config.Localization[confEmail.Language].Messages.RegistrationSubject
			} else {
				subject = config.Localization[confEmail.Language].Messages.EmailSubject
			}

			buf, err := templaterepository.GenerateText(models.NewDtoTemplate(confEmail.Email, confEmail.Language,
				request.Host, confEmail.Code), services.TEMPLATE_EMAIL, services.TEMPLATE_LAYOUT)
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return err
			}

			err = emailrepository.SendEmail(confEmail.Email, subject, buf.String())
			if err != nil {
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
				return err
			}
		}
	}

	return nil
}
