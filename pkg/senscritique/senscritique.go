package senscritique

import (
	"net/url"
	"regexp"

	"github.com/gocolly/colly"
)

const (
	defaultBaseURL = "https://www.senscritique.com/"
)

// A Scraper represents a scraper for the SensCritique website
type Scraper struct {
	BaseURL *url.URL

	collector *colly.Collector
}

// NewScraper creates a new Scraper instance with default configuration
func NewScraper() *Scraper {
	baseURL, _ := url.Parse(defaultBaseURL)
	return &Scraper{
		BaseURL:   baseURL,
		collector: colly.NewCollector(),
	}
}

func trimString(value string) *string {
	if value == "" {
		return nil
	}
	space := regexp.MustCompile(`\s+`)
	s := space.ReplaceAllString(value, " ")
	return &s
}
