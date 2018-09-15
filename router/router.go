package router

import (
	"reflect"
	"strings"
)

var mapping map[string]reflect.Type = make(map[string]reflect.Type)

func router(pattern string, t reflect.Type) {
	mapping[strings.ToLower(pattern)] = t
}

func Router(pattern string, app IApp) {
	refV := reflect.ValueOf(app)
	refT := reflect.Indirect(refV).Type()
	router(pattern, refT)
}

func AutoRouter(app IApp) {
	refV := reflect.ValueOf(app)
	refT := reflect.Indirect(refV).Type()
	refName := strings.TrimSuffix(strings.ToLower(refT.Name()), "controller")
	router(refName, refT)
}
