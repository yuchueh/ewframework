package context

import (
	"net/http"
	"net"
	"bufio"
	"errors"
)

//Response is a wrapper for the http.ResponseWriter
//started set to true if response was written to then don't execute other handler
type Response struct {
	http.ResponseWriter
	Started bool
	Status  int
}

func (r *Response) reset(rw http.ResponseWriter) {
	r.ResponseWriter = rw
	r.Status = 0
	r.Started = false
}

// Write writes the data to the connection as part of an HTTP reply,
// and sets `started` to true.
// started means the response has sent out.
func (r *Response) Write(p []byte) (int, error) {
	r.Started = true
	return r.ResponseWriter.Write(p)
}

// WriteHeader sends an HTTP response header with status code,
// and sets `started` to true.
func (r *Response) WriteHeader(code int) {
	if r.Status > 0 {
		//prevent multiple response.WriteHeader calls
		return
	}
	r.Status = code
	r.Started = true
	r.ResponseWriter.WriteHeader(code)
}

// Hijack hijacker for http
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("webserver doesn't support hijacking")
	}
	return hj.Hijack()
}

// Flush http.Flusher
func (r *Response) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// CloseNotify http.CloseNotifier
func (r *Response) CloseNotify() <-chan bool {
	if cn, ok := r.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}
