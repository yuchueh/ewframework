package context

import (
	"net/http"
	"strings"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"encoding/base64"
	"strconv"
	"time"
	r "github.com/yuchueh/ewframework/utils/rand"
)

// Context Http request context struct including HttpInput, HttpOutput, http.Request and http.ResponseWriter.
// HttpInput and HttpOutput provides some api to operate request and response more easily.
type Context struct {
	Input          *HttpInput
	Output         *HttpOutput
	Request        *http.Request
	ResponseWriter *Response
	_xsrfToken     string
}

// NewContext return the Context with Input and Output
func NewContext() *Context {
	return &Context{
		Input:  NewInput(),
		Output: NewOutput(),
	}
}

// Reset init Context
func (ctx *Context) Reset(rw http.ResponseWriter, r *http.Request) {
	ctx.Request = r
	if ctx.ResponseWriter == nil {
		ctx.ResponseWriter = &Response{}
	}
	ctx.ResponseWriter.reset(rw)
	ctx.Input.Reset(ctx)
	ctx.Output.Reset(ctx)
	ctx._xsrfToken = ""
}

// Redirect does redirection to localurl with http header status code.
func (ctx *Context) Redirect(status int, localurl string) {
	http.Redirect(ctx.ResponseWriter, ctx.Request, localurl, status)
}

// Abort stops this request.
// if beego.ErrorMaps exists, panic body.
func (ctx *Context) Abort(status int, body string) {
	ctx.Output.SetStatus(status)
	panic(body)
}

// WriteString Write string to response body.
// it sends response body.
func (ctx *Context) WriteString(content string) {
	ctx.ResponseWriter.Write([]byte(content))
}

// GetCookie Get cookie from request by a given key.
// It's alias of HttpInput.Cookie.
func (ctx *Context) GetCookie(key string) string {
	return ctx.Input.Cookie(key)
}

// SetCookie Set cookie for response.
// It's alias of HttpOutput.Cookie.
func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	ctx.Output.Cookie(name, value, others...)
}

// GetSecureCookie Get secure cookie from request by a given key.
func (ctx *Context) GetSecureCookie(Secret, key string) (string, bool) {
	val := ctx.Input.Cookie(key)
	if val == "" {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3)

	if len(parts) != 3 {
		return "", false
	}

	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]

	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)

	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}

// SetSecureCookie Set Secure cookie for response.
func (ctx *Context) SetSecureCookie(Secret, name, value string, others ...interface{}) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	ctx.Output.Cookie(name, cookie, others...)
}

// XSRFToken creates a xsrf token string and returns.
func (ctx *Context) XSRFToken(key string, expire int64) string {
	if ctx._xsrfToken == "" {
		token, ok := ctx.GetSecureCookie(key, "_xsrf")
		if !ok {
			token = string(r.RandomCreateBytes(32))
			ctx.SetSecureCookie(key, "_xsrf", token, expire)
		}
		ctx._xsrfToken = token
	}
	return ctx._xsrfToken
}

// CheckXSRFCookie checks xsrf token in this request is valid or not.
// the token can provided in request header "X-Xsrftoken" and "X-CsrfToken"
// or in form field value named as "_xsrf".
func (ctx *Context) CheckXSRFCookie() bool {
	token := ctx.Input.Query("_xsrf")
	if token == "" {
		token = ctx.Request.Header.Get("X-Xsrftoken")
	}
	if token == "" {
		token = ctx.Request.Header.Get("X-Csrftoken")
	}
	if token == "" {
		ctx.Abort(403, "'_xsrf' argument missing from POST")
		return false
	}
	if ctx._xsrfToken != token {
		ctx.Abort(403, "XSRF cookie does not match POST argument")
		return false
	}
	return true
}

// RenderMethodResult renders the return value of a controller method to the output
func (ctx *Context) RenderMethodResult(result interface{}) {
	if result != nil {
		renderer, ok := result.(Renderer)
		if !ok {
			err, ok := result.(error)
			if ok {
				renderer = errorRenderer(err)
			} else {
				renderer = jsonRenderer(result)
			}
		}
		renderer.Render(ctx)
	}
}