package workflows

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"errors"
	"time"
)

type HeaderWorkflow struct {
	OrderRepository           services.OrderRepository
	FacilityRepository        services.FacilityRepository
	HeaderFacilityRepository  services.HeaderFacilityRepository
	OrderStatusRepository     services.OrderStatusRepository
	InvoiceRepository         services.InvoiceRepository
	CompanyRepository         services.CompanyRepository
	OperationRepository       services.OperationRepository
	TransactionTypeRepository services.TransactionTypeRepository
	TableColumnRepository     services.TableColumnRepository
	UnitRepository            services.UnitRepository
	TableRowRepository        services.TableRowRepository
	PriceRepository           services.PriceRepository
	HeaderProductRepository   services.HeaderProductRepository
	TemplateRepository        services.TemplateRepository
	EmailRepository           services.EmailRepository
}

func NewHeaderWorkflow(orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	headerfacilityrepository services.HeaderFacilityRepository, orderstatusrepository services.OrderStatusRepository,
	invoicerepository services.InvoiceRepository, companyrepository services.CompanyRepository, operationrepository services.OperationRepository,
	transactiontyperepository services.TransactionTypeRepository, tablecolumnrepository services.TableColumnRepository,
	unitrepository services.UnitRepository, tablerowrepository services.TableRowRepository, pricerepository services.PriceRepository,
	headerproductrepository services.HeaderProductRepository, templaterepository services.TemplateRepository,
	emailrepository services.EmailRepository) *HeaderWorkflow {
	return &HeaderWorkflow{
		OrderRepository:           orderrepository,
		FacilityRepository:        facilityrepository,
		HeaderFacilityRepository:  headerfacilityrepository,
		OrderStatusRepository:     orderstatusrepository,
		InvoiceRepository:         invoicerepository,
		CompanyRepository:         companyrepository,
		OperationRepository:       operationrepository,
		TransactionTypeRepository: transactiontyperepository,
		TableColumnRepository:     tablecolumnrepository,
		UnitRepository:            unitrepository,
		TableRowRepository:        tablerowrepository,
		PriceRepository:           pricerepository,
		HeaderProductRepository:   headerproductrepository,
		TemplateRepository:        templaterepository,
		EmailRepository:           emailrepository,
	}
}

func (headerworkflow *HeaderWorkflow) CheckHeaderOrder(dtoheaderfacility *models.DtoHeaderFacility) (err error) {
	return nil
}

func (headerworkflow *HeaderWorkflow) CalculateCost(
	dtoorder *models.DtoOrder, dtoheaderfacility *models.DtoHeaderFacility) (cost float64, err error) {
	headerprices, err := helpers.GetHeaderPrices(models.SERVICE_TYPE_HEADER, dtoorder.Supplier_ID, nil, headerworkflow.PriceRepository,
		headerworkflow.TableColumnRepository, headerworkflow.TableRowRepository, headerworkflow.HeaderProductRepository, "", true)
	if err != nil {
		return 0, err
	}
	cost = 0

	found := false
	for _, headerprice := range *headerprices {
		if !headerprice.Increase {
			cost += headerprice.Price
			found = true
		}
	}
	if !found {
		log.Error("Price list doesn't contain positions for supplier %v", dtoorder.Supplier_ID)
		return 0, errors.New("Empty price list")
	}

	for _, headerprice := range *headerprices {
		if headerprice.Increase {
			cost *= headerprice.PriceIncrease
		}
	}

	return cost, nil
}

func (headerworkflow *HeaderWorkflow) PayAndInvoice(dtoorder *models.DtoOrder, dtoheaderfacility *models.DtoHeaderFacility, balance float64) (err error) {
	if balance < dtoheaderfacility.Cost {
		log.Error("Not enough money at unit balance %v to pay for order %v execution", dtoorder.Unit_ID, dtoorder.ID)
		return errors.New("Not enough money")
	}
	dtocompany, err := headerworkflow.CompanyRepository.GetPrimaryByUnit(dtoorder.Unit_ID)
	if err != nil {
		return
	}
	dtounit, err := headerworkflow.UnitRepository.Get(config.Configuration.SystemAccount)
	if err != nil {
		return err
	}
	dtotransactiontype, err := headerworkflow.TransactionTypeRepository.Get(models.TRANSACTION_TYPE_SERVICE_FEE_HEADER)
	if err != nil {
		return err
	}

	dtotransaction := new(models.DtoTransaction)
	dtotransaction.Source_ID = dtoorder.Unit_ID
	dtotransaction.Destination_ID = dtounit.ID
	dtotransaction.Type_ID = dtotransactiontype.ID

	dtoinvoice := new(models.DtoInvoice)
	dtoinvoice.Company_ID = dtocompany.ID
	if dtocompany.VAT != 0 {
		dtoinvoice.VAT = (dtoheaderfacility.Cost / (1 + float64(dtocompany.VAT)/100)) * float64(dtocompany.VAT) / 100
	}
	dtoinvoice.Total = dtoheaderfacility.Cost
	dtoinvoice.Paid = true
	dtoinvoice.Created = time.Now()
	dtoinvoice.Active = true
	dtoinvoice.InvoiceItems = []models.DtoInvoiceItem{*models.NewDtoInvoiceItem(0, 0, dtotransactiontype.Name, models.INVOICE_ITEM_TYPE_ROUBLE,
		1, dtoinvoice.Total, dtoinvoice.Total)}
	dtoinvoice.PaidAt = time.Now()

	err = headerworkflow.InvoiceRepository.PayForOrder(dtoorder, dtoinvoice, dtotransaction, true)
	if err != nil {
		return err
	}

	return nil
}

func (headerworkflow *HeaderWorkflow) SendReuest(dtoorder *models.DtoOrder, dtoheaderfacility *models.DtoHeaderFacility) (err error) {
	content := config.Localization[config.Configuration.Server.DefaultLanguage].Messages.HeaderRequest + " " + dtoheaderfacility.Name
	buf, err := headerworkflow.TemplateRepository.GenerateText(models.NewDtoHTMLTemplate(content, config.Configuration.Server.DefaultLanguage),
		services.TEMPLATE_FEEDBACK, services.TEMPLATE_DIRECTORY_EMAILS, "")
	if err != nil {
		return err
	}

	header := config.Localization[config.Configuration.Server.DefaultLanguage].Messages.OrderHeader + " " + dtoorder.Name
	err = headerworkflow.EmailRepository.SendHTML(config.Configuration.Mail.Receiver, header, buf.String(),
		"", config.Configuration.Mail.Sender)
	if err != nil {
		return err
	}

	return nil
}

func (headerworkflow *HeaderWorkflow) SetStatus(dtoorder *models.DtoOrder, orderstatus models.OrderStatus, active bool) (err error) {
	dtoorderstatus := models.NewDtoOrderStatus(dtoorder.ID, orderstatus, active, "", time.Now())
	err = headerworkflow.OrderStatusRepository.Save(dtoorderstatus, nil)
	if err != nil {
		return err
	}

	return nil
}

func (headerworkflow *HeaderWorkflow) ExecuteOrder(order_id int64) {
	log.Info("Starting order %v execution at %v", order_id, time.Now())
	log.Info("Checking order type ...")
	dtoorder, err := headerworkflow.OrderRepository.Get(order_id)
	if err != nil {
		return
	}
	dtofacility, err := headerworkflow.FacilityRepository.Get(dtoorder.Facility_ID)
	if err != nil {
		return
	}
	if dtofacility.Alias != models.SERVICE_TYPE_HEADER {
		log.Error("Order service is not macthed to the service method %v", dtoorder.Facility_ID)
		return
	}
	if !dtofacility.Active {
		log.Error("Service is not active %v", dtofacility.ID)
		return
	}
	log.Info("Checking service type ...")
	dtoheaderfacility, err := headerworkflow.HeaderFacilityRepository.Get(dtoorder.ID)
	if err != nil {
		return
	}
	log.Info("Checking order status ...")
	dtoorderstatuses, err := headerworkflow.OrderStatusRepository.GetByOrder(dtoorder.ID)
	if err != nil {
		return
	}

	order := models.NewApiLongOrderFromDto(dtoorder, dtoorderstatuses)
	if order.IsAssembled && order.IsConfirmed && !order.IsOpen && !order.IsCancelled && !order.IsExecuted && !order.IsArchived && !order.IsDeleted {
		/* 1 */

		/* 2 */ err = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_OPEN, true)
		if err != nil {
			_ = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Checking order data ...")
		err = headerworkflow.CheckHeaderOrder(dtoheaderfacility)
		if err != nil {
			_ = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}

		/* 3 */ err = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_SUPPLIER_COST_NEW, true)
		if err != nil {
			_ = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Calculating order cost ...")
		dtoheaderfacility.Cost, err = headerworkflow.CalculateCost(dtoorder, dtoheaderfacility)
		if err != nil {
			_ = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		dtoheaderfacility.CostFactual = dtoheaderfacility.Cost
		err = headerworkflow.HeaderFacilityRepository.Update(dtoheaderfacility)
		if err != nil {
			_ = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}

		/* 4 */ err = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED, true)
		if err != nil {
			_ = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Start order processing ...")
		/* 5 */ dtoorder.Begin_Date = time.Now()
		err = headerworkflow.OrderRepository.Update(dtoorder, &[]models.DtoOrderStatus{
			*models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_MODERATOR_BEGIN, true, "", time.Now())}, nil, true)
		if err != nil {
			_ = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Calculating unit balance ...")
		balance, err := headerworkflow.OperationRepository.CalculateBalance(dtoorder.Unit_ID)
		if err != nil {
			_ = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Order payment and invoicement ...")
		/* 6 */ err = headerworkflow.PayAndInvoice(dtoorder, dtoheaderfacility, balance)
		if err != nil {
			_ = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Order notification ...")
		/* 6 */ err = headerworkflow.SendReuest(dtoorder, dtoheaderfacility)
		if err != nil {
			_ = headerworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Completing order execution %v", time.Now())
	}
}
