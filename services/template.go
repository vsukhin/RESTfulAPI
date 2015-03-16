package services

import (
	"application/models"
	"bytes"
	"net/http"
)

const (
	TEMPLATE_LAYOUT   = "layout.txt.tmpl"
	TEMPLATE_EMAIL    = "registration.txt.tmpl"
	TEMPLATE_PASSWORD = "user_password_reset.txt.tmpl"
)

type TemplateRepository interface {
	GenerateText(dtotemplate *models.DtoTemplate, name string, layout string) (buf *bytes.Buffer, err error)
	GenerateHTML(name string, w http.ResponseWriter, object interface{}) (err error)
}

type TemplateService struct {
}

func NewTemplateService() *TemplateService {
	return &TemplateService{}
}
