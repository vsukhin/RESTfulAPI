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

const (
	NEWS_NUMBER                 = 10
	NEWS_VERSION                = "2.0"
	METHOD_NAME_SUBSCRIPTION    = "/api/v1.0/subscriptions/news/:email/"
	METHOD_TIMEOUT_SUBSCRIPTION = time.Second
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
	requestrepository services.RequestRepository, sessionrepository services.SessionRepository) {
	var user_id int64 = 0
	var language = config.Configuration.Server.DefaultLanguage

	session, _, err := sessionrepository.GetAndSaveSession(request, r, nil, false, false, true)
	if err == nil {
		user_id = session.UserID
		language = session.Language
	}

	if user_id == 0 {
		if helpers.CheckFrequence(METHOD_NAME_SUBSCRIPTION, METHOD_TIMEOUT_SUBSCRIPTION, request, r, requestrepository, language) != nil {
			return
		}
	}

	email, err := url.QueryUnescape(params[helpers.PARAMETER_NAME_SUBSCRIBТION_EMAIL])
	if err != nil {
		log.Error("Can't unescape %v url data", err)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return
	}
	if email == "" || len([]rune(email)) > helpers.PARAM_LENGTH_MAX {
		log.Error("Wrong parameter length %v", email)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return
	}

	found, err := subscriptionrepository.Exists(email)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return
	}

	confirmed := false
	if found {
		dtosubscription, err := subscriptionrepository.Get(email)
		if err != nil {
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return
		}
		confirmed = dtosubscription.Confirmed
	}

	valid := strings.Contains(email, "@")

	r.JSON(http.StatusOK, models.NewApiLongSubscription(email, confirmed, valid))
}

// post /api/v1.0/subscriptions/news/
func CreateSubscription(errors binding.Errors, viewsubscription models.ViewSubscription, request *http.Request, r render.Render,
	subscriptionrepository services.SubscriptionRepository, captcharepository services.CaptchaRepository,
	sessionrepository services.SessionRepository, emailrepository services.EmailRepository, templaterepository services.TemplateRepository,
	accesslogrepository services.AccessLogRepository) {
	if helpers.CheckValidation(errors, r, config.Configuration.Server.DefaultLanguage) != nil {
		return
	}
	if helpers.Check(viewsubscription.CaptchaHash, viewsubscription.CaptchaValue, r, captcharepository) != nil {
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

	dtoaccesslog, err := helpers.CreateAccessLog(request.RequestURI, request, r, accesslogrepository, config.Configuration.Server.DefaultLanguage)
	if err != nil {
		return
	}
	dtosubscription.Subscr_AccessLog_ID = dtoaccesslog.ID

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
	host := request.Header.Get(helpers.REQUEST_HEADER_X_FORWARDED_FOR)
	if host == "" {
		host, _, err = net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			log.Error("Can't detect ip address %v from %v", err, request.RemoteAddr)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
			return
		}
	}
	subject := config.Localization[dtosubscription.Language].Messages.SubscriptionSubject
	buf, err := templaterepository.GenerateText(models.NewDtoDualCodeTemplate(models.NewDtoTemplate(dtosubscription.Email, dtosubscription.Language,
		request.Host, time.Now(), host), dtosubscription.Subscr_Code, dtosubscription.Unsubscr_Code), services.TEMPLATE_SUBSCRIPTION,
		services.TEMPLATE_DIRECTORY_EMAILS, "")
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Object_NotExist})
		return
	}

	headers := "List-Unsubscribe: <http://" + request.Host + "/subscriptions/unsubscribe/" + dtosubscription.Unsubscr_Code + "/>"
	err = emailrepository.SendHTML(dtosubscription.Email, subject, buf.String(), headers, "")
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[config.Configuration.Server.DefaultLanguage].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, models.NewApiMiddleSubscription(config.Configuration.Mail.Sender, dtosubscription.Email))
}

// patch /api/v1.0/subscriptions/news/
func ConfirmSubscription(errors binding.Errors, confirm models.SubscriptionConfirm, request *http.Request, r render.Render,
	subscriptionrepository services.SubscriptionRepository, accesslogrepository services.AccessLogRepository) {
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

	dtoaccesslog, err := helpers.CreateAccessLog(request.RequestURI, request, r, accesslogrepository, config.Configuration.Server.DefaultLanguage)
	if err != nil {
		return
	}
	dtosubscription.Conf_AccessLog_ID = dtoaccesslog.ID

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

// delete /api/v1.0/subscriptions/news/:email/
func DeleteSubscription(r render.Render, params martini.Params, subscriptionrepository services.SubscriptionRepository,
	emailrepository services.EmailRepository, session *models.DtoSession) {
	email, err := url.QueryUnescape(params[helpers.PARAMETER_NAME_SUBSCRIBТION_EMAIL])
	if err != nil {
		log.Error("Can't unescape %v url data", err)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}
	if email == "" || len([]rune(email)) > helpers.PARAM_LENGTH_MAX {
		log.Error("Wrong parameter length %v", email)
		r.JSON(http.StatusBadRequest, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	found, err := subscriptionrepository.Exists(email)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if !found {
		log.Error("Can't find email in subscription %v", email)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	dtoemail, err := emailrepository.Get(email)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}
	if dtoemail.UserID != session.UserID {
		log.Error("Email %v doesn't belong to user %v", email, session.UserID)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	err = subscriptionrepository.Delete(email)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[session.Language].Errors.Api.Data_Wrong})
		return
	}

	r.JSON(http.StatusOK, types.ResponseOK{Message: config.Localization[session.Language].Messages.OK})
}

// get /api/v1.0/news/
func GetNews(w http.ResponseWriter, request *http.Request, r render.Render, newsrepository services.NewsRepository, session *models.DtoSession) {
	query := ""
	var filters *[]models.FilterExp
	filters, err := helpers.GetFilterArray(new(models.NewsSearch), nil, request, r, session.Language)
	if err != nil {
		return
	}
	if len(*filters) != 0 {
		var masks []string
		for _, filter := range *filters {
			var exps []string
			for _, field := range filter.Fields {
				exps = append(exps, "`"+field+"` "+filter.Op+" "+filter.Value)
			}
			masks = append(masks, "("+strings.Join(exps, " or ")+")")
		}
		query += " and "
		query += strings.Join(masks, " and ")
	}

	var sorts *[]models.OrderExp
	sorts, err = helpers.GetOrderArray(new(models.NewsSearch), request, r, session.Language)
	if err != nil {
		return
	}
	if len(*sorts) != 0 {
		var orders []string
		for _, sort := range *sorts {
			orders = append(orders, " `"+sort.Field+"` "+sort.Order)
		}
		query += " order by"
		query += strings.Join(orders, ",")
	}

	var limit string
	limit, err = helpers.GetLimitQuery(request, r, session.Language)
	if err != nil {
		return
	}
	query += limit

	news, err := newsrepository.GetAny(query)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[session.Language].Errors.Api.Object_NotExist})
		return
	}

	helpers.RenderJSONArray(news, len(*news), w, r)
}
