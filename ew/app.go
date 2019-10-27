package ew

import (
	"net/http"
	"github.com/yuchueh/ewframework/router"
	"fmt"
	"net"
	"net/http/fcgi"
	"os"
	"time"
	"path"
	"github.com/yuchueh/ewframework/config"
	"github.com/yuchueh/ewframework/logs"
	"github.com/yuchueh/ewframework/utils/osext"
	"github.com/yuchueh/ewframework/controller"
	"github.com/yuchueh/ewframework/filter"
)

var (
	// EwApp is an application instance
	EwApp *App
)

func init() {
	// create ew application
	EwApp = NewApp()
}

// App defines application with a new PatternServeMux.
type App struct {
	Handlers *router.ControllerRegister
	Server   *http.Server
}

// NewApp returns a new application.
func NewApp() *App {
	logs.Debug("app.NewApp")
	cr := router.NewControllerRegister()
	app := &App{Handlers: cr, Server: &http.Server{}}
	return app
}

// Run ew application.
func (app *App) Run() {
	logs.Debug("app.Run()")
	addr := config.BConfig.Listen.HTTPAddr

	if config.BConfig.Listen.HTTPPort != 0 {
		addr = fmt.Sprintf("%s:%d", config.BConfig.Listen.HTTPAddr, config.BConfig.Listen.HTTPPort)
	}
	logs.Debug("app.Run.addr=", addr)

	var (
		err        error
		l          net.Listener
		endRunning = make(chan bool, 1)
	)

	// run cgi server
	if config.BConfig.Listen.EnableFcgi {
		logs.Informational("Listen.EnableFcgi=true")

		if config.BConfig.Listen.EnableStdIo {
			if err = fcgi.Serve(nil, app.Handlers); err == nil { // standard I/O
				logs.Informational("Use FCGI via standard I/O")
			} else {
				logs.Critical("Cannot use FCGI via standard I/O", err)
			}
			return
		}
		if config.BConfig.Listen.HTTPPort == 0 {
			// remove the Socket file before start
			if osext.FileExists(addr) {
				os.Remove(addr)
			}
			l, err = net.Listen("unix", addr)
		} else {
			l, err = net.Listen("tcp", addr)
		}
		if err != nil {
			logs.Critical("Listen: ", err)
		}
		if err = fcgi.Serve(l, app.Handlers); err != nil {
			logs.Critical("fcgi.Serve: ", err)
		}
		return
	}

	//Bind Server Handler
	logs.Debug("Bind Server Handler:", app.Handlers)
	app.Server.Handler = app.Handlers
	app.Server.ReadTimeout = time.Duration(config.BConfig.Listen.ServerTimeOut) * time.Second
	app.Server.WriteTimeout = time.Duration(config.BConfig.Listen.ServerTimeOut) * time.Second
	app.Server.ErrorLog = logs.GetLogger("HTTP")

	// run graceful mode
	//if config.BConfig.Listen.Graceful {
	//	httpsAddr := config.BConfig.Listen.HTTPSAddr
	//	app.Server.Addr = httpsAddr
	//	if config.BConfig.Listen.EnableHTTPS {
	//		go func() {
	//			time.Sleep(20 * time.Microsecond)
	//			if config.BConfig.Listen.HTTPSPort != 0 {
	//				httpsAddr = fmt.Sprintf("%s:%d", config.BConfig.Listen.HTTPSAddr, config.BConfig.Listen.HTTPSPort)
	//				app.Server.Addr = httpsAddr
	//			}
	//			server := grace.NewServer(httpsAddr, app.Handlers)
	//			server.Server.ReadTimeout = app.Server.ReadTimeout
	//			server.Server.WriteTimeout = app.Server.WriteTimeout
	//			if err := server.ListenAndServeTLS(config.BConfig.Listen.HTTPSCertFile, config.BConfig.Listen.HTTPSKeyFile); err != nil {
	//				logs.Critical("ListenAndServeTLS: ", err, fmt.Sprintf("%d", os.Getpid()))
	//				time.Sleep(100 * time.Microsecond)
	//				endRunning <- true
	//			}
	//		}()
	//	}
	//	if config.BConfig.Listen.EnableHTTP {
	//		go func() {
	//			server := grace.NewServer(addr, app.Handlers)
	//			server.Server.ReadTimeout = app.Server.ReadTimeout
	//			server.Server.WriteTimeout = app.Server.WriteTimeout
	//			if config.BConfig.Listen.ListenTCP4 {
	//				server.Network = "tcp4"
	//			}
	//			if err := server.ListenAndServe(); err != nil {
	//				logs.Critical("ListenAndServe: ", err, fmt.Sprintf("%d", os.Getpid()))
	//				time.Sleep(100 * time.Microsecond)
	//				endRunning <- true
	//			}
	//		}()
	//	}
	//	<-endRunning
	//	return
	//}

	// run normal mode EnableHTTPS
	if config.BConfig.Listen.EnableHTTPS {
		go func() {
			time.Sleep(20 * time.Microsecond)
			if config.BConfig.Listen.HTTPSPort != 0 {
				app.Server.Addr = fmt.Sprintf("%s:%d", config.BConfig.Listen.HTTPSAddr, config.BConfig.Listen.HTTPSPort)
			} else if config.BConfig.Listen.EnableHTTP {
				logs.Informational("Start https server error, confict with http.Please reset https port")
				return
			}
			logs.Informational("https server Running on https://%s", app.Server.Addr)
			if err := app.Server.ListenAndServeTLS(config.BConfig.Listen.HTTPSCertFile, config.BConfig.Listen.HTTPSKeyFile); err != nil {
				logs.Critical("ListenAndServeTLS: ", err)
				time.Sleep(100 * time.Microsecond)
				endRunning <- true
			}
		}()
	}

	//EnableHTTP
	if config.BConfig.Listen.EnableHTTP {
		go func() {
			app.Server.Addr = addr

			//Listen.HTTPSAddr
			if len(config.BConfig.Listen.HTTPSAddr) == 0 {
				//localhost
				logs.Informational("http server Running on http://localhost:", config.BConfig.Listen.HTTPPort)
			} else {
				logs.Informational("http server Running on http://", app.Server.Addr)
			}

			if config.BConfig.Listen.ListenTCP4 {
				ln, err := net.Listen("tcp4", app.Server.Addr)
				if err != nil {
					logs.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
					endRunning <- true
					return
				}
				if err = app.Server.Serve(ln); err != nil {
					logs.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
					endRunning <- true
					return
				}
			} else {
				if err := app.Server.ListenAndServe(); err != nil {
					logs.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
					endRunning <- true
				}
			}
		}()
	}
	<-endRunning
}

// Router adds a patterned controller handler to EwApp.
// it's an alias method of App.Router.
// usage:
//  simple router
//  ew.Router("/admin", &admin.UserController{})
//  ew.Router("/admin/index", &admin.ArticleController{})
//
//  regex router
//
//  ew.Router("/api/:id([0-9]+)", &controllers.RController{})
//
//  custom rules
//  ew.Router("/api/list",&RestController{},"*:ListFood")
//  ew.Router("/api/create",&RestController{},"post:CreateFood")
//  ew.Router("/api/update",&RestController{},"put:UpdateFood")
//  ew.Router("/api/delete",&RestController{},"delete:DeleteFood")
func Router(rootpath string, c controller.ControllerInterface, mappingMethods ...string) *App {
	EwApp.Handlers.Add(rootpath, c, mappingMethods...)
	return EwApp
}

// Include will generate router file in the router/xxx.go from the controller's comments
// usage:
// ew.Include(&BankAccount{}, &OrderController{},&RefundController{},&ReceiptController{})
// type BankAccount struct{
//   ew.Controller
// }
//
// register the function
// func (b *BankAccount)Mapping(){
//  b.Mapping("ShowAccount" , b.ShowAccount)
//  b.Mapping("ModifyAccount", b.ModifyAccount)
//}
//
// //@router /account/:id  [get]
// func (b *BankAccount) ShowAccount(){
//    //logic
// }
//
//
// //@router /account/:id  [post]
// func (b *BankAccount) ModifyAccount(){
//    //logic
// }
//
// the comments @router url methodlist
// url support all the function Router's pattern
// methodlist [get post head put delete options *]
func Include(cList ...controller.ControllerInterface) *App {
	EwApp.Handlers.Include(cList...)
	return EwApp
}

// RESTRouter adds a restful controller handler to EwApp.
// its' controller implements ew.ControllerInterface and
// defines a param "pattern/:objectId" to visit each resource.
func RESTRouter(rootpath string, c controller.ControllerInterface) *App {
	Router(rootpath, c)
	Router(path.Join(rootpath, ":objectId"), c)
	return EwApp
}

// AutoRouter adds defined controller handler to EwApp.
// it's same to App.AutoRouter.
// if ew.AddAuto(&MainContorlller{}) and MainController has methods List and Page,
// visit the url /main/list to exec List function or /main/page to exec Page function.
func AutoRouter(c controller.ControllerInterface) *App {
	EwApp.Handlers.AddAuto(c)
	return EwApp
}

// AutoPrefix adds controller handler to EwApp with prefix.
// it's same to App.AutoRouterWithPrefix.
// if ew.AutoPrefix("/admin",&MainContorlller{}) and MainController has methods List and Page,
// visit the url /admin/main/list to exec List function or /admin/main/page to exec Page function.
func AutoPrefix(prefix string, c controller.ControllerInterface) *App {
	EwApp.Handlers.AddAutoPrefix(prefix, c)
	return EwApp
}

// Get used to register router for Get method
// usage:
//    ew.Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Get(rootpath string, f filter.FilterFunc) *App {
	EwApp.Handlers.Get(rootpath, f)
	return EwApp
}

// Post used to register router for Post method
// usage:
//    ew.Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Post(rootpath string, f filter.FilterFunc) *App {
	EwApp.Handlers.Post(rootpath, f)
	return EwApp
}

// Delete used to register router for Delete method
// usage:
//    ew.Delete("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Delete(rootpath string, f filter.FilterFunc) *App {
	EwApp.Handlers.Delete(rootpath, f)
	return EwApp
}

// Put used to register router for Put method
// usage:
//    ew.Put("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Put(rootpath string, f filter.FilterFunc) *App {
	EwApp.Handlers.Put(rootpath, f)
	return EwApp
}

// Head used to register router for Head method
// usage:
//    ew.Head("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Head(rootpath string, f filter.FilterFunc) *App {
	EwApp.Handlers.Head(rootpath, f)
	return EwApp
}

// Options used to register router for Options method
// usage:
//    ew.Options("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Options(rootpath string, f filter.FilterFunc) *App {
	EwApp.Handlers.Options(rootpath, f)
	return EwApp
}

// Patch used to register router for Patch method
// usage:
//    ew.Patch("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Patch(rootpath string, f filter.FilterFunc) *App {
	EwApp.Handlers.Patch(rootpath, f)
	return EwApp
}

// Any used to register router for all methods
// usage:
//    ew.Any("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Any(rootpath string, f filter.FilterFunc) *App {
	EwApp.Handlers.Any(rootpath, f)
	return EwApp
}

// Handler used to register a Handler router
// usage:
//    ew.Handler("/api", http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
//          fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
//    }))
func Handler(rootpath string, h http.Handler, options ...interface{}) *App {
	EwApp.Handlers.Handler(rootpath, h, options...)
	return EwApp
}

// InsertFilter adds a FilterFunc with pattern condition and action constant.
// The pos means action constant including
// ew.BeforeStatic, ew.BeforeRouter, ew.BeforeExec, ew.AfterExec and ew.FinishRouter.
// The bool params is for setting the returnOnOutput value (false allows multiple filters to execute)
func InsertFilter(pattern string, pos int, filter filter.FilterFunc, params ...bool) *App {
	EwApp.Handlers.InsertFilter(pattern, pos, filter, params...)
	return EwApp
}