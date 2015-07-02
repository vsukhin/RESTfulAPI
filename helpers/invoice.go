package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
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
