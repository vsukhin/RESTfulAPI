package services

import (
	"application/config"
	"github.com/yosssi/ace"
	"html/template"
	"net/http"
	"path/filepath"
)

func (templateservice *TemplateService) GenerateHTML(name string, w http.ResponseWriter, object interface{}) (err error) {
	var tpl *template.Template
	tpl, err = ace.Load(filepath.Join(config.Configuration.Server.TemplateStorage, name), "", nil)
	if err != nil {
		log.Error("Error during loading jade template %v  with value %v", err, name)
		return err
	}
	if err = tpl.Execute(w, object); err != nil {
		log.Error("Error during executing jade template %v", err)
		return err
	}

	return nil
}
