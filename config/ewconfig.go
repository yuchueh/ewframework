package config

import (
	"path/filepath"
	"os"
	"github.com/yuchueh/ewframework/utils/osext"
	"errors"
	"github.com/yuchueh/ewframework/utils"
)

// Config is the main struct for BConfig
type EwConfig struct {
	AppName             string //Application name
	RunMode             string //Running Mode: dev | prod
	RouterCaseSensitive bool
	ServerName          string
	RecoverPanic        bool
	//RecoverFunc         func(*context.Context)
	CopyRequestBody     bool
	EnableGzip          bool
	MaxMemory           int64
	EnableErrorsShow    bool
	EnableErrorsRender  bool
	Listen              Listen
	WebConfig           WebConfig
	Log                 LogConfig
}

// Listen holds for http and https related config
type Listen struct {
	Graceful      bool // Graceful means use graceful module to start the server
	ServerTimeOut int64
	ListenTCP4    bool
	EnableHTTP    bool
	HTTPAddr      string
	HTTPPort      int
	EnableHTTPS   bool
	HTTPSAddr     string
	HTTPSPort     int
	HTTPSCertFile string
	HTTPSKeyFile  string
	EnableAdmin   bool
	AdminAddr     string
	AdminPort     int
	EnableFcgi    bool
	EnableStdIo   bool // EnableStdIo works with EnableFcgi Use FCGI via standard I/O
}

// WebConfig holds web related config
type WebConfig struct {
	AutoRender             bool
	EnableDocs             bool
	FlashName              string
	FlashSeparator         string
	DirectoryIndex         bool
	StaticDir              map[string]string
	StaticExtensionsToGzip []string
	TemplateLeft           string
	TemplateRight          string
	ViewsPath              string
	EnableXSRF             bool
	XSRFKey                string
	XSRFExpire             int
	Session                SessionConfig
}

// SessionConfig holds session related config
type SessionConfig struct {
	SessionOn                    bool
	SessionProvider              string
	SessionName                  string
	SessionGCMaxLifetime         int64
	SessionProviderConfig        string
	SessionCookieLifeTime        int
	SessionAutoSetCookie         bool
	SessionDomain                string
	SessionDisableHTTPOnly       bool // used to allow for cross domain cookies/javascript cookies.
	SessionEnableSidInHTTPHeader bool //	enable store/get the sessionId into/from http headers
	SessionNameInHTTPHeader      string
	SessionEnableSidInURLQuery   bool //	enable get the sessionId from Url Query params
}

// LogConfig holds Log related config
type LogConfig struct {
	AccessLogs  bool
	FileLineNum bool
	Outputs     map[string]string // Store Adaptor : config
}

type ewAppConfig struct {
	innerConfig *Config
	filename	string
}

var (
	// BConfig is the default config for Application
	BConfig *EwConfig
	// AppConfig is the instance of Config, store the config information from file
	AppConfig *ewAppConfig
	// AppPath is the absolute path to the app
	AppPath string
	// GlobalSessions is the instance for the session manager
	//GlobalSessions *session.Manager

	// appConfigPath is the path to the config files
	appConfigPath string
	// appConfigProvider is the provider for the config, default is ini
	appConfigProvider = "ini"
)

//func recoverPanic(ctx *context.Context) {
//	if err := recover(); err != nil {
//		if err == utils.ErrAbort {
//			return
//		}
//		if !BConfig.RecoverPanic {
//			panic(err)
//		}
//		if BConfig.EnableErrorsShow {
//			if _, ok := ErrorMaps[fmt.Sprint(err)]; ok {
//				icode, _ := strconv.Atoi(fmt.Sprint(err))
//				Exception(uint64(icode), ctx)
//				return
//			}
//		}
//		var stack string
//		logs.Critical("the request url is ", ctx.Input.URL())
//		logs.Critical("Handler crashed with error", err)
//		for i := 1; ; i++ {
//			_, file, line, ok := runtime.Caller(i)
//			if !ok {
//				break
//			}
//			logs.Critical(fmt.Sprintf("%s:%d", file, line))
//			stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
//		}
//		if BConfig.RunMode == DEV && BConfig.EnableErrorsRender {
//			ShowErr(err, ctx, stack)
//		}
//	}
//}

func newBConfig() *EwConfig {
	return &EwConfig{
		AppName:             "eweb",
		RunMode:             utils.DEV,
		RouterCaseSensitive: true,
		ServerName:          "EwGoServer:" + utils.VERSION,
		RecoverPanic:        true,
		//RecoverFunc:         recoverPanic,
		CopyRequestBody:     false,
		EnableGzip:          false,
		MaxMemory:           1 << 26, //64MB
		EnableErrorsShow:    true,
		EnableErrorsRender:  true,
		Listen: Listen{
			Graceful:      false,
			ServerTimeOut: 0,
			ListenTCP4:    false,
			EnableHTTP:    true,
			HTTPAddr:      "",
			HTTPPort:      8080,
			EnableHTTPS:   false,
			HTTPSAddr:     "",
			HTTPSPort:     10443,
			HTTPSCertFile: "",
			HTTPSKeyFile:  "",
			EnableAdmin:   false,
			AdminAddr:     "",
			AdminPort:     8088,
			EnableFcgi:    false,
			EnableStdIo:   false,
		},
		WebConfig: WebConfig{
			AutoRender:             true,
			EnableDocs:             false,
			FlashName:              "EWGO_FLASH",
			FlashSeparator:         "EWGOFLASH",
			DirectoryIndex:         false,
			StaticDir:              map[string]string{"/static": "static"},
			StaticExtensionsToGzip: []string{".css", ".js"},
			TemplateLeft:           "{{",
			TemplateRight:          "}}",
			ViewsPath:              "views",
			EnableXSRF:             false,
			XSRFKey:                "ewgoxsrf",
			XSRFExpire:             0,
			Session: SessionConfig{
				SessionOn:                    false,
				SessionProvider:              "memory",
				SessionName:                  "ewgosessionID",
				SessionGCMaxLifetime:         3600,
				SessionProviderConfig:        "",
				SessionDisableHTTPOnly:       false,
				SessionCookieLifeTime:        0, //set cookie default is the browser life
				SessionAutoSetCookie:         true,
				SessionDomain:                "",
				SessionEnableSidInHTTPHeader: false, //	enable store/get the sessionId into/from http headers
				SessionNameInHTTPHeader:      "Ewgosessionid",
				SessionEnableSidInURLQuery:   false, //	enable get the sessionId from Url Query params
			},
		},
		Log: LogConfig{
			AccessLogs:  false,
			FileLineNum: true,
			Outputs:     map[string]string{"console": ""},
		},
	}
}

func init() {
	BConfig = newBConfig()
	var err error
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	appConfigPath = filepath.Join(workPath, "conf", "app.conf")
	if !osext.FileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppPath, "conf", "app.conf")
		if !osext.FileExists(appConfigPath) {
			AppConfig = &ewAppConfig{ }
			cf, _ := ReadDefault(appConfigPath)
			AppConfig.innerConfig = cf
			//file name
			AppConfig.filename = appConfigPath
			return
		}
	}
}

func (b *ewAppConfig) Set(section, key, value string) error {
	if !b.innerConfig.AddOption(section, key, value) {
		return errors.New("AddOption error!")
	} else {
		return nil
	}
}

func (b *ewAppConfig) String(section, key string) string {
	val, err := b.innerConfig.String(section, key)
	if err != nil {
		return ""
	}
	return val
}

func (b *ewAppConfig) Int(section, key string) (int, error) {
	return b.innerConfig.Int(section, key)
}

func (b *ewAppConfig) Bool(section, key string) (bool, error) {
	return b.innerConfig.Bool(section, key)
}

func (b *ewAppConfig) Float(section, key string) (float64, error) {
	return b.innerConfig.Float(section, key)
}

func (b *ewAppConfig) DefaultString(section, key string, defaultVal string) string {
	if v := b.String(section, key); v != "" {
		return v
	}
	return defaultVal
}

func (b *ewAppConfig) DefaultInt(section, key string, defaultVal int) int {
	if v, err := b.Int(section, key); err == nil {
		return v
	}
	return defaultVal
}

func (b *ewAppConfig) DefaultBool(section, key string, defaultVal bool) bool {
	if v, err := b.Bool(section, key); err == nil {
		return v
	}
	return defaultVal
}

func (b *ewAppConfig) DefaultFloat(section, key string, defaultVal float64) float64 {
	if v, err := b.Float(section, key); err == nil {
		return v
	}
	return defaultVal
}

func (b *ewAppConfig) HasSection(section string) bool {
	return b.innerConfig.HasSection(section)
}

func (b *ewAppConfig) GetSection(section string) ([]string, error) {
	return b.innerConfig.Options(section)
}

func (b *ewAppConfig) LoadConfig(filename string)  {
	if osext.FileExists(filename) {
		cf, _ := ReadDefault(filename)
		AppConfig.innerConfig = cf
	}
}

func (b *ewAppConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.WriteFile(b.filename, 0644, "")
}
