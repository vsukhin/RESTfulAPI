package middlewares

import (
	"net/http"
)

// Not found error for wrong URI
func Default(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}
