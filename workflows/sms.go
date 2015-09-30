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
	"math"
	"time"
	"unicode"
)

const (
	SMS_LENGTH_ASCII  = 153
	SMS_LENGTH_UICODE = 67

	COLUMN_NAME_SMS_ID           = "SmsId"
	COLUMN_NAME_SMS_ERROR        = "SmsError"
	COLUMN_NAME_SMS_STATUS_ID    = "SmsStatusId"
	COLUMN_NAME_SMS_STATUS_ERROR = "SmsStatusError"
)

type SMSWorkflow struct {
	OrderRepository           services.OrderRepository
	FacilityRepository        services.FacilityRepository
	SMSFacilityRepository     services.SMSFacilityRepository
	OrderStatusRepository     services.OrderStatusRepository
	CustomerTableRepository   services.CustomerTableRepository
	SMSTableRepository        services.SMSTableRepository
	SMSSenderRepository       services.SMSSenderRepository
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

func NewSMSWorkflow(orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	smsfacilityrepository services.SMSFacilityRepository, orderstatusrepository services.OrderStatusRepository,
	customertablerepository services.CustomerTableRepository, smstablerepository services.SMSTableRepository,
	smssenderrepository services.SMSSenderRepository, resulttablerepository services.ResultTableRepository,
	worktablerepository services.WorkTableRepository, invoicerepository services.InvoiceRepository,
	companyrepository services.CompanyRepository, operationrepository services.OperationRepository,
	transactiontyperepository services.TransactionTypeRepository, tablecolumnrepository services.TableColumnRepository,
	unitrepository services.UnitRepository, tablerowrepository services.TableRowRepository,
	pricerepository services.PriceRepository, mobileoperatorrepository services.MobileOperatorRepository,
	columntyperepository services.ColumnTypeRepository) *SMSWorkflow {
	return &SMSWorkflow{
		OrderRepository:           orderrepository,
		FacilityRepository:        facilityrepository,
		SMSFacilityRepository:     smsfacilityrepository,
		OrderStatusRepository:     orderstatusrepository,
		CustomerTableRepository:   customertablerepository,
		SMSTableRepository:        smstablerepository,
		SMSSenderRepository:       smssenderrepository,
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

func SendSMS(supplier *libTypes.Supplier, sms *[]libTypes.Sms) (smsresponse *libTypes.SmsResponse, err error) {
	// Отправка SMS
	suppliers.WaitClientReady()                         // Функция блокируется до момента пока связь с сервером не будет установлена
	response, err := suppliers.SendSms(*supplier, *sms) // Функция не блокируется. Если связи с сервером нет, то вернётся ошибка
	if err != nil {
		log.Error("Can't send SMS messages %v for supplier %v", err, supplier.Name)
		return nil, err
	}

	return &response, nil
}

func GetSMSStatus(smsresponse *libTypes.SmsResponse) (smsstatuses map[gocql.UUID]libTypes.SmsStatus, err error) {
	smsstatuses = make(map[gocql.UUID]libTypes.SmsStatus)
	// Получение статуса по UUID запроса
	for {
		final := true
		for i := range smsresponse.Ids {
			var status libTypes.SmsStatus
			status, err = suppliers.StatusSms(smsresponse.Ids[i])
			if err != nil {
				log.Error("Can't get SMS statuses %v", err)
				return map[gocql.UUID]libTypes.SmsStatus{}, err
			}
			if status.Final {
				smsstatuses[status.Id] = status
			} else {
				final = false
			}
		}
		if final {
			break
		}
		time.Sleep(time.Second)
	}

	return smsstatuses, nil
}

func CalculateSMSQuantity(sms string) (count int) {
	count = 0
	runes := []rune(sms)
	isASCII := true
	for _, current_rune := range runes {
		if current_rune > unicode.MaxASCII {
			isASCII = false
			break
		}
	}
	length := 0
	if isASCII {
		length = SMS_LENGTH_ASCII
		count = len(runes) / length

	} else {
		length = SMS_LENGTH_UICODE
		count = len(runes) / length
	}
	if math.Mod(float64(len(runes)), float64(length)) != 0 || sms == "" {
		count++
	}

	return count
}

func (smsworkflow *SMSWorkflow) CheckSMSOrder(dtoorder *models.DtoOrder, dtosmsfacility *models.DtoSMSFacility) (columnmobilephone_id int64, err error) {
	smstable, err := smsworkflow.SMSTableRepository.Get(dtosmsfacility.DeliveryDataId)
	if err != nil {
		return 0, err
	}
	if dtosmsfacility.MessageFromInColumnId != 0 {
		found := false
		for _, smssender := range smstable.SMSSenders {
			if smssender.ID == dtosmsfacility.MessageFromInColumnId {
				found = true
				break
			}
		}
		if !found {
			log.Error("Can't find sms sender column %v in table %v", dtosmsfacility.MessageFromInColumnId, dtosmsfacility.DeliveryDataId)
			return 0, errors.New("Wrong sms sender column")
		}
		dtotablecolumn, err := smsworkflow.TableColumnRepository.Get(dtosmsfacility.MessageFromInColumnId)
		if err != nil {
			return 0, err
		}
		found, err = smsworkflow.SMSSenderRepository.Belongs(dtotablecolumn, dtoorder.Unit_ID, dtoorder.Supplier_ID)
		if err != nil {
			return 0, err
		}
		if !found {
			log.Error("SMS sender column %v data doesn't belong to unit %v", dtotablecolumn.ID, dtoorder.Unit_ID)
			return 0, errors.New("SMS sender column not for unit")
		}
	} else {
		dtosmsender, err := smsworkflow.SMSSenderRepository.Get(dtosmsfacility.MessageFromId)
		if err != nil {
			return 0, err
		}
		if dtosmsender.Unit_ID != dtoorder.Unit_ID {
			log.Error("SMS sender %v unit doesn't match order %v unit", dtosmsender.ID, dtoorder.ID)
			return 0, errors.New("Wrong sms sender unit")
		}
		if !dtosmsender.Active {
			log.Error("SMS sender is not active %v", dtosmsender.ID)
			return 0, errors.New("SMS sender not active")
		}
		if !dtosmsender.Registered {
			log.Error("SMS sender is not registered %v", dtosmsender.ID)
			return 0, errors.New("SMS sender not registered")
		}
	}

	if dtosmsfacility.MessageToInColumnId != 0 {
		found := false
		for _, mobilephone := range smstable.MobilePhones {
			if mobilephone.ID == dtosmsfacility.MessageToInColumnId {
				found = true
				columnmobilephone_id = dtosmsfacility.MessageToInColumnId
				break
			}
		}
		if !found {
			log.Error("Can't find mobile phone column %v in table %v", dtosmsfacility.MessageToInColumnId, dtosmsfacility.DeliveryDataId)
			return 0, errors.New("Wrong mobile phone column")
		}
	} else {
		if len(smstable.MobilePhones) == 0 {
			log.Error("Can't find mobile column in table %v", dtosmsfacility.DeliveryDataId)
			return 0, errors.New("Missed mobile phone column")
		} else {
			columnmobilephone_id = smstable.MobilePhones[0].ID
		}
	}

	if dtosmsfacility.MessageBodyInColumnId != 0 {
		found := false
		for _, message := range smstable.Messages {
			if message.ID == dtosmsfacility.MessageBodyInColumnId {
				found = true
				break
			}
		}
		if !found {
			log.Error("Can't find message column %v in table %v", dtosmsfacility.MessageBodyInColumnId, dtosmsfacility.DeliveryDataId)
			return 0, errors.New("Wrong message column")
		}
	}

	return columnmobilephone_id, nil
}

func (smsworkflow *SMSWorkflow) GetSMSTableColumns(dtosmsfacility *models.DtoSMSFacility, columnmobilephone_id int64) (
	columnmessage, columnmobilephone, columnsmssender *models.DtoTableColumn, tablecolumns *[]models.DtoTableColumn, err error) {
	tablecolumns = new([]models.DtoTableColumn)
	if dtosmsfacility.MessageBodyInColumnId != 0 {
		columnmessage, err = smsworkflow.TableColumnRepository.Get(dtosmsfacility.MessageBodyInColumnId)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		*tablecolumns = append(*tablecolumns, *columnmessage)
	}
	if dtosmsfacility.MessageFromInColumnId != 0 {
		columnsmssender, err = smsworkflow.TableColumnRepository.Get(dtosmsfacility.MessageFromInColumnId)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		*tablecolumns = append(*tablecolumns, *columnsmssender)
	}
	columnmobilephone, err = smsworkflow.TableColumnRepository.Get(columnmobilephone_id)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	*tablecolumns = append(*tablecolumns, *columnmobilephone)

	return columnmessage, columnmobilephone, columnsmssender, tablecolumns, nil
}

func (smsworkflow *SMSWorkflow) CalculateCost(apitablerows *[]models.ApiInfoTableRow, columnmessage, columnmobilephone, columnsender *models.DtoTableColumn,
	dtoorder *models.DtoOrder, dtosmsfacility *models.DtoSMSFacility) (cost float64, err error) {
	dtomobileoperators, err := smsworkflow.MobileOperatorRepository.FindAll()
	if err != nil {
		return 0, err
	}
	mobileoperators_uuid := make(map[string]*models.DtoMobileOperator)
	mobileoperators_id := make(map[int]*models.DtoMobileOperator)
	for index := range *dtomobileoperators {
		mobileoperators_uuid[(*dtomobileoperators)[index].UUID] = &(*dtomobileoperators)[index]
		mobileoperators_id[(*dtomobileoperators)[index].ID] = &(*dtomobileoperators)[index]
	}
	dtocolumntype, err := smsworkflow.ColumnTypeRepository.Get(columnmobilephone.Column_Type_ID)
	if err != nil {
		return 0, err
	}
	unitedsmscount := 0
	if columnmessage == nil {
		unitedsmscount = CalculateSMSQuantity(dtosmsfacility.MessageBody)
	}
	unitedsmssender := ""
	if columnsender == nil {
		dtosmsender, err := smsworkflow.SMSSenderRepository.Get(dtosmsfacility.MessageFromId)
		if err != nil {
			return 0, err
		}
		unitedsmssender = dtosmsender.Name
	}

	mobilephones := []uint64{}
	mobilesmses := []int{}
	smssenders := []string{}
	for _, apitablerow := range *apitablerows {
		var smscount int = 0
		var mobilephone uint64 = 0
		var smssender = ""
		if columnmessage == nil {
			smscount = unitedsmscount
		}
		if columnsender == nil {
			smssender = unitedsmssender
		}
		for _, apitablecell := range apitablerow.Cells {
			if columnmessage != nil {
				if apitablecell.Table_Column_ID == columnmessage.ID {
					smscount = CalculateSMSQuantity(apitablecell.Value)
				}
			}
			if apitablecell.Table_Column_ID == columnmobilephone.ID {
				mobilephone, err = helpers.CheckMobilePhone(apitablecell.Value, dtocolumntype, smsworkflow.ColumnTypeRepository)
				if err != nil {
					return 0, err
				}
			}
			if columnsender != nil {
				if apitablecell.Table_Column_ID == columnsender.ID {
					smssender = apitablecell.Value
				}
			}
		}
		mobilephones = append(mobilephones, mobilephone)
		mobilesmses = append(mobilesmses, smscount)
		smssenders = append(smssenders, smssender)
	}

	mobileoperatoruuids, err := suppliers.MobileOperator(mobilephones)
	if err != nil {
		log.Error("Error during detecting mobile operators %v", err)
		return 0, err
	}
	defaultmobileoperator, err := smsworkflow.MobileOperatorRepository.GetDefault()
	if err != nil {
		return 0, err
	}
	mobileoperatorsmses := make(map[int]map[string]int)
	for index := range mobileoperatoruuids {
		if mobileoperatoruuids[index].Id.String() == models.MOBILE_OPERATOR_UUID_UNKNOWN {
			log.Error("Not existed mobile operator for phone %v", mobilephones[index])
			return 0, errors.New("Not existed mobile operator")
		}
		mobileoperator, ok := mobileoperators_uuid[mobileoperatoruuids[index].Id.String()]
		if !ok {
			mobileoperator = defaultmobileoperator
		}
		mobileoperatorsenders, ok := mobileoperatorsmses[mobileoperator.ID]
		if !ok {
			mobileoperatorsenders = make(map[string]int)
			mobileoperatorsmses[mobileoperator.ID] = mobileoperatorsenders
		}
		mobileoperatorsenders[smssenders[index]] += mobilesmses[index]
	}

	smshlrprices, err := helpers.GetSMSHLRPrices(models.SERVICE_TYPE_SMS, dtoorder.Supplier_ID, nil, smsworkflow.PriceRepository,
		smsworkflow.TableColumnRepository, smsworkflow.TableRowRepository,
		smsworkflow.MobileOperatorRepository, "", true)
	if err != nil {
		return 0, err
	}

	pricemobileoperators := make(map[int]int)
	for _, smshlrprice := range *smshlrprices {
		pricemobileoperators[smshlrprice.Mobile_Operator_ID] = smshlrprice.Mobile_Operator_ID
	}

	cost = 0
	for mobileoperator_id, mobileoperatorsenders := range mobileoperatorsmses {
		pricemobileoperator_id, ok := pricemobileoperators[mobileoperator_id]
		if !ok {
			pricemobileoperator_id, ok = pricemobileoperators[defaultmobileoperator.ID]
			if !ok {
				log.Error("Can't find default mobile operator %v in price list for supplier %v", defaultmobileoperator.ID, dtoorder.Supplier_ID)
				return 0, errors.New("Missed default mobile operator in price list")
			}
		}

		foundprices := make(map[string]bool)
		for smssender := range mobileoperatorsenders {
			foundprices[smssender] = false
		}
		for _, smshlrprice := range *smshlrprices {
			if smshlrprice.Mobile_Operator_ID == pricemobileoperator_id {
				if mobileoperators_id[pricemobileoperator_id].SMSBillingModel == models.BILLING_MODEL_RANGE {
					count := 0
					for _, smscount := range mobileoperatorsenders {
						count += smscount
					}
					if count >= smshlrprice.AmountRange.Begin && (count <= smshlrprice.AmountRange.End || smshlrprice.AmountRange.End == 0) {
						cost += float64(count) * smshlrprice.Price
						for smssender := range foundprices {
							foundprices[smssender] = true
						}
						break
					}
				}
				if mobileoperators_id[pricemobileoperator_id].SMSBillingModel == models.BILLING_MODEL_CUMULATIVE_PERHEADER {
					for smssender, count := range mobileoperatorsenders {
						if count >= smshlrprice.AmountRange.Begin && (count > smshlrprice.AmountRange.End && smshlrprice.AmountRange.End != 0) {
							cost += float64(smshlrprice.AmountRange.End-smshlrprice.AmountRange.Begin+1) * smshlrprice.Price
						}
						if count >= smshlrprice.AmountRange.Begin && (count <= smshlrprice.AmountRange.End || smshlrprice.AmountRange.End == 0) {
							cost += float64(count-smshlrprice.AmountRange.Begin+1) * smshlrprice.Price
							foundprices[smssender] = true
						}
					}
				}
			}
		}

		for _, found := range foundprices {
			if !found {
				log.Error("Can't find  mobile operator %v in price list for supplier %v", pricemobileoperator_id, dtoorder.Supplier_ID)
				return 0, errors.New("Missed obile operator in price list")
			}
		}
	}

	return cost, nil
}

func (smsworkflow *SMSWorkflow) PayAndInvoice(dtoorder *models.DtoOrder, dtosmsfacility *models.DtoSMSFacility, balance float64) (err error) {
	if balance < dtosmsfacility.Cost {
		log.Error("Not enough money at unit balance %v to pay for order %v execution", dtoorder.Unit_ID, dtoorder.ID)
		return errors.New("Not enough money")
	}
	dtocompany, err := smsworkflow.CompanyRepository.GetPrimaryByUnit(dtoorder.Unit_ID)
	if err != nil {
		return
	}
	dtounit, err := smsworkflow.UnitRepository.Get(config.Configuration.SystemAccount)
	if err != nil {
		return err
	}
	dtotransactiontype, err := smsworkflow.TransactionTypeRepository.Get(models.TRANSACTION_TYPE_SERVICE_FEE_SMS)
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
		dtoinvoice.VAT = (dtosmsfacility.Cost / (1 + float64(dtocompany.VAT)/100)) * float64(dtocompany.VAT) / 100
	}
	dtoinvoice.Total = dtosmsfacility.Cost
	dtoinvoice.Paid = true
	dtoinvoice.Created = time.Now()
	dtoinvoice.Active = true
	dtoinvoice.InvoiceItems = []models.DtoInvoiceItem{*models.NewDtoInvoiceItem(0, 0, dtotransactiontype.Name, models.INVOICE_ITEM_TYPE_ROUBLE,
		1, dtoinvoice.Total, dtoinvoice.Total)}
	dtoinvoice.PaidAt = time.Now()

	err = smsworkflow.InvoiceRepository.PayForOrder(dtoorder, dtoinvoice, dtotransaction, true)
	if err != nil {
		return err
	}

	return nil
}

func (smsworkflow *SMSWorkflow) SendSMS(apitablerows *[]models.ApiInfoTableRow, dtoorder *models.DtoOrder, dtosmsfacility *models.DtoSMSFacility,
	columnsmssender, columnmessage, columnmobilephone *models.DtoTableColumn) (smsresponse *libTypes.SmsResponse, err error) {
	dtosupplier, err := smsworkflow.UnitRepository.Get(dtoorder.Supplier_ID)
	if err != nil {
		return nil, err
	}
	smssupplier, err := GetSupplier(dtosupplier.UUID)
	if err != nil {
		return nil, err
	}
	smssender := ""
	if columnsmssender == nil {
		dtosmsender, err := smsworkflow.SMSSenderRepository.Get(dtosmsfacility.MessageFromId)
		if err != nil {
			return nil, err
		}
		smssender = dtosmsender.Name
	}
	dtocolumntype, err := smsworkflow.ColumnTypeRepository.Get(columnmobilephone.Column_Type_ID)
	if err != nil {
		return nil, err
	}

	sms := new([]libTypes.Sms)
	for _, apitablerow := range *apitablerows {
		obj := new(libTypes.Sms)
		obj.Flash = false
		if columnmessage == nil {
			obj.Message = []byte(dtosmsfacility.MessageBody)
		}
		if columnsmssender == nil {
			obj.Sender = smssender
		}
		for _, apitablecell := range apitablerow.Cells {
			if columnmessage != nil {
				if apitablecell.Table_Column_ID == columnmessage.ID {
					obj.Message = []byte(apitablecell.Value)
				}
			}
			if apitablecell.Table_Column_ID == columnmobilephone.ID {
				obj.Recipient, err = helpers.CheckMobilePhone(apitablecell.Value, dtocolumntype, smsworkflow.ColumnTypeRepository)
				if err != nil {
					return nil, err
				}
			}
			if columnsmssender != nil {
				if apitablecell.Table_Column_ID == columnsmssender.ID {
					obj.Sender = apitablecell.Value
				}
			}
		}
		*sms = append(*sms, *obj)
	}
	smsresponse, err = SendSMS(smssupplier, sms)
	if err != nil {
		return nil, err
	}

	return smsresponse, nil
}

func (smsworkflow *SMSWorkflow) CopyData(dtoorder *models.DtoOrder,
	dtodatatable *models.DtoCustomerTable) (dtoworkdatatable *models.DtoCustomerTable, err error) {
	dtoworkdatatable, err = smsworkflow.CustomerTableRepository.Copy(dtodatatable, true)
	if err != nil {
		return nil, err
	}
	dtodatatable.TypeID = models.TABLE_TYPE_READONLY
	err = smsworkflow.CustomerTableRepository.Update(dtodatatable)
	if err != nil {
		return nil, err
	}
	dtoworkdatatable.TypeID = models.TABLE_TYPE_HIDDEN
	err = smsworkflow.CustomerTableRepository.Update(dtoworkdatatable)
	if err != nil {
		return nil, err
	}

	dtoworktable := models.NewDtoWorkTable(dtoorder.ID, dtoworkdatatable.ID)
	err = smsworkflow.WorkTableRepository.Create(dtoworktable, nil)
	if err != nil {
		return nil, err
	}

	return dtoworkdatatable, nil
}

func (smsworkflow *SMSWorkflow) SaveSMS(dtoworkdatatable *models.DtoCustomerTable, smsresponse *libTypes.SmsResponse) (err error) {
	position, err := smsworkflow.TableColumnRepository.GetDefaultPosition(dtoworkdatatable.ID)
	if err != nil {
		return err
	}
	tablecolumns := []models.DtoTableColumn{}
	columnnames := []string{COLUMN_NAME_SMS_ID, COLUMN_NAME_SMS_ERROR}
	for _, columnname := range columnnames {
		position++
		fieldnum, err := helpers.FindFreeColumnInternal(dtoworkdatatable.ID, 0, smsworkflow.TableColumnRepository)
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

		err = smsworkflow.TableColumnRepository.Create(dtotablecolumn, nil)
		if err != nil {
			return err
		}
		tablecolumns = append(tablecolumns, *dtotablecolumn)
	}

	apitablerows, err := smsworkflow.TableRowRepository.GetAll("", "", dtoworkdatatable.ID, &tablecolumns)
	if err != nil {
		return err
	}

	workdatatablecolumns, err := smsworkflow.TableColumnRepository.GetByTable(dtoworkdatatable.ID)
	if err != nil {
		return err
	}

	for index, apitablerow := range *apitablerows {
		tablerow, err := smsworkflow.TableRowRepository.Get(apitablerow.ID)
		if err != nil {
			return err
		}

		tablecells, err := tablerow.TableRowToDtoTableCells(workdatatablecolumns)
		if err != nil {
			return err
		}

		for i, _ := range *tablecells {
			for j, _ := range tablecolumns {
				FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_SMS_ID, smsresponse.Ids[index].String())
				if smsresponse.Errors[index] != nil {
					FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_SMS_ERROR, smsresponse.Errors[index].Error())
				}
			}
		}

		err = tablerow.TableCellsToTableRow(tablecells, workdatatablecolumns)
		if err != nil {
			return err
		}

		err = smsworkflow.TableRowRepository.Update(tablerow, nil, true, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (smsworkflow *SMSWorkflow) SaveSMSStatus(dtoworkdatatable *models.DtoCustomerTable, smsstatuses map[gocql.UUID]libTypes.SmsStatus) (err error) {
	var columnsmsid *models.DtoTableColumn
	alltablecolumns, err := smsworkflow.TableColumnRepository.GetByTable(dtoworkdatatable.ID)
	if err != nil {
		return err
	}
	found := false
	for i := range *alltablecolumns {
		if (*alltablecolumns)[i].Name == COLUMN_NAME_SMS_ID {
			columnsmsid = &(*alltablecolumns)[i]
			found = true
			break
		}
	}
	if !found {
		log.Error("Can't find sms id column for table %v", dtoworkdatatable.ID)
		return errors.New("Missed sms id column")
	}

	position, err := smsworkflow.TableColumnRepository.GetDefaultPosition(dtoworkdatatable.ID)
	if err != nil {
		return err
	}
	tablecolumns := []models.DtoTableColumn{*columnsmsid}
	columnnames := []string{COLUMN_NAME_SMS_STATUS_ID, COLUMN_NAME_SMS_STATUS_ERROR}
	for _, columnname := range columnnames {
		position++
		fieldnum, err := helpers.FindFreeColumnInternal(dtoworkdatatable.ID, 0, smsworkflow.TableColumnRepository)
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

		err = smsworkflow.TableColumnRepository.Create(dtotablecolumn, nil)
		if err != nil {
			return err
		}
		tablecolumns = append(tablecolumns, *dtotablecolumn)
	}

	apitablerows, err := smsworkflow.TableRowRepository.GetAll("", "", dtoworkdatatable.ID, &tablecolumns)
	if err != nil {
		return err
	}

	workdatatablecolumns, err := smsworkflow.TableColumnRepository.GetByTable(dtoworkdatatable.ID)
	if err != nil {
		return err
	}

	for _, apitablerow := range *apitablerows {
		tablerow, err := smsworkflow.TableRowRepository.Get(apitablerow.ID)
		if err != nil {
			return err
		}

		tablecells, err := tablerow.TableRowToDtoTableCells(workdatatablecolumns)
		if err != nil {
			return err
		}

		var smsid gocql.UUID
		found := false
		for i, _ := range *tablecells {
			if (*tablecells)[i].Table_Column_ID == columnsmsid.ID {
				smsid, err = gocql.ParseUUID((*tablecells)[i].Value)
				if err != nil {
					log.Error("Can't parse sms id from column %v, %v", err, columnsmsid.ID)
					return err
				}
				found = true
				break
			}
		}
		if !found {
			log.Error("Can't find sms id value for column %v", columnsmsid.ID)
			return errors.New("Missed sms id value")
		}

		for i, _ := range *tablecells {
			for j, _ := range tablecolumns {
				FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_SMS_STATUS_ID, fmt.Sprintf("%v", smsstatuses[smsid].Status))
				FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_SMS_STATUS_ERROR, fmt.Sprintf("%v", smsstatuses[smsid].Error))
			}
		}

		err = tablerow.TableCellsToTableRow(tablecells, workdatatablecolumns)
		if err != nil {
			return err
		}

		err = smsworkflow.TableRowRepository.Update(tablerow, nil, true, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (smsworkflow *SMSWorkflow) ClearTables(dtoorder *models.DtoOrder, dtosmsfacility *models.DtoSMSFacility,
	dtodatatable *models.DtoCustomerTable) (err error) {
	if dtosmsfacility.DeliveryDataDelete {
		err = smsworkflow.CustomerTableRepository.Deactivate(dtodatatable)
		if err != nil {
			return err
		}
	} else {
		dtodatatable.TypeID = models.TABLE_TYPE_DEFAULT
		err = smsworkflow.CustomerTableRepository.Update(dtodatatable)
		if err != nil {
			return err
		}
	}

	worktables, err := smsworkflow.WorkTableRepository.GetByOrder(dtoorder.ID)
	if err != nil {
		return err
	}
	for _, worktable := range *worktables {
		dtocustomertable, err := smsworkflow.CustomerTableRepository.Get(worktable.Customer_Table_ID)
		if err != nil {
			return err
		}
		err = smsworkflow.CustomerTableRepository.Deactivate(dtocustomertable)
		if err != nil {
			return err
		}
	}

	return nil
}

func (smsworkflow *SMSWorkflow) SetStatus(dtoorder *models.DtoOrder, orderstatus models.OrderStatus, active bool) (err error) {
	dtoorderstatus := models.NewDtoOrderStatus(dtoorder.ID, orderstatus, active, "", time.Now())
	err = smsworkflow.OrderStatusRepository.Save(dtoorderstatus, nil)
	if err != nil {
		return err
	}

	return nil
}

func (smsworkflow *SMSWorkflow) ExecuteOrder(order_id int64) {
	log.Info("Starting order %v execution at %v", order_id, time.Now())
	log.Info("Checking order type ...")
	dtoorder, err := smsworkflow.OrderRepository.Get(order_id)
	if err != nil {
		return
	}
	dtofacility, err := smsworkflow.FacilityRepository.Get(dtoorder.Facility_ID)
	if err != nil {
		return
	}
	if dtofacility.Alias != models.SERVICE_TYPE_SMS {
		log.Error("Order service is not macthed to the service method %v", dtoorder.Facility_ID)
		return
	}
	if !dtofacility.Active {
		log.Error("Service is not active %v", dtofacility.ID)
		return
	}
	log.Info("Checking service type ...")
	dtosmsfacility, err := smsworkflow.SMSFacilityRepository.Get(dtoorder.ID)
	if err != nil {
		return
	}
	log.Info("Checking order status ...")
	dtoorderstatuses, err := smsworkflow.OrderStatusRepository.GetByOrder(dtoorder.ID)
	if err != nil {
		return
	}

	order := models.NewApiLongOrderFromDto(dtoorder, dtoorderstatuses)
	if order.IsAssembled && order.IsConfirmed && !order.IsOpen && !order.IsCancelled && !order.IsExecuted && !order.IsArchived && !order.IsDeleted {
		/* 1 */
		if (dtosmsfacility.DeliveryType == models.TYPE_DELIVERY_ONCE && !dtosmsfacility.DeliveryTime) ||
			(dtosmsfacility.DeliveryTime && ((!dtosmsfacility.DeliveryTimeStart.IsZero() && dtosmsfacility.DeliveryTimeStart.Sub(time.Now()) <= 0) &&
				(!dtosmsfacility.DeliveryTimeEnd.IsZero() && dtosmsfacility.DeliveryTimeEnd.Sub(time.Now()) > 0))) {
			schhour, schmin, _ := dtosmsfacility.DeliveryBaseTime.Local().Clock()
			hour, min, _ := time.Now().Clock()
			if schhour == hour && schmin == min {
				log.Info("Checking data table ...")
				dtodatatable, err := smsworkflow.CustomerTableRepository.Get(dtosmsfacility.DeliveryDataId)
				if err != nil {
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}
				if !dtodatatable.Active {
					log.Error("Data table is not active %v", dtosmsfacility.DeliveryDataId)
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}
				if !dtodatatable.Permanent {
					log.Error("Data table is not permanent %v", dtosmsfacility.DeliveryDataId)
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}
				/* 2 */ err = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_OPEN, true)
				if err != nil {
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}
				log.Info("Checking order data ...")
				columnmobilephone_id, err := smsworkflow.CheckSMSOrder(dtoorder, dtosmsfacility)
				if err != nil {
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}

				/* 3 */ err = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_SUPPLIER_COST_NEW, true)
				if err != nil {
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}
				log.Info("Checking table columns ...")
				columnmessage, columnmobilephone, columnsmssender, tablecolumns, err := smsworkflow.GetSMSTableColumns(dtosmsfacility, columnmobilephone_id)
				if err != nil {
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}
				log.Info("Checking table data ...")
				apitablerows, err := smsworkflow.TableRowRepository.GetAll("", "", dtosmsfacility.DeliveryDataId, tablecolumns)
				if err != nil {
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}
				log.Info("Calculating order cost ...")
				dtosmsfacility.Cost, err = smsworkflow.CalculateCost(apitablerows, columnmessage, columnmobilephone, columnsmssender, dtoorder, dtosmsfacility)
				if err != nil {
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}
				dtosmsfacility.CostFactual = dtosmsfacility.Cost
				err = smsworkflow.SMSFacilityRepository.Update(dtosmsfacility, true, false)
				if err != nil {
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}

				/* 4 */ err = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED, true)
				if err != nil {
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}
				log.Info("Start order processing ...")
				/* 5 */ dtoorder.Begin_Date = time.Now()
				err = smsworkflow.OrderRepository.Update(dtoorder, &[]models.DtoOrderStatus{
					*models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_MODERATOR_BEGIN, true, "", time.Now())}, nil, true)
				if err != nil {
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}
				log.Info("Calculating unit balance ...")
				balance, err := smsworkflow.OperationRepository.CalculateBalance(dtoorder.Unit_ID)
				if err != nil {
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}
				log.Info("Order payment and invoicement ...")
				/* 6 */ err = smsworkflow.PayAndInvoice(dtoorder, dtosmsfacility, balance)
				if err != nil {
					_ = smsworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
					return
				}
				return
				log.Info("Copying table data ...")
				/* 7 */ dtoworkdatatable, err := smsworkflow.CopyData(dtoorder, dtodatatable)
				if err != nil {
					return
				}
				log.Info("Sending data to supplier ...")
				/* 8 */ smsresponse, err := smsworkflow.SendSMS(apitablerows, dtoorder, dtosmsfacility, columnsmssender, columnmessage, columnmobilephone)
				if err != nil {
					return
				}
				log.Info("Saving supplier response ...")
				/* 9 */ err = smsworkflow.SaveSMS(dtoworkdatatable, smsresponse)
				if err != nil {
					return
				}
				log.Info("Getting supplier results ...")
				/* 10 */ smsstatuses, err := GetSMSStatus(smsresponse)
				if err != nil {
					return
				}
				log.Info("Saving supplier results ...")
				/* 11 */ err = smsworkflow.SaveSMSStatus(dtoworkdatatable, smsstatuses)
				if err != nil {
					return
				}
				log.Info("Finsing order processing and clearing data ...")
				/* 12 */ err = smsworkflow.ClearTables(dtoorder, dtosmsfacility, dtodatatable)
				if err != nil {
					return
				}
				log.Info("Completing order execution %v", time.Now())
				/* 13 */ dtoorder.End_Date = time.Now()
				err = smsworkflow.OrderRepository.Update(dtoorder, &[]models.DtoOrderStatus{
					*models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_SUPPLIER_CLOSE, true, "", time.Now())}, nil, true)
				if err != nil {
					return
				}

			}
		}
	}
}
