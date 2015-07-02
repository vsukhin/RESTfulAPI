package services

import (
	"application/models"
	"regexp"
	"strings"
)

var (
	NumberRegExp = regexp.MustCompile("[0-9]+")
)

type ColumnTypeRepository interface {
	Validate(dtocolumntype *models.DtoColumnType, columnRegExp *regexp.Regexp, value string) (valid bool, corrected string, err error)
	Get(id int) (columntype *models.DtoColumnType, err error)
	GetAll() (columntypes *[]models.ApiColumnType, err error)
	GetByTable(tableid int64) (columntypes map[int]models.DtoColumnType, err error)
	FindByName(name string) (id int, err error)
}

type ColumnTypeService struct {
	*Repository
}

func NewColumnTypeService(repository *Repository) *ColumnTypeService {
	repository.DbContext.AddTableWithName(models.DtoColumnType{}, repository.Table).SetKeys(true, "id")
	return &ColumnTypeService{
		repository,
	}
}

func (columntypeservice *ColumnTypeService) Validate(dtocolumntype *models.DtoColumnType,
	columnRegExp *regexp.Regexp, value string) (valid bool, corrected string, err error) {
	valid = true
	switch dtocolumntype.ID {
	case models.COLUMN_TYPE_BIRTHDAY:
		_, err = models.ParseDate(value)
		if err == nil {
			valid = true
		} else {
			valid = false
		}
	case models.COLUMN_TYPE_MOBILE_PHONE:
		numbers := NumberRegExp.FindAllString(value, -1)
		value = strings.Join(numbers, "")
		runes := []rune(value)
		if len(runes) != 0 {
			if runes[0] == '8' {
				runes[0] = '7'
			}
		}
		value = string(runes)
		fallthrough
	default:
		if dtocolumntype.Regexp != "" {
			if columnRegExp == nil {
				valid, err = regexp.MatchString(dtocolumntype.Regexp, value)
				if err != nil {
					log.Error("Error during running reg exp %v with value %v", err, dtocolumntype.Regexp)
					return false, value, err
				}
			} else {
				valid = columnRegExp.MatchString(value)
			}
		}
	}
	if dtocolumntype.Required {
		if value == "" {
			valid = false
		}
	}

	return valid, value, nil
}

func (columntypeservice *ColumnTypeService) Get(id int) (columntype *models.DtoColumnType, err error) {
	columntype = new(models.DtoColumnType)
	err = columntypeservice.DbContext.SelectOne(columntype, "select * from "+columntypeservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting column type object from database %v with value %v", err, id)
		return nil, err
	}

	return columntype, nil
}

func (columntypeservice *ColumnTypeService) GetAll() (columntypes *[]models.ApiColumnType, err error) {
	columntypes = new([]models.ApiColumnType)
	_, err = columntypeservice.DbContext.Select(columntypes,
		"select id, name, description, required, `regexp`, "+
			"case horAlignmentHead when 1 then 'left' when 2 then 'center' when 3 then 'right' end as alignmentHead, "+
			"case horAlignmentBody when 1 then 'left' when 2 then 'center' when 3 then 'right' end as alignmentBody from "+
			columntypeservice.Table+" where active = 1")
	if err != nil {
		log.Error("Error during getting all column type object from database %v", err)
		return nil, err
	}

	return columntypes, nil
}

func (columntypeservice *ColumnTypeService) GetByTable(tableid int64) (columntypes map[int]models.DtoColumnType, err error) {
	columntypes = make(map[int]models.DtoColumnType)

	tempcolumntypes := new([]models.DtoColumnType)
	_, err = columntypeservice.DbContext.Select(tempcolumntypes,
		"select c.* from "+columntypeservice.Table+" c inner join table_columns t on c.id = t.column_type_id where c.active = 1"+
			" and t.active = 1 and t.customer_table_id = ? order by t.position asc", tableid)
	if err != nil {
		log.Error("Error during getting all column type object from database %v with value %v", err, tableid)
		return nil, err
	}
	for _, columntype := range *tempcolumntypes {
		columntypes[columntype.ID] = columntype
	}

	return columntypes, nil
}

func (columntypeservice *ColumnTypeService) FindByName(name string) (id int, err error) {
	err = columntypeservice.DbContext.SelectOne(&id, "select id from "+columntypeservice.Table+" where name = ?", name)
	if err != nil {
		log.Error("Error during finding column type object from database %v with value %v", err, name)
		return 0, err
	}

	return id, nil
}
