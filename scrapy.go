package scrapy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"launchpad.net/xmlpath"
)

type Selector struct {
	Name           string
	TypeOfSelector string
	Value          string
}

type ScrapeResult struct {
	Name           string
	TypeOfSelector string
	Value          string
	Result         string
	Error          error
	Time           time.Time
}

// Scrape will use the selector to scrape the endpoint and return the result of the selector or an error.
func Scrape(endpoint string, selectors []Selector) []ScrapeResult {
	var scrapeResults = []ScrapeResult{}
	var err error

	body, err := getEndpointBody(endpoint)

	if err == nil {
		for i := range selectors {
			scrapedValue, err := scrapeValue(body, selectors[i])

			scrapeResults = append(scrapeResults, ScrapeResult{
				Name:           selectors[i].Name,
				TypeOfSelector: selectors[i].TypeOfSelector,
				Value:          selectors[i].Value,
				Result:         scrapedValue,
				Error:          err,
				Time:           time.Now(),
			})
		}
	}

	return scrapeResults
}

func getEndpointBody(endpoint string) ([]byte, error) {
	var body []byte
	var err error
	response, err := http.Get(endpoint)

	if err == nil {
		body, err = ioutil.ReadAll(response.Body)
		defer response.Body.Close()
	}

	return body, err
}

func scrapeValue(document []byte, selector Selector) (string, error) {
	var result string
	var err error

	switch selector.TypeOfSelector {
	case "xpath":
		path := xmlpath.MustCompile(selector.Value)
		root, err := xmlpath.ParseHTML(bytes.NewReader(document))
		if err == nil {
			result, _ = path.String(root)
		}
	case "regex":
		regex := regexp.MustCompile(selector.Value)
		matches := regex.FindSubmatch(document)
		if len(matches) == 2 {
			result = string(matches[1])
		} else {
			err = fmt.Errorf("Regex selector resulted into multiple group matches ('%v'), only one is allowed", len(matches)-1)
		}
	default:
		err = fmt.Errorf("Selector type '%s' is not supported", selector.TypeOfSelector)
	}

	return result, err
}
