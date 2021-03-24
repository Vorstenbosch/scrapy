package scrapy

import "testing"

func TestScrape(t *testing.T) {
	// Given
	// TODO: create a webapp stub for testing
	endpoint := "http://aap.nl/"
	selector := Selector{
		typeOfSelector: "xpath",
		value:          "//h2[@class='heading-lg']",
	}

	expectedResult := "AAP geeft dieren weer een toekomst"

	// When
	result, err := Scrape(endpoint, selector)

	// Then
	if err != nil {
		t.Errorf("Scrape test failed due to unexpected error '%v'", err)
	}

	if result != expectedResult {
		t.Errorf("Scrape test failed as the result did not match '%s' but was '%s'", expectedResult, result)
	}

}
