package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"bytes"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

const (
	PARAMETER_NAME_CODE = "code"
)

func SendPassword(language string, email *models.DtoEmail, user *models.DtoUser, request *http.Request, r render.Render,
	emailservice *services.EmailService, templateservice *services.TemplateService) (err error) {
	var buf *bytes.Buffer

	subject := config.Localization[email.Language].Messages.PasswordSubject
	buf, err = templateservice.GenerateText(models.NewDtoTemplate(email.Email, email.Language, request.Host, user.Code),
		services.TEMPLATE_PASSWORD, services.TEMPLATE_LAYOUT)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return err
	}

	err = emailservice.SendEmail(email.Email, subject, buf.String())
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return err
	}

	return nil
}
