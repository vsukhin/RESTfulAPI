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
	COLUMN_NAME_VERIFY_ID                             = "DataId"
	COLUMN_NAME_VERIFY_ERROR                          = "DataError"
	COLUMN_NAME_VERIFY_STATUS_ID                      = "DataStatusId"
	COLUMN_NAME_VERIFY_STATUS_ERROR                   = "DataStatusError"
	COLUMN_NAME_VERIFY_POSTADDRESS_RESULT             = "DataPostAddressResult"
	COLUMN_NAME_VERIFY_POSTADDRESS_POSTALCODE         = "DataPostAddressPostalCode"
	COLUMN_NAME_VERIFY_POSTADDRESS_COUNTRY            = "DataPostAddressCountry"
	COLUMN_NAME_VERIFY_POSTADDRESS_REGIONTYPE         = "DataPostAddressRegionType"
	COLUMN_NAME_VERIFY_POSTADDRESS_REGIONTYPEFULL     = "DataPostAddressRegionTypeFull"
	COLUMN_NAME_VERIFY_POSTADDRESS_REGION             = "DataPostAddressRegion"
	COLUMN_NAME_VERIFY_POSTADDRESS_AREATYPE           = "DataPostAddressAreaType"
	COLUMN_NAME_VERIFY_POSTADDRESS_AREATYPEFULL       = "DataPostAddressAreTypeFull"
	COLUMN_NAME_VERIFY_POSTADDRESS_AREA               = "DataPostAddressArea"
	COLUMN_NAME_VERIFY_POSTADDRESS_CITYTYPE           = "DataPostAddressCityType"
	COLUMN_NAME_VERIFY_POSTADDRESS_CITYTYPEFULL       = "DataPostAddressCityTypeFull"
	COLUMN_NAME_VERIFY_POSTADDRESS_CITY               = "DataPostAddressCity"
	COLUMN_NAME_VERIFY_POSTADDRESS_SETTLEMENTTYPE     = "DataPostAddressSettlementType"
	COLUMN_NAME_VERIFY_POSTADDRESS_SETTLEMENTTYPEFULL = "DataPostAddressSettlementTypeFull"
	COLUMN_NAME_VERIFY_POSTADDRESS_SETTLEMENT         = "DataPostAddressSettlement"
	COLUMN_NAME_VERIFY_POSTADDRESS_STREETTYPE         = "DataPostAddressStreetType"
	COLUMN_NAME_VERIFY_POSTADDRESS_STREETTYPEFULL     = "DataPostAddressStreetTypeFull"
	COLUMN_NAME_VERIFY_POSTADDRESS_STREET             = "DataPostAddressStreet"
	COLUMN_NAME_VERIFY_POSTADDRESS_HOUSETYPE          = "DataPostAddressHouseType"
	COLUMN_NAME_VERIFY_POSTADDRESS_HOUSETYPEFULL      = "DataPostAddressHouseTypeFull"
	COLUMN_NAME_VERIFY_POSTADDRESS_HOUSE              = "DataPostAddressHouse"
	COLUMN_NAME_VERIFY_POSTADDRESS_BLOCKTYPE          = "DataPostAddressBlockType"
	COLUMN_NAME_VERIFY_POSTADDRESS_BLOCKTYPEFULL      = "DataPostAddressBlockTypeFull"
	COLUMN_NAME_VERIFY_POSTADDRESS_BLOCK              = "DataPostAddressBlock"
	COLUMN_NAME_VERIFY_POSTADDRESS_FLATTYPE           = "DataPostAddressFlatType"
	COLUMN_NAME_VERIFY_POSTADDRESS_FLAT               = "DataPostAddressFlat"
	COLUMN_NAME_VERIFY_POSTADDRESS_FLATAREA           = "DataPostAddressFlatArea"
	COLUMN_NAME_VERIFY_POSTADDRESS_SQUAREMETERPRICE   = "DataPostAddressSquareMeterPrice"
	COLUMN_NAME_VERIFY_POSTADDRESS_FLATPRICE          = "DataPostAddressFlatPrice"
	COLUMN_NAME_VERIFY_POSTADDRESS_POSTALBOX          = "DataPostAddressPostalBox"
	COLUMN_NAME_VERIFY_POSTADDRESS_FIASID             = "DataPostAddressFiasId"
	COLUMN_NAME_VERIFY_POSTADDRESS_KLADRID            = "DataPostAddressKladrId"
	COLUMN_NAME_VERIFY_POSTADDRESS_OKATO              = "DataPostAddressOkato"
	COLUMN_NAME_VERIFY_POSTADDRESS_OKTMO              = "DataPostAddressOktmo"
	COLUMN_NAME_VERIFY_POSTADDRESS_TAXOFFICE          = "DataPostAddressTaxOffice"
	COLUMN_NAME_VERIFY_POSTADDRESS_TAXOFFICELEGAL     = "DataPostAddressTaxOfficeLegal"
	COLUMN_NAME_VERIFY_POSTADDRESS_TIMEZONE           = "DataPostAddressTimezone"
	COLUMN_NAME_VERIFY_POSTADDRESS_GEOLAT             = "DataPostAddressGeoLat"
	COLUMN_NAME_VERIFY_POSTADDRESS_GEOLON             = "DataPostAddressGeoLon"
	COLUMN_NAME_VERIFY_POSTADDRESS_QCGEO              = "DataPostAddressQcGeo"
	COLUMN_NAME_VERIFY_POSTADDRESS_QCCOMPLETE         = "DataPostAddressQcComplete"
	COLUMN_NAME_VERIFY_POSTADDRESS_QCHOUSE            = "DataPostAddressQcHouse"
	COLUMN_NAME_VERIFY_POSTADDRESS_QUALITYCODE        = "DataPostAddressQualityCode"
	COLUMN_NAME_VERIFY_POSTADDRESS_UNPARSEDPARTS      = "DataPostAddressUnparsedParts"
	COLUMN_NAME_VERIFY_VEHICLE_RESULT                 = "DataVehicleResult"
	COLUMN_NAME_VERIFY_VEHICLE_BRAND                  = "DataVehicleBrand"
	COLUMN_NAME_VERIFY_VEHICLE_MODEL                  = "DataVehicleModel"
	COLUMN_NAME_VERIFY_VEHICLE_QUALITYCODE            = "DataVehicleQualityCode"
	COLUMN_NAME_VERIFY_PHONE_TYPE                     = "DataPhoneType"
	COLUMN_NAME_VERIFY_PHONE_RESULT                   = "DataPhoneResult"
	COLUMN_NAME_VERIFY_PHONE_COUNTRYCODE              = "DataPhoneCountryCode"
	COLUMN_NAME_VERIFY_PHONE_CITYCODE                 = "DataPhoneCityCode"
	COLUMN_NAME_VERIFY_PHONE_NUMBER                   = "DataPhoneNumber"
	COLUMN_NAME_VERIFY_PHONE_EXTENSION                = "DataPhoneExtension"
	COLUMN_NAME_VERIFY_PHONE_PROVIDER                 = "DataPhoneProvider"
	COLUMN_NAME_VERIFY_PHONE_REGION                   = "DataPhoneRegion"
	COLUMN_NAME_VERIFY_PHONE_TIMEZONE                 = "DataPhoneTimezone"
	COLUMN_NAME_VERIFY_PHONE_QCCONFLICT               = "DataPhoneQCConflict"
	COLUMN_NAME_VERIFY_PHONE_QUALITYCODE              = "DataPhoneQualityCode"
	COLUMN_NAME_VERIFY_FULLNAME_RESULT                = "DataFullnameResult"
	COLUMN_NAME_VERIFY_FULLNAME_SURNAME               = "DataFullnameSurname"
	COLUMN_NAME_VERIFY_FULLNAME_NAME                  = "DataFullnameName"
	COLUMN_NAME_VERIFY_FULLNAME_PATRONYMIC            = "DataFullnamePatronymic"
	COLUMN_NAME_VERIFY_FULLNAME_GENDER                = "DataFullnameGender"
	COLUMN_NAME_VERIFY_FULLNAME_QUALITYCODE           = "DataFullnameQualityCode"
	COLUMN_NAME_VERIFY_EMAIL_RESULT                   = "DataEmailResult"
	COLUMN_NAME_VERIFY_PASSPORT_CODE                  = "DataPassportCode"
	COLUMN_NAME_VERIFY_PASSPORT_NUMBER                = "DataPassportNumber"
	COLUMN_NAME_VERIFY_PASSPORT_QUALITYCODE           = "DataPassportQualityCode"
)

type VerifyWorkflow struct {
	OrderRepository           services.OrderRepository
	FacilityRepository        services.FacilityRepository
	VerifyFacilityRepository  services.VerifyFacilityRepository
	OrderStatusRepository     services.OrderStatusRepository
	CustomerTableRepository   services.CustomerTableRepository
	VerifyTableRepository     services.VerifyTableRepository
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
	VerifyProductRepository   services.VerifyProductRepository
	DataColumnRepository      services.DataColumnRepository
}

func NewVerifyWorkflow(orderrepository services.OrderRepository, facilityrepository services.FacilityRepository,
	verifyfacilityrepository services.VerifyFacilityRepository, orderstatusrepository services.OrderStatusRepository,
	customertablerepository services.CustomerTableRepository, verifytablerepository services.VerifyTableRepository,
	resulttablerepository services.ResultTableRepository, worktablerepository services.WorkTableRepository,
	invoicerepository services.InvoiceRepository, companyrepository services.CompanyRepository,
	operationrepository services.OperationRepository, transactiontyperepository services.TransactionTypeRepository,
	tablecolumnrepository services.TableColumnRepository, unitrepository services.UnitRepository,
	tablerowrepository services.TableRowRepository, pricerepository services.PriceRepository,
	verifyproductrepository services.VerifyProductRepository, datacolumnrepository services.DataColumnRepository) *VerifyWorkflow {
	return &VerifyWorkflow{
		OrderRepository:           orderrepository,
		FacilityRepository:        facilityrepository,
		VerifyFacilityRepository:  verifyfacilityrepository,
		OrderStatusRepository:     orderstatusrepository,
		CustomerTableRepository:   customertablerepository,
		VerifyTableRepository:     verifytablerepository,
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
		VerifyProductRepository:   verifyproductrepository,
		DataColumnRepository:      datacolumnrepository,
	}
}

func SendVerify(supplier *libTypes.Supplier, verify *[]libTypes.VerifyData) (verifyresponse *libTypes.VerifyDataResponse, err error) {
	// Отправка Verify
	suppliers.WaitClientReady()                                   // Функция блокируется до момента пока связь с сервером не будет установлена
	response, err := suppliers.SendVerifyData(*supplier, *verify) // Функция не блокируется. Если связи с сервером нет, то вернётся ошибка
	if err != nil {
		log.Error("Can't send verify requests %v for supplier %v", err, supplier.Name)
		return nil, err
	}
	return &response, nil
}

func GetVerifyStatus(verifyresponse *libTypes.VerifyDataResponse) (verifystatuses map[gocql.UUID]libTypes.VerifyDataStatus, err error) {
	verifystatuses = make(map[gocql.UUID]libTypes.VerifyDataStatus)
	// Получение статуса по UUID запроса
	for {
		final := true
		for i := range verifyresponse.Ids {
			var status libTypes.VerifyDataStatus
			status, err = suppliers.StatusVerifyData(verifyresponse.Ids[i])
			if err != nil {
				log.Error("Can't get verify statuses %v", err)
				return map[gocql.UUID]libTypes.VerifyDataStatus{}, err
			}
			if status.Final {
				verifystatuses[status.Id] = status
			} else {
				final = false
			}
		}
		if final {
			break
		}
		time.Sleep(time.Second)
	}

	return verifystatuses, nil
}

func (verifyworkflow *VerifyWorkflow) CheckVerifyOrder(dtoorder *models.DtoOrder,
	dtoverifyfacility *models.DtoVerifyFacility) (verifytable *models.ApiVerifyTable, err error) {
	verifytable, err = verifyworkflow.VerifyTableRepository.Get(dtoverifyfacility.TablesDataId)
	if err != nil {
		return nil, err
	}
	datacolumns, err := verifyworkflow.DataColumnRepository.GetByOrder(dtoorder.ID)
	if err != nil {
		return nil, err
	}

	if len(*datacolumns) != 0 {
		for _, datacolumn := range *datacolumns {
			found := false
			for _, tablecolumn := range verifytable.Verification {
				if datacolumn.Table_Column_ID == tablecolumn.ID {
					found = true
					break
				}
			}
			if !found {
				log.Error("Can find data column %v in table %v", datacolumn.Table_Column_ID, dtoverifyfacility.TablesDataId)
				return nil, errors.New("Wrong data column")
			}
		}
	} else {
		log.Error("Can't find data columns for order %v", dtoorder.ID)
		return nil, errors.New("Missed data columns")
	}

	return verifytable, nil
}

func (verifyworkflow *VerifyWorkflow) GetVerifyTableColumns(dtoorder *models.DtoOrder,
	dtoverifyfacility *models.DtoVerifyFacility) (tablecolumns *[]models.DtoTableColumn,
	verifycolumns map[int]int64, err error) {
	datacolumns, err := verifyworkflow.DataColumnRepository.GetByOrder(dtoorder.ID)
	if err != nil {
		return nil, map[int]int64{}, err
	}

	tablecolumns = new([]models.DtoTableColumn)
	verifycolumns = make(map[int]int64)
	for _, datacolumn := range *datacolumns {
		dtotablecolumn, err := verifyworkflow.TableColumnRepository.Get(datacolumn.Table_Column_ID)
		if err != nil {
			return nil, map[int]int64{}, err
		}
		if verifycolumns[dtotablecolumn.Column_Type_ID] != 0 {
			log.Error("Only one verify column type %v allowed for table %v", dtotablecolumn.Column_Type_ID, dtotablecolumn.Customer_Table_ID)
			return nil, map[int]int64{}, errors.New("Duplicated verify column type")
		}
		verifycolumns[dtotablecolumn.Column_Type_ID] = dtotablecolumn.ID
		*tablecolumns = append(*tablecolumns, *dtotablecolumn)
	}

	return tablecolumns, verifycolumns, nil
}

func (verifyworkflow *VerifyWorkflow) CalculateCost(apitablerows *[]models.ApiInfoTableRow,
	dtoorder *models.DtoOrder, dtoverifyfacility *models.DtoVerifyFacility) (cost float64, err error) {
	verifyprices, err := helpers.GetVerifyPrices(models.SERVICE_TYPE_VERIFY, dtoorder.Supplier_ID, nil, verifyworkflow.PriceRepository,
		verifyworkflow.TableColumnRepository, verifyworkflow.TableRowRepository, verifyworkflow.VerifyProductRepository, "", true)
	if err != nil {
		return 0, err
	}
	cost = 0
	count := 0
	for _, verifyprice := range *verifyprices {
		cost = float64(len(*apitablerows)) * verifyprice.Price
		count++
	}
	if count < 1 {
		log.Error("Price list doesn't contain positions for supplier %v", dtoorder.Supplier_ID)
		return 0, errors.New("Empty price list")
	}
	if count > 1 {
		log.Error("Price list contains multiple positions for supplier %v", dtoorder.Supplier_ID)
		return 0, errors.New("Multiple price list")
	}

	return cost, nil
}

func (verifyworkflow *VerifyWorkflow) PayAndInvoice(dtoorder *models.DtoOrder, dtoverifyfacility *models.DtoVerifyFacility, balance float64) (err error) {
	if balance < dtoverifyfacility.Cost {
		log.Error("Not enough money at unit balance %v to pay for order %v execution", dtoorder.Unit_ID, dtoorder.ID)
		return errors.New("Not enough money")
	}
	dtocompany, err := verifyworkflow.CompanyRepository.GetPrimaryByUnit(dtoorder.Unit_ID)
	if err != nil {
		return err
	}
	dtounit, err := verifyworkflow.UnitRepository.Get(config.Configuration.SystemAccount)
	if err != nil {
		return err
	}
	dtotransactiontype, err := verifyworkflow.TransactionTypeRepository.Get(models.TRANSACTION_TYPE_SERVICE_FEE_VERIFY)
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
		dtoinvoice.VAT = (dtoverifyfacility.Cost / (1 + float64(dtocompany.VAT)/100)) * float64(dtocompany.VAT) / 100
	}
	dtoinvoice.Total = dtoverifyfacility.Cost
	dtoinvoice.Paid = true
	dtoinvoice.Created = time.Now()
	dtoinvoice.Active = true
	dtoinvoice.InvoiceItems = []models.DtoInvoiceItem{*models.NewDtoInvoiceItem(0, 0, dtotransactiontype.Name, models.INVOICE_ITEM_TYPE_ROUBLE,
		1, dtoinvoice.Total, dtoinvoice.Total)}
	dtoinvoice.PaidAt = time.Now()

	err = verifyworkflow.InvoiceRepository.PayForOrder(dtoorder, dtoinvoice, dtotransaction, true)
	if err != nil {
		return err
	}

	return nil
}

func (verifyworkflow *VerifyWorkflow) SendVerify(apitablerows *[]models.ApiInfoTableRow, dtoorder *models.DtoOrder,
	dtoverifyfacility *models.DtoVerifyFacility) (verifyresponse *libTypes.VerifyDataResponse, err error) {
	dtosupplier, err := verifyworkflow.UnitRepository.Get(dtoorder.Supplier_ID)
	if err != nil {
		return nil, err
	}
	verifysupplier, err := GetSupplier(dtosupplier.UUID)
	if err != nil {
		return nil, err
	}

	verify := new([]libTypes.VerifyData)
	for _, apitablerow := range *apitablerows {
		obj := new(libTypes.VerifyData)
		for _, apitablecell := range apitablerow.Cells {
			var datatype libTypes.VerifyDataType

			switch apitablecell.Column_Type_ID {
			case models.COLUMN_TYPE_SOURCE_ADDRESS:
				datatype = libTypes.CVerifyDataPostalAddress
			case models.COLUMN_TYPE_SOURCE_PHONE:
				datatype = libTypes.CVerifyDataPhone
			case models.COLUMN_TYPE_SOURCE_PASSPORT:
				datatype = libTypes.CVerifyDataPassport
			case models.COLUMN_TYPE_SOURCE_FIO:
				datatype = libTypes.CVerifyDataFullName
			case models.COLUMN_TYPE_SOURCE_EMAIL:
				datatype = libTypes.CVerifyDataEmail
			case models.COLUMN_TYPE_SOURCE_DATE:
				datatype = libTypes.CVerifyDataIgnore
			case models.COLUMN_TYPE_SOURCE_AUTOMOBILE:
				datatype = libTypes.CVerifyDataVehicleMakeModel
			default:
				log.Error("Can't find corresponding data type for column %v", apitablecell.Table_Column_ID)
				return nil, errors.New("Wrong data type")
			}

			obj.Data = append(obj.Data, libTypes.VerifyDataBody{
				Type: datatype,
				Data: []string{apitablecell.Value},
			})
		}
		*verify = append(*verify, *obj)
	}
	verifyresponse, err = SendVerify(verifysupplier, verify)
	if err != nil {
		return nil, err
	}

	return verifyresponse, nil
}

func (verifyworkflow *VerifyWorkflow) CopyData(dtoorder *models.DtoOrder,
	dtodatatable *models.DtoCustomerTable) (dtoworkdatatable *models.DtoCustomerTable, err error) {
	dtoworkdatatable, err = verifyworkflow.CustomerTableRepository.Copy(dtodatatable, true)
	if err != nil {
		return nil, err
	}
	dtodatatable.TypeID = models.TABLE_TYPE_READONLY
	err = verifyworkflow.CustomerTableRepository.Update(dtodatatable)
	if err != nil {
		return nil, err
	}
	dtoworkdatatable.TypeID = models.TABLE_TYPE_HIDDEN
	err = verifyworkflow.CustomerTableRepository.Update(dtoworkdatatable)
	if err != nil {
		return nil, err
	}

	dtoworktable := models.NewDtoWorkTable(dtoorder.ID, dtoworkdatatable.ID)
	err = verifyworkflow.WorkTableRepository.Create(dtoworktable, nil)
	if err != nil {
		return nil, err
	}

	return dtoworkdatatable, nil
}

func (verifyworkflow *VerifyWorkflow) SaveVerify(dtoworkdatatable *models.DtoCustomerTable, verifyresponse *libTypes.VerifyDataResponse) (err error) {
	position, err := verifyworkflow.TableColumnRepository.GetDefaultPosition(dtoworkdatatable.ID)
	if err != nil {
		return err
	}
	tablecolumns := []models.DtoTableColumn{}
	columnnames := []string{COLUMN_NAME_VERIFY_ID, COLUMN_NAME_VERIFY_ERROR}
	for _, columnname := range columnnames {
		position++
		fieldnum, err := helpers.FindFreeColumnInternal(dtoworkdatatable.ID, 0, verifyworkflow.TableColumnRepository)
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

		err = verifyworkflow.TableColumnRepository.Create(dtotablecolumn, nil)
		if err != nil {
			return err
		}
		tablecolumns = append(tablecolumns, *dtotablecolumn)
	}

	apitablerows, err := verifyworkflow.TableRowRepository.GetAll("", "", dtoworkdatatable.ID, &tablecolumns)
	if err != nil {
		return err
	}

	workdatatablecolumns, err := verifyworkflow.TableColumnRepository.GetByTable(dtoworkdatatable.ID)
	if err != nil {
		return err
	}

	for index, apitablerow := range *apitablerows {
		tablerow, err := verifyworkflow.TableRowRepository.Get(apitablerow.ID)
		if err != nil {
			return err
		}

		tablecells, err := tablerow.TableRowToDtoTableCells(workdatatablecolumns)
		if err != nil {
			return err
		}

		for i, _ := range *tablecells {
			for j, _ := range tablecolumns {
				FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_ID, verifyresponse.Ids[index].String())
				if verifyresponse.Errors[index] != nil {
					FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_ERROR, verifyresponse.Errors[index].Error())
				}
			}
		}

		err = tablerow.TableCellsToTableRow(tablecells, workdatatablecolumns)
		if err != nil {
			return err
		}

		err = verifyworkflow.TableRowRepository.Update(tablerow, nil, true, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (verifyworkflow *VerifyWorkflow) SaveVerifyStatus(dtoworkdatatable *models.DtoCustomerTable,
	verifystatuses map[gocql.UUID]libTypes.VerifyDataStatus, verifycolumns map[int]int64) (err error) {
	var columnverifyid *models.DtoTableColumn
	alltablecolumns, err := verifyworkflow.TableColumnRepository.GetByTable(dtoworkdatatable.ID)
	if err != nil {
		return err
	}
	found := false
	for i := range *alltablecolumns {
		if (*alltablecolumns)[i].Name == COLUMN_NAME_VERIFY_ID {
			columnverifyid = &(*alltablecolumns)[i]
			found = true
			break
		}
	}
	if !found {
		log.Error("Can't find verify id column for table %v", dtoworkdatatable.ID)
		return errors.New("Missed verify id column")
	}

	position, err := verifyworkflow.TableColumnRepository.GetDefaultPosition(dtoworkdatatable.ID)
	if err != nil {
		return err
	}
	tablecolumns := []models.DtoTableColumn{*columnverifyid}

	columnnames := []string{COLUMN_NAME_VERIFY_STATUS_ID, COLUMN_NAME_VERIFY_STATUS_ERROR}
	columntypes := []int{models.COLUMN_TYPE_DEFAULT, models.COLUMN_TYPE_DEFAULT}
	if verifycolumns[models.COLUMN_TYPE_SOURCE_ADDRESS] != 0 {
		columnnames = append(columnnames, COLUMN_NAME_VERIFY_POSTADDRESS_RESULT, COLUMN_NAME_VERIFY_POSTADDRESS_POSTALCODE,
			COLUMN_NAME_VERIFY_POSTADDRESS_COUNTRY, COLUMN_NAME_VERIFY_POSTADDRESS_REGIONTYPE, COLUMN_NAME_VERIFY_POSTADDRESS_REGIONTYPEFULL,
			COLUMN_NAME_VERIFY_POSTADDRESS_REGION, COLUMN_NAME_VERIFY_POSTADDRESS_AREATYPE, COLUMN_NAME_VERIFY_POSTADDRESS_AREATYPEFULL,
			COLUMN_NAME_VERIFY_POSTADDRESS_AREA, COLUMN_NAME_VERIFY_POSTADDRESS_CITYTYPE, COLUMN_NAME_VERIFY_POSTADDRESS_CITYTYPEFULL,
			COLUMN_NAME_VERIFY_POSTADDRESS_CITY, COLUMN_NAME_VERIFY_POSTADDRESS_SETTLEMENTTYPE, COLUMN_NAME_VERIFY_POSTADDRESS_SETTLEMENTTYPEFULL,
			COLUMN_NAME_VERIFY_POSTADDRESS_SETTLEMENT, COLUMN_NAME_VERIFY_POSTADDRESS_STREETTYPE, COLUMN_NAME_VERIFY_POSTADDRESS_STREETTYPEFULL,
			COLUMN_NAME_VERIFY_POSTADDRESS_STREET, COLUMN_NAME_VERIFY_POSTADDRESS_HOUSETYPE, COLUMN_NAME_VERIFY_POSTADDRESS_HOUSETYPEFULL,
			COLUMN_NAME_VERIFY_POSTADDRESS_HOUSE, COLUMN_NAME_VERIFY_POSTADDRESS_BLOCKTYPE, COLUMN_NAME_VERIFY_POSTADDRESS_BLOCKTYPEFULL,
			COLUMN_NAME_VERIFY_POSTADDRESS_BLOCK, COLUMN_NAME_VERIFY_POSTADDRESS_FLATTYPE, COLUMN_NAME_VERIFY_POSTADDRESS_FLAT,
			COLUMN_NAME_VERIFY_POSTADDRESS_FLATAREA, COLUMN_NAME_VERIFY_POSTADDRESS_SQUAREMETERPRICE, COLUMN_NAME_VERIFY_POSTADDRESS_FLATPRICE,
			COLUMN_NAME_VERIFY_POSTADDRESS_POSTALBOX, COLUMN_NAME_VERIFY_POSTADDRESS_FIASID, COLUMN_NAME_VERIFY_POSTADDRESS_KLADRID,
			COLUMN_NAME_VERIFY_POSTADDRESS_OKATO, COLUMN_NAME_VERIFY_POSTADDRESS_OKTMO, COLUMN_NAME_VERIFY_POSTADDRESS_TAXOFFICE,
			COLUMN_NAME_VERIFY_POSTADDRESS_TAXOFFICELEGAL, COLUMN_NAME_VERIFY_POSTADDRESS_TIMEZONE, COLUMN_NAME_VERIFY_POSTADDRESS_GEOLAT,
			COLUMN_NAME_VERIFY_POSTADDRESS_GEOLON, COLUMN_NAME_VERIFY_POSTADDRESS_QCGEO, COLUMN_NAME_VERIFY_POSTADDRESS_QCCOMPLETE,
			COLUMN_NAME_VERIFY_POSTADDRESS_QCHOUSE, COLUMN_NAME_VERIFY_POSTADDRESS_QUALITYCODE, COLUMN_NAME_VERIFY_POSTADDRESS_UNPARSEDPARTS)
		columntypes = append(columntypes, models.COLUMN_TYPE_ANSWER_POSTADDRESS_RESULT, models.COLUMN_TYPE_ANSWER_POSTADDRESS_POSTALCODE,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_COUNTRY, models.COLUMN_TYPE_ANSWER_POSTADDRESS_REGIONTYPE,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_REGIONTYPEFULL, models.COLUMN_TYPE_ANSWER_POSTADDRESS_REGION, models.COLUMN_TYPE_ANSWER_POSTADDRESS_AREATYPE,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_AREATYPEFULL, models.COLUMN_TYPE_ANSWER_POSTADDRESS_AREA, models.COLUMN_TYPE_ANSWER_POSTADDRESS_CITYTYPE,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_CITYTYPEFULL, models.COLUMN_TYPE_ANSWER_POSTADDRESS_CITY, models.COLUMN_TYPE_ANSWER_POSTADDRESS_SETTLEMENTTYPE,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_SETTLEMENTTYPEFULL, models.COLUMN_TYPE_ANSWER_POSTADDRESS_SETTLEMENT,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_STREETTYPE, models.COLUMN_TYPE_ANSWER_POSTADDRESS_STREETTYPEFULL, models.COLUMN_TYPE_ANSWER_POSTADDRESS_STREET,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_HOUSETYPE, models.COLUMN_TYPE_ANSWER_POSTADDRESS_HOUSETYPEFULL, models.COLUMN_TYPE_ANSWER_POSTADDRESS_HOUSE,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_BLOCKTYPE, models.COLUMN_TYPE_ANSWER_POSTADDRESS_BLOCKTYPEFULL, models.COLUMN_TYPE_ANSWER_POSTADDRESS_BLOCK,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_FLATTYPE, models.COLUMN_TYPE_ANSWER_POSTADDRESS_FLAT, models.COLUMN_TYPE_ANSWER_POSTADDRESS_FLATAREA,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_SQUAREMETERPRICE, models.COLUMN_TYPE_ANSWER_POSTADDRESS_FLATPRICE,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_POSTALBOX, models.COLUMN_TYPE_ANSWER_POSTADDRESS_FIASID, models.COLUMN_TYPE_ANSWER_POSTADDRESS_KLADRID,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_OKATO, models.COLUMN_TYPE_ANSWER_POSTADDRESS_OKTMO, models.COLUMN_TYPE_ANSWER_POSTADDRESS_TAXOFFICE,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_TAXOFFICELEGAL, models.COLUMN_TYPE_ANSWER_POSTADDRESS_TIMEZONE, models.COLUMN_TYPE_ANSWER_POSTADDRESS_GEOLAT,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_GEOLON, models.COLUMN_TYPE_ANSWER_POSTADDRESS_QCGEO, models.COLUMN_TYPE_ANSWER_POSTADDRESS_QCCOMPLETE,
			models.COLUMN_TYPE_ANSWER_POSTADDRESS_QCHOUSE, models.COLUMN_TYPE_ANSWER_POSTADDRESS_QUALITYCODE, models.COLUMN_TYPE_ANSWER_POSTADDRESS_UNPARSEDPARTS)
	}
	if verifycolumns[models.COLUMN_TYPE_SOURCE_PHONE] != 0 {
		columnnames = append(columnnames, COLUMN_NAME_VERIFY_PHONE_TYPE, COLUMN_NAME_VERIFY_PHONE_RESULT, COLUMN_NAME_VERIFY_PHONE_COUNTRYCODE,
			COLUMN_NAME_VERIFY_PHONE_CITYCODE, COLUMN_NAME_VERIFY_PHONE_NUMBER, COLUMN_NAME_VERIFY_PHONE_EXTENSION, COLUMN_NAME_VERIFY_PHONE_PROVIDER,
			COLUMN_NAME_VERIFY_PHONE_REGION, COLUMN_NAME_VERIFY_PHONE_TIMEZONE, COLUMN_NAME_VERIFY_PHONE_QCCONFLICT, COLUMN_NAME_VERIFY_PHONE_QUALITYCODE)
		columntypes = append(columntypes, models.COLUMN_TYPE_ANSWER_PHONE_TYPE, models.COLUMN_TYPE_ANSWER_PHONE_RESULT,
			models.COLUMN_TYPE_ANSWER_PHONE_COUNTRYCODE, models.COLUMN_TYPE_ANSWER_PHONE_CITYCODE, models.COLUMN_TYPE_ANSWER_PHONE_NUMBER,
			models.COLUMN_TYPE_ANSWER_PHONE_EXTENSION, models.COLUMN_TYPE_ANSWER_PHONE_PROVIDER, models.COLUMN_TYPE_ANSWER_PHONE_REGION,
			models.COLUMN_TYPE_ANSWER_PHONE_TIMEZONE, models.COLUMN_TYPE_ANSWER_PHONE_QCCONFLICT, models.COLUMN_TYPE_ANSWER_PHONE_QUALITYCODE)
	}
	if verifycolumns[models.COLUMN_TYPE_SOURCE_PASSPORT] != 0 {
		columnnames = append(columnnames, COLUMN_NAME_VERIFY_PASSPORT_CODE, COLUMN_NAME_VERIFY_PASSPORT_NUMBER, COLUMN_NAME_VERIFY_PASSPORT_QUALITYCODE)
		columntypes = append(columntypes, models.COLUMN_TYPE_ANSWER_PASSPORT_CODE, models.COLUMN_TYPE_ANSWER_PASSPORT_NUMBER,
			models.COLUMN_TYPE_ANSWER_PASSPORT_QUALITYCODE)
	}
	if verifycolumns[models.COLUMN_TYPE_SOURCE_FIO] != 0 {
		columnnames = append(columnnames, COLUMN_NAME_VERIFY_FULLNAME_RESULT, COLUMN_NAME_VERIFY_FULLNAME_SURNAME, COLUMN_NAME_VERIFY_FULLNAME_NAME,
			COLUMN_NAME_VERIFY_FULLNAME_PATRONYMIC, COLUMN_NAME_VERIFY_FULLNAME_GENDER, COLUMN_NAME_VERIFY_FULLNAME_QUALITYCODE)
		columntypes = append(columntypes, models.COLUMN_TYPE_ANSWER_FULLNAME_RESULT, models.COLUMN_TYPE_ANSWER_FULLNAME_SURNAME,
			models.COLUMN_TYPE_ANSWER_FULLNAME_NAME, models.COLUMN_TYPE_ANSWER_FULLNAME_PATRONYMIC, models.COLUMN_TYPE_ANSWER_FULLNAME_GENDER,
			models.COLUMN_TYPE_ANSWER_FULLNAME_QUALITYCODE)
	}
	if verifycolumns[models.COLUMN_TYPE_SOURCE_EMAIL] != 0 {
		columnnames = append(columnnames, COLUMN_NAME_VERIFY_EMAIL_RESULT)
		columntypes = append(columntypes, models.COLUMN_TYPE_ANSWER_EMAIL_RESULT)
	}
	if verifycolumns[models.COLUMN_TYPE_SOURCE_AUTOMOBILE] != 0 {
		columnnames = append(columnnames, COLUMN_NAME_VERIFY_VEHICLE_RESULT, COLUMN_NAME_VERIFY_VEHICLE_BRAND, COLUMN_NAME_VERIFY_VEHICLE_MODEL)
		columntypes = append(columntypes, models.COLUMN_TYPE_ANSWER_VEHICLE_RESULT, models.COLUMN_TYPE_ANSWER_VEHICLE_BRAND,
			models.COLUMN_TYPE_ANSWER_VEHICLE_MODEL)
	}
	for index := range columnnames {
		position++
		fieldnum, err := helpers.FindFreeColumnInternal(dtoworkdatatable.ID, 0, verifyworkflow.TableColumnRepository)
		if err != nil {
			return err
		}

		dtotablecolumn := new(models.DtoTableColumn)
		dtotablecolumn.Created = time.Now()
		dtotablecolumn.Position = position
		dtotablecolumn.Name = columnnames[index]
		dtotablecolumn.Customer_Table_ID = dtoworkdatatable.ID
		dtotablecolumn.Column_Type_ID = columntypes[index]
		dtotablecolumn.Prebuilt = true
		dtotablecolumn.FieldNum = fieldnum
		dtotablecolumn.Active = true
		dtotablecolumn.Edition = 0

		err = verifyworkflow.TableColumnRepository.Create(dtotablecolumn, nil)
		if err != nil {
			return err
		}
		tablecolumns = append(tablecolumns, *dtotablecolumn)
	}

	apitablerows, err := verifyworkflow.TableRowRepository.GetAll("", "", dtoworkdatatable.ID, &tablecolumns)
	if err != nil {
		return err
	}

	workdatatablecolumns, err := verifyworkflow.TableColumnRepository.GetByTable(dtoworkdatatable.ID)
	if err != nil {
		return err
	}

	for _, apitablerow := range *apitablerows {
		tablerow, err := verifyworkflow.TableRowRepository.Get(apitablerow.ID)
		if err != nil {
			return err
		}

		tablecells, err := tablerow.TableRowToDtoTableCells(workdatatablecolumns)
		if err != nil {
			return err
		}

		var verifyid gocql.UUID
		found := false
		for i, _ := range *tablecells {
			if (*tablecells)[i].Table_Column_ID == columnverifyid.ID {
				verifyid, err = gocql.ParseUUID((*tablecells)[i].Value)
				if err != nil {
					log.Error("Can't parse verify id from column %v, %v", err, columnverifyid.ID)
					return err
				}
				found = true
				break
			}
		}
		if !found {
			log.Error("Can't find verify id value for column %v", columnverifyid.ID)
			return errors.New("Missed verify id value")
		}

		for i, _ := range *tablecells {
			for j, _ := range tablecolumns {
				FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_STATUS_ID, fmt.Sprintf("%v", verifystatuses[verifyid].Status))
				FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_STATUS_ERROR, fmt.Sprintf("%v", verifystatuses[verifyid].Error))
				if verifycolumns[models.COLUMN_TYPE_SOURCE_ADDRESS] != 0 {
					if verifystatuses[verifyid].DataPostalAddress != nil {
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_RESULT, verifystatuses[verifyid].DataPostalAddress.Result)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_POSTALCODE, verifystatuses[verifyid].DataPostalAddress.PostalCode)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_COUNTRY, verifystatuses[verifyid].DataPostalAddress.Country)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_REGIONTYPE, verifystatuses[verifyid].DataPostalAddress.RegionType)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_REGIONTYPEFULL, verifystatuses[verifyid].DataPostalAddress.RegionTypeFull)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_REGION, verifystatuses[verifyid].DataPostalAddress.Region)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_AREATYPE, verifystatuses[verifyid].DataPostalAddress.AreaType)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_AREATYPEFULL, verifystatuses[verifyid].DataPostalAddress.AreaTypeFull)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_AREA, verifystatuses[verifyid].DataPostalAddress.Area)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_CITYTYPE, verifystatuses[verifyid].DataPostalAddress.CityType)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_CITYTYPEFULL, verifystatuses[verifyid].DataPostalAddress.CityTypeFull)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_CITY, verifystatuses[verifyid].DataPostalAddress.City)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_SETTLEMENTTYPE, verifystatuses[verifyid].DataPostalAddress.SettlementType)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_SETTLEMENTTYPEFULL, verifystatuses[verifyid].DataPostalAddress.SettlementTypeFull)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_SETTLEMENT, verifystatuses[verifyid].DataPostalAddress.Settlement)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_STREETTYPE, verifystatuses[verifyid].DataPostalAddress.StreetType)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_STREETTYPEFULL, verifystatuses[verifyid].DataPostalAddress.StreetTypeFull)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_STREET, verifystatuses[verifyid].DataPostalAddress.Street)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_HOUSETYPE, verifystatuses[verifyid].DataPostalAddress.HouseType)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_HOUSETYPEFULL, verifystatuses[verifyid].DataPostalAddress.HouseTypeFull)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_HOUSE, verifystatuses[verifyid].DataPostalAddress.House)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_BLOCKTYPE, verifystatuses[verifyid].DataPostalAddress.BlockType)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_BLOCKTYPEFULL, verifystatuses[verifyid].DataPostalAddress.BlockTypeFull)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_BLOCK, verifystatuses[verifyid].DataPostalAddress.Block)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_FLATTYPE, verifystatuses[verifyid].DataPostalAddress.FlatType)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_FLAT, verifystatuses[verifyid].DataPostalAddress.Flat)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_FLATAREA, verifystatuses[verifyid].DataPostalAddress.FlatArea)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_SQUAREMETERPRICE, verifystatuses[verifyid].DataPostalAddress.SquareMeterPrice)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_FLATPRICE, verifystatuses[verifyid].DataPostalAddress.FlatPrice)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_POSTALBOX, verifystatuses[verifyid].DataPostalAddress.PostalBox)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_FIASID, verifystatuses[verifyid].DataPostalAddress.FiasId)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_KLADRID, verifystatuses[verifyid].DataPostalAddress.KladrId)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_OKATO, verifystatuses[verifyid].DataPostalAddress.Okato)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_OKTMO, verifystatuses[verifyid].DataPostalAddress.Oktmo)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_TAXOFFICE, verifystatuses[verifyid].DataPostalAddress.TaxOffice)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_TAXOFFICELEGAL, verifystatuses[verifyid].DataPostalAddress.TaxOfficeLegal)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_TIMEZONE, verifystatuses[verifyid].DataPostalAddress.Timezone)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_GEOLAT, verifystatuses[verifyid].DataPostalAddress.GeoLat)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_GEOLON, verifystatuses[verifyid].DataPostalAddress.GeoLon)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_QCGEO, verifystatuses[verifyid].DataPostalAddress.QcGeo)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_QCCOMPLETE, verifystatuses[verifyid].DataPostalAddress.QcComplete)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_QCHOUSE, verifystatuses[verifyid].DataPostalAddress.QcHouse)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_QUALITYCODE, fmt.Sprintf("%v", verifystatuses[verifyid].DataPostalAddress.QualityCode))
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_POSTADDRESS_UNPARSEDPARTS, verifystatuses[verifyid].DataPostalAddress.UnparsedParts)
					}
				}
				if verifycolumns[models.COLUMN_TYPE_SOURCE_PHONE] != 0 {
					if verifystatuses[verifyid].DataPhone != nil {
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PHONE_TYPE, verifystatuses[verifyid].DataPhone.Type)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PHONE_RESULT, verifystatuses[verifyid].DataPhone.Phone)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PHONE_COUNTRYCODE, verifystatuses[verifyid].DataPhone.CountryCode)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PHONE_CITYCODE, verifystatuses[verifyid].DataPhone.CityCode)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PHONE_NUMBER, verifystatuses[verifyid].DataPhone.Number)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PHONE_EXTENSION, verifystatuses[verifyid].DataPhone.Extension)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PHONE_PROVIDER, verifystatuses[verifyid].DataPhone.Provider)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PHONE_REGION, verifystatuses[verifyid].DataPhone.Region)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PHONE_TIMEZONE, verifystatuses[verifyid].DataPhone.Timezone)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PHONE_QCCONFLICT, verifystatuses[verifyid].DataPhone.QcConflict)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PHONE_QUALITYCODE, fmt.Sprintf("%v", verifystatuses[verifyid].DataPhone.QualityCode))
					}
				}
				if verifycolumns[models.COLUMN_TYPE_SOURCE_PASSPORT] != 0 {
					if verifystatuses[verifyid].DataPassport != nil {
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PASSPORT_CODE, verifystatuses[verifyid].DataPassport.Series)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PASSPORT_NUMBER, verifystatuses[verifyid].DataPassport.Number)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_PASSPORT_QUALITYCODE, fmt.Sprintf("%v", verifystatuses[verifyid].DataPassport.QualityCode))

					}
				}
				if verifycolumns[models.COLUMN_TYPE_SOURCE_FIO] != 0 {
					if verifystatuses[verifyid].DataFullName != nil {
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_FULLNAME_RESULT, verifystatuses[verifyid].DataFullName.Result)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_FULLNAME_SURNAME, verifystatuses[verifyid].DataFullName.Surname)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_FULLNAME_NAME, verifystatuses[verifyid].DataFullName.Name)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_FULLNAME_PATRONYMIC, verifystatuses[verifyid].DataFullName.Patronymic)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_FULLNAME_GENDER, fmt.Sprintf("%v", verifystatuses[verifyid].DataFullName.Gender))
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_FULLNAME_QUALITYCODE, fmt.Sprintf("%v", verifystatuses[verifyid].DataFullName.QualityCode))
					}
				}
				if verifycolumns[models.COLUMN_TYPE_SOURCE_EMAIL] != 0 {
					if verifystatuses[verifyid].DataEmail != nil {
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_EMAIL_RESULT, verifystatuses[verifyid].DataEmail.Email)
					}
				}
				if verifycolumns[models.COLUMN_TYPE_SOURCE_AUTOMOBILE] != 0 {
					if verifystatuses[verifyid].DataVehicleMakeModel != nil {
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_VEHICLE_RESULT, verifystatuses[verifyid].DataVehicleMakeModel.Result)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_VEHICLE_BRAND, verifystatuses[verifyid].DataVehicleMakeModel.Brand)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_VEHICLE_MODEL, verifystatuses[verifyid].DataVehicleMakeModel.Model)
						FillTableCell(&(*tablecells)[i], &tablecolumns[j], COLUMN_NAME_VERIFY_VEHICLE_QUALITYCODE, fmt.Sprintf("%v", verifystatuses[verifyid].DataVehicleMakeModel.QualityCode))
					}
				}
			}
		}

		err = tablerow.TableCellsToTableRow(tablecells, workdatatablecolumns)
		if err != nil {
			return err
		}

		err = verifyworkflow.TableRowRepository.Update(tablerow, nil, true, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (verifyworkflow *VerifyWorkflow) ClearTables(dtoorder *models.DtoOrder, dtoverifyfacility *models.DtoVerifyFacility,
	dtodatatable *models.DtoCustomerTable) (err error) {
	if dtoverifyfacility.TablesDataDelete {
		err = verifyworkflow.CustomerTableRepository.Deactivate(dtodatatable)
		if err != nil {
			return err
		}
	} else {
		dtodatatable.TypeID = models.TABLE_TYPE_DEFAULT
		err = verifyworkflow.CustomerTableRepository.Update(dtodatatable)
		if err != nil {
			return err
		}
	}

	worktables, err := verifyworkflow.WorkTableRepository.GetByOrder(dtoorder.ID)
	if err != nil {
		return err
	}
	for _, worktable := range *worktables {
		dtocustomertable, err := verifyworkflow.CustomerTableRepository.Get(worktable.Customer_Table_ID)
		if err != nil {
			return err
		}
		err = verifyworkflow.CustomerTableRepository.Deactivate(dtocustomertable)
		if err != nil {
			return err
		}
	}

	return nil
}

func (verifyworkflow *VerifyWorkflow) SetStatus(dtoorder *models.DtoOrder, orderstatus models.OrderStatus, active bool) (err error) {
	dtoorderstatus := models.NewDtoOrderStatus(dtoorder.ID, orderstatus, active, "", time.Now())
	err = verifyworkflow.OrderStatusRepository.Save(dtoorderstatus, nil)
	if err != nil {
		return err
	}

	return nil
}

func (verifyworkflow *VerifyWorkflow) ExecuteOrder(order_id int64) {
	log.Info("Starting order %v execution at %v", order_id, time.Now())
	log.Info("Checking order type ...")
	dtoorder, err := verifyworkflow.OrderRepository.Get(order_id)
	if err != nil {
		return
	}
	dtofacility, err := verifyworkflow.FacilityRepository.Get(dtoorder.Facility_ID)
	if err != nil {
		return
	}
	if dtofacility.Alias != models.SERVICE_TYPE_VERIFY {
		log.Error("Order service is not macthed to the service method %v", dtoorder.Facility_ID)
		return
	}
	if !dtofacility.Active {
		log.Error("Service is not active %v", dtofacility.ID)
		return
	}
	log.Info("Checking service type ...")
	dtoverifyfacility, err := verifyworkflow.VerifyFacilityRepository.Get(dtoorder.ID)
	if err != nil {
		return
	}
	log.Info("Checking order status ...")
	dtoorderstatuses, err := verifyworkflow.OrderStatusRepository.GetByOrder(dtoorder.ID)
	if err != nil {
		return
	}

	order := models.NewApiLongOrderFromDto(dtoorder, dtoorderstatuses)
	if order.IsAssembled && order.IsConfirmed && !order.IsOpen && !order.IsCancelled && !order.IsExecuted && !order.IsArchived && !order.IsDeleted {
		/* 1 */
		log.Info("Checking data table ...")
		dtodatatable, err := verifyworkflow.CustomerTableRepository.Get(dtoverifyfacility.TablesDataId)
		if err != nil {
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		if !dtodatatable.Active {
			log.Error("Data table is not active %v", dtoverifyfacility.TablesDataId)
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		if !dtodatatable.Permanent {
			log.Error("Data table is not permanent %v", dtoverifyfacility.TablesDataId)
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		/* 2 */ err = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_OPEN, true)
		if err != nil {
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Checking order data ...")
		_, err = verifyworkflow.CheckVerifyOrder(dtoorder, dtoverifyfacility)
		if err != nil {
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}

		/* 3 */ err = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_SUPPLIER_COST_NEW, true)
		if err != nil {
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Checking table columns ...")
		tablecolumns, verifycolumns, err := verifyworkflow.GetVerifyTableColumns(dtoorder, dtoverifyfacility)
		if err != nil {
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Checking table data ...")
		apitablerows, err := verifyworkflow.TableRowRepository.GetAll("", "", dtoverifyfacility.TablesDataId, tablecolumns)
		if err != nil {
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Calculating order cost ...")
		dtoverifyfacility.Cost, err = verifyworkflow.CalculateCost(apitablerows, dtoorder, dtoverifyfacility)
		if err != nil {
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		dtoverifyfacility.CostFactual = dtoverifyfacility.Cost
		err = verifyworkflow.VerifyFacilityRepository.Update(dtoverifyfacility, true, false)
		if err != nil {
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}

		/* 4 */ err = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED, true)
		if err != nil {
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Start order processing ...")
		/* 5 */ dtoorder.Begin_Date = time.Now()
		err = verifyworkflow.OrderRepository.Update(dtoorder, &[]models.DtoOrderStatus{
			*models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_MODERATOR_BEGIN, true, "", time.Now())}, nil, true)
		if err != nil {
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Calculating unit balance ...")
		balance, err := verifyworkflow.OperationRepository.CalculateBalance(dtoorder.Unit_ID)
		if err != nil {
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Order payment and invoicement ...")
		/* 6 */ err = verifyworkflow.PayAndInvoice(dtoorder, dtoverifyfacility, balance)
		if err != nil {
			_ = verifyworkflow.SetStatus(dtoorder, models.ORDER_STATUS_COMPLETED, false)
			return
		}
		log.Info("Copying table data ...")
		/* 7 */ dtoworkdatatable, err := verifyworkflow.CopyData(dtoorder, dtodatatable)
		if err != nil {
			return
		}
		log.Info("Sending data to supplier ...")
		/* 8 */ verifyresponse, err := verifyworkflow.SendVerify(apitablerows, dtoorder, dtoverifyfacility)
		if err != nil {
			return
		}
		log.Info("Saving supplier response ...")
		/* 9 */ err = verifyworkflow.SaveVerify(dtoworkdatatable, verifyresponse)
		if err != nil {
			return
		}
		log.Info("Getting supplier results ...")
		/* 10 */ verifystatuses, err := GetVerifyStatus(verifyresponse)
		if err != nil {
			return
		}
		log.Info("Saving supplier results ...")
		/* 11 */ err = verifyworkflow.SaveVerifyStatus(dtoworkdatatable, verifystatuses, verifycolumns)
		if err != nil {
			return
		}
		log.Info("Finsing order processing and clearing data ...")
		/* 12 */ err = verifyworkflow.ClearTables(dtoorder, dtoverifyfacility, dtodatatable)
		if err != nil {
			return
		}
		log.Info("Completing order execution %v", time.Now())
		/* 13 */ dtoorder.End_Date = time.Now()
		err = verifyworkflow.OrderRepository.Update(dtoorder, &[]models.DtoOrderStatus{
			*models.NewDtoOrderStatus(dtoorder.ID, models.ORDER_STATUS_SUPPLIER_CLOSE, true, "", time.Now())}, nil, true)
		if err != nil {
			return
		}
	}
}
