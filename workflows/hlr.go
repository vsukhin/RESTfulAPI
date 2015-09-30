package workflows

import (
	"application/communication/suppliers"
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	libTypes "lib/suppliers/types"
	"time"
)

const (
	HLR_LENGTH_ASCII  = 153
	HLR_LENGTH_UICODE = 67

	COLUMN_NAME_HLR_ID           = "HlrId"
	COLUMN_NAME_HLR_ERROR        = "HlrError"
	COLUMN_NAME_HLR_STATUS_ID    = "HlrStatusId"
	COLUMN_NAME_HLR_STATUS_ERROR = "HlrStatusError"
	COLUMN_NAME_HLR_ORN          = "HlrOrn"
	COLUMN_NAME_HLR_ROAMING      = "HlrRoaming"
)

type HLRWorkflow struct {
	OrderRepository           services.OrderRepository
	FacilityRepository        services.FacilityRepository
	HLRFacilityRepository     services.HLRFacilityRepository
	OrderStatusRepository     services.OrderStatusRepository
	CustomerTableRepository   services.CustomerTableRepository
	HLRTableRepository        services.HLRTableRepository
	ResultTableRepository     services.ResultTableRepository
	WorkTableRepository       services.WorkTableRepository
	InvoiceRepository         services.InvoiceRepository
	CompanyRepository         services.CompanyRepository
	OperationRepository       services.OperationRepository
	TransactionTypeRepository services.TransactionTypeRepository
	TableColumnRepository     services.TableColumnRepository
	UnitRepository            services.UnitRepository
	TableRowRepository        services.TableRowRepository
	PriceRepository           services.PriceRepository
	MobileOperatorRepository  services.MobileOperatorRepository
	ColumnTypeRepository      services.ColumnTypeRepository
}

func NewHLRWorkflow(orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	hlrfacilityrepository services.HLRFacilityRepository, orderstatusrepository services.OrderStatusRepository,
	customertablerepository services.CustomerTableRepository, hlrtablerepository services.HLRTableRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository,
	invoicerepository services.InvoiceRepository, companyrepository services.CompanyRepository,
	operationrepository services.OperationRepository, transactiontyperepository services.TransactionTypeRepository,
	tablecolumnrepository services.TableColumnRepository, unitrepository services.UnitRepository,
	tablerowrepository services.TableRowRepository, pricerepository services.PriceRepository,
	mobileoperatorrepository services.MobileOperatorRepository, columntyperepository services.ColumnTypeRepository) *HLRWorkflow {
	return &HLRWorkflow{
		OrderRepository:           orderrepository,
		FacilityRepository:        facilityrepository,
		HLRFacilityRepository:     hlrfacilityrepository,
		OrderStatusRepository:     orderstatusrepository,
		CustomerTableRepository:   customertablerepository,
		HLRTableRepository:        hlrtablerepository,
		ResultTableRepository:     resulttablerepository,
		WorkTableRepository:       worktablerepository,
		InvoiceRepository:         invoicerepository,
		CompanyRepository:         companyrepository,
		OperationRepository:       operationrepository,
		TransactionTypeRepository: transactiontyperepository,
		TableColumnRepository:     tablecolumnrepository,
		UnitRepository:            unitrepository,
		TableRowRepository:        tablerowrepository,
		PriceRepository:           pricerepository,
		MobileOperatorRepository:  mobileoperatorrepository,
		ColumnTypeRepository:      columntyperepository,
	}
}

func SendHLR(supplier *libTypes.Supplier, hlr *[]libTypes.Hlr) (hlrresponse *libTypes.HlrResponse, err error) {
	// Отправка HLR
	suppliers.WaitClientReady()                         // Функция блокируется до момента пока связь с сервером не будет установлена
	response, err := suppliers.SendHlr(*supplier, *hlr) // Функция не блокируется. Если связи с сервером нет, то вернётся ошибка
	if err != nil {
		log.Error("Can't send HLR requests %v for supplier %v", err, supplier.Name)
		return nil, err
	}

	return &response, nil
}

func GetHLRStatus(hlrresponse *libTypes.HlrResponse) (hlrstatuses map[gocql.UUID]libTypes.HlrStatus, err error) {
	hlrstatuses = make(map[gocql.UUID]libTypes.HlrStatus)
	// Получение статуса по UUID запроса
	for {
		final := true
		for i := range hlrresponse.Ids {
			var status libTypes.HlrStatus
			status, err = suppliers.StatusHlr(hlrresponse.Ids[i])
			if err != nil {
				log.Error("Can't get HLR statuses %v", err)
				return map[gocql.UUID]libTypes.HlrStatus{}, err
			}
			if status.Final {
				hlrstatuses[status.Id] = status
			} else {
				final = false
			}
		}
		if final {
			break
		}
		time.Sleep(time.Second)
	}

	return hlrstatuses, nil
}

func (hlrworkflow *HLRWorkflow) CheckHLROrder(dtohlrfacility *models.DtoHLRFacility) (columnmobilephone_id int64, err error) {
	hlrtable, err := hlrworkflow.HLRTableRepository.Get(dtohlrfacility.DeliveryDataId)
	if err != nil {
		return 0, err
	}

	if dtohlrfacility.MessageToInColumnId != 0 {
		found := false
		for _, mobilephone := range hlrtable.MobilePhones {
			if mobilephone.ID == dtohlrfacility.MessageToInColumnId {
				found = true
				columnmobilephone_id = dtohlrfacility.MessageToInColumnId
				break
			}
		}
		if !found {
			log.Error("Can't find mobile phone column %v in table %v", dtohlrfacility.MessageToInColumnId, dtohlrfacility.DeliveryDataId)
			return 0, errors.New("Wrong mobile phone column")
		}
	} else {
		if len(hlrtable.MobilePhones) == 0 {
			log.Error("Can't find mobile column in table %v", dtohlrfacility.DeliveryDataId)
			return 0, errors.New("Missed mobile phone column")
		} else {
			columnmobilephone_id = hlrtable.MobilePhones[0].ID
		}
	}

	return columnmobilephone_id, nil
}

func (hlrworkflow *HLRWorkflow) GetHLRTableColumns(dtohlrfacility *models.DtoHLRFacility, columnmobilephone_id int64) (
	columnmobilephone *models.DtoTableColumn, tablecolumns *[]models.DtoTableColumn, err error) {
	tablecolumns = new([]models.DtoTableColumn)
	columnmobilephone, err = hlrworkflow.TableColumnRepository.Get(columnmobilephone_id)
	if err != nil {
		return nil, nil, err
	}
	*tablecolumns = append(*tablecolumns, *columnmobilephone)

	return columnmobilephone, tablecolumns, nil
}

func (hlrworkflow *HLRWorkflow) CalculateCost(apitablerows *[]models.ApiInfoTableRow, columnmobilephone *models.DtoTableColumn,
	dtoorder *models.DtoOrder, dtohlrfacility *models.DtoHLRFacility) (cost float64, err error) {
	dtomobileoperators, err := hlrworkflow.MobileOperatorRepository.FindAll()
	if err != nil {
		return 0, err
	}
	mobileoperators_uuid := make(map[string]*models.DtoMobileOperator)
	mobileoperators_id := make(map[int]*models.DtoMobileOperator)
	for index := range *dtomobileoperators {
		mobileoperators_uuid[(*dtomobileoperators)[index].UUID] = &(*dtomobileoperators)[index]
		mobileoperators_id[(*dtomobileoperators)[index].ID] = &(*dtomobileoperators)[index]
	}
	dtocolumntype, err := hlrworkflow.ColumnTypeRepository.Get(columnmobilephone.Column_Type_ID)
	if err != nil {
		return 0, err
	}

	mobilephones := []uint64{}
	for _, apitablerow := range *apitablerows {
		var mobilephone uint64 = 0
		for _, apitablecell := range apitablerow.Cells {
			if apitablecell.Table_Column_ID == columnmobilephone.ID {
				mobilephone, err = helpers.CheckMobilePhone(apitablecell.Value, dtocolumntype, hlrworkflow.ColumnTypeRepository)
				if err != nil {
					return 0, err
				}
			}
		}
		mobilephones = append(mobilephones, mobilephone)
	}

	mobileoperatoruuids, err := suppliers.MobileOperator(mobilephones)
	if err != nil {
		log.Error("Error during detecting mobile operators %v", err)
		return 0, errors.New("Mobile operator detection error")
	}
	defaultmobileoperator, err := hlrworkflow.MobileOperatorRepository.GetDefault()
	if err != nil {
		return 0, err
	}

	mobileoperatorhlres := make(map[int]int)
	for index := range mobileoperatoruuids {
		if mobileoperatoruuids[index].Id.String() == models.MOBILE_OPERATOR_UUID_UNKNOWN {
			log.Error("Not existed mobile operator for phone %v", mobilephones[index])
			return 0, errors.New("Not existed mobile operator")
		}
		if mobileoperatoruuids[index].Id.String() == models.MOBILE_OPERATOR_UUID_MEGAFON {
			log.Error("Not HLR enabled mobile operator for phone %v", mobilephones[index])
			return 0, errors.New("Not HLR enabled mobile operator")
		}
		mobileoperator, ok := mobileoperators_uuid[mobileoperatoruuids[index].Id.String()]
		if !ok {
			mobileoperator = defaultmobileoperator
		}
		mobileoperatorhlres[mobileoperator.ID]++
	}

	smshlrprices, err := helpers.GetSMSHLRPrices(models.SERVICE_TYPE_HLR, dtoorder.Supplier_ID, nil, hlrworkflow.PriceRepository,
		hlrworkflow.TableColumnRepository, hlrworkflow.TableRowRepository,
		hlrworkflow.MobileOperatorRepository, "", true)
	if err != nil {
		return 0, err
	}

	pricemobileoperators := make(map[int]int)
	for _, smshlrprice := range *smshlrprices {
		pricemobileoperators[smshlrprice.Mobile_Operator_ID] = smshlrprice.Mobile_Operator_ID
	}

	cost = 0
	for mobileoperator_id, count := range mobileoperatorhlres {
		pricemobileoperator_id, ok := pricemobileoperators[mobileoperator_id]
		if !ok {
			pricemobileoperator_id, ok = pricemobileoperators[defaultmobileoperator.ID]
			if !ok {
				log.Error("Can't find default mobile operator %v in price list for supplier %v", defaultmobileoperator.ID, dtoorder.Supplier_ID)
				return 0, errors.New("Missed default mobile operator in price list")
			}
		}

		found := false
		if mobileoperators_id[pricemobileoperator_id].HLRBillingModel == models.BILLING_MODEL_RANGE {
			for _, smshlrprice := range *smshlrprices {
				if smshlrprice.Mobile_Operator_ID == pricemobileoperator_id &&
					count >= smshlrprice.AmountRange.Begin && (count <= smshlrprice.AmountRange.End || smshlrprice.AmountRange.End == 0) {
					cost += float64(count) * smshlrprice.Price
					found = true
					break
				}
			}
		}
		if !found {
			log.Error("Can't find  mobile operator %v in price list for supplier %v", pricemobileoperator_id, dtoorder.Supplier_ID)
			return 0, errors.New("Missed obile operator in price list")
		}
	}

	return cost, nil
}

func (hlrworkflow *HLRWorkflow) PayAndInvoice(dtoorder *models.DtoOrder, dtohlrfacility *models.DtoHLRFacility, balance float64) (err error) {
	if balance < dtohlrfacility.Cost {
		log.Error("Not enough money at unit balance %v to pay for order %v execution", dtoorder.Unit_ID, dtoorder.ID)
		return errors.New("Not enough money")
	}
	dtocompany, err := hlrworkflow.CompanyRepository.GetPrimaryByUnit(dtoorder.Unit_ID)
	if err != nil {
		return
	}
	dtounit, err := hlrworkflow.UnitRepository.Get(config.Configuration.SystemAccount)
	if err != nil {
		return err
	}
	dtotransactiontype, err := hlrworkflow.TransactionTypeRepository.Get(models.TRANSACTION_TYPE_SERVICE_FEE_HLR)
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
		dtoinvoice.VAT = (dtohlrfacility.Cost / (1 + float64(dtocompany.VAT)/100)) * float64(dtocompany.VAT) / 100
	}
	dtoinvoice.Total = dtohlrfacility.Cost
	dtoinvoice.Paid = true
	dtoinvoice.Created = time.Now()
	dtoinvoice.Active = true
	dtoinvoice.InvoiceItems = []models.DtoInvoiceItem{*models.NewDtoInvoiceItem(0, 0, dtotransactiontype.Name, models.INVOICE_ITEM_TYPE_ROUBLE,
		1, dtoinvoice.Total, dtoinvoice.Total)}
	dtoinvoice.PaidAt = time.Now()

	err = hlrworkflow.InvoiceRepository.PayForOrder(dtoorder, dtoinvoice, dtotransaction, true)
	if err != nil {
		return err
	}

	return nil
}

func (hlrworkflow *HLRWorkflow) SendHLR(apitablerows *[]models.ApiInfoTableRow, dtoorder *models.DtoOrder, dtohlrfacility *models.DtoHLRFacility,
	columnmobilephone *models.DtoTableColumn) (hlrresponse *libTypes.HlrResponse, err error) {
	dtosupplier, err := hlrworkflow.UnitRepository.Get(dtoorder.Supplier_ID)
	if err != nil {
		return nil, err
	}
	hlrsupplier, err := GetSupplier(dtosupplier.UUID)
	if err != nil {
		return nil, err
	}

	dtocolumntype, err := hlrworkflow.ColumnTypeRepository.Get(columnmobilephone.Column_Type_ID)
	if err != nil {
		return nil, err
	}

	hlr := new([]libTypes.Hlr)
	for _, apitablerow := range *apitablerows {
		obj := new(libTypes.Hlr)
		for _, apitablecell := range apitablerow.Cells {
			if apitablecell.Table_Column_ID == columnmobilephone.ID {
				obj.Recipient, err = helpers.CheckMobilePhone(apitablecell.Value, dtocolumntype, hlrworkflow.ColumnTypeRepository)
				if err != nil {
					return nil, err
				}
			}
		}
		*hlr = append(*hlr, *obj)
	}
	hlrresponse, err = SendHLR(hlrsupplier, hlr)
	if err != nil {
		return nil, err
	}

	return hlrresponse, nil
}

func (hlrworkflow *HLRWorkflow) CopyData(dtoorder *models.DtoOrder,
	dtodatatable *models.DtoCustomerTable) (dtoworkdatatable *models.DtoCustomerTable, err error) {
	dtoworkdatatable, err = hlrworkflow.CustomerTableRepository.Copy(dtodatatable, true)
	if err != nil {
		return nil, err
	}
	dtodatatable.TypeID = models.TABLE_TYPE_READONLY
	err = hlrworkflow.CustomerTableRepository.Update(dtodatatable)
	if err != nil {
		return nil, err
	}
	dtoworkdatatable.TypeID = models.TABLE_TYPE_HIDDEN
	err = hlrworkflow.CustomerTableRepository.Update(dtoworkdatatable)
	if err != nil {
		return nil, err
	}

	dtoworktable := models.NewDtoWorkTable(dtoorder.ID, dtoworkdatatable.ID)
	err = hlrworkflow.WorkTableRepository.Create(dtoworktable, nil)
	if err != nil {
		return nil, err
	}

	return dtoworkdatatable, nil
}

func (hlrworkflow *HLRWorkflow) SaveHLR(dtoworkdatatable *models.DtoCustomerTable, hlrresponse *libTypes.HlrResponse) (err error) {
	position, err := hlrworkflow.TableColumnRepository.GetDefaultPosition(dtoworkdatatable.ID)
	if err != nil {
		return err
	}
	tablecolumns := []models.DtoTableColumn{}
	columnnames := []string{COLUMN_NAME_HLR_ID, COLUMN_NAME_HLR_ERROR}
	for _, columnname := range columnnames {
		position++
		fieldnum, err := helpers.FindFreeColumnInternal(dtoworkdatatable.ID, 0, hlrworkflow.TableColumnRepository)
		if err != nil {
			return err
		}

		dtotablecolumn := new(models.DtoTableColumn)
		dtotablecolumn.Created = time.Now()
		dtotablecolumn.Position = position
		dtotablecolumn.Name = columnname
		dtotablecolumn.Customer_Table_ID = dtoworkdatatable.ID
		dtotablecolumn.Column_Type_ID = models.COLUMN_TYPE_DEFAULT
		dtotablecolumn.Prebuilt = true
		dtotablecolumn.FieldNum = fieldnum
		dtotablecolumn.Active = true
		dtotablecolumn.Edition = 0

		err = hlrworkflow.TableColumnRepository.Create(dtotablecolumn, nil)
		if err != nil {
			return err
		}
		tablecolumns = append(tablecolumns, *dtotablecolumn)
	}

	apitablerows, err := hlrworkflow.TableRowRepository.GetAll("", "", dtoworkdatatable.ID, &tablecolumns)
	if err != nil {
		return err
	}

	workdatatablecolumns, err := hlrworkflow.TableColumnRepository.GetByTable(dtoworkdatatable.ID)
	if err != nil {
		return err
	}

	for index, apitablerow := range *apitablerows {
		tablerow, err := hlrworkflow.TableRowRepository.Get(apitablerow.ID)
		if err != nil {
			return err
		}

		tablecells, err := tablerow.TableRowToDtoTableCells(workdatatablecolumns)
		if err != nil {
			return err
		}

		for i, _ := range *tablecells {
			for j, _ := range tablecolumns {
				FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_HLR_ID, hlrresponse.Ids[index].String())
				FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_HLR_ERROR, hlrresponse.Errors[index].Error())
			}
		}

		err = tablerow.TableCellsToTableRow(tablecells, workdatatablecolumns)
		if err != nil {
			return err
		}

		err = hlrworkflow.TableRowRepository.Update(tablerow, nil, true, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (hlrworkflow *HLRWorkflow) SaveHLRStatus(dtoworkdatatable *models.DtoCustomerTable, hlrstatuses map[gocql.UUID]libTypes.HlrStatus) (err error) {
	var columnhlrid *models.DtoTableColumn
	alltablecolumns, err := hlrworkflow.TableColumnRepository.GetByTable(dtoworkdatatable.ID)
	if err != nil {
		return err
	}
	found := false
	for i := range *alltablecolumns {
		if (*alltablecolumns)[i].Name == COLUMN_NAME_HLR_ID {
			columnhlrid = &(*alltablecolumns)[i]
			found = true
			break
		}
	}
	if !found {
		log.Error("Can't find hlr id column for table %v", dtoworkdatatable.ID)
		return errors.New("Missed hlr id column")
	}

	position, err := hlrworkflow.TableColumnRepository.GetDefaultPosition(dtoworkdatatable.ID)
	if err != nil {
		return err
	}
	tablecolumns := []models.DtoTableColumn{*columnhlrid}
	columnnames := []string{COLUMN_NAME_HLR_STATUS_ID, COLUMN_NAME_HLR_STATUS_ERROR, COLUMN_NAME_HLR_ORN, COLUMN_NAME_HLR_ROAMING}
	for _, columnname := range columnnames {
		position++
		fieldnum, err := helpers.FindFreeColumnInternal(dtoworkdatatable.ID, 0, hlrworkflow.TableColumnRepository)
		if err != nil {
			return err
		}

		dtotablecolumn := new(models.DtoTableColumn)
		dtotablecolumn.Created = time.Now()
		dtotablecolumn.Position = position
		dtotablecolumn.Name = columnname
		dtotablecolumn.Customer_Table_ID = dtoworkdatatable.ID
		dtotablecolumn.Column_Type_ID = models.COLUMN_TYPE_DEFAULT
		dtotablecolumn.Prebuilt = true
		dtotablecolumn.FieldNum = fieldnum
		dtotablecolumn.Active = true
		dtotablecolumn.Edition = 0

		err = hlrworkflow.TableColumnRepository.Create(dtotablecolumn, nil)
		if err != nil {
			return err
		}
		tablecolumns = append(tablecolumns, *dtotablecolumn)
	}

	apitablerows, err := hlrworkflow.TableRowRepository.GetAll("", "", dtoworkdatatable.ID, &tablecolumns)
	if err != nil {
		return err
	}

	workdatatablecolumns, err := hlrworkflow.TableColumnRepository.GetByTable(dtoworkdatatable.ID)
	if err != nil {
		return err
	}

	for _, apitablerow := range *apitablerows {
		tablerow, err := hlrworkflow.TableRowRepository.Get(apitablerow.ID)
		if err != nil {
			return err
		}

		tablecells, err := tablerow.TableRowToDtoTableCells(workdatatablecolumns)
		if err != nil {
			return err
		}

		var hlrid gocql.UUID
		found := false
		for i, _ := range *tablecells {
			if (*tablecells)[i].Table_Column_ID == columnhlrid.ID {
				hlrid, err = gocql.ParseUUID((*tablecells)[i].Value)
				if err != nil {
					log.Error("Can't parse hlr id from column %v, %v", err, columnhlrid.ID)
					return err
				}
				found = true
				break
			}
		}
		if !found {
			log.Error("Can't find hlr id value for column %v", columnhlrid.ID)
			return errors.New("Missed hlr id value")
		}

		for i, _ := range *tablecells {
			for j, _ := range tablecolumns {
				FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_HLR_STATUS_ID, fmt.Sprintf("%v", hlrstatuses[hlrid].Status))
				FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_HLR_STATUS_ERROR, fmt.Sprintf("%v", hlrstatuses[hlrid].Error))
				FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_HLR_ORN, fmt.Sprintf("%v", hlrstatuses[hlrid].ResponseOrn))
				FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_HLR_ROAMING, fmt.Sprintf("%v", hlrstatuses[hlrid].ResponseIsRoaming))
			}
		}

		err = tablerow.TableCellsToTableRow(tablecells, workdatatablecolumns)
		if err != nil {
			return err
		}

		err = hlrworkflow.TableRowRepository.Update(tablerow, nil, true, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (hlrworkflow *HLRWorkflow) ClearTables(dtoorder *models.DtoOrder, dtohlrfacility *models.DtoHLRFacility,
	dtodatatable *models.DtoCustomerTable) (err error) {
	if dtohlrfacility.DeliveryDataDelete {
		err = hlrworkflow.CustomerTableRepository.Deactivate(dtodatatable)
		if err != nil {
			return err
		}
	} else {
		dtodatatable.TypeID = models.TABLE_TYPE_DEFAULT
		err = hlrworkflow.CustomerTableRepository.Update(dtodatatable)
		if err != nil {
			return err
		}
	}

	worktables, err := hlrworkflow.WorkTableRepository.GetByOrder(dtoorder.ID)
	if err != nil {
		return err
	}
	for _, worktable := range *worktables {
		dtocustomertable, err := hlrworkflow.CustomerTableRepository.Get(worktable.Customer_Table_ID)
		if err != nil {
			return err
		}
		err = hlrworkflow.CustomerTableRepository.Deactivate(dtocustomertable)
		if err != nil {
			return err
		}
	}

	return nil
}

func (hlrworkflow *HLRWorkflow) SetStatus(dtoorder *models.DtoOrder, orderstatus models.OrderStatus, active bool) (err error) {
	dtoorderstatus := models.NewDtoOrderStatus(dtoorder.ID, orderstatus, active, "", time.Now())
	err = hlrworkflow.OrderStatusRepository.Save(dtoorderstatus, nil)
	if err != nil {
		return err
	}

	return nil
}

func (hlrworkflow *HLRWorkflow) ExecuteOrder(order_id int64) {
	log.Info("Starting order %v execution at %v", order_id, time.Now())
	log.Info("Checking order type ...")
	dtoorder, err := hlrworkflow.OrderRepository.Get(order_id)
	if err != nil {
		return
	}
	dtofacility, err := hlrworkflow.FacilityRepository.Get(dtoorder.Facility_ID)
	if err != nil {
		return
	}
	if dtofacility.Alias != models.SERVICE_TYPE_HLR {
		log.Error("Order service is not macthed to the service method %v", dtoorder.Facility_ID)
		return
	}
	if !dtofacility.Active {
		log.Error("Service is not active %v", dtofacility.ID)
		return
	}
	log.Info("Checking service type ...")
	dtohlrfacility, err := hlrworkflow.HLRFacilityRepository.Get(dtoorder.ID)
	if err != nil {
		return
	}
	log.Info("Checking order status ...")
	dtoorderstatuses, err := hlrworkflow.OrderStatusRepository.GetByOrder(dtoorder.ID)
	if err != nil {
		return
	}

	order := models.NewApiLongOrderFromDto(dtoorder, dtoorderstatuses)
	if order.IsAssembled && order.IsConfirmed && !order.IsOpen && !order.IsCancelled && !order.IsExecuted && !order.IsArchived && !order.IsDeleted {
		/* 1 */
		log.Info("Checking data table ...")
		dtodatatable, err := hlrworkflow.CustomerTableRepository.Get(dtohlrfacility.DeliveryDataId)
		if err != nil {
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		if !dtodatatable.Active {
			log.Error("Data table is not active %v", dtohlrfacility.DeliveryDataId)
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		if !dtodatatable.Permanent {
			log.Error("Data table is not permanent %v", dtohlrfacility.DeliveryDataId)
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		/* 2 */ err = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_OPEN, true)
		if err != nil {
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Checking order data ...")
		columnmobilephone_id, err := hlrworkflow.CheckHLROrder(dtohlrfacility)
		if err != nil {
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}

		/* 3 */ err = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_SUPPLIER_COST_NEW, true)
		if err != nil {
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Checking table columns ...")
		columnmobilephone, tablecolumns, err := hlrworkflow.GetHLRTableColumns(dtohlrfacility, columnmobilephone_id)
		if err != nil {
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Checking table data ...")
		apitablerows, err := hlrworkflow.TableRowRepository.GetAll("", "", dtohlrfacility.DeliveryDataId, tablecolumns)
		if err != nil {
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Calculating order cost ...")
		dtohlrfacility.Cost, err = hlrworkflow.CalculateCost(apitablerows, columnmobilephone, dtoorder, dtohlrfacility)
		if err != nil {
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		dtohlrfacility.CostFactual = dtohlrfacility.Cost
		err = hlrworkflow.HLRFacilityRepository.Update(dtohlrfacility, true, false)
		if err != nil {
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}

		/* 4 */ err = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED, true)
		if err != nil {
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Start order processing ...")
		/* 5 */ dtoorder.Begin_Date = time.Now()
		err = hlrworkflow.OrderRepository.Update(dtoorder, &[]models.DtoOrderStatus{
			*models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_MODERATOR_BEGIN, true, "", time.Now())}, nil, true)
		if err != nil {
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Calculating unit balance ...")
		balance, err := hlrworkflow.OperationRepository.CalculateBalance(dtoorder.Unit_ID)
		if err != nil {
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Order payment and invoicement ...")
		/* 6 */ err = hlrworkflow.PayAndInvoice(dtoorder, dtohlrfacility, balance)
		if err != nil {
			_ = hlrworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Copying table data ...")
		/* 7 */ dtoworkdatatable, err := hlrworkflow.CopyData(dtoorder, dtodatatable)
		if err != nil {
			return
		}
		log.Info("Sending data to supplier ...")
		/* 8 */ hlrresponse, err := hlrworkflow.SendHLR(apitablerows, dtoorder, dtohlrfacility, columnmobilephone)
		if err != nil {
			return
		}
		log.Info("Saving supplier response ...")
		/* 9 */ err = hlrworkflow.SaveHLR(dtoworkdatatable, hlrresponse)
		if err != nil {
			return
		}
		log.Info("Getting supplier results ...")
		/* 10 */ hlrstatuses, err := GetHLRStatus(hlrresponse)
		if err != nil {
			return
		}
		log.Info("Saving supplier results ...")
		/* 11 */ err = hlrworkflow.SaveHLRStatus(dtoworkdatatable, hlrstatuses)
		if err != nil {
			return
		}
		log.Info("Finsing order processing and clearing data ...")
		/* 12 */ err = hlrworkflow.ClearTables(dtoorder, dtohlrfacility, dtodatatable)
		if err != nil {
			return
		}
		log.Info("Completing order execution %v", time.Now())
		/* 13 */ dtoorder.End_Date = time.Now()
		err = hlrworkflow.OrderRepository.Update(dtoorder, &[]models.DtoOrderStatus{
			*models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_SUPPLIER_CLOSE, true, "", time.Now())}, nil, true)
		if err != nil {
			return
		}
	}
}
