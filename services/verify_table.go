package services

import (
	"application/models"
)

type VerifyTableRepository interface {
	GetAll(user_id int64) (tables *[]models.ApiVerifyTable, err error)
}

type VerifyTableService struct {
	FacilityTableRepository FacilityTableRepository
	*Repository
}

func NewVerifyTableService(repository *Repository) *VerifyTableService {
	return &VerifyTableService{Repository: repository}
}

func (verifytableservice *VerifyTableService) GetAll(user_id int64) (verifytables *[]models.ApiVerifyTable, err error) {
	verifytables = new([]models.ApiVerifyTable)

	tables := new([]models.DtoCustomerTable)
	_, err = verifytableservice.DbContext.Select(tables, "select * from "+verifytableservice.Table+
		" where unit_id = (select unit_id from users where id = ?) and active = 1 and permanent = 1", user_id)
	if err != nil {
		log.Error("Error during getting all verify facility tables in database %v with value %v", err, user_id)
		return nil, err
	}

	verification, err := verifytableservice.FacilityTableRepository.GetColumnsByTypes(user_id, []int{models.COLUMN_TYPE_SOURCE_ADDRESS,
		models.COLUMN_TYPE_SOURCE_PHONE, models.COLUMN_TYPE_SOURCE_FIO, models.COLUMN_TYPE_SOURCE_EMAIL, models.COLUMN_TYPE_SOURCE_DATE,
		models.COLUMN_TYPE_SOURCE_AUTOMOBILE})
	if err != nil {
		return nil, err
	}

	for _, table := range *tables {
		verifytable := new(models.ApiVerifyTable)
		verifytable.ID = table.ID
		verifytable.Name = table.Name
		verifytable.UnitID = table.UnitID
		verifytable.TypeID = table.TypeID
		for _, verify := range *verification {
			if verify.Customer_Table_ID == table.ID {
				verifytable.Verification = append(verifytable.Verification, *models.NewApiTableColumn(verify.ID, verify.Name,
					verify.Column_Type_ID, verify.Position))
			}
		}
		if len(verifytable.Verification) == 0 {
			continue
		}

		*verifytables = append(*verifytables, *verifytable)
	}

	return verifytables, nil
}
