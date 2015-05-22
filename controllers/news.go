package controllers

import (
	"application/config"
	"application/helpers"
	"application/models"
	"application/services"
	"encoding/xml"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	"types"
)

// get /subscriptions/news/rss/
func GetNewsRss(w http.ResponseWriter, request *http.Request, r render.Render, newsrepository services.NewsRepository) {
	news, err := newsrepository.GetAll(config.Configuration.Server.DefaultLanguage, NEWS_NUMBER)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	for i := range *news {
		(*news)[i].URL = "http://" + request.Host
	}
	apichannel := models.NewApiChannel(config.Localization[config.Configuration.Server.DefaultLanguage].Messages.NewsHeader,
		config.Localization[config.Configuration.Server.DefaultLanguage].Messages.NewsHeader, "http://"+request.Host,
		config.Configuration.Server.DefaultLanguage, news)
	apiRSSFeed := models.NewApiRSSFeed(NEWS_VERSION, &([]models.ApiChannel{*apichannel}))

	data, err := xml.Marshal(apiRSSFeed)
	if err == nil {
		w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
		w.Write([]byte("<?xml version=\"1.0\"?>"))
		w.Write(data)
	} else {
		log.Error("Can't marshal xml data %v", err)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
}

// get /api/v1.0/subscriptions/news/:email/
func GetNewsSubscription(request *http.Request, r render.Render, params martini.Params, subscriptionrepository services.SubscriptionRepository,
	requestrepository services.RequestRepository) {
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	exists, err := requestrepository.Exists(host, SUBSCRIPTION_METHOD_NAME)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	var hits int64 = 0
	if exists {
		request, err := requestrepository.Get(host, SUBSCRIPTION_METHOD_NAME)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
			return
		}
		if time.Now().Sub(request.LastUpdated) < time.Second {
			log.Error("Requests are too frequent for method %v", request.Method)
			r.JSON(http.StatusServiceUnavailable, types.Error{Code: types.TYPE_ERROR_REQUEST_TOOFREQUENT,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Request_Too_Often})
			return
		}
		hits = request.Hits
	}
	dtorequest := models.NewDtoRequest(host, SUBSCRIPTION_METHOD_NAME, time.Now(), hits+1)
	err = requestrepository.Save(dtorequest)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	email, err := url.QueryUnescape(params[helpers.PARAMETER_NAME_SUBSCRIBТION_EMAIL])
	if err != nil {
		log.Error("Can't unescape %v url data", err)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}
	if email == "" || len([]rune(email)) > helpers.PARAM_LENGTH_MAX {
		log.Error("Wrong parameter length %v", email)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	found, err := subscriptionrepository.Exists(email)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	confirmed := false
	if found {
		dtosubscription, err := subscriptionrepository.Get(email)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
			return
		}
		confirmed = dtosubscription.Confirmed
	}

	r.JSON(http.StatusOK, models.NewApiLongSubscription(email, confirmed, found))
}

// patch /api/v1.0/subscriptions/news/
func CreateSubscription(errors binding.Errors, viewsubscription models.ViewSubscription, request *http.Request, r render.Render,
	subscriptionrepository services.SubscriptionRepository, captcharepository services.CaptchaRepository,
	sessionrepository services.SessionRepository, emailrepository services.EmailRepository, templaterepository services.TemplateRepository) {
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	if captcharepository.Check(viewsubscription.CaptchaHash, viewsubscription.CaptchaValue, r) != nil {
		return
	}

	found, err := subscriptionrepository.Exists(viewsubscription.Email)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	if found {
		subscription, err := subscriptionrepository.Get(viewsubscription.Email)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
			return
		}
		if subscription.Confirmed {
			log.Error("Subscription is already confirmed for email %v", subscription.Email)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
			return
		}
	}
	var language string
	if viewsubscription.Language != "" {
		language = strings.ToLower(viewsubscription.Language)
	} else {
		language = config.Configuration.Server.DefaultLanguage
	}

	subscr_code, err := sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}
	unsubscr_code, err := sessionrepository.GenerateToken(helpers.TOKEN_LENGTH)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}
	dtosubscription := new(models.DtoSubscription)
	dtosubscription.Subscr_Code = subscr_code
	dtosubscription.Unsubscr_Code = unsubscr_code
	dtosubscription.Email = strings.ToLower(viewsubscription.Email)
	dtosubscription.Language = language
	dtosubscription.Confirmed = false
	dtosubscription.Subscr_Created = time.Now()
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	dtosubscription.Subscr_IP_Address = host
	dnses, err := net.LookupAddr(dtosubscription.Subscr_IP_Address)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	dtosubscription.Subscr_Reverse_DNS = strings.Join(dnses, ",")
	dtosubscription.Subscr_User_Agent = request.UserAgent()
	dtosubscription.Conf_Created = time.Now()
	if !found {
		err = subscriptionrepository.Create(dtosubscription)
	} else {
		err = subscriptionrepository.Update(dtosubscription)
	}
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	subject := config.Localization[dtosubscription.Language].Messages.SubscriptionSubject
	buf, err := templaterepository.GenerateText(models.NewDtoDualCodeTemplate(models.NewDtoTemplate(dtosubscription.Email, dtosubscription.Language,
		request.Host), dtosubscription.Subscr_Code, dtosubscription.Unsubscr_Code), services.TEMPLATE_SUBSCRIPTION, services.TEMPLATE_LAYOUT)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	headers := "List-Unsubscribe: <http://" + request.Host + "/subscriptions/unsubscribe/" + dtosubscription.Unsubscr_Code + "/>"
	err = emailrepository.SendEmail(dtosubscription.Email, subject, buf.String(), headers)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiMiddleSubscription(config.Configuration.Mail.Sender, dtosubscription.Email))
}

// patch /api/v1.0/subscriptions/news/
func ConfirmSubscription(errors binding.Errors, confirm models.SubscriptionConfirm, request *http.Request, r render.Render,
	subscriptionrepository services.SubscriptionRepository) {
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}

	dtosubscription, err := subscriptionrepository.FindBySubscrCode(confirm.Code)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	if dtosubscription.Confirmed {
		log.Error("Subscription is already confirmed for code %v", confirm.Code)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	dtosubscription.Confirmed = true
	dtosubscription.Conf_Created = time.Now()
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	dtosubscription.Conf_IP_Address = host
	dnses, err := net.LookupAddr(dtosubscription.Conf_IP_Address)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}
	dtosubscription.Conf_Reverse_DNS = strings.Join(dnses, ",")
	dtosubscription.Conf_User_Agent = request.UserAgent()
	err = subscriptionrepository.Update(dtosubscription)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiShortSubscription(dtosubscription.Email))
}

// get /subscriptions/unsubscribe/:unsubscribeCode/
func UnsubscribeFromNews(r render.Render, params martini.Params, subscriptionrepository services.SubscriptionRepository) {
	code := params[helpers.PARAMETER_NAME_UNSUBSCRIBE_CODE]
	if code == "" || len(code) > helpers.PARAM_LENGTH_MAX {
		log.Error("Wrong parameter length %v", code)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	dtosubscription, err := subscriptionrepository.FindByUnsubscrCode(code)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	err = subscriptionrepository.Delete(dtosubscription.Email)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[config.Configuration.Server.DefaultLanguage].Messages.OK})
}
