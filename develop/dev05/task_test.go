package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func createTempFile(t *testing.T, content string) string {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name()
}

func removeTempFile(t *testing.T, name string) {
	if err := os.Remove(name); err != nil {
		t.Fatal(err)
	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputCh := make(chan string)
	go func() {
		var buf strings.Builder
		io.Copy(&buf, r)
		outputCh <- buf.String()
	}()

	f()

	w.Close()
	output := <-outputCh
	os.Stdout = old
	return output
}

func TestGrep(t *testing.T) {
	content := "Hello World\nThis is a test\nAnother test line\nHello again"
	tmpfile := createTempFile(t, content)
	defer removeTempFile(t, tmpfile)

	tests := []struct {
		options  grepOptions
		pattern  string
		expected string
	}{
		{grepOptions{}, "Hello", "Hello World\nHello again\n"},
		{grepOptions{ignoreCase: true}, "hello", "Hello World\nHello again\n"},
		{grepOptions{invert: true}, "Hello", "This is a test\nAnother test line\n"},
		{grepOptions{fixed: true}, "Hello World", "Hello World\n"},
		{grepOptions{lineNum: true}, "Hello", "1:Hello World\n4:Hello again\n"},
		{grepOptions{count: true}, "Hello", "2\n"},
		{grepOptions{after: 1}, "Hello", "Hello World\nThis is a test\nHello again\n"},
	}

	for _, test := range tests {
		output := captureOutput(func() {
			pattern := compilePattern(test.pattern, test.options)
			grepFile(tmpfile, pattern, test.options)
		})

		if output != test.expected {
			t.Errorf("For pattern222 %s with options %+v, expected %q but got %q", test.pattern, test.options, test.expected, output)
		}
	}
}
