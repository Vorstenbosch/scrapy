package scrapyboss

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/vorstenbosch/scrapy/scrapy"
)

var server http.Server

func TestScrapyBossCreation(t *testing.T) {
	// Given
	config := ScrapyBossConfig{
		ScrapeIntervalInSeconds: 10,
		ScrapeEndpoints: []ScrapeEndpoint{
			ScrapeEndpoint{
				Endpoint: "http://localhost:5555",
			},
		},
	}

	// When
	scrapyBoss := NewScrapyBoss(config)

	// Then
	if scrapyBoss.GetConfig().ScrapeIntervalInSeconds != 10 {
		t.Errorf("Config 'ScrapeIntervalInSeconds' was not '10', as expected, but was '%v'", scrapyBoss.GetConfig().ScrapeIntervalInSeconds)
	}
}

func TestScrapyBoss(t *testing.T) {
	// Given
	config := ScrapyBossConfig{
		ScrapeIntervalInSeconds: 10,
		ScrapeEndpoints: []ScrapeEndpoint{
			ScrapeEndpoint{
				Endpoint: "http://localhost:5555",
				Selectors: []scrapy.Selector{
					scrapy.Selector{
						Name:           "test",
						Value:          "//div",
						TypeOfSelector: "xpath",
					},
				},
			},
		},
	}

	scrapyBoss := NewScrapyBoss(config)

	// When
	scrapyBoss.Start()
	defer scrapyBoss.Stop()

	time.Sleep(5 * time.Second)

	// Then
	if scrapyBoss.GetConfig().ScrapeIntervalInSeconds != 10 {
		t.Errorf("Config 'ScrapeIntervalInSeconds' was not '10', as expected, but was '%v'", scrapyBoss.GetConfig().ScrapeIntervalInSeconds)
	}

	if len(scrapyBoss.GetScrapeData()) == 0 {
		t.Errorf("Scrape data was empty")
	}
}

func TestScrapyBossIterations(t *testing.T) {
	// Given
	config := ScrapyBossConfig{
		ScrapeIntervalInSeconds: 5,
		ScrapeEndpoints: []ScrapeEndpoint{
			ScrapeEndpoint{
				Endpoint: "http://localhost:5555",
				Selectors: []scrapy.Selector{
					scrapy.Selector{
						Name:           "test",
						Value:          "//div",
						TypeOfSelector: "xpath",
					},
				},
			},
		},
	}

	scrapyBoss := NewScrapyBoss(config)

	// When
	scrapyBoss.Start()
	defer scrapyBoss.Stop()

	time.Sleep(15 * time.Second)

	// Then
	if !scrapyBoss.IsRunning() {
		t.Errorf("ScrapyBoss should be running but is indicating it is not")
	}

	if scrapyBoss.GetConfig().ScrapeIntervalInSeconds != 5 {
		t.Errorf("Config 'ScrapeIntervalInSeconds' was not '5', as expected, but was '%v'", scrapyBoss.GetConfig().ScrapeIntervalInSeconds)
	}

	if len(scrapyBoss.GetScrapeData()) > 1 {
		t.Errorf("Did not found multiple scrape results")
	}
}

func setup() {
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

func TestMain(m *testing.M) {
	setup()
	tearDown()
	os.Exit(m.Run())
}
