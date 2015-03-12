package services

import (
	"application/config"
	"application/models"
	"bytes"
	"path/filepath"
	"text/template"
)

func (templateservice *TemplateService) GenerateText(dtotemplate *models.DtoTemplate, name string, layout string) (buf *bytes.Buffer, err error) {
	var tpl *template.Template
	buf = new(bytes.Buffer)
	path := filepath.Join(config.Configuration.Server.TemplateStorage, "/mailers")

	if layout == "" {
		tpl, err = template.ParseFiles(filepath.Join(path, name))
	} else {
		tpl, err = template.ParseFiles(filepath.Join(path, layout), filepath.Join(path, name))
	}

	if err != nil {
		log.Error("Error during loading go template %v with value %v", err, name)
		return nil, err
	}

	if err = tpl.Execute(buf, dtotemplate); err != nil {
		log.Error("Error during executing go template %v", err)
		return nil, err
	}

	return buf, nil
}
