package scrapy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"launchpad.net/xmlpath"
)

type Selector struct {
	typeOfSelector string
	value          string
}

// Scrape will use the selector to scrape the endpoint and return the result of the selector or an error.
func Scrape(endpoint string, selector Selector) (string, error) {
	var scrapedValue string
	var err error

	body, err := getEndpointBody(endpoint)

	if err == nil {
		scrapedValue, err = scrapeValue(body, selector)
	}

	return scrapedValue, err
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
	var value string
	var err error

	switch selector.typeOfSelector {
	case "xpath":
		path := xmlpath.MustCompile(selector.value)
		root, err := xmlpath.ParseHTML(bytes.NewReader(document))
		if err == nil {
			value, _ = path.String(root)
		}
	default:
		err = fmt.Errorf("Selector type '%s' is not supported", selector.typeOfSelector)
	}

	return value, err
}
