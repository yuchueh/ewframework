package main

import (
	"fmt"
	"github.com/yuchueh/ewframework/ew"
	"github.com/yuchueh/ewframework/controller"
)

type myController struct {
	controller.Controller
}

func main() {
	fmt.Println("this is ewframework start")

	ew.AutoRouter(&myController{})
	//ew.Router("/", &myController{}, "*:Index")
	ew.Run()

	fmt.Println("this is ewframework end")
}
