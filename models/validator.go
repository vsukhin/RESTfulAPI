package models

import (
	"application/config"
	"fmt"
	"github.com/martini-contrib/binding"
	logging "github.com/op/go-logging"
	"gopkg.in/validator.v2"
	"net/http"
	"strconv"
	"strings"
	"types"
)

var (
	log config.Logger = logging.MustGetLogger("models")
)

const (
	VALIDATE_FIELD_EMPTY = iota + 1
	VALIDATE_FIELD_SMALL
	VALIDATE_FIELD_BIG
	VALIDATE_LANGUAGE_WRONG
	VALIDATE_FIELD_REGEXP
)

func InitLogger(logger config.Logger) {
	log = logger
}

func Validate(object interface{}, errors binding.Errors, req *http.Request) binding.Errors {
	err := validator.Validate(object)
	if err != nil {
		errormap := err.(validator.ErrorMap)
		for f, errarray := range errormap {
			for _, e := range errarray {
				code := ""
				value := ""
				switch e {
				case validator.ErrZeroValue:
					code = strconv.Itoa(VALIDATE_FIELD_EMPTY)
				case validator.ErrMin:
					code = strconv.Itoa(VALIDATE_FIELD_SMALL)
				case validator.ErrMax:
					code = strconv.Itoa(VALIDATE_FIELD_BIG)
				case validator.ErrRegexp:
					code = strconv.Itoa(VALIDATE_FIELD_REGEXP)
				}
				fieldname := GetJsonTag(f, object)
				if fieldname == "" {
					fieldname = f
				}
				v, found := GetFieldValue(fieldname, object)
				if found {
					value = fmt.Sprintf("%v", v)
				}
				log.Error("Error during validation field %s with value %v (%v)", f, v, e)
				errors = append(errors, binding.Error{
					FieldNames:     []string{fieldname},
					Classification: code,
					Message:        value,
				})
			}
		}
	}

	return errors
}

func ValidateWithLanguage(object interface{}, errors binding.Errors, req *http.Request, language string) binding.Errors {
	if language != "" {
		found := false
		for _, lang := range config.Configuration.Server.AvailableLanguages {
			if strings.ToLower(language) == strings.ToLower(lang) {
				found = true
				break
			}
		}
		if !found {
			log.Error("Error during looking up for existing language %v", language)
			errors = append(errors, binding.Error{
				FieldNames:     []string{"language"},
				Classification: strconv.Itoa(VALIDATE_LANGUAGE_WRONG),
			})
		}
	}

	return Validate(object, errors, req)
}

func ConvertErrors(language string, errors binding.Errors) (fielderrors *[]types.FieldError, err int) {
	fielderrors = new([]types.FieldError)

	for _, bindingerror := range errors {
		fielderror := new(types.FieldError)

		if len(bindingerror.FieldNames) > 0 {
			fielderror.Field = strings.Join(bindingerror.FieldNames, ",")
			fielderror.FieldValue = bindingerror.Message
			message := ""
			coderr, errconv := strconv.Atoi(bindingerror.Classification)
			if errconv == nil {
				switch coderr {
				case VALIDATE_FIELD_EMPTY:
					message = config.Localization[language].Errors.Binding.Field_Empty
				case VALIDATE_FIELD_SMALL:
					message = config.Localization[language].Errors.Binding.Field_Small
				case VALIDATE_FIELD_BIG:
					message = config.Localization[language].Errors.Binding.Field_Big
				case VALIDATE_FIELD_REGEXP:
					message = config.Localization[language].Errors.Binding.Field_Regexp
				case VALIDATE_LANGUAGE_WRONG:
					return nil, types.TYPE_ERROR_LANGUAGE_NOTSUPPORTED
				}
			} else {
				message = bindingerror.Message
			}
			fielderror.Message = message
		} else {
			fielderror.Message = bindingerror.Message
			log.Error("Error during object binding %v", fielderror.Message)
		}
		*fielderrors = append(*fielderrors, *fielderror)
	}

	return fielderrors, types.TYPE_ERROR_NONE
}
