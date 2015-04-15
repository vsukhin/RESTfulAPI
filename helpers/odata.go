package helpers

import (
	"application/config"
	"application/models"
	"errors"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
	"strings"
	"types"
)

const (
	PARAM_QUERY_LIMIT  = "limit"
	PARAM_QUERY_ORDER  = "order"
	PARAM_QUERY_FILTER = "filter"
	PARAM_QUERY_NUMBER = 3

	PARAM_LIMIT_LOW    = 0
	PARAM_LIMIT_HIGH   = 1
	PARAM_LIMIT_NUMBER = 2

	PARAM_SORT_ASC    = "asc"
	PARAM_SORT_DESC   = "desc"
	PARAM_SORT_FIELD  = 0
	PARAM_SORT_ORDER  = 1
	PARAM_SORT_NUMBER = 2

	PARAM_FILTER_OP_EQ = "eq"
	PARAM_FILTER_OP_LT = "lt"
	PARAM_FILTER_OP_LE = "le"
	PARAM_FILTER_OP_GT = "gt"
	PARAM_FILTER_OP_GE = "ge"
	PARAM_FILTER_OP_NE = "ne"
	PARAM_FILTER_OP_LK = "lk"
	PARAM_FILTER_FIELD = 0
	PARAM_FILTER_OP    = 1
	PARAM_FILTER_VALUE = 2
)

func GetLimitQuery(request *http.Request, r render.Render, language string) (query string, err error) {
	query = ""
	limit := request.URL.Query().Get(PARAM_QUERY_LIMIT)
	if limit != "" {
		var offset int64
		var count int64

		limits := strings.Split(limit, ":")
		if len(limits) != PARAM_LIMIT_NUMBER {
			log.Error("Wrong number of limit parameter elements %v", len(limits))
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return "", errors.New("Wrong number of parameters")
		}

		offset, err = strconv.ParseInt(limits[PARAM_LIMIT_LOW], 0, 64)
		if err != nil {
			log.Error("Wrong limit offset %v", limits[PARAM_LIMIT_LOW])
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return "", err
		}
		if offset < 0 {
			log.Error("Wrong limit offset %v", limits[PARAM_LIMIT_LOW])
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return "", errors.New("Wrong offset")
		}

		count, err = strconv.ParseInt(limits[PARAM_LIMIT_HIGH], 0, 64)
		if err != nil {
			log.Error("Wrong limit count %v", limits[PARAM_LIMIT_HIGH])
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return "", err
		}

		if count > 0 {
			query = " limit " + limits[PARAM_LIMIT_LOW] + ", " + limits[PARAM_LIMIT_HIGH]
		}
	} else {
		query = " limit 0, 100"
	}

	return query, nil
}

func GetOrderArray(checker models.Checker, request *http.Request, r render.Render, language string) (sorts *[]models.OrderExp, err error) {
	order := request.URL.Query().Get(PARAM_QUERY_ORDER)
	sorts = new([]models.OrderExp)
	if order != "" {
		orders := strings.Split(order, ",")
		for _, element := range orders {
			elements := strings.Split(element, ":")
			if len(elements) != PARAM_SORT_NUMBER {
				log.Error("Wrong number of sort parameter elements %v", len(elements))
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return nil, errors.New("Wrong parameter number")
			}

			var valid bool
			valid, err = checker.Check(elements[PARAM_SORT_FIELD])
			if !valid || err != nil {
				log.Error("Unknown field name %v", elements[PARAM_SORT_FIELD])
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return nil, errors.New("Uknown field")
			}

			if strings.ToLower(elements[PARAM_SORT_ORDER]) != PARAM_SORT_ASC &&
				strings.ToLower(elements[PARAM_SORT_ORDER]) != PARAM_SORT_DESC {
				log.Error("Unknown sort operation %v", elements[PARAM_SORT_ORDER])
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return nil, errors.New("Uknown sort")
			}
			*sorts = append(*sorts, models.OrderExp{Field: elements[PARAM_SORT_FIELD], Order: elements[PARAM_SORT_ORDER]})
		}
		if len(*sorts) == 0 {
			log.Error("Sort is not found")
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return nil, errors.New("Sort not found")
		}
	}

	return sorts, nil
}

func GetFilterArray(extractor models.Extractor, parameter interface{}, request *http.Request, r render.Render,
	language string) (filters *[]models.FilterExp, err error) {
	filter := request.URL.Query().Get(PARAM_QUERY_FILTER)
	filters = new([]models.FilterExp)
	if filter != "" {
		cons := strings.Split(filter, ",")
		for _, element := range cons {
			elements := strings.Split(element, ":")
			if len(elements) != PARAM_QUERY_NUMBER {
				log.Error("Wrong number of filter parameter elements %v", len(elements))
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return nil, errors.New("Wrong parameter number")
			}
			var allfields bool
			var field string
			var value string

			allfields = false
			if elements[PARAM_FILTER_FIELD] == "*" {
				allfields = true
			}

			if allfields {
				if strings.Contains(elements[PARAM_FILTER_VALUE], "'") {
					log.Error("Wrong field value %v for %v", elements[PARAM_FILTER_VALUE], elements[PARAM_FILTER_FIELD])
					r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[language].Errors.Api.Data_Wrong})
					return nil, errors.New("Wrong value")
				}
				field = ""
				value = "'" + elements[PARAM_FILTER_VALUE] + "'"
			} else {
				var errField error
				var errValue error

				field, value, errField, errValue = extractor.Extract(elements[PARAM_FILTER_FIELD], elements[PARAM_FILTER_VALUE])
				if errField != nil {
					log.Error("Unknown field name %v", elements[PARAM_FILTER_FIELD])
					r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[language].Errors.Api.Data_Wrong})
					return nil, errField
				}
				if errValue != nil {
					log.Error("Wrong field value %v for %v", elements[PARAM_FILTER_VALUE], elements[PARAM_FILTER_FIELD])
					r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[language].Errors.Api.Data_Wrong})
					return nil, errValue
				}
			}

			op := ""
			switch strings.ToLower(elements[PARAM_FILTER_OP]) {
			case PARAM_FILTER_OP_EQ:
				op = "="
			case PARAM_FILTER_OP_LT:
				op = "<"
			case PARAM_FILTER_OP_LE:
				op = "<="
			case PARAM_FILTER_OP_GT:
				op = ">"
			case PARAM_FILTER_OP_GE:
				op = ">="
			case PARAM_FILTER_OP_NE:
				op = "!="
			case PARAM_FILTER_OP_LK:
				op = "like"
				value = strings.Replace(value, "*", "%", -1)
			default:
				log.Error("Unknown filter operation %v", elements[PARAM_FILTER_OP])
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return nil, errors.New("Uknown filter")
			}

			if allfields {
				var fields []string
				for _, field = range *extractor.GetAllFields(parameter) {
					fields = append(fields, field)
				}
				*filters = append(*filters, models.FilterExp{Fields: fields, Op: op, Value: value})
			} else {
				*filters = append(*filters, models.FilterExp{Fields: []string{field}, Op: op, Value: value})
			}
		}
		if len(*filters) == 0 {
			log.Error("Filter is not found")
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return nil, errors.New("Filter not found")
		}
	}

	return filters, nil
}
