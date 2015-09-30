package workflows

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"errors"
	"time"
)

const (
	TABLE_WORK_NAME   = "Recognize work data"
	TABLE_RESULT_NAME = "Recognize result data"
)

type RecognizeWorkflow struct {
	OrderRepository             services.OrderRepository
	FacilityRepository          services.FacilityRepository
	RecognizeFacilityRepository services.RecognizeFacilityRepository
	OrderStatusRepository       services.OrderStatusRepository
	CustomerTableRepository     services.CustomerTableRepository
	ResultTableRepository       services.ResultTableRepository
	WorkTableRepository         services.WorkTableRepository
	InvoiceRepository           services.InvoiceRepository
	CompanyRepository           services.CompanyRepository
	OperationRepository         services.OperationRepository
	TransactionTypeRepository   services.TransactionTypeRepository
	TableColumnRepository       services.TableColumnRepository
	UnitRepository              services.UnitRepository
	TableRowRepository          services.TableRowRepository
	PriceRepository             services.PriceRepository
	RecognizeProductRepository  services.RecognizeProductRepository
	InputFieldRepository        services.InputFieldRepository
	InputProductRepository      services.InputProductRepository
	SupplierRequestRepository   services.SupplierRequestRepository
}

func NewRecognizeWorkflow(orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	recognizefacilityrepository services.RecognizeFacilityRepository, orderstatusrepository services.OrderStatusRepository,
	customertablerepository services.CustomerTableRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, invoicerepository services.InvoiceRepository, companyrepository services.CompanyRepository,
	operationrepository services.OperationRepository, transactiontyperepository services.TransactionTypeRepository,
	tablecolumnrepository services.TableColumnRepository, unitrepository services.UnitRepository,
	tablerowrepository services.TableRowRepository, pricerepository services.PriceRepository,
	recognizeproductrepository services.RecognizeProductRepository, inputfieldrepository services.InputFieldRepository,
	inputproductrepository services.InputProductRepository, supplierrequestrepository services.SupplierRequestRepository) *RecognizeWorkflow {
	return &RecognizeWorkflow{
		OrderRepository:             orderrepository,
		FacilityRepository:          facilityrepository,
		RecognizeFacilityRepository: recognizefacilityrepository,
		OrderStatusRepository:       orderstatusrepository,
		CustomerTableRepository:     customertablerepository,
		ResultTableRepository:       resulttablerepository,
		WorkTableRepository:         worktablerepository,
		InvoiceRepository:           invoicerepository,
		CompanyRepository:           companyrepository,
		OperationRepository:         operationrepository,
		TransactionTypeRepository:   transactiontyperepository,
		TableColumnRepository:       tablecolumnrepository,
		UnitRepository:              unitrepository,
		TableRowRepository:          tablerowrepository,
		PriceRepository:             pricerepository,
		RecognizeProductRepository:  recognizeproductrepository,
		InputFieldRepository:        inputfieldrepository,
		InputProductRepository:      inputproductrepository,
		SupplierRequestRepository:   supplierrequestrepository,
	}
}

func (recognizeworkflow *RecognizeWorkflow) CheckRecognizeOrder(dtorecognizefacility *models.DtoRecognizeFacility) (err error) {
	return nil
}

func (recognizeworkflow *RecognizeWorkflow) CalculateCost(
	dtoorder *models.DtoOrder, dtorecognizefacility *models.DtoRecognizeFacility) (cost float64, err error) {
	recognizeprices, err := helpers.GetRecognizePrices(models.SERVICE_TYPE_RECOGNIZE, dtoorder.Supplier_ID, nil, recognizeworkflow.PriceRepository,
		recognizeworkflow.TableColumnRepository, recognizeworkflow.TableRowRepository, recognizeworkflow.RecognizeProductRepository, "", true)
	if err != nil {
		return 0, err
	}
	cost = 0

	if dtorecognizefacility.EstimatedCalculationOnFields {
		inputfields, err := recognizeworkflow.InputFieldRepository.GetByOrder(dtoorder.ID)
		if err != nil {
			return 0, err
		}
		inputproducts, err := recognizeworkflow.InputProductRepository.GetByOrder(dtoorder.ID)
		if err != nil {
			return 0, err
		}

		for _, inputfield := range *inputfields {
			count := 0
			for _, recognizeprice := range *recognizeprices {
				if inputfield.Product_ID == recognizeprice.Product_ID && !recognizeprice.Increase {
					cost += float64(inputfield.Count) * recognizeprice.Price
					count++
				}
			}
			if count < 1 {
				log.Error("Price list doesn't contain field position %v for supplier %v", inputfield.Product_ID, dtoorder.Supplier_ID)
				return 0, errors.New("Empty field price list")
			}
			if count > 1 {
				log.Error("Price list contains multiple field position %v for supplier %v", inputfield.Product_ID, dtoorder.Supplier_ID)
				return 0, errors.New("Multiple field price list")
			}
		}
		for _, inputproduct := range *inputproducts {
			count := 0
			for _, recognizeprice := range *recognizeprices {
				if inputproduct.Product_ID == recognizeprice.Product_ID && recognizeprice.Increase {
					cost *= recognizeprice.PriceIncrease
					count++
				}
			}
			if count < 1 {
				log.Error("Price list doesn't contain discount position %v for supplier %v", inputproduct.Product_ID, dtoorder.Supplier_ID)
				return 0, errors.New("Empty discount price list")
			}
			if count > 1 {
				log.Error("Price list contains multiple discount position %v for supplier %v", inputproduct.Product_ID, dtoorder.Supplier_ID)
				return 0, errors.New("Multiple discount price list")
			}
		}
		cost *= float64(dtorecognizefacility.EstimatedNumbersForm)
	} else {
		supplierrequests, err := recognizeworkflow.SupplierRequestRepository.GetByOrder(dtoorder.ID)
		if err != nil {
			return 0, err
		}
		found := false
		for _, supplierrequest := range *supplierrequests {
			if supplierrequest.MyChoice && supplierrequest.Supplier_ID == dtoorder.Supplier_ID {
				cost = supplierrequest.EstimatedCost
				found = true
				break
			}
		}
		if !found {
			log.Error("Can't find estimated supplier cost for order %v", dtoorder.ID)
			return 0, errors.New("Not available cost")
		}
	}

	return cost, nil
}

func (recognizeworkflow *RecognizeWorkflow) PayAndInvoice(dtoorder *models.DtoOrder, dtorecognizefacility *models.DtoRecognizeFacility, balance float64) (err error) {
	if balance < dtorecognizefacility.Cost {
		log.Error("Not enough money at unit balance %v to pay for order %v execution", dtoorder.Unit_ID, dtoorder.ID)
		return errors.New("Not enough money")
	}
	dtocompany, err := recognizeworkflow.CompanyRepository.GetPrimaryByUnit(dtoorder.Unit_ID)
	if err != nil {
		return
	}
	dtounit, err := recognizeworkflow.UnitRepository.Get(config.Configuration.SystemAccount)
	if err != nil {
		return err
	}
	dtotransactiontype, err := recognizeworkflow.TransactionTypeRepository.Get(models.TRANSACTION_TYPE_SERVICE_FEE_RECOGNIZE)
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
		dtoinvoice.VAT = (dtorecognizefacility.Cost / (1 + float64(dtocompany.VAT)/100)) * float64(dtocompany.VAT) / 100
	}
	dtoinvoice.Total = dtorecognizefacility.Cost
	dtoinvoice.Paid = true
	dtoinvoice.Created = time.Now()
	dtoinvoice.Active = true
	dtoinvoice.InvoiceItems = []models.DtoInvoiceItem{*models.NewDtoInvoiceItem(0, 0, dtotransactiontype.Name, models.INVOICE_ITEM_TYPE_ROUBLE,
		1, dtoinvoice.Total, dtoinvoice.Total)}
	dtoinvoice.PaidAt = time.Now()

	err = recognizeworkflow.InvoiceRepository.PayForOrder(dtoorder, dtoinvoice, dtotransaction, true)
	if err != nil {
		return err
	}

	return nil
}

func (recognizeworkflow *RecognizeWorkflow) CopyData(dtoorder *models.DtoOrder) (err error) {
	dtoworkdatatable := new(models.DtoCustomerTable)
	dtoworkdatatable.Name = TABLE_WORK_NAME
	dtoworkdatatable.Created = time.Now()
	dtoworkdatatable.TypeID = models.TABLE_TYPE_HIDDEN
	dtoworkdatatable.UnitID = dtoorder.Unit_ID
	dtoworkdatatable.Active = true
	dtoworkdatatable.Permanent = true
	dtoworkdatatable.Signature = models.CUSTOMER_TABLE_SIGNATURE_DEFAULT
	err = recognizeworkflow.CustomerTableRepository.Create(dtoworkdatatable)
	if err != nil {
		return err
	}
	dtoworktable := models.NewDtoWorkTable(dtoorder.ID, dtoworkdatatable.ID)
	err = recognizeworkflow.WorkTableRepository.Create(dtoworktable, nil)
	if err != nil {
		return err
	}

	dtoresultdatatable := new(models.DtoCustomerTable)
	dtoresultdatatable.Name = TABLE_RESULT_NAME
	dtoresultdatatable.Created = time.Now()
	dtoresultdatatable.TypeID = models.TABLE_TYPE_HIDDEN
	dtoresultdatatable.UnitID = dtoorder.Supplier_ID
	dtoresultdatatable.Active = true
	dtoresultdatatable.Permanent = true
	dtoresultdatatable.Signature = models.CUSTOMER_TABLE_SIGNATURE_DEFAULT
	err = recognizeworkflow.CustomerTableRepository.Create(dtoresultdatatable)
	if err != nil {
		return err
	}
	dtoresulttable := models.NewDtoResultTable(dtoorder.ID, dtoresultdatatable.ID)
	err = recognizeworkflow.ResultTableRepository.Create(dtoresulttable, nil)
	if err != nil {
		return err
	}

	return nil
}

func (recognizeworkflow *RecognizeWorkflow) SetStatus(dtoorder *models.DtoOrder, orderstatus models.OrderStatus, active bool) (err error) {
	dtoorderstatus := models.NewDtoOrderStatus(dtoorder.ID, orderstatus, active, "", time.Now())
	err = recognizeworkflow.OrderStatusRepository.Save(dtoorderstatus, nil)
	if err != nil {
		return err
	}

	return nil
}

func (recognizeworkflow *RecognizeWorkflow) ExecuteOrder(order_id int64) {
	dtoorder, err := recognizeworkflow.OrderRepository.Get(order_id)
	if err != nil {
		return
	}
	dtofacility, err := recognizeworkflow.FacilityRepository.Get(dtoorder.Facility_ID)
	if err != nil {
		return
	}
	if dtofacility.Alias != models.SERVICE_TYPE_RECOGNIZE {
		log.Error("Order service is not macthed to the service method %v", dtoorder.Facility_ID)
		return
	}
	dtorecognizefacility, err := recognizeworkflow.RecognizeFacilityRepository.Get(dtoorder.ID)
	if err != nil {
		return
	}
	dtoorderstatuses, err := recognizeworkflow.OrderStatusRepository.GetByOrder(dtoorder.ID)
	if err != nil {
		return
	}

	order := models.NewApiLongOrderFromDto(dtoorder, dtoorderstatuses)
	if order.IsAssembled && order.IsConfirmed && !order.IsOpen && !order.IsCancelled && !order.IsExecuted && !order.IsArchived && !order.IsDeleted {
		/* 1 */

		/* 2 */ err = recognizeworkflow.SetStatus(dtoorder, models.ORDER_STATUS_OPEN, true)
		if err != nil {
			_ = recognizeworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		err = recognizeworkflow.CheckRecognizeOrder(dtorecognizefacility)
		if err != nil {
			_ = recognizeworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}

		/* 3 */ err = recognizeworkflow.SetStatus(dtoorder, models.ORDER_STATUS_SUPPLIER_COST_NEW, true)
		if err != nil {
			_ = recognizeworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		dtorecognizefacility.Cost, err = recognizeworkflow.CalculateCost(dtoorder, dtorecognizefacility)
		if err != nil {
			_ = recognizeworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		dtorecognizefacility.CostFactual = dtorecognizefacility.Cost
		err = recognizeworkflow.RecognizeFacilityRepository.Update(dtorecognizefacility, true, false)
		if err != nil {
			_ = recognizeworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}

		/* 4 */ err = recognizeworkflow.SetStatus(dtoorder, models.ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED, true)
		if err != nil {
			_ = recognizeworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}

		/* 5 */ dtoorder.Begin_Date = time.Now()
		err = recognizeworkflow.OrderRepository.Update(dtoorder, &[]models.DtoOrderStatus{
			*models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_MODERATOR_BEGIN, true, "", time.Now())}, nil, true)
		if err != nil {
			_ = recognizeworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		balance, err := recognizeworkflow.OperationRepository.CalculateBalance(dtoorder.Unit_ID)
		if err != nil {
			_ = recognizeworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}

		/* 6 */ err = recognizeworkflow.PayAndInvoice(dtoorder, dtorecognizefacility, balance)
		if err != nil {
			_ = recognizeworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}

		/* 7 */ err = recognizeworkflow.CopyData(dtoorder)
		if err != nil {
			return
		}

		/* 8 */ dtoorder.End_Date = time.Now()
		err = recognizeworkflow.OrderRepository.Update(dtoorder, &[]models.DtoOrderStatus{
			*models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_SUPPLIER_CLOSE, true, "", time.Now())}, nil, true)
		if err != nil {
			return
		}
	}
}
