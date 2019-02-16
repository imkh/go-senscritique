package scscraper

import (
	"github.com/gocolly/colly"
)

// SensCritiqueScraper represents a scraper for SensCritique
type SensCritiqueScraper struct {
	baseURL string

	c *colly.Collector
}

// New creates a new SensCritiqueScraper instance with default configuration
func New() *SensCritiqueScraper {
	return &SensCritiqueScraper{
		baseURL: "https://www.senscritique.com/",
		c:       colly.NewCollector(),
	}
}
