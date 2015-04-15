package middlewares

import (
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	
	"application/config"
)

// Redirect to public address
func Default(r render.Render, context martini.Context) {
	r.Redirect(config.Configuration.Server.PublicAddress, http.StatusMovedPermanently)
}
