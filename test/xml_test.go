package test

import (
	"testing"
	"encoding/xml"
)

type post struct {
	XMLName		xml.Name	`xml:"post"`
	Id			string		`xml:"id,attr"`
	Content 	string 		`xml:"content"`
	Author 		Author 		`xml:"author"`
	Xml 		string 		`xml:",innerxml"`
}

type Author struct {
	Id 			string 		`xml:"id,attr"`
	Name 		string 		`xml:",chardata"`
}

func Test_XML(t *testing.T)  {

}
