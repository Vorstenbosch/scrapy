package scrapy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

var server http.Server

func TestScrapeXpath(t *testing.T) {
	// Given
	endpoint := "http://localhost:5555/"
	selector := Selector{
		Name:           "xpath-scrape",
		TypeOfSelector: "xpath",
		Value:          "//div",
	}

	expectedResult := "Hello world"
	c := NewScrapeClient()

	// When
	var result []ScrapeResult
	c.Scrape(endpoint, []Selector{selector}, &result)

	// Then
	if len(result) != 1 {
		t.Errorf("Scrape test failed due to incorrect result list length")
	}

	if result[0].Error != nil {
		t.Errorf("Scrape test failed due to unexpected error '%v'", result[0].Error)
	}

	if result[0].Result != expectedResult {
		t.Errorf("Scrape test failed as the result did not match '%s' but was '%s'", expectedResult, result)
	}

	if result[0].Name != "xpath-scrape" {
		t.Errorf("Scrape test failed as the name did not match '%s' but was '%s'", "xpath-scrape", result)
	}
}

// The xpath library failed on invalid html pages. Therefore we are now 'fixing' the scrape target
//   before applying the xpatch selector. This test is testing xpath scraping on invalid html.
func TestScrapeXpathInvalidHtml(t *testing.T) {
	// Given
	endpoint := "http://localhost:5555/invalid-html"
	selector := Selector{
		Name:           "xpath-scrape",
		TypeOfSelector: "xpath",
		Value:          "//div",
	}

	expectedResult := "Hello world"
	c := NewScrapeClient()

	// When
	var result []ScrapeResult
	c.Scrape(endpoint, []Selector{selector}, &result)

	// Then
	if len(result) != 1 {
		t.Errorf("Scrape test failed due to incorrect result list length")
	}

	if result[0].Error != nil {
		t.Errorf("Scrape test failed due to unexpected error '%v'", result[0].Error)
	}

	if result[0].Result != expectedResult {
		t.Errorf("Scrape test failed as the result did not match '%s' but was '%s'", expectedResult, result)
	}

	if result[0].Name != "xpath-scrape" {
		t.Errorf("Scrape test failed as the name did not match '%s' but was '%s'", "xpath-scrape", result)
	}
}

func TestScrapeXpathNotFound(t *testing.T) {
	// Given
	endpoint := "http://localhost:5555/invalid-html"
	selector := Selector{
		Name:           "xpath-scrape",
		TypeOfSelector: "xpath",
		Value:          "//table",
	}

	expectedResult := fmt.Errorf("Unable to find '//table'")
	c := NewScrapeClient()

	// When
	var result []ScrapeResult
	c.Scrape(endpoint, []Selector{selector}, &result)

	// Then
	if len(result) != 1 {
		t.Errorf("Scrape test failed due to incorrect result list length")
	}

	if result[0].Error == fmt.Errorf("Unable to find '//table'") {
		t.Errorf("Scrape test failed, expected error '%v' but received '%v'", expectedResult, result[0].Error)
	}

	if result[0].Name != "xpath-scrape" {
		t.Errorf("Scrape test failed as the name did not match '%s' but was '%s'", "xpath-scrape", result)
	}
}

func TestScrapeRegex(t *testing.T) {
	// Given
	endpoint := "http://localhost:5555/"
	selector := Selector{
		TypeOfSelector: "regex",
		Value:          "H(.*?) world",
	}

	expectedResult := "ello"
	c := NewScrapeClient()

	// When
	var result []ScrapeResult
	c.Scrape(endpoint, []Selector{selector}, &result)

	// Then
	if len(result) != 1 {
		t.Errorf("Scrape test failed due to incorrect result list length")
	}

	if result[0].Error != nil {
		t.Errorf("Scrape test failed due to unexpected error '%v'", result[0].Error)
	}

	if result[0].Result != expectedResult {
		t.Errorf("Scrape test failed as the result did not match '%s' but was '%s'", expectedResult, result)
	}
}

func setup() {
	server := &http.Server{Addr: ":5555"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<div>Hello world</div>")
	})

	http.HandleFunc("/invalid-html", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<div>Hello world</div></ul>")
	})

	go func() {
		// always returns error. ErrServerClosed on graceful close
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
}

func tearDown() {
	// now close the server gracefully ("shutdown")
	// timeout could be given with a proper context
	// (in real world you shouldn't use TODO()).
	if err := server.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

}

func TestMain(m *testing.M) {
	setup()
	tearDown()
	os.Exit(m.Run())
}
