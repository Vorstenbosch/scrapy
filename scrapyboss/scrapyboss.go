package scrapyboss

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vorstenbosch/scrapy/scrapy"
)

type scrapyBoss struct {
	config       ScrapyBossConfig
	scrapeData   map[string]*[]scrapy.ScrapeResult
	running      bool
	scrapeClient *scrapy.ScrapeClient
}

func NewScrapyBoss(c ScrapyBossConfig) *scrapyBoss {
	errorList := ValidateConfig(c)
	if len(errorList) > 0 {
		panic(fmt.Sprintf("Configuration is invalid due to '%v'", errorList))

	}

	// Default pool size
	idleConnectionPool := c.IdleConnectionPool
	if c.IdleConnectionPool == 0 {
		idleConnectionPool = 10
	}

	transport := &http.Transport{
		MaxIdleConns:       idleConnectionPool,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(c.ScrapeTimeoutInSeconds) * time.Second,
	}

	scrapeClient := scrapy.ScrapeClient{HttpClient: client}

	return &scrapyBoss{
		config:       c,
		scrapeData:   map[string]*[]scrapy.ScrapeResult{},
		running:      false,
		scrapeClient: &scrapeClient,
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

func (s *scrapyBoss) IsRunning() bool {
	return s.running
}

func (s *scrapyBoss) Stop() {
	s.running = false
}

func (s *scrapyBoss) collect(endpoints []ScrapeEndpoint, interval int) {
	for s.running {
		for i := range endpoints {
			if s.scrapeData[endpoints[i].Endpoint] == nil {
				s.scrapeData[endpoints[i].Endpoint] = &[]scrapy.ScrapeResult{}
			}

			go s.scrapeClient.Scrape(endpoints[i].Endpoint, endpoints[i].Selectors, s.scrapeData[endpoints[i].Endpoint])
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
