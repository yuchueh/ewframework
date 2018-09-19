package param

import (
	"fmt"
	"strings"
)

//MethodParam keeps param information to be auto passed to controller methods
type MethodParam struct {
	name         string
	in           paramType
	required     bool
	defaultValue string
}

type paramType byte

const (
	param paramType = iota
	path
	body
	header
)


func (mp *MethodParam) String() string {
	options := []string{}
	result := "param.New(\"" + mp.name + "\""
	if mp.required {
		options = append(options, "param.IsRequired")
	}
	switch mp.in {
	case path:
		options = append(options, "param.InPath")
	case body:
		options = append(options, "param.InBody")
	case header:
		options = append(options, "param.InHeader")
	}
	if mp.defaultValue != "" {
		options = append(options, fmt.Sprintf(`param.Default("%s")`, mp.defaultValue))
	}
	if len(options) > 0 {
		result += ", "
	}
	result += strings.Join(options, ", ")
	result += ")"
	return result
}

