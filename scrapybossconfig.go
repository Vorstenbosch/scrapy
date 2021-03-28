package scrapy

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// ScrapeConfig is the struct representation of the config yaml file
// TODO: add general configurations like; retries, storage, api protection etc.
type ScrapyBossConfig struct {
	ScrapeEndpoints         []ScrapeEndpoint
	ScrapeIntervalInSeconds int
}

type ScrapeEndpoint struct {
	Endpoint  string
	Selectors []Selector
}

func ParseConfig(b []byte) (ScrapyBossConfig, error) {
	var config ScrapyBossConfig
	var err error

	err = yaml.Unmarshal(b, &config)

	if config.ScrapeIntervalInSeconds < 1 {
		config = ScrapyBossConfig{}
		err = fmt.Errorf("Config setting of 'ScrapeIntervalInSeconds' must be higher than 0")
	}

	return config, err
}
