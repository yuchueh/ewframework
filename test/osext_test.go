package test

import (
	"testing"
	"github.com/yuchueh/ewframework/utils/osext"
)

func Test_osext(t *testing.T)  {
	s, err := osext.Executable()
	if err != nil {
		t.Fatal("Bad test case Executable:", err)
	} else {
		t.Log("Executable:", s)
	}

	b := osext.HaveReadPermission(s)
	t.Log(s, "HaveReadPermission:", b)

	b = osext.HaveWritePermission(s)
	t.Log(s, "HaveWritePermission:", b)

	b = osext.HaveRWPermission(s)
	t.Log(s, "HaveRWPermission:", b)

	s, err = osext.ExecutableFolder()
	if err != nil {
		t.Fatal("Bad test case ExecutableFolder:", err)
	} else {
		t.Log("ExecutableFolder:", s)
	}

	b = osext.Exits(s)
	t.Log(s, "osext.Exits:", b)

	fi := osext.SearchDir(s)
	t.Log(fi)
	for i,v := range fi {
		t.Log(i, v)
	}

	s, err = osext.GetCurrentPath()
	if err != nil {
		t.Fatal("Bad test case GetCurrentPath:", err)
	} else {
		t.Log("GetCurrentPath:", s)
	}

	s, err = osext.GetWd()
	if err != nil {
		t.Fatal("Bad test case GetWd:", err)
	} else {
		t.Log("GetWd:", s)
	}
}

func Benchmark_osext(b *testing.B)  {
	s, err := osext.ExecutableFolder()
	if err != nil {
		b.Fatal("Bad test case ExecutableFolder:", err)
	} else {
		b.Log("ExecutableFolder:", s)
	}

	bb := osext.Exits(s)
	b.Log(s, "osext.Exits:", bb)

	fi := osext.SearchDir(s)
	b.Log(fi)
	for i,v := range fi {
		b.Log(i, v)
	}
}

func Benchmark_GetCurrentFileName(b *testing.B)  {
	b.Log(osext.GetCurrentFileName(true))
	b.Log(osext.GetCurrentFileName(false))
}