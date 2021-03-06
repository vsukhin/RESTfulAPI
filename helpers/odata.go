package helpers

import (
	"application/config"
	"application/models"
	"errors"
	"github.com/martini-contrib/render"
	"net/http"
	"net/url"
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

func ParseRawQuery(query string) (m map[string][]string) {
	m = make(map[string][]string)
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		m[key] = append(m[key], value)
	}

	return m
}

func GetLimitQuery(request *http.Request, r render.Render, language string) (query string, err error) {
	query = ""
	limit, err := url.QueryUnescape(request.URL.Query().Get(PARAM_QUERY_LIMIT))
	if err != nil {
		log.Error("Can't unescape %v url data", err)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return "", errors.New("Wrong data")
	}
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

		param_low := limits[PARAM_LIMIT_LOW]
		param_high := limits[PARAM_LIMIT_HIGH]
		offset, err = strconv.ParseInt(param_low, 0, 64)
		if err != nil {
			log.Error("Wrong limit offset %v", param_low)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return "", err
		}
		if offset < 0 {
			log.Error("Wrong limit offset %v", param_low)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return "", errors.New("Wrong offset")
		}

		count, err = strconv.ParseInt(param_high, 0, 64)
		if err != nil {
			log.Error("Wrong limit count %v", param_high)
			r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return "", err
		}

		if count > 0 {
			query = " limit " + param_low + ", " + param_high
		}
	} else {
		query = " limit 0, 100"
	}

	return query, nil
}

func GetOrderArray(checker models.Checker, request *http.Request, r render.Render, language string) (sorts *[]models.OrderExp, err error) {
	order, err := url.QueryUnescape(request.URL.Query().Get(PARAM_QUERY_ORDER))
	if err != nil {
		log.Error("Can't unescape %v url data", err)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, errors.New("Wrong data")
	}
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
			param_field := elements[PARAM_SORT_FIELD]
			param_order := elements[PARAM_SORT_ORDER]
			valid, err = checker.Check(param_field)
			if !valid || err != nil {
				log.Error("Unknown field name %v", param_field)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return nil, errors.New("Unknown field")
			}

			if strings.ToLower(param_order) != PARAM_SORT_ASC &&
				strings.ToLower(param_order) != PARAM_SORT_DESC {
				log.Error("Unknown sort operation %v", param_order)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return nil, errors.New("Unknown sort")
			}
			*sorts = append(*sorts, models.OrderExp{Field: param_field, Order: param_order})
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
	filter := strings.Join(ParseRawQuery(request.URL.RawQuery)[PARAM_QUERY_FILTER], ",")
	filters = new([]models.FilterExp)
	if filter != "" {
		cons := strings.Split(filter, ",")
		for _, element := range cons {
			elements := strings.Split(element, ":")
			if len(elements) < PARAM_QUERY_NUMBER {
				log.Error("Wrong number of filter parameter elements %v", len(elements))
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return nil, errors.New("Wrong parameter number")
			}
			for i := range elements {
				elements[i], err = url.QueryUnescape(elements[i])
				if err != nil {
					log.Error("Can't unescape %v url data", err)
					r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[language].Errors.Api.Data_Wrong})
					return nil, errors.New("Wrong data")
				}
			}

			var allfields bool
			var field string
			var value string

			param_field := elements[PARAM_FILTER_FIELD]
			param_value := strings.Join(elements[PARAM_FILTER_VALUE:], ":")
			param_op := elements[PARAM_FILTER_OP]

			allfields = false
			if param_field == "*" {
				allfields = true
			}

			if allfields {
				if strings.Contains(param_value, "'") {
					log.Error("Wrong field value %v for %v", param_value, param_field)
					r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[language].Errors.Api.Data_Wrong})
					return nil, errors.New("Wrong value")
				}
				field = ""
				value = "'" + param_value + "'"
			} else {
				var errField error
				var errValue error

				field, value, errField, errValue = extractor.Extract(param_field, param_value)
				if errField != nil {
					log.Error("Unknown field name %v", param_field)
					r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[language].Errors.Api.Data_Wrong})
					return nil, errField
				}
				if errValue != nil {
					log.Error("Wrong field value %v for %v", param_value, param_field)
					r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[language].Errors.Api.Data_Wrong})
					return nil, errValue
				}
			}

			op := ""
			switch strings.ToLower(param_op) {
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
				log.Error("Unknown filter operation %v", param_op)
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return nil, errors.New("Unknown filter")
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
