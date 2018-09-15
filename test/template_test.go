package test

import (
	"github.com/yuchueh/ewframework/utils/osext"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
)

func processTemplate01(w http.ResponseWriter, r *http.Request) {
	sHtmlFile := osext.GetWd() + "html" + string(os.PathSeparator) + "tmpl01.html"
	//fmt.Println(sHtmlFile)
	temp, _ := template.ParseFiles(sHtmlFile)
	temp.Execute(w, "Hello world!")
	//temp.ExecuteTemplate(w, "html" + string(os.PathSeparator) + "tmpl01.html", "Hello world!")
}

func processTemplate02(w http.ResponseWriter, r *http.Request) {
	sHtmlFile := osext.GetWd() + "html" + string(os.PathSeparator) + "tmpl02.html"
	//fmt.Println(sHtmlFile)
	temp, _ := template.ParseFiles(sHtmlFile)
	temp.Execute(w, rand.Intn(10) > 5)
}

func processTemplate03(w http.ResponseWriter, r *http.Request) {
	sHtmlFile := osext.GetWd() + "html" + string(os.PathSeparator) + "tmpl03.html"
	//fmt.Println(sHtmlFile)
	temp, _ := template.ParseFiles(sHtmlFile)
	month := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	temp.Execute(w, month)
}

func processTemplate04(w http.ResponseWriter, r *http.Request) {
	sHtmlFile := osext.GetWd() + "html" + string(os.PathSeparator) + "tmpl04.html"
	//fmt.Println(sHtmlFile)
	temp, _ := template.ParseFiles(sHtmlFile)
	month := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	temp.Execute(w, month)
}

func processTemplate05(w http.ResponseWriter, r *http.Request) {
	sHtmlFile := osext.GetWd() + "html" + string(os.PathSeparator) + "tmpl05.html"
	sHtmlFile_1 := osext.GetWd() + "html" + string(os.PathSeparator) + "tmpl05_1.html"
	//fmt.Println(sHtmlFile)
	temp, _ := template.ParseFiles(sHtmlFile, sHtmlFile_1)
	temp.Execute(w, "Hello")
}

func formatDate(t time.Time) string {
	layout := "2006-01-02"
	return t.Format(layout)
}

func processTemplate06(w http.ResponseWriter, r *http.Request) {
	sHtmlFile := osext.GetWd() + "html" + string(os.PathSeparator) + "tmpl06.html"
	//fmt.Println(sHtmlFile)
	funcsMap := template.FuncMap { "fdate": formatDate }
	t := template.New("tmpl06.html").Funcs(funcsMap)
	//temp, _ := template.ParseFiles(sHtmlFile, sHtmlFile)
	//temp.Funcs(funcsMap)
	//temp.Execute(w, time.Now())
	t, _ = t.ParseFiles(sHtmlFile)
	t.Execute(w, time.Now())
}

func processTemplate07(w http.ResponseWriter, r *http.Request) {
	sHtmlFile := osext.GetWd() + "html" + string(os.PathSeparator) + "tmpl07.html"
	temp, _ := template.ParseFiles(sHtmlFile)
	context := `I asked: <i>"What's up?'"</i>`
	temp.Execute(w, context)
}

func form(w http.ResponseWriter, r *http.Request) {
	sHtmlFile := osext.GetWd() + "html" + string(os.PathSeparator) + "form.html"
	temp, _ := template.ParseFiles(sHtmlFile)
	temp.Execute(w, nil)
}

func processform(w http.ResponseWriter, r *http.Request) {
	//stop the browser from protecting you from XSS attacks, you
	//simply need to set a response header in our handler
	//0： 表示关闭浏览器的XSS防护机制
	//1：删除检测到的恶意代码， 如果响应报文中没有看到X-XSS-Protection 字段，那么浏览器就认为X-XSS-Protection配置为1，这是浏览器的默认设置
	//1; mode=block：如果检测到恶意代码，在不渲染恶意代码
	w.Header().Set("X-XSS-Protection", "0")
	sHtmlFile := osext.GetWd() + "html" + string(os.PathSeparator) + "processform.html"
	temp, _ := template.ParseFiles(sHtmlFile)
	temp.Execute(w, template.HTML(r.FormValue("comment")))
	//temp.Execute(w, r.FormValue("comment"))
}

func processlayout(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().Unix())
	sHtmlFile := osext.GetWd() + "html" + string(os.PathSeparator) + "layout.html"
	sHtmlFile_red := osext.GetWd() + "html" + string(os.PathSeparator) + "red_hello.html"
	sHtmlFile_blue := osext.GetWd() + "html" + string(os.PathSeparator) + "blud_hello.html"

	var t *template.Template
	if rand.Intn(10) > 5 {
		t, _ = template.ParseFiles(sHtmlFile, sHtmlFile_red)
	} else {
		t, _ = template.ParseFiles(sHtmlFile, sHtmlFile_blue)
	}
	//temp, _ := template.ParseFiles(sHtmlFile)
	//temp.ExecuteTemplate(w, "layout", "")
	t.ExecuteTemplate(w, "layout", "")
}

func Test_Template(t *testing.T) {
	svr := http.Server{
		Addr: "",
	}

	http.HandleFunc("/process01", processTemplate01)
	http.HandleFunc("/process02", processTemplate02)
	http.HandleFunc("/process03", processTemplate03)
	http.HandleFunc("/process04", processTemplate04)
	http.HandleFunc("/process05", processTemplate05)
	http.HandleFunc("/process06", processTemplate06)
	http.HandleFunc("/process07", processTemplate07)
	http.HandleFunc("/form", form)
	http.HandleFunc("/process", processform)
	http.HandleFunc("/processlayout", processlayout)
	svr.ListenAndServe()
}
