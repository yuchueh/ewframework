package test

import (
	"testing"
	"github.com/yuchueh/ewframework/config"
	"strings"
)

func TestConfig(t *testing.T)  {
	c, err := config.ReadDefault("testdata/config.cfg")
	if err != nil {
		t.Error(err)
	} else {
		s, _ := c.String("DEFAULT", "url")
		t.Log("url:", s)
		s, _ = c.String("DEFAULT", "protocol")
		t.Log("protocol:", s)

		b := c.AddSection("NewSection")
		t.Log("AddSection NewSection = ", b)
	}

	//config.LoadContext
	ctx, err := config.LoadContext("app.conf", []string{"testdata"})
	if err != nil {
		t.Errorf("Error: %v", err)
		t.FailNow()
	}

	result, _ := ctx.String("one")
	if !strings.EqualFold("source1", result) {
		t.Errorf("Expected '[X] x.two' to be 'override-conf2-sourcex2' but instead it was '%s'", result)
	}

	result, _ = ctx.String("two")
	t.Log(result)

	ctx.SetSection("Y")
	result, found := ctx.String("y.three")

	_, found = ctx.String("y.notexists")
	if found {
		t.Error("Config 'y.notexists' shouldn't be found")
	}

	//two_+_four
	result, _ = ctx.String("two_+_four")
	t.Log("two_+_four：", result)

	result, _ = ctx.String("test")
	t.Log("test：", result)

	ctx.SetSection("X")
	result, _ = ctx.String("x.one")
	t.Log("x.one：", result)

}