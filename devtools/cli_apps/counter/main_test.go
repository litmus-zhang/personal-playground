package counter

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("hello world of songs\n")
	exp := 4

	res := count(b, false, false)
	if res != exp {
		t.Errorf("Expected %v, got %v", exp, res)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 \n line2 \n line3 word1")
	exp := 3

	res := count(b, true, false)
	if res != exp {
		t.Errorf("Expected %v, got %v", exp, res)
	}
}

func TestBytes(t *testing.T) {
	b := bytes.NewBufferString("hello world of songs")
	exp := 20

	res := count(b, false, true)
	if res != exp {
		t.Errorf("Expected %v, got %v", exp, res)
	}
}
