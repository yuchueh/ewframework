package test

import (
	"github.com/yuchueh/ewframework/utils/osext"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"testing"
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

func Test_Template(t *testing.T) {
	svr := http.Server{
		Addr: "",
	}

	http.HandleFunc("/process01", processTemplate01)
	http.HandleFunc("/process02", processTemplate02)
	http.HandleFunc("/process03", processTemplate03)
	http.HandleFunc("/process04", processTemplate04)
	http.HandleFunc("/process05", processTemplate05)
	svr.ListenAndServe()
}
