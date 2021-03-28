package scrapy

import (
	"fmt"
)

type scrapyBoss struct {
	config       ScrapyBossConfig
	scrapingData map[string][]ScrapeResult
	running      bool
}

func NewScrapyBoss(c ScrapyBossConfig) *scrapyBoss {
	return &scrapyBoss{
		config:       c,
		scrapingData: map[string][]ScrapeResult{},
		running:      false,
	}
}

func (s *scrapyBoss) Start() error {
	var err error

	if s.running {
		err = fmt.Errorf("This ScrapyBoss is already running")
	} else {
		s.running = true
		go s.collect(s.config.ScrapeEndpoints, s.config.ScrapeIntervalInSeconds)
	}

	return err
}

func (s *scrapyBoss) Stop() {
	s.running = false
}

func (s *scrapyBoss) collect(endpoints []ScrapeEndpoint, interval int) {
	for s.running {
		for i := range endpoints {
			Scrape(endpoints[i].Endpoint, endpoints[i].Selectors)
		}
	}
}
