package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"errors"
	"github.com/martini-contrib/render"
	"net/http"
	"types"
)

const (
	PARAM_NAME_CLASSIFIER_ID = "id"
)

func CheckClassifier(classifierid int, r render.Render, classifierrepository services.ClassifierRepository,
	language string) (dtoclassifier *models.DtoClassifier, err error) {
	dtoclassifier, err = classifierrepository.Get(classifierid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	if !dtoclassifier.Active {
		log.Error("Classifier is not active %v", dtoclassifier.ID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, errors.New("Classifier not active")
	}

	return dtoclassifier, nil
}

func CheckPrimaryEmail(user *models.ViewApiUserFull, language string, r render.Render) (err error) {
	count := 0
	for _, checkEmail := range user.Emails {
		if checkEmail.Primary {
			if user.Confirmed != checkEmail.Confirmed {
				log.Error("Confirmation statuses for user and email are different")
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return errors.New("Mismatched statuses")
			}
			count++
		}
	}
	if count != 1 {
		log.Error("Only one primary email is allowed")
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return errors.New("Wrong primary emails amount")
	}

	return nil
}

func CheckPrimaryMobilePhone(user *models.ViewApiUserFull, language string, r render.Render) (err error) {
	count := 0
	for _, checkMobilePhone := range user.MobilePhones {
		if checkMobilePhone.Primary {
			if user.Confirmed != checkMobilePhone.Confirmed {
				log.Error("Confirmation statuses for user and mobile phone are different")
				r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
					Message: config.Localization[language].Errors.Api.Data_Wrong})
				return errors.New("Mismatched statuses")
			}
			count++
		}
	}
	if count != 1 {
		log.Error("Only one primary mobile phone is allowed")
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return errors.New("Wrong primary mobile phones amount")
	}

	return nil
}

func CheckEmailAvailability(value string, language string, r render.Render,
	emailrepository services.EmailRepository) (emailExists bool, err error) {
	emailExists, err = emailrepository.Exists(value)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return emailExists, err
	}

	if emailExists {
		email, err := emailrepository.Get(value)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return emailExists, err
		}
		if email.Confirmed {
			log.Error("Email exists in database %v", value)
			r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_EMAIL_INUSE,
				Message: config.Localization[language].Errors.Api.Email_InUse})
			return emailExists, errors.New("Email exists")
		}
	}

	return emailExists, nil
}

func CheckMobilePhoneAvailability(value string, language string, r render.Render,
	mobilephonerepository services.MobilePhoneRepository) (phoneExists bool, err error) {
	phoneExists, err = mobilephonerepository.Exists(value)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return phoneExists, err
	}

	if phoneExists {
		mobilephone, err := mobilephonerepository.Get(value)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return phoneExists, err
		}
		if mobilephone.Confirmed {
			log.Error("Mobile phone exists in database %v", value)
			r.JSON(http.StatusConflict, types.Error{Code: types.TYPE_ERROR_MOBILEPHONE_INUSE,
				Message: config.Localization[language].Errors.Api.MobilePhone_InUse})
			return phoneExists, errors.New("Mobile phone exists")
		}
	}

	return phoneExists, nil
}

func SendConfirmations(dtouser *models.DtoUser, session *models.DtoSession, request *http.Request, r render.Render,
	emailrepository services.EmailRepository, templaterepository services.TemplateRepository, sendpassword bool) (err error) {
	for _, confEmail := range *dtouser.Emails {
		if !confEmail.Confirmed {
			if confEmail.Primary && sendpassword {
				err = SendPasswordRegistration(confEmail.Language, &confEmail, dtouser, request, r, emailrepository, templaterepository)
				if err != nil {
					return err
				}
			} else {
				subject := ""
				if confEmail.Primary {
					subject = config.Localization[confEmail.Language].Messages.RegistrationSubject
				} else {
					subject = config.Localization[confEmail.Language].Messages.EmailSubject
				}

				buf, err := templaterepository.GenerateText(models.NewDtoCodeTemplate(models.NewDtoTemplate(confEmail.Email, confEmail.Language,
					request.Host), confEmail.Code), services.TEMPLATE_EMAIL_CONFIRMATION, services.TEMPLATE_LAYOUT)
				if err != nil {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
						Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
					return err
				}

				err = emailrepository.SendHTML(confEmail.Email, subject, buf.String(), "", "")
				if err != nil {
					r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
						Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
					return err
				}
			}
		}
	}

	return nil
}

func GetUnitDependences(unitid int64, r render.Render, userrepository services.UserRepository,
	customertablerepository services.CustomerTableRepository, projectrepository services.ProjectRepository, orderrepository services.OrderRepository,
	facilityrepository services.FacilityRepository, companyrepository services.CompanyRepository, smssenderrepository services.SMSSenderRepository,
	invoicerepository services.InvoiceRepository, language string) (users *[]models.ApiUserTiny, tables *[]models.ApiMiddleCustomerTable,
	projects *[]models.ApiShortProject, orders *[]models.ApiBriefOrder, facilities *[]models.ApiLongFacility,
	companies *[]models.ApiShortCompany, smssenders *[]models.ApiLongSMSSender, invoices *[]models.ApiShortInvoice, err error) {
	users, err = userrepository.GetByUnit(unitid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	tables, err = customertablerepository.GetByUnit(unitid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	projects, err = projectrepository.GetByUnit(unitid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	orders, err = orderrepository.GetByUnit(unitid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	facilities, err = facilityrepository.GetByUnit(unitid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	companies, err = companyrepository.GetByUnit(unitid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	smssenders, err = smssenderrepository.GetByUnit(unitid)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	invoices, err = invoicerepository.GetByUnit(unitid, "")
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return users, tables, projects, orders, facilities, companies, smssenders, invoices, nil
}
