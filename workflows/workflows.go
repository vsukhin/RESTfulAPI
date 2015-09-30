/* Workflows package provides methods and data structures for workflow layer implementation */

package workflows

import (
	"application/config"
	"application/models"
	logging "github.com/op/go-logging"
	libSuppliers "lib/suppliers"
	libTypes "lib/suppliers/types"
	"lib/uuid"
)

type Executor interface {
	ExecuteOrder(order_id int64)
}

var (
	log config.Logger = logging.MustGetLogger("workflows")
)

func InitLogger(logger config.Logger) {
	log = logger
}

// Получение поставщика из библиотеки по UUID в виде строки
func GetSupplier(uuidstr string) (sp *libTypes.Supplier, err error) {
	var supp *libSuppliers.Suppliers
	var uuidObj uuid.UUID

	uuidObj, err = uuid.ParseUUID(uuidstr)
	if err != nil {
		log.Error("Can't parse UUID %v, %v", uuidstr, err)
		return nil, err
	}
	supp = libSuppliers.New()
	supplier, err := supp.SuppliersByUUID(uuidObj)
	if err != nil {
		log.Error("Can't find supplier by UUID %v, %v", uuidstr, err)
		return nil, err
	}

	return &supplier, nil
}

func FillTableCell(dtotablecell *models.DtoTableCell, dtotablecolumn *models.DtoTableColumn, column string, value string) {
	if (dtotablecell.Table_Column_ID == dtotablecolumn.ID) && (dtotablecolumn.Name == column) && dtotablecolumn.Prebuilt {
		dtotablecell.Value = value
		dtotablecell.Checked = true
		dtotablecell.Valid = true
	}
}
