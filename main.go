package main

import (
	"fmt"
	"net/http"
)

type MyStruct struct {
}

func (s *MyStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func main() {
	fmt.Println("this is ewframework start")

	handler := MyStruct{}
	//http.ListenAndServe("", nil)
	svr := http.Server{
		Addr: "",
		//Handler:&handler,
	}

	http.Handle("/hello", &handler)
	svr.ListenAndServe()

	fmt.Println("this is ewframework end")
}
