package controllers

import (
	"github.com/martini-contrib/render"
	"html/template"
	"net/http"
	"types"
)

type Renderer struct {
	ErrorValue  types.Error
	StatusValue int
}

func (r *Renderer) JSON(status int, v interface{}) {
	r.StatusValue = status
	r.ErrorValue, _ = v.(types.Error)
}
func (r *Renderer) Data(status int, v []byte) {

}
func (r *Renderer) Error(status int) {

}
func (r *Renderer) HTML(status int, name string, binding interface{}, htmlOpt ...render.HTMLOptions) {

}
func (r *Renderer) XML(status int, v interface{}) {

}
func (r *Renderer) Status(status int) {

}
func (r *Renderer) Redirect(location string, status ...int) {

}
func (r *Renderer) Header() http.Header {
	return nil
}
func (r *Renderer) Template() *template.Template {
	return nil
}
