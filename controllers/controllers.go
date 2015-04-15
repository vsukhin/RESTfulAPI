/* Controllers package provides methods responsible for RESTFul API business logic implementation */

package controllers

import (
	"application/config"
	logging "github.com/op/go-logging"
)

var (
	log config.Logger = logging.MustGetLogger("controllers")
)

func InitLogger(logger config.Logger) {
	log = logger
}
