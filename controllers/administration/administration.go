/* Administration package provides methods responsible for RESTFul API administrative operations implementation */

package administration

import (
	"application/config"
	logging "github.com/op/go-logging"
)

var (
	log config.Logger = logging.MustGetLogger("administration")
)

func InitLogger(logger config.Logger) {
	log = logger
}
