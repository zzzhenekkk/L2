package main

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wget_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	filepath := filepath.Join(tempDir, "testfile.txt")
	err = ioutil.WriteFile(filepath, []byte("This is a test file for wget utility.\n"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestExtractLinks(t *testing.T) {
	htmlContent := `
	<html>
		<body>
			<a href="/test1">link1</a>
			<img src="/image1.jpg">
			<script src="/script1.js"></script>
		</body>
	</html>`

	baseURL, _ := url.Parse("http://example.com")

	links, err := extractLinks(baseURL, strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedLinks := []string{
		"http://example.com/test1",
		"http://example.com/image1.jpg",
		"http://example.com/script1.js",
	}

	for i, link := range links {
		if link != expectedLinks[i] {
			t.Fatalf("expected %q, got %q", expectedLinks[i], link)
		}
	}
}

func TestDownloadPage(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wget_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	filepath := filepath.Join(tempDir, "index.html")
	err = ioutil.WriteFile(filepath, []byte("<html><body>Test</body></html>"), 0644)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
