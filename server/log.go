package server

import (
	"application/config"
	"fmt"
	"github.com/go-martini/martini"
	logging "github.com/op/go-logging"
	"net/http"
)

var (
	log config.Logger = logging.MustGetLogger("server")
)

func InitLogger(logger config.Logger) {
	log = logger
}

func LogRequest(context martini.Context, request *http.Request, response http.ResponseWriter) {
	responseWriter := response.(martini.ResponseWriter)
	log.Info("Request: %s %s", request.Method, request.URL.String())
	context.Next()
	responseInfo := fmt.Sprintf("Response: %d %s", responseWriter.Status(), http.StatusText(responseWriter.Status()))
	if responseWriter.Status() < 400 {
		log.Info(responseInfo)
	} else if responseWriter.Status() < 500 {
		log.Warning(responseInfo)
	} else {
		log.Error(responseInfo)
	}
}
