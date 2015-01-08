package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestByLength(t *testing.T) {
	expected := []string{
		"1",
		"123",
		"1234",
		"12345",
		"123456",
		"1234567",
		"12345678",
		"123456789",
	}

	testArr := []string{
		"123456",
		"123456789",
		"1234",
		"12345",
		"1",
		"12345678",
		"123",
		"1234567",
	}

	sort.Sort(ByLength(testArr))
	for idx, val := range testArr {
		if val != expected[idx] {
			t.Errorf("Expected ordered list, got %v", testArr)
			break
		}
	}
}

func TestPossibleVersions(t *testing.T) {
	tests := []struct {
		url      string
		versions []string
	}{
		{"github.com/rsenk330/gogetver", []string{"master"}},
		{"github.com/rsenk330/gogetver.develop", []string{"develop"}},
		{"github.com/rsenk330/gogetver.v1", []string{"v1"}},
		{"github.com/rsenk330/gogetver.v1.0", []string{"v1.0", "0"}},
		{"github.com/rsenk330/gogetver.v1.0a", []string{"v1.0a", "0a"}},
	}

	for _, test := range tests {
		versions := PossibleVersions(test.url)

		for idx, eVer := range test.versions {
			if versions[idx] != eVer {
				t.Errorf("Expected versions to equal %v, got %v", test.versions, versions)
				break
			}
		}
	}
}

func TestApp_Home(t *testing.T) {
	app := NewApp(NewAppConfig())

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	// Test normal case
	w := httptest.NewRecorder()
	app.Home(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 response code, got %v.", w.Code)
	}

	doc, err := goquery.NewDocumentFromReader(w.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	if !doc.Find("body").HasClass("home") {
		t.Error("Expected <body class=\"home\"></body> to exist.")
	}

	// Test when template does not exist
	config := NewAppConfig()
	config.TemplatesDir = "/tmp/gdags"
	app = NewApp(config)

	w = httptest.NewRecorder()
	app.Home(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 response code, got %v", w.Code)
	}
}

func TestApp_Package(t *testing.T) {
	app := NewApp(NewAppConfig())
	r := Router(app)
	server := httptest.NewServer(r)
	defer server.Close()

	url, err := r.Get("package").URL("pkg", "github.com/rsenk330/gogetver")
	if err != nil {
		t.Error(err)
	}

	// Test non-go get request (render package page)
	res, err := http.Get(fmt.Sprintf("%v%v", server.URL, url))
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 response code, got %v", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	if !doc.Find("body").HasClass("package") {
		t.Error("Expected <body class=\"package\"></body> to exist.")
	}

	// Test go get request
	res, err = http.Get(fmt.Sprintf("%v%v?go-get=1", server.URL, url))
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 response code, got %v", res.StatusCode)
	}

	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	importPath, exists := doc.Find("meta[name=go-import]").Attr("content")
	if exists && importPath != "gogetver.com/github.com/rsenk330/gogetver git https://gogetver.com/github.com/rsenk330/gogetver" {
		t.Errorf("Expected go-import meta content to be 'gogetver.com/github.com/rsenk330/gogetver git https://gogetver.com/github.com/rsenk330/gogetver', got '%v'", importPath)
	}
}
