package helpers

import (
	"application/models"
	"application/services"
	"errors"
	"strconv"
)

func CheckMobilePhone(value string, dtocolumntype *models.DtoColumnType, columntyperepository services.ColumnTypeRepository) (mobilephone uint64, err error) {
	valid, corrected, _ := columntyperepository.Validate(dtocolumntype, nil, value)
	if !valid {
		log.Error("Mobile phone is not valid %v", value)
		return 0, errors.New("Wrong mobile phone")
	}
	mobilephone, err = strconv.ParseUint(corrected, 0, 64)
	if err != nil {
		log.Error("Can't convert mobile phone %v", err)
		return 0, err
	}

	return mobilephone, nil
}
