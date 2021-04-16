package scrapyboss

import (
	"testing"

	"github.com/vorstenbosch/scrapy/scrapy"
)

func TestParseScrapyConfig(t *testing.T) {
	// Given
	yamlConfig := `
scrapeintervalinseconds: 2
idleconnectionpool: 1
scrapetimeoutinseconds: 1
scrapeendpoints: 
  - endpoint: http://localhost:5555
    selectors:
      - name: test
        typeofselector: xpath
        value: //div
`

	expectedConfig := ScrapyBossConfig{
		ScrapeEndpoints: []ScrapeEndpoint{
			ScrapeEndpoint{
				Endpoint: "http://localhost:5555",
				Selectors: []scrapy.Selector{
					scrapy.Selector{
						Name:           "test",
						TypeOfSelector: "xpath",
						Value:          "//div",
					},
				},
			},
		},
		ScrapeIntervalInSeconds: 2,
		ScrapeTimeoutInSeconds:  1,
		IdleConnectionPool:      1,
	}

	// When
	config, err := ParseConfig([]byte(yamlConfig))

	// Then
	if err != nil {
		t.Errorf("ParseConfig test failed due to unexpected error '%v'", err)
	}

	if len(config.ScrapeEndpoints) != len(expectedConfig.ScrapeEndpoints) {
		t.Errorf("ParseConfig test failed because resulting config '%v' is not the same as what was expected '%v'", config, expectedConfig)
	}

	if len(config.ScrapeEndpoints[0].Selectors) != len(expectedConfig.ScrapeEndpoints[0].Selectors) {
		t.Errorf("ParseConfig test failed because resulting config '%v' is not the same as what was expected '%v'", config, expectedConfig)
	}

	if config.ScrapeIntervalInSeconds != expectedConfig.ScrapeIntervalInSeconds {
		t.Errorf("ParseConfig test failed because resulting config '%v' is not the same as what was expected '%v'", config, expectedConfig)
	}

	if config.ScrapeTimeoutInSeconds != expectedConfig.ScrapeTimeoutInSeconds {
		t.Errorf("ParseConfig test failed because resulting config '%v' is not the same as what was expected '%v'", config, expectedConfig)
	}

	if config.IdleConnectionPool != expectedConfig.IdleConnectionPool {
		t.Errorf("ParseConfig test failed because resulting config '%v' is not the same as what was expected '%v'", config, expectedConfig)
	}
}

func TestValidateConfig(t *testing.T) {
	// Given
	config := ScrapyBossConfig{
		ScrapeEndpoints: []ScrapeEndpoint{
			ScrapeEndpoint{
				Endpoint: "http://localhost:5555",
				Selectors: []scrapy.Selector{
					scrapy.Selector{
						Name:           "test",
						TypeOfSelector: "xpath",
						Value:          "//div",
					},
				},
			},
		},
		ScrapeIntervalInSeconds: 2,
		ScrapeTimeoutInSeconds:  1,
		IdleConnectionPool:      1,
	}

	// When
	errorList := ValidateConfig(config)

	// Then
	if len(errorList) != 0 {
		t.Errorf("Expecting that the config was valid but errors '%v' where found", errorList)
	}
}

func TestInvalidConfig(t *testing.T) {
	// Given
	config := ScrapyBossConfig{
		ScrapeEndpoints: []ScrapeEndpoint{
			ScrapeEndpoint{
				Endpoint: "http://localhost:5555",
				Selectors: []scrapy.Selector{
					scrapy.Selector{
						Name:           "test",
						TypeOfSelector: "xpath",
						Value:          "//div",
					},
				},
			},
		},
		ScrapeIntervalInSeconds: -1,
		ScrapeTimeoutInSeconds:  1,
		IdleConnectionPool:      -4,
	}

	// When
	errorList := ValidateConfig(config)

	// Then
	if len(errorList) != 3 {
		t.Errorf("Expected '3' errors but found only '%v' when validating the config", len(errorList))
	}
}
