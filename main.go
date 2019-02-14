package main

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/hokaccha/go-prettyjson"
)

// DiaryEntry represent an entry in a SensCritique user's diary.
type DiaryEntry struct {
	FrenchTitle      *string `json:"french_title"`
	TitleDate        *string `json:"title_date"`
	OriginalTitle    *string `json:"original_title"`
	TitleDescription *string `json:"title_description"`
	Date             *string `json:"date"`
	Score            *string `json:"score"`
}

func trimString(value string) *string {
	if value == "" {
		return nil
	}
	space := regexp.MustCompile(`\s+`)
	s := space.ReplaceAllString(value, " ")
	return &s
}

func main() {
	c := colly.NewCollector()

	journal := make([]DiaryEntry, 0)

	c.OnHTML("div.eldi-collection", func(e *colly.HTMLElement) {
		e.ForEach("li.eldi-list-item", func(i int, e *colly.HTMLElement) {
			if date := e.Attr("data-sc-datedone"); date != "" {
				score := e.ChildText("div.epri-score")
				if score == "" { // TODO: check "done" state (no score)
					score = "✓" // e.DOM.Find("span.eins-done")
				}
				journal = append(journal, DiaryEntry{
					FrenchTitle:      trimString(e.ChildText("[id^=product-title]")),
					TitleDate:        trimString(e.ChildText("span.elco-date")),
					OriginalTitle:    trimString(e.ChildText("p.elco-original-title")),
					TitleDescription: trimString(e.ChildText("p.elco-baseline")),
					Date:             trimString(date),
					Score:            trimString(score),
				})
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		x, _ := prettyjson.Marshal(journal)
		fmt.Println(string(x))
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit("https://www.senscritique.com/iMkh/journal")
}
