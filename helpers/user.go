package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

const (
	PARAMETER_NAME_CODE = "code"
)

func SendPassword(language string, email *models.DtoEmail, user *models.DtoUser, request *http.Request, r render.Render,
	emailrepository services.EmailRepository, templaterepository services.TemplateRepository) (err error) {
	subject := config.Localization[email.Language].Messages.PasswordSubject
	buf, err := templaterepository.GenerateText(models.NewDtoTemplate(email.Email, email.Language, request.Host, user.Code),
		services.TEMPLATE_PASSWORD, services.TEMPLATE_LAYOUT)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return err
	}

	err = emailrepository.SendEmail(email.Email, subject, buf.String())
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return err
	}

	return nil
}
