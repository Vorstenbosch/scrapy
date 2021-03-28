package scrapyboss

import (
	"testing"

	"github.com/vorstenbosch/scrapy/scrapy"
)

func TestScrapyConfig(t *testing.T) {
	// Given
	yamlConfig := `
scrapeintervalinseconds: 1
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
		ScrapeIntervalInSeconds: 1,
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

	// TODO: add other assetions
}

func TestInvalidConfig(t *testing.T) {
	// Given
	yamlConfig := `
scrapeintervalinseconds: 0
scrapeendpoints: 
  - endpoint: http://localhost:5555
    selectors:
      - name: test
        typeofselector: xpath
        value: //div
`

	// When
	_, err := ParseConfig([]byte(yamlConfig))

	// Then
	if err == nil {
		t.Errorf("ParseConfig test failed due to missing expected error")
	}

	if err.Error() != "Config setting of 'ScrapeIntervalInSeconds' must be higher than 0" {
		t.Errorf("ParseConfig test failed because expected error was not as expected '%v'", err)
	}
}
