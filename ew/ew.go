package ew

import (
	"strings"
	"strconv"
	"github.com/yuchueh/ewframework/config"
	"github.com/yuchueh/ewframework/logs"
)

//hook function to run
type hookfunc func() error

var (
	hooks = make([]hookfunc, 0) //hook function slice to store the hookfunc
)

// AddAPPStartHook is used to register the hookfunc
// The hookfuncs will run in ew.Run()
// such as initiating session , starting middleware , building template, starting admin control and so on.
func AddAPPStartHook(hf ...hookfunc) {
	hooks = append(hooks, hf...)
}

// Run ew application.
// ew.Run() default run on HttpPort
// ew.Run("localhost")
// ew.Run(":8089")
// ew.Run("127.0.0.1:8089")
func Run(params ...string) {

	logs.GLogger.Debug("ew.Run params:", params)

	initBeforeHTTPRun()

	if len(params) > 0 && params[0] != "" {
		strs := strings.Split(params[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			config.BConfig.Listen.HTTPAddr = strs[0]
		}
		if len(strs) > 1 && strs[1] != "" {
			config.BConfig.Listen.HTTPPort, _ = strconv.Atoi(strs[1])
		}
	}

	EwApp.Run()
}

func initBeforeHTTPRun() {
	//init hooks
	AddAPPStartHook(
		registerMime,
		registerDefaultErrorHandler,
		registerSession,
		registerTemplate,
		registerAdmin,
		registerGzip,
	)

	for _, hk := range hooks {
		if err := hk(); err != nil {
			panic(err)
		}
	}
}