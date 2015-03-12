package services

const (
	TEMPLATE_LAYOUT   = "layout.txt.tmpl"
	TEMPLATE_EMAIL    = "registration.txt.tmpl"
	TEMPLATE_PASSWORD = "user_password_reset.txt.tmpl"
)

type TemplateService struct {
}

func NewTemplateService() *TemplateService {
	return &TemplateService{}
}
