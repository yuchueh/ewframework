package logs

import (
	"time"
	"bytes"
	"testing"
)

func TestFormatHeader_0(t *testing.T) {
	tm := time.Now()
	if tm.Year() >= 2100 {
		t.FailNow()
	}
	dur := time.Second
	for {
		if tm.Year() >= 2100 {
			break
		}
		h, _ := formatTimeHeader(tm)
		if tm.Format("2006/01/02 15:04:05 ") != string(h) {
			t.Log(tm)
			t.FailNow()
		}
		tm = tm.Add(dur)
		dur *= 2
	}
}

func TestFormatHeader_1(t *testing.T) {
	tm := time.Now()
	year := tm.Year()
	dur := time.Second
	for {
		if tm.Year() >= year+1 {
			break
		}
		h, _ := formatTimeHeader(tm)
		if tm.Format("2006/01/02 15:04:05 ") != string(h) {
			t.Log(tm)
			t.FailNow()
		}
		tm = tm.Add(dur)
	}
}

func TestNewAnsiColor1(t *testing.T) {
	inner := bytes.NewBufferString("")
	w := NewAnsiColorWriter(inner)
	if w == inner {
		t.Errorf("Get %#v, want %#v", w, inner)
	}
}

func TestNewAnsiColor2(t *testing.T) {
	inner := bytes.NewBufferString("")
	w1 := NewAnsiColorWriter(inner)
	w2 := NewAnsiColorWriter(w1)
	if w1 != w2 {
		t.Errorf("Get %#v, want %#v", w1, w2)
	}
}

