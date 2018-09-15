package main

import (
	"net/http"
	"html/template"
	//"github.com/yuchueh/ewframework/router"
)

type IApp interface {
	//Init(ctx *Context)
	W() http.ResponseWriter
	R() *http.Request
	Display(tpls ...string)
	DisplayWithFuncs(funcs template.FuncMap, tpls ...string)
}
