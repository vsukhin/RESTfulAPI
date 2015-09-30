package helpers

import (
	"application/config"
	"application/models"
	"application/services"
	"github.com/martini-contrib/render"
	"net"
	"net/http"
	"strings"
	"time"
	"types"
)

const (
	URL_LENGTH_MAX                 = 1000
	REQUEST_HEADER_X_FORWARDED_FOR = "X-Forwarded-For"
)

func CreateAccessLog(url string, request *http.Request, r render.Render, accesslogrepository services.AccessLogRepository,
	language string) (dtoaccesslog *models.DtoAccessLog, err error) {
	dtoaccesslog = new(models.DtoAccessLog)
	host := request.Header.Get(REQUEST_HEADER_X_FORWARDED_FOR)
	if host == "" {
		host, _, err = net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			log.Error("Can't detect ip address %v from %v", err, request.RemoteAddr)
			r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
				Message: config.Localization[language].Errors.Api.Object_NotExist})
			return nil, err
		}
	}
	dtoaccesslog.IP_Address = host
	dnses, err := net.LookupAddr(dtoaccesslog.IP_Address)
	if err != nil {
		log.Error("Can't detect reverse dns %v for %v", err, dtoaccesslog.IP_Address)
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_OBJECT_NOTEXIST,
			Message: config.Localization[language].Errors.Api.Object_NotExist})
		return nil, err
	}
	dtoaccesslog.Reverse_DNS = strings.Join(dnses, ",")
	dtoaccesslog.User_Agent = request.UserAgent()
	var urllen int
	urllen = len([]rune(request.Referer()))
	if urllen > URL_LENGTH_MAX {
		urllen = URL_LENGTH_MAX
	}
	if urllen > 0 {
		dtoaccesslog.Referer = string([]rune(request.Referer())[:urllen-1])
	}
	urllen = len([]rune(url))
	if urllen > URL_LENGTH_MAX {
		urllen = URL_LENGTH_MAX
	}
	if urllen > 0 {
		dtoaccesslog.URL = string([]rune(url)[:urllen-1])
	}
	dtoaccesslog.Created = time.Now()

	err = accesslogrepository.Create(dtoaccesslog)
	if err != nil {
		r.JSON(http.StatusNotFound, types.Error{Code: types.TYPE_ERROR_DATA_WRONG,
			Message: config.Localization[language].Errors.Api.Data_Wrong})
		return nil, err
	}

	return dtoaccesslog, nil
}
