package scrapy

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

var server http.Server
var handlerInitiated bool = false

func TestScrapeXpath(t *testing.T) {
	// Given
	endpoint := "http://localhost:5555/"
	selector := Selector{
		Name:           "xpath-scrape",
		TypeOfSelector: "xpath",
		Value:          "//div",
	}

	expectedResult := "Hello world"

	// When
	result := Scrape(endpoint, []Selector{selector})

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

func TestScrapeRegex(t *testing.T) {
	// Given
	endpoint := "http://localhost:5555/"
	selector := Selector{
		TypeOfSelector: "regex",
		Value:          "H(.*?) world",
	}

	expectedResult := "ello"

	// When
	result := Scrape(endpoint, []Selector{selector})

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

func TestScrapeRegexMultipleMatches(t *testing.T) {
	// Given
	endpoint := "http://localhost:5555/"
	selector := Selector{
		TypeOfSelector: "regex",
		Value:          "^(.)(.)+$",
	}

	// When
	result := Scrape(endpoint, []Selector{selector})

	// Then
	if len(result) != 1 {
		t.Errorf("Scrape test failed due to incorrect result list length")
	}

	if result[0].Error == nil {
		t.Errorf("Scrape test failed due to missing error")
	}

	if result[0].Error.Error() != "Regex selector resulted into multiple group matches ('2'), only one is allowed" {
		t.Errorf("Error is not as expected but is '%s'", result[0].Error.Error())
	}
}

func TestInvalidSelectorType(t *testing.T) {
	// Given
	endpoint := "http://localhost:5555/"
	selector := Selector{
		TypeOfSelector: "NOT_VALID",
		Value:          "//div",
	}

	// When
	result := Scrape(endpoint, []Selector{selector})

	// Then
	if len(result) != 1 {
		t.Errorf("Scrape test failed due to incorrect result list length")
	}

	if result[0].Error == nil {
		t.Errorf("Scrape test failed due to missing error")
	}

	if result[0].Error.Error() != "Selector type 'NOT_VALID' is not supported" {
		t.Errorf("Error is not as expected but is '%s'", result[0].Error.Error())
	}
}

func setup() {
	server := &http.Server{Addr: ":5555"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<div>Hello world</div>")
	})

	handlerInitiated = true

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
	os.Exit(m.Run())
	tearDown()
}
