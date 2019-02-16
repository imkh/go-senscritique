package main

import "github.com/imkh/scscraper/pkg/scscraper"

func main() {
	scs := scscraper.New()

	scs.ScrapeDiary()
}
