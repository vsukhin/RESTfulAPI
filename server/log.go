package server

import (
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	logging "github.com/op/go-logging"
)

var (
	logger = logging.MustGetLogger("server")
)

func logRequest(context martini.Context, request *http.Request, response http.ResponseWriter) {
	responseWriter := response.(martini.ResponseWriter)
	logger.Info("Request: %s %s", request.Method, request.URL.String())
	context.Next()
	responseInfo := fmt.Sprintf("Response: %d %s", responseWriter.Status(), http.StatusText(responseWriter.Status()))
	if responseWriter.Status() < 400 {
		logger.Info(responseInfo)
	} else if responseWriter.Status() < 500 {
		logger.Warning(responseInfo)
	} else {
		logger.Error(responseInfo)
	}
}
