package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("hello world of songs")
	exp := 4

	res := count(b, false)
	if res != exp {
		t.Errorf("Expected %v, got %v", exp, res)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("hello\nworld\nof\nsongs")
	exp := 4

	res := count(b, true)
	if res != exp {
		t.Errorf("Expected %v, got %v", exp, res)
	}
}
