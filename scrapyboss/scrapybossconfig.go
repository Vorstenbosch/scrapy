package scrapyboss

import (
	"fmt"

	"github.com/vorstenbosch/scrapy/scrapy"
	"gopkg.in/yaml.v2"
)

// ScrapeConfig is the struct representation of the config yaml file
// TODO: add general configurations like; retries, storage, api protection etc.
type ScrapyBossConfig struct {
	ScrapeEndpoints         []ScrapeEndpoint
	ScrapeIntervalInSeconds int
	IdleConnectionPool      int
	ScrapeTimeoutInSeconds  int
}

type ScrapeEndpoint struct {
	Endpoint  string
	Selectors []scrapy.Selector
}

func ParseConfig(b []byte) (ScrapyBossConfig, error) {
	var config ScrapyBossConfig
	var err error

	err = yaml.Unmarshal(b, &config)

	return config, err
}

// ValidateConfig validates the config on a functional level
// It prevents unwanted behaviour (e.g. DOS-ing a scrape target due to a invalid scrape interval)
func ValidateConfig(c ScrapyBossConfig) []error {
	var errorList []error

	if c.ScrapeIntervalInSeconds < 1 {
		errorList = append(errorList, fmt.Errorf("Config setting of 'ScrapeIntervalInSeconds' must be higher than 0"))
	}

	if c.ScrapeIntervalInSeconds < c.ScrapeTimeoutInSeconds {
		errorList = append(errorList, fmt.Errorf("Config setting of 'ScrapeIntervalInSeconds' must be higher than '%v'", c.ScrapeTimeoutInSeconds))
	}

	if c.IdleConnectionPool < 1 {
		errorList = append(errorList, fmt.Errorf("Config setting of 'IdleConnectionPool' must be higher than '0' but '%v' was configured", c.IdleConnectionPool))
	}

	return errorList
}
