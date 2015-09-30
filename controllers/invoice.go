package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"types"
)

// options /api/v1.0/customers/invoices/
func GetMetaInvoices(request *http.Request, r render.Render, invoicerepository services.InvoiceRepository, session *models.DtoSession) {
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

	invoice, err := invoicerepository.GetMeta(session.UserID, query)
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
func CreateInvoice(errors binding.Errors, viewinvoice models.ViewInvoice, r render.Render,
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
	dtoinvoice.VAT = (dtoinvoice.Total / (1 + float64(company.VAT)/100)) * float64(company.VAT) / 100
	dtoinvoice.InvoiceItems = []models.DtoInvoiceItem{*models.NewDtoInvoiceItem(0, 0, models.INVOICE_ITEM_NAME_DEFAULT, models.INVOICE_ITEM_TYPE_ROUBLE,
		1, viewinvoice.Total, viewinvoice.Total)}

	err = invoicerepository.Create(dtoinvoice, nil, true)
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

	r.JSON(http.StatusOK, models.NewApiFullInvoice(dtoinvoice.ID, dtoinvoice.Company_ID, dtoinvoice.Created, dtoinvoice.VAT, dtoinvoice.Total,
		*invoiceitems, dtoinvoice.Paid, dtoinvoice.PaidAt, !dtoinvoice.Active))
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

	r.JSON(http.StatusOK, models.NewApiLongInvoice(dtoinvoice.Company_ID, dtoinvoice.Created, dtoinvoice.VAT, dtoinvoice.Total,
		*invoiceitems, dtoinvoice.Paid, dtoinvoice.PaidAt, !dtoinvoice.Active))
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
	dtoinvoice.VAT = (dtoinvoice.Total / (1 + float64(company.VAT)/100)) * float64(company.VAT) / 100
	dtoinvoice.InvoiceItems = []models.DtoInvoiceItem{*models.NewDtoInvoiceItem(0, dtoinvoice.ID, models.INVOICE_ITEM_NAME_DEFAULT, models.INVOICE_ITEM_TYPE_ROUBLE,
		1, viewinvoice.Total, viewinvoice.Total)}

	err = invoicerepository.Update(dtoinvoice, nil, true)
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

	r.JSON(http.StatusOK, models.NewApiLongInvoice(dtoinvoice.Company_ID, dtoinvoice.Created, dtoinvoice.VAT, dtoinvoice.Total,
		*invoiceitems, dtoinvoice.Paid, dtoinvoice.PaidAt, !dtoinvoice.Active))
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

// get /api/v1.0/customers/invoices/:iid/export/
func GetExportInvoice(request *http.Request, r render.Render, params martini.Params, companyrepository services.CompanyRepository,
	invoicerepository services.InvoiceRepository, unitrepository services.UnitRepository, templaterepository services.TemplateRepository,
	invoiceitemrepository services.InvoiceItemRepository, companycoderepository services.CompanyCodeRepository,
	companyaddressrepository services.CompanyAddressRepository, companybankrepository services.CompanyBankRepository,
	companyemployeerepository services.CompanyEmployeeRepository, filerepository services.FileRepository,
	contractrepository services.ContractRepository, session *models.DtoSession) {
	format, err := url.QueryUnescape(request.URL.Query().Get(helpers.PARAM_QUERY_FORMAT))
	if err != nil {
		log.Error("Can't unescape %v url data", err)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	if strings.ToLower(format) != models.DATA_FORMAT_EXPORT_PDF {
		log.Error("Export is not available for format %v", format)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

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

	dtounit, err := unitrepository.Get(config.Configuration.SystemAccount)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	dtoseller, err := companyrepository.GetPrimaryByUnit(dtounit.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	apiseller, err := helpers.LoadCompany(dtoseller, r, companycoderepository, companyaddressrepository, companybankrepository,
		companyemployeerepository, session.Language)
	if err != nil {
		return
	}
	found := false
	var bankindex int
	for bankindex = range apiseller.CompanyBanks {
		if apiseller.CompanyBanks[bankindex].Primary && !apiseller.CompanyBanks[bankindex].Deleted {
			found = true
			break
		}
	}
	if !found {
		log.Error("Primary bank is not found for company %v", dtoseller.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	templateseller, err := helpers.PrepareCompanyTemplate(dtoseller.ID, apiseller, r, session.Language)
	if err != nil {
		return
	}

	dtobuyer, err := companyrepository.Get(dtoinvoice.Company_ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	apibuyer, err := helpers.LoadCompany(dtobuyer, r, companycoderepository, companyaddressrepository, companybankrepository,
		companyemployeerepository, session.Language)
	if err != nil {
		return
	}
	contracts, err := contractrepository.GetByUnit(dtobuyer.Unit_ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	found = false
	var contractindex int
	for contractindex = range *contracts {
		if (*contracts)[contractindex].Company_ID == dtobuyer.ID && (*contracts)[contractindex].Signed {
			found = true
			break
		}
	}
	if !found {
		log.Error("Contract is not found for company %v", dtobuyer.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	templatebuyer, err := helpers.PrepareCompanyTemplate(dtobuyer.ID, apibuyer, r, session.Language)
	if err != nil {
		return
	}

	buf, err := templaterepository.GenerateText(models.NewDtoInvoiceTemplate(
		*models.NewApiFullInvoice(dtoinvoice.ID, dtoinvoice.Company_ID, dtoinvoice.Created, dtoinvoice.VAT, dtoinvoice.Total,
			*invoiceitems, dtoinvoice.Paid, dtoinvoice.PaidAt, !dtoinvoice.Active), apiseller.CompanyBanks[bankindex],
		*templateseller, *templatebuyer, (*contracts)[contractindex]),
		services.TEMPLATE_INVOICE, "", "")
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	file := new(models.DtoFile)
	file.Created = time.Now()
	file.Name = fmt.Sprintf("invoice_%v.pdf", dtoinvoice.ID)
	file.Path = "/" + fmt.Sprintf("%04d/%02d/%02d/", file.Created.Year(), file.Created.Month(), file.Created.Day())
	file.Permanent = false
	file.Export_Ready = false
	file.Export_Percentage = 0
	file.Export_Object_ID = dtoinvoice.ID
	file.Export_Error = false
	file.Export_ErrorDescription = ""

	err = filerepository.Create(file, nil)
	if err != nil {
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	filename := filepath.Join(config.Configuration.FileStorage, file.Path, fmt.Sprintf("%08d", file.ID))
	err = os.Rename(filename, filename+".html")
	if err != nil {
		log.Error("Can't rename file %v with value %v", err, filename)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	err = ioutil.WriteFile(filename+".html", buf.Bytes(), 0666)
	if err != nil {
		log.Error("Can't save invoice html %v", err)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	absfilepath, err := filepath.Abs(config.Configuration.FileStorage)
	if err != nil {
		log.Error("Can't make an absolute path for %v, %v", config.Configuration.FileStorage, err)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	go helpers.HTMLtoPDF(absfilepath, file, filerepository)

	r.JSON(http.StatusOK, models.ApiFile{ID: file.ID})
}

// options /api/v1.0/customers/invoices/:iid/export/:fid/
func GetExportInvoiceStatus(r render.Render, params martini.Params, filerepository services.FileRepository,
	invoicerepository services.InvoiceRepository, session *models.DtoSession) {
	dtoinvoice, err := helpers.CheckInvoice(r, params, invoicerepository, session.Language)
	if err != nil {
		return
	}

	fileid, err := helpers.CheckParameterInt(r, params[helpers.PARAM_NAME_FILE_ID], session.Language)
	if err != nil {
		return
	}

	file, err := filerepository.Get(fileid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	if dtoinvoice.ID != file.Export_Object_ID {
		log.Error("Linked file object %v and exported invoice %v don't match", file.Export_Object_ID, dtoinvoice.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	r.JSON(http.StatusOK, models.NewApiExportStatus(file.Export_Ready, file.Export_Percentage,
		fmt.Sprintf("%v", file.Created.Add(config.Configuration.FileTimeout)), file.Export_Error, file.Export_ErrorDescription))
}
