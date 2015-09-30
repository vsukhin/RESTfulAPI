package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"regexp"
	"strings"
	"types"
)

const (
	PARAM_NAME_COMPANY_ID = "orgid"
)

func CheckCompany(r render.Render, params martini.Params, companyrepository services.CompanyRepository,
	language string) (dtocompany *models.DtoCompany, err error) {
	company_id, err := CheckParameterInt(r, params[PARAM_NAME_COMPANY_ID], language)
	if err != nil {
		return nil, err
	}

	dtocompany, err = companyrepository.Get(company_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return dtocompany, nil
}

func CheckCompanyAvailability(company_id int64, user_id int64, r render.Render, companyrepository services.CompanyRepository,
	language string) (dtocompany *models.DtoCompany, err error) {
	if company_id != 0 {
		dtocompany, err = companyrepository.Get(company_id)
	} else {
		dtocompany, err = companyrepository.GetPrimaryByUser(user_id)
	}
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtocompany.Active {
		log.Error("Company is not active %v", dtocompany.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Company is not active")
	}
	allowed, err := companyrepository.CheckUserAccess(user_id, dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !allowed {
		log.Error("Company %v is not accessible for user %v", company_id, user_id)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Company is not accessible")
	}

	return dtocompany, nil
}

func CheckCompanyType(companytype_id int, r render.Render, companytyperepository services.CompanyTypeRepository,
	language string) (dtocompanytype *models.DtoCompanyType, err error) {
	dtocompanytype, err = companytyperepository.Get(companytype_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtocompanytype.Active {
		log.Error("Company type is not active %v", dtocompanytype.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Company type is not active")
	}

	return dtocompanytype, nil
}

func CheckCompanyClass(companyclass_id int, r render.Render, companyclassrepository services.CompanyClassRepository,
	language string) (dtocompanyclass *models.DtoCompanyClass, err error) {
	dtocompanyclass, err = companyclassrepository.Get(companyclass_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtocompanyclass.Active {
		log.Error("Company class is not active %v", dtocompanyclass.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Company class is not active")
	}

	return dtocompanyclass, nil
}

func CheckAddressType(addreetype_id int, r render.Render, addresstyperepository services.AddressTypeRepository,
	language string) (dtoaddresstype *models.DtoAddressType, err error) {
	dtoaddresstype, err = addresstyperepository.Get(addreetype_id)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtoaddresstype.Active {
		log.Error("Address type is not active %v", dtoaddresstype.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Address type is not active")
	}

	return dtoaddresstype, nil
}

func FillCompany(viewcompany *models.ViewCompany, dtocompany *models.DtoCompany, r render.Render,
	companytyperepository services.CompanyTypeRepository, companyclassrepository services.CompanyClassRepository,
	addresstyperepository services.AddressTypeRepository, language string) (err error) {
	_, err = CheckCompanyType(viewcompany.Company_Type_ID, r, companytyperepository, language)
	if err != nil {
		return err
	}
	dtocompany.Company_Type_ID = viewcompany.Company_Type_ID
	dtocompany.Primary = viewcompany.Primary
	dtocompany.FullName_Rus = viewcompany.FullName_Rus
	dtocompany.FullName_Eng = viewcompany.FullName_Eng
	dtocompany.ShortName_Rus = viewcompany.ShortName_Rus
	dtocompany.ShortName_Eng = viewcompany.ShortName_Eng
	dtocompany.Resident = viewcompany.Resident
	dtocompany.VAT = viewcompany.VAT

	codeclasses := make(map[int]int)
	for _, viewcode := range viewcompany.CompanyCodes {
		for _, code := range strings.Split(viewcode.Codes, ",") {
			dtocompanycode := new(models.DtoCompanyCode)
			companyclass, err := CheckCompanyClass(viewcode.Company_Class_ID, r, companyclassrepository, language)
			if err != nil {
				return err
			}
			dtocompanycode.Company_Class_ID = companyclass.ID
			_, ok := codeclasses[companyclass.ID]
			if !ok {
				codeclasses[companyclass.ID] = 0
			}
			codeclasses[companyclass.ID]++
			if codeclasses[companyclass.ID] > 1 && !companyclass.Multiple {
				log.Error("Code class could not have multiple values %v", companyclass.ID)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return errors.New("Multiple codes not allowed")
			}
			if len([]rune(code)) > models.CODE_FIELD_MAX_LENGTH_VALUE {
				log.Error("Wrong length of code field %v", len([]rune(code)))
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return errors.New("Wrong code length")
			}
			if companyclass.Format != "" {
				valid, err := regexp.MatchString(companyclass.Format, code)
				if err != nil {
					log.Error("Error during running reg exp %v with value %v", err, companyclass.Format)
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[language].Errors.Api.Data_Wrong})
					return errors.New("Wrong reg exp")
				}
				if !valid {
					log.Error("Error during checking code %v with format %v", code, companyclass.Format)
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[language].Errors.Api.Data_Wrong})
					return errors.New("Wrong value")
				}
			}
			if code == "" && companyclass.Required {
				log.Error("Code could not be empty for class %v", companyclass.ID)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return errors.New("Empty code")
			}
			dtocompanycode.Code = code
			dtocompany.CompanyCodes = append(dtocompany.CompanyCodes, *dtocompanycode)
		}
	}

	primaryaddresses := make(map[int]int)
	legaladdresses := 0
	for _, viewaddress := range viewcompany.CompanyAddresses {
		dtocompanyaddress := new(models.DtoCompanyAddress)
		dtocompanyaddress.Primary = viewaddress.Primary
		if viewaddress.Ditto != 0 {
			dittotype, err := CheckAddressType(viewaddress.Ditto, r, addresstyperepository, language)
			if err != nil {
				return err
			}
			dtocompanyaddress.Ditto = dittotype.ID
		}
		addresstype, err := CheckAddressType(viewaddress.Address_Type_ID, r, addresstyperepository, language)
		if err != nil {
			return err
		}
		dtocompanyaddress.Address_Type_ID = addresstype.ID
		if addresstype.ID == models.ADDRESS_TYPE_LEGAL {
			legaladdresses++
		}
		_, ok := primaryaddresses[addresstype.ID]
		if !ok {
			primaryaddresses[addresstype.ID] = 0
		}
		if dtocompanyaddress.Primary {
			primaryaddresses[addresstype.ID]++
		}
		if viewaddress.Ditto == 0 {
			if addresstype.Required {
				if viewaddress.Full == "" {
					if viewaddress.Zip == "" {
						log.Error("Zip field could not be empty")
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
							Message: config.Localization[language].Errors.Api.Data_Wrong})
						return errors.New("Empty zip")
					}
					if viewaddress.Country == "" {
						log.Error("Country field could not be empty")
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
							Message: config.Localization[language].Errors.Api.Data_Wrong})
						return errors.New("Empty country")
					}
					if viewaddress.City == "" {
						log.Error("City field could not be empty")
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
							Message: config.Localization[language].Errors.Api.Data_Wrong})
						return errors.New("Empty city")
					}
					if viewaddress.Street == "" {
						log.Error("Street field could not be empty")
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
							Message: config.Localization[language].Errors.Api.Data_Wrong})
						return errors.New("Empty street")
					}
					if viewaddress.Building == "" {
						log.Error("Building field could not be empty")
						r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
							Message: config.Localization[language].Errors.Api.Data_Wrong})
						return errors.New("Empty building")
					}
				}
			}
			dtocompanyaddress.Full = viewaddress.Full
			dtocompanyaddress.Zip = viewaddress.Zip
			dtocompanyaddress.Country = viewaddress.Country
			dtocompanyaddress.City = viewaddress.City
			dtocompanyaddress.Street = viewaddress.Street
			dtocompanyaddress.Building = viewaddress.Building
			dtocompanyaddress.Region = viewaddress.Region
			dtocompanyaddress.Postbox = viewaddress.Postbox
			dtocompanyaddress.Company = viewaddress.Company
		}
		dtocompanyaddress.Comments = viewaddress.Comments
		dtocompany.CompanyAddresses = append(dtocompany.CompanyAddresses, *dtocompanyaddress)
	}
	if legaladdresses != 1 {
		log.Error("Exactly one legal address must be provided %v", legaladdresses)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return errors.New("Wrong legal address")
	}
	for addresstype_id, count := range primaryaddresses {
		if count != 1 {
			log.Error("Exactly one primary address must be provided for %v, %v", addresstype_id, count)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return errors.New("Wrong primary count")
		}
	}

	primarybanks := 0
	for _, viewbank := range viewcompany.CompanyBanks {
		dtocompanybank := new(models.DtoCompanyBank)
		dtocompanybank.Primary = viewbank.Primary
		if viewbank.Primary && !viewbank.Deleted {
			primarybanks++
		}
		dtocompanybank.Bik = viewbank.Bik
		dtocompanybank.Name = viewbank.Name
		dtocompanybank.CheckingAccount = viewbank.CheckingAccount
		dtocompanybank.CorrespondingAccount = viewbank.CorrespondingAccount
		dtocompanybank.Active = !viewbank.Deleted
		dtocompany.CompanyBanks = append(dtocompany.CompanyBanks, *dtocompanybank)
	}
	if primarybanks != 1 {
		log.Error("Exactly one primary bank must be provided %v", primarybanks)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return errors.New("Wrong primary banks")
	}

	ceocount := 0
	accountantcount := 0
	for _, viewemployee := range viewcompany.CompanyStaff {
		dtocompanyemployee := new(models.DtoCompanyEmployee)
		if viewemployee.Employee_Type != models.EMPLOYEE_TYPE_CEO && viewemployee.Employee_Type != models.EMPLOYEE_TYPE_ACCOUNTANT {
			log.Error("Unknown employee type %v", viewemployee.Employee_Type)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
				Message: config.Localization[language].Errors.Api.Data_Wrong})
			return errors.New("Unknown employee type")
		}
		dtocompanyemployee.Employee_Type = viewemployee.Employee_Type
		if viewemployee.Employee_Type == models.EMPLOYEE_TYPE_CEO && !viewemployee.Deleted {
			ceocount++
		}
		if viewemployee.Employee_Type == models.EMPLOYEE_TYPE_ACCOUNTANT && !viewemployee.Deleted {
			accountantcount++
		}
		if viewemployee.Ditto != "" {
			if viewemployee.Ditto != models.EMPLOYEE_TYPE_CEO && viewemployee.Ditto != models.EMPLOYEE_TYPE_ACCOUNTANT {
				log.Error("Unknown employee ditto %v", viewemployee.Ditto)
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return errors.New("Unknown employee ditto")
			}
			dtocompanyemployee.Ditto = viewemployee.Ditto
		} else {
			if viewemployee.Surname == "" {
				log.Error("Surname field could not be empty")
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return errors.New("Empty surname")
			}
			dtocompanyemployee.Surname = viewemployee.Surname
			if viewemployee.Name == "" {
				log.Error("Name field could not be empty")
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return errors.New("Empty name")
			}
			dtocompanyemployee.Name = viewemployee.Name
			if viewemployee.MiddleName == "" {
				log.Error("Middlename field could not be empty")
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return errors.New("Empty middlename")
			}
			dtocompanyemployee.MiddleName = viewemployee.MiddleName
			if viewemployee.Base == "" {
				log.Error("Base field could not be empty")
				r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return errors.New("Empty base")
			}
			dtocompanyemployee.Base = viewemployee.Base
		}
		dtocompanyemployee.Active = !viewemployee.Deleted
		dtocompany.CompanyStaff = append(dtocompany.CompanyStaff, *dtocompanyemployee)
	}
	if ceocount != 1 {
		log.Error("Exactly one ceo must be provided %v", ceocount)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return errors.New("Wrong ceo count")
	}
	if accountantcount != 1 {
		log.Error("Exactly one accountant must be provided %v", accountantcount)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return errors.New("Wrong accountant count")
	}

	return nil
}

func LoadCompany(dtocompany *models.DtoCompany, r render.Render, companycoderepository services.CompanyCodeRepository,
	companyaddressrepository services.CompanyAddressRepository, companybankrepository services.CompanyBankRepository,
	companyemployeerepository services.CompanyEmployeeRepository, language string) (apicompany *models.ApiMiddleCompany, err error) {
	companycodes, err := companycoderepository.GetByCompany(dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	companyaddresses, err := companyaddressrepository.GetByCompany(dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	companybanks, err := companybankrepository.GetByCompany(dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	companystaff, err := companyemployeerepository.GetByCompany(dtocompany.ID)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}

	return models.NewApiMiddleCompany(dtocompany.Primary, dtocompany.Company_Type_ID,
		dtocompany.FullName_Rus, dtocompany.FullName_Eng, dtocompany.ShortName_Rus, dtocompany.ShortName_Eng, dtocompany.Resident,
		*companycodes, *companyaddresses, *companybanks, *companystaff, dtocompany.VAT, dtocompany.Locked, !dtocompany.Active), nil
}

func PrepareCompanyTemplate(company_id int64, apicompany *models.ApiMiddleCompany, r render.Render,
	language string) (dtocompanytemplate *models.DtoCompanyTemplate, err error) {
	dtocompanytemplate = new(models.DtoCompanyTemplate)
	dtocompanytemplate.Name = apicompany.FullName_Rus
	for _, company_class_id := range []int{models.CODE_TYPE_INN, models.CODE_TYPE_KPP} {
		found := false
		for _, companycode := range apicompany.CompanyCodes {
			if companycode.Company_Class_ID == company_class_id {
				found = true
				switch company_class_id {
				case models.CODE_TYPE_INN:
					dtocompanytemplate.INN = companycode.Codes
				case models.CODE_TYPE_KPP:
					dtocompanytemplate.KPP = companycode.Codes
				}
				break
			}
		}
		if !found {
			log.Error("Company class is not found %v with value %v", company_class_id, company_id)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, errors.New("Not found company class")
		}
	}
	found := false
	for _, companyaddress := range apicompany.CompanyAddresses {
		if companyaddress.Primary && companyaddress.Address_Type_ID == models.ADDRESS_TYPE_LEGAL {
			if companyaddress.Full != "" {
				found = true
				dtocompanytemplate.Address = companyaddress.Full
				break
			} else if companyaddress.Country != "" && companyaddress.Zip != "" && companyaddress.City != "" && companyaddress.Street != "" &&
				companyaddress.Building != "" {
				found = true
				dtocompanytemplate.Address = companyaddress.Country + ", " + companyaddress.Zip + ", "
				if companyaddress.Region != "" {
					dtocompanytemplate.Address += companyaddress.Region + ", "
				}
				dtocompanytemplate.Address += companyaddress.City + ", " + companyaddress.Street + ", " + companyaddress.Building
				break
			}
		}
	}
	if !found {
		log.Error("Company address is not found %v with value %v", models.ADDRESS_TYPE_LEGAL, company_id)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Not found company address")
	}
	for _, employee_type := range []string{models.EMPLOYEE_TYPE_CEO, models.EMPLOYEE_TYPE_ACCOUNTANT} {
		found := false
		for _, companyemployee := range apicompany.CompanyStaff {
			if companyemployee.Employee_Type == employee_type && !companyemployee.Deleted {
				found = true
				fio := ""
				for _, name := range []string{companyemployee.Name, companyemployee.MiddleName} {
					runes := []rune(name)
					if len(runes) != 0 {
						fio += " " + string(runes[0]) + "."
					}
				}
				switch employee_type {
				case models.EMPLOYEE_TYPE_CEO:
					dtocompanytemplate.CEO = companyemployee.Surname + fio
					if companyemployee.Base != "" {
						dtocompanytemplate.CEO += " (" + companyemployee.Base + ")"
					}
				case models.EMPLOYEE_TYPE_ACCOUNTANT:
					dtocompanytemplate.Accountant = companyemployee.Surname + fio
					if companyemployee.Base != "" {
						dtocompanytemplate.Accountant += " (" + companyemployee.Base + ")"
					}
				}
				break
			}
		}
		if !found {
			log.Error("Company employee is not found %v with value %v", employee_type, company_id)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, errors.New("Not found company employee")
		}
	}

	return dtocompanytemplate, nil
}
