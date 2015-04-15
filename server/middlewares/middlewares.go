/* Middlewares package provides methods for preprocessing requests to the system */

package middlewares

import (
	"application/config"
	logging "github.com/op/go-logging"
)

var (
	log config.Logger = logging.MustGetLogger("middlewares")
)

func InitLogger(logger config.Logger) {
	log = logger
}
