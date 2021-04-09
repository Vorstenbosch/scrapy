# scrapy

## TL;DR
Scrapy is a easy and simple scraping library.

## How to use

### Scrapy
To use Scrapy directly you need to call the 'Scrape' method with an endpoint to scrape and a list of selectors, e.g.:
```go
endpoint := "http://localhost:5555/"
	selector := Selector{
		Name:           "xpath-scrape",
		TypeOfSelector: "xpath",
		Value:          "//div",
	}

var result []ScrapeResult
Scrape(endpoint, []Selector{selector}, &result)
```

### ScrapyBoss
ScrapyBoss serves as a lightweight scheduler for scrapes done by Scrapy. It can be configured by providing a yaml file or by creating the configuration directly from code. Once starting it will start scraping the endpoints according to the configured selectors based on the 'ScrapeIntervalInSeconds'.

#### Config from code
```go
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
scrapyBoss.Start()
```
#### Config from a configuration file
```go
data, err := ioutil.ReadFile("/path/to/the/config/file.yaml")
if err != nil {
    log.Fatal(err)
}

config, err := scrapyboss.ParseConfig(data)
if err != nil {
    log.Fatal(err)
}

scrapyBoss := scrapyboss.NewScrapyBoss(config)
scrapyBoss.Start()
```

## Selectors
Scrapy supports the following selector types:
- xpath
- regex

## Scrapi project
The Scrapi project (https://github.com/Vorstenbosch/scrapi) provides an example of an implementation of the scrapy library.

## Roadmap to releasing version 1.0.0
- Improve test cases with:
    - invalid html pages
    - more complex html pages
    - assertions for the regex selector
    - assertions on error messages (e.g. when the selector does not find anything)