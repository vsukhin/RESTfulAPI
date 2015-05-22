package services

import (
	"application/models"
)

type HLRTableRepository interface {
	GetAll(user_id int64) (tables *[]models.ApiHLRTable, err error)
}

type HLRTableService struct {
	FacilityTableRepository FacilityTableRepository
	*Repository
}

func NewHLRTableService(repository *Repository) *HLRTableService {
	return &HLRTableService{Repository: repository}
}

func (hlrtableservice *HLRTableService) GetAll(user_id int64) (hlrtables *[]models.ApiHLRTable, err error) {
	hlrtables = new([]models.ApiHLRTable)

	tables := new([]models.DtoCustomerTable)
	_, err = hlrtableservice.DbContext.Select(tables, "select * from "+hlrtableservice.Table+
		" where unit_id = (select unit_id from users where id = ?) and active = 1 and permanent = 1", user_id)
	if err != nil {
		log.Error("Error during getting all hlr facility tables in database %v with value %v", err, user_id)
		return nil, err
	}

	mobilephones, goodmobilephones, err := hlrtableservice.FacilityTableRepository.GetColumnsByType(user_id, models.COLUMN_TYPE_MOBILE_PHONE)
	if err != nil {
		return nil, err
	}

	for _, table := range *tables {
		hlrtable := new(models.ApiHLRTable)
		hlrtable.ID = table.ID
		hlrtable.Name = table.Name
		hlrtable.UnitID = table.UnitID
		hlrtable.TypeID = table.TypeID
		for _, mobilephone := range *mobilephones {
			if mobilephone.Customer_Table_ID == table.ID && goodmobilephones[mobilephone.ID] != 0 {
				hlrtable.MobilePhones = append(hlrtable.MobilePhones, *models.NewApiTableColumn(mobilephone.ID, mobilephone.Name,
					mobilephone.Column_Type_ID, mobilephone.Position))
			}
		}
		if len(hlrtable.MobilePhones) == 0 {
			continue
		}

		*hlrtables = append(*hlrtables, *hlrtable)
	}

	return hlrtables, nil
}
