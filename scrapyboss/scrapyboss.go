package scrapyboss

import (
	"fmt"
	"time"

	"github.com/vorstenbosch/scrapy/scrapy"
)

type scrapyBoss struct {
	config     ScrapyBossConfig
	scrapeData map[string]*[]scrapy.ScrapeResult
	running    bool
}

func NewScrapyBoss(c ScrapyBossConfig) *scrapyBoss {
	return &scrapyBoss{
		config:     c,
		scrapeData: map[string]*[]scrapy.ScrapeResult{},
		running:    false,
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
			var d []scrapy.ScrapeResult
			s.scrapeData[endpoints[i].Endpoint] = &d
			go scrapy.Scrape(endpoints[i].Endpoint, endpoints[i].Selectors, &d)
		}
		time.Sleep(time.Duration(s.config.ScrapeIntervalInSeconds) * time.Second)
	}
}

func (s *scrapyBoss) GetScrapeData() map[string]*[]scrapy.ScrapeResult {
	return s.scrapeData
}

func (s *scrapyBoss) GetConfig() ScrapyBossConfig {
	return s.config
}
