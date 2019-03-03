package senscritique

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

const (
	defaultBaseURL = "https://www.senscritique.com"
)

// A Scraper represents a scraper for the SensCritique website.
type Scraper struct {
	// SensCritique base URL.
	baseURL *url.URL

	// Colly's main entity: Collector, which provides
	// the actual scraper instance for a scraping job.
	collector *colly.Collector

	// Reuse a single struct instead of allocating one for each service on the heap.
	common service

	// Services used for talking to different parts of the SensCritique scraper.
	Diary *DiaryService
}

type service struct {
	scraper *Scraper
}

// NewScraper creates a new Scraper instance with default configuration.
func NewScraper() *Scraper {
	baseURL, _ := url.Parse(defaultBaseURL)

	s := &Scraper{
		baseURL:   baseURL,
		collector: colly.NewCollector(),
	}

	s.common.scraper = s
	s.Diary = (*DiaryService)(&s.common)

	return s
}

// SetBaseURL sets the base URL to a custom endpoint.
func (s *Scraper) SetBaseURL(urlStr string) error {
	if strings.HasSuffix(urlStr, "/") {
		urlStr = strings.TrimSuffix(urlStr, "/")
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	s.baseURL = baseURL

	return nil
}

func trimString(s string) string {
	space := regexp.MustCompile(`\s+`)
	ts := space.ReplaceAllString(s, " ")
	return ts
}

// Bool is a helper method to return a pointer to a provided bool value.
func Bool(v bool) *bool { return &v }

// Int is a helper method to return a pointer to a provided int value.
func Int(v int) *int { return &v }

// String is a helper method to return a pointer to a provided string value.
func String(v string) *string { return &v }
