package ew

import (
	"encoding/json"
	"mime"
	"net/http"
	"path/filepath"
	"github.com/yuchueh/ewframework/utils"
	"github.com/yuchueh/ewframework/context"
	"github.com/yuchueh/ewframework/config"
	"github.com/yuchueh/ewframework/session"
	"github.com/yuchueh/ewframework/logs"
	ewtemplate "github.com/yuchueh/ewframework/template"
)

//register Mime
func registerMime() error {
	for k, v := range utils.Mimemaps {
		mime.AddExtensionType(k, v)
	}
	return nil
}

// register default error http handlers, 404,401,403,500 and 503.
func registerDefaultErrorHandler() error {
	m := map[string]func(http.ResponseWriter, *http.Request){
		"401": context.Unauthorized,
		"402": context.PaymentRequired,
		"403": context.Forbidden,
		"404": context.NotFound,
		"405": context.MethodNotAllowed,
		"500": context.InternalServerError,
		"501": context.NotImplemented,
		"502": context.BadGateway,
		"503": context.ServiceUnavailable,
		"504": context.GatewayTimeout,
		"417": context.Invalidxsrf,
		"422": context.Missingxsrf,
	}
	for e, h := range m {
		if _, ok := context.ErrorMaps[e]; !ok {
			context.ErrorHandler(e, h)
		}
	}
	return nil
}

func registerSession() error {
	if config.BConfig.WebConfig.Session.SessionOn {
		var err error
		sessionConfig := config.AppConfig.String("session","sessionConfig")
		conf := new(session.ManagerConfig)
		if sessionConfig == "" {
			conf.CookieName = config.BConfig.WebConfig.Session.SessionName
			conf.EnableSetCookie = config.BConfig.WebConfig.Session.SessionAutoSetCookie
			conf.Gclifetime = config.BConfig.WebConfig.Session.SessionGCMaxLifetime
			conf.Secure = config.BConfig.Listen.EnableHTTPS
			conf.CookieLifeTime = config.BConfig.WebConfig.Session.SessionCookieLifeTime
			conf.ProviderConfig = filepath.ToSlash(config.BConfig.WebConfig.Session.SessionProviderConfig)
			conf.DisableHTTPOnly = config.BConfig.WebConfig.Session.SessionDisableHTTPOnly
			conf.Domain = config.BConfig.WebConfig.Session.SessionDomain
			conf.EnableSidInHTTPHeader = config.BConfig.WebConfig.Session.SessionEnableSidInHTTPHeader
			conf.SessionNameInHTTPHeader = config.BConfig.WebConfig.Session.SessionNameInHTTPHeader
			conf.EnableSidInURLQuery = config.BConfig.WebConfig.Session.SessionEnableSidInURLQuery
		} else {
			if err = json.Unmarshal([]byte(sessionConfig), conf); err != nil {
				return err
			}
		}
		if session.GlobalSessions, err = session.NewManager(config.BConfig.WebConfig.Session.SessionProvider, conf); err != nil {
			return err
		}
		go session.GlobalSessions.GC()
	}
	return nil
}

func registerTemplate() error {
	defer ewtemplate.LockViewPaths()
	if err := ewtemplate.AddViewPath(config.BConfig.WebConfig.ViewsPath); err != nil {
		if config.BConfig.RunMode == utils.DEV {
			logs.Warning(err)
		}
		return err
	}
	return nil
}

//TODO registerAdmin
func registerAdmin() error {
	if config.BConfig.Listen.EnableAdmin {
		//go beeAdminApp.Run()
	}
	return nil
}

func registerGzip() error {
	if config.BConfig.EnableGzip {
		context.InitGzip(
			config.AppConfig.DefaultInt("","gzipMinLength", -1),
			config.AppConfig.DefaultInt("","gzipCompressLevel", -1),
			config.AppConfig.DefaultStrings("","includedMethods", []string{"GET"}),
		)
	}
	return nil
}
