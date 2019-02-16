package main

import "github.com/imkh/senscritique-scraper/pkg/scscraper"

func main() {
	scs := scscraper.New()

	scs.ScrapeDiary()
}
