package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"os/exec"
	"path/filepath"
	"types"
)

const (
	PARAM_NAME_INVOICE_ID = "iid"
)

func CheckInvoice(r render.Render, params martini.Params, invoicerepository services.InvoiceRepository,
	language string) (dtoinvoice *models.DtoInvoice, err error) {
	invoice_id, err := CheckParameterInt(r, params[PARAM_NAME_INVOICE_ID], language)
	if err != nil {
		return nil, err
	}

	dtoinvoice, err = invoicerepository.Get(invoice_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return dtoinvoice, nil
}

func HTMLtoPDF(absfilepath string, file *models.DtoFile, filerepository services.FileRepository) {
	command := exec.Command("/bin/sh", "-c", "docker run --rm -m=\"500m\" -v /usr/share/fonts/:/usr/share/fonts/truetype/pdf "+
		"-v "+filepath.Join(absfilepath, file.Path)+":/data webdeskltd/wkhtmltopdf --print-media-type file:///data/"+
		fmt.Sprintf("%08d", file.ID)+".html "+"/data/"+fmt.Sprintf("%08d", file.ID))
	err := command.Run()
	if err == nil {
		file.Export_Ready = true
		file.Export_Percentage = 100
	} else {
		log.Error("Can't convert invoice html to pdf format %v", err)
		file.Export_Error = true
		file.Export_ErrorDescription = err.Error()
	}

	err = filerepository.Update(file)
	if err != nil {
		return
	}
}
