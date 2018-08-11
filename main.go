package main

import (
	"fmt"
	"github.com/yuchueh/ewframework/config"
	"runtime"
	"path"
	"path/filepath"
	"strings"
)

func main()  {
	fmt.Println("this is ewframework")

	c, err := config.ReadDefault("config.cfg")
	if err != nil {
		fmt.Println(err)
	} else {
		s, _ := c.String("DEFAULT", "url")
		fmt.Println(s)
	}

	_, filename, line, _ := runtime.Caller(0)
	fmt.Println(filename, line, path.Base(filename), path.Ext(filename), filepath.Ext(filename), strings.Replace(path.Base(filename), filepath.Ext(filename), "", -1))


	//filepath.Ext()

}
