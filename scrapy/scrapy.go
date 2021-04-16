package scrapy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
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

type ScrapeClient struct {
	HttpClient *http.Client
}

func NewScrapeClient() *ScrapeClient {
	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return &ScrapeClient{HttpClient: client}
}

// Scrape will use the selector to scrape the endpoint and place the result on the 'out'
func (c *ScrapeClient) Scrape(endpoint string, selectors []Selector, out *[]ScrapeResult) {
	var err error

	body, err := c.getEndpointBody(endpoint)

	sr := ScrapeResult{
		Error: err,
	}

	if err == nil {
		for i := range selectors {
			scrapedValue, err := scrapeValue(body, selectors[i])

			sr.Name = selectors[i].Name
			sr.TypeOfSelector = selectors[i].TypeOfSelector
			sr.Value = selectors[i].Value
			sr.Result = scrapedValue
			sr.Error = err
		}
	}

	sr.Time = time.Now()
	*out = append(*out, sr)
}

func (c *ScrapeClient) getEndpointBody(endpoint string) ([]byte, error) {
	var body []byte
	var err error

	response, err := c.HttpClient.Get(endpoint)

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
		var xmlPathRoot *xmlpath.Node
		var htmlRoot *html.Node

		path := xmlpath.MustCompile(selector.Value)

		// In case we run into invalid xml we need to clean it first
		reader := bytes.NewReader(document)
		htmlRoot, err = html.Parse(reader)

		var b bytes.Buffer
		html.Render(&b, htmlRoot)

		// We now have the cleaned version of the html page
		xmlPathRoot, err = xmlpath.ParseHTML(strings.NewReader(b.String()))
		if err == nil {
			var found bool
			result, found = path.String(xmlPathRoot)
			if !found {
				err = fmt.Errorf("Unable to find '%s'", selector.Value)
			}
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
