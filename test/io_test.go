package test

import (
	"testing"
	"io/ioutil"
	"fmt"
	"os"
	"bytes"
	"encoding/gob"
)

func Test_IO(t *testing.T)  {
	data := []byte("Hello world!\n")
	err := ioutil.WriteFile("data/data1", data, 0644)
	if err != nil {
		panic(err)
	}

	dat, _ := ioutil.ReadFile("data/data1")
	fmt.Println(string(dat))

	f, _ := os.Create("data/data2")
	defer f.Close()
	bytes, _ := f.Write(data)
	fmt.Printf("Write %d bytes to file\n", bytes)

	read1 := make([]byte, bytes)
	f, _ = os.Open("data/data2")
	defer f.Close()
	f.Read(read1)
	fmt.Println(string(read1))

}

type PostData struct {
	Id		int
	Content	string
	Author	string
}

func store(data interface{}, filename string)  {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)

	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filename, buffer.Bytes(), 600)
	if err != nil {
		panic(err)
	}
}

func load(data interface{}, filename string)  {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		panic(err)
	}
}

func Test_gob(t *testing.T)  {
	post := PostData{Id: 1, Content: "Hello World!", Author: "hai he"}
	store(post, "data/post1")

	var postdata PostData
	load(&postdata, "data/post1")
	fmt.Println(postdata)
}