package scrapy

import (
	"context"
	"io"
	"log"
	"net/http"
	"testing"
)

var server http.Server

func TestScrape(t *testing.T) {
	// Given
	setup()
	endpoint := "http://localhost:5555/"
	selector := Selector{
		typeOfSelector: "xpath",
		value:          "//div",
	}

	expectedResult := "Hello world"

	// When
	result, err := Scrape(endpoint, selector)

	// Then
	if err != nil {
		t.Errorf("Scrape test failed due to unexpected error '%v'", err)
	}

	if result != expectedResult {
		t.Errorf("Scrape test failed as the result did not match '%s' but was '%s'", expectedResult, result)
	}

}

func setup() {
	defer tearDown()
	server := &http.Server{Addr: ":5555"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<div>Hello world</div>")
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
