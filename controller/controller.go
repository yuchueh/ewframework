package controller

import "github.com/yuchueh/ewframework/context/param"

// ControllerInterface is an interface to uniform all controller handler.
type ControllerInterface interface {
	//Init(ct *context.Context, controllerName, actionName string, app interface{})
	Prepare()
	Get()
	Post()
	Delete()
	Put()
	Head()
	Patch()
	Options()
	Finish()
	Render() error
	XSRFToken() string
	CheckXSRFCookie() bool
	HandlerFunc(fn string) bool
	URLMapping()
}

// ControllerComments store the comment for the controller method
type ControllerComments struct {
	Method           string
	Router           string
	AllowHTTPMethods []string
	Params           []map[string]string
	MethodParams     []*param.MethodParam
}

// ControllerCommentsSlice implements the sort interface
type ControllerCommentsSlice []ControllerComments

func (p ControllerCommentsSlice) Len() int           { return len(p) }
func (p ControllerCommentsSlice) Less(i, j int) bool { return p[i].Router < p[j].Router }
func (p ControllerCommentsSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
