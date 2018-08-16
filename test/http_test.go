package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime"
	"testing"
	"time"
)

type MyStruct struct {
}

func (s *MyStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func World(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "World")
}

func Test_Http(t *testing.T) {
	t.Log("Test_Http Start")
	handler := MyStruct{}
	//http.ListenAndServe("", nil)
	svr := http.Server{
		Addr: "",
		//Handler:&handler,
	}

	http.Handle("/", &handler)
	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/world", World)
	svr.ListenAndServe()

	t.Log("Test_Http End")
}

func log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		fmt.Println("Handler function called - " + name)
		h(w, r)
	}
}

func Test_ChainHttp(t *testing.T) {
	t.Log("Test_ChainHttp Start")

	svr := http.Server{
		Addr: "",
	}
	http.HandleFunc("/hello", log(Hello))
	svr.ListenAndServe()
}

type HelloStruct struct {
}

func (c HelloStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("logHandler called - %T\n", h)
		h.ServeHTTP(w, r)
	})
}

func protectHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("protectHandler called - %T\n", h)
		h.ServeHTTP(w, r)
	})
}

func Test_ChainHttp_handler(t *testing.T) {
	hello := HelloStruct{}

	svr := http.Server{
		Addr: "",
		//Handler:&hello,
	}

	http.Handle("/hello", protectHandler(logHandler(hello)))
	svr.ListenAndServe()
}

func httpHeader(w http.ResponseWriter, r *http.Request) {
	h := r.Header
	fmt.Fprintln(w, h)
}

func httpBody(w http.ResponseWriter, r *http.Request) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	fmt.Fprintln(w, string(body))
}

func Test_HttpHeader(t *testing.T) {
	svr := http.Server{
		Addr: "",
	}

	http.HandleFunc("/header", httpHeader)
	http.HandleFunc("/body", httpBody)
	svr.ListenAndServe()
}

func process(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintln(w, r.Form)
}

func processHtml_Post(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>Test Post</title>
			</head>
			<body>
				<form action="/process" method="post"	enctype="application/x-www-form-urlencoded">
					<input type="text" name="first_name"/>
					<input type="text" name="last_name"/>
					<input type="submit"/>
				</form>
			</body>
			</html>
	`
	fmt.Fprintln(w, html)
}

func Test_HttpForm_Post(t *testing.T) {
	svr := http.Server{
		Addr: "",
	}

	http.HandleFunc("/process", process)
	http.HandleFunc("/", processHtml_Post)
	svr.ListenAndServe()
}

func processHtml_Get(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>Test Post</title>
			</head>
			<body>
				<form action="/process" method="get">
					<input type="text" name="first_name"/>
					<input type="text" name="last_name"/>
					<input type="submit"/>
				</form>
			</body>
			</html>
	`
	fmt.Fprintln(w, html)
}

func Test_HttpForm_Get(t *testing.T) {
	svr := http.Server{
		Addr: "",
	}

	http.HandleFunc("/process", process)
	http.HandleFunc("/", processHtml_Get)
	svr.ListenAndServe()
}

func processFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	fileHeader := r.MultipartForm.File["uploaded"][0]
	file, err := fileHeader.Open()
	if err == nil {
		data, err := ioutil.ReadAll(file)
		if err == nil {
			fmt.Fprintln(w, string(data))
		}
	}
}

type Post struct {
	User    string
	Threads []string
}

func processHtml_File(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
				<title>Go Web Programming</title>
				</head>
			<body>
				<form action="/process?hello=world&thread=123"	method="post" enctype="multipart/form-data">
					<input type="text" name="hello" value="sau sheong"/>
					<input type="text" name="post" value="456"/>
					<input type="file" name="uploaded">
					<input type="submit">
				</form>
			</body>
			</html>
	`
	fmt.Fprintln(w, html)
}

func jsonExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	post := &Post{
		User:    "Sau Sheong",
		Threads: []string{"first", "second", "third"},
	}
	json, _ := json.Marshal(post)
	w.Write(json)
}

func redirectExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "http://baidu.com")
	w.WriteHeader(302)
}

func Test_HttpForm_File(t *testing.T) {
	svr := http.Server{
		Addr: "",
	}

	http.HandleFunc("/process", processFile)
	http.HandleFunc("/", processHtml_File)
	http.HandleFunc("/json", jsonExample)
	http.HandleFunc("/redirect", redirectExample)

	svr.ListenAndServe()
}

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	Raw        string
	Unparsed   []string
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	c1 := http.Cookie{
		Name:     "first_cookie",
		Value:    "Go Web Programming",
		HttpOnly: true,
	}

	c2 := http.Cookie{
		Name:     "second_cookie",
		Value:    "Manning Publications Co",
		HttpOnly: true,
	}
	w.Header().Set("Set-Cookie", c1.String())
	w.Header().Add("Set-Cookie", c2.String())
}

func getCookie(w http.ResponseWriter, r *http.Request) {
	h1, _ := r.Cookie("first_cookie")
	h := r.Header["Cookie"]
	fmt.Fprintln(w, h1)
	fmt.Fprintln(w, h)
}

func Test_Http_Cookie(t *testing.T) {
	svr := http.Server{
		Addr: "",
	}

	http.HandleFunc("/setcookie", setCookie)
	http.HandleFunc("/getcookie", getCookie)

	svr.ListenAndServe()
}
