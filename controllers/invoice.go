package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
	"time"
	"types"
)

// options /api/v1.0/customers/invoices/
func GetMetaInvoices(r render.Render, invoicerepository services.InvoiceRepository, session *models.DtoSession) {
	invoice, err := invoicerepository.GetMeta(session.UserID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, invoice)
}

// get /api/v1.0/customers/invoices/
func GetInvoices(w http.ResponseWriter, request *http.Request, r render.Render, invoicerepository services.InvoiceRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.InvoiceSearch), nil, request, r, session.Language)
	if err != nil {
		return
	}
	if len(*filters) != 0 {
		var masks []string
		for _, filter := range *filters {
			var exps []string
			for _, field := range filter.Fields {
				exps = append(exps, field+" "+filter.Op+" "+filter.Value)
			}
			masks = append(masks, "("+strings.Join(exps, " or ")+")")
		}
		query += " and "
		query += strings.Join(masks, " and ")
	}

	var sorts *[]models.OrderExp
	sorts, err = helpers.GetOrderArray(new(models.InvoiceSearch), request, r, session.Language)
	if err != nil {
		return
	}
	if len(*sorts) != 0 {
		var orders []string
		for _, sort := range *sorts {
			orders = append(orders, " "+sort.Field+" "+sort.Order)
		}
		query += " order by"
		query += strings.Join(orders, ",")
	}

	var limit string
	limit, err = helpers.GetLimitQuery(request, r, session.Language)
	if err != nil {
		return
	}
	query += limit

	invoices, err := invoicerepository.GetByUser(session.UserID, query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(invoices, len(*invoices), w, r)
}

// post /api/v1.0/customers/invoices/
func CreateInvoice(errors binding.Errors, viewinvoice models.ViewInvoice, r render.Render, params martini.Params,
	invoicerepository services.InvoiceRepository, companyrepository services.CompanyRepository,
	invoiceitemrepository services.InvoiceItemRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	company, err := helpers.CheckCompanyAvailability(viewinvoice.Company_ID, session.UserID, r, companyrepository, session.Language)
	if err != nil {
		return
	}

	dtoinvoice := new(models.DtoInvoice)
	dtoinvoice.Company_ID = company.ID
	dtoinvoice.Paid = false
	dtoinvoice.Created = time.Now()
	dtoinvoice.Active = true
	dtoinvoice.Total = viewinvoice.Total
	dtoinvoice.VAT = (dtoinvoice.Total / (1 + helpers.FIELD_VAT_RATE)) * helpers.FIELD_VAT_RATE
	dtoinvoice.InvoiceItems = []models.DtoInvoiceItem{*models.NewDtoInvoiceItem(0, 0, "Оплата по договору", "руб", 1, viewinvoice.Total, viewinvoice.Total)}

	err = invoicerepository.Create(dtoinvoice, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	invoiceitems, err := invoiceitemrepository.GetByInvoice(dtoinvoice.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiFullInvoice(dtoinvoice.ID, dtoinvoice.Company_ID, dtoinvoice.VAT, dtoinvoice.Total,
		*invoiceitems, dtoinvoice.Paid, !dtoinvoice.Active))
}

//get /api/v1.0/customers/invoices/:iid/
func GetInvoice(r render.Render, params martini.Params, invoicerepository services.InvoiceRepository,
	invoiceitemrepository services.InvoiceItemRepository, session *models.DtoSession) {
	dtoinvoice, err := helpers.CheckInvoice(r, params, invoicerepository, session.Language)
	if err != nil {
		return
	}

	invoiceitems, err := invoiceitemrepository.GetByInvoice(dtoinvoice.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongInvoice(dtoinvoice.Company_ID, dtoinvoice.VAT, dtoinvoice.Total,
		*invoiceitems, dtoinvoice.Paid, !dtoinvoice.Active))
}

// patch /api/v1.0/customers/invoices/
func UpdateInvoice(errors binding.Errors, viewinvoice models.ViewInvoice, r render.Render, params martini.Params,
	invoicerepository services.InvoiceRepository, companyrepository services.CompanyRepository,
	invoiceitemrepository services.InvoiceItemRepository, session *models.DtoSession) {
	if helpers.CheckValidation(errors, r, session.Language) != nil {
		return
	}
	if viewinvoice.Company_ID == 0 {
		log.Error("Company is not provided")
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	dtoinvoice, err := helpers.CheckInvoice(r, params, invoicerepository, session.Language)
	if err != nil {
		return
	}
	if !dtoinvoice.Active {
		log.Error("Invoice is not active %v", dtoinvoice.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if dtoinvoice.Paid {
		log.Error("Invoice is paid %v", dtoinvoice.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	company, err := helpers.CheckCompanyAvailability(viewinvoice.Company_ID, session.UserID, r, companyrepository, session.Language)
	if err != nil {
		return
	}

	dtoinvoice.Company_ID = company.ID
	dtoinvoice.Total = viewinvoice.Total
	dtoinvoice.VAT = (dtoinvoice.Total / (1 + helpers.FIELD_VAT_RATE)) * helpers.FIELD_VAT_RATE
	dtoinvoice.InvoiceItems = []models.DtoInvoiceItem{*models.NewDtoInvoiceItem(0, dtoinvoice.ID, "Оплата по договору", "руб", 1,
		viewinvoice.Total, viewinvoice.Total)}

	err = invoicerepository.Update(dtoinvoice, true)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	invoiceitems, err := invoiceitemrepository.GetByInvoice(dtoinvoice.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiLongInvoice(dtoinvoice.Company_ID, dtoinvoice.VAT, dtoinvoice.Total,
		*invoiceitems, dtoinvoice.Paid, !dtoinvoice.Active))
}

// delete /api/v1.0/customers/invoices/:iid/
func DeleteInvoice(r render.Render, params martini.Params, invoicerepository services.InvoiceRepository,
	session *models.DtoSession) {
	dtoinvoice, err := helpers.CheckInvoice(r, params, invoicerepository, session.Language)
	if err != nil {
		return
	}
	if !dtoinvoice.Active {
		log.Error("Invoice is not active %v", dtoinvoice.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	err = invoicerepository.Deactivate(dtoinvoice)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}
