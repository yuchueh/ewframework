package router

import "net/http"

type IContext interface {
	Config(w http.ResponseWriter, r *http.Request)
}

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func (c *Context) Config(w http.ResponseWriter, r *http.Request) {
	c.w = w
	c.r = r
}
