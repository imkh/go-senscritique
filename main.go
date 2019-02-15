package main

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly"
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

	diary := make([]DiaryEntry, 0)

	username := "iMkh"
	category := "all"
	page := 1

	fmt.Println(category)

	c.OnHTML("div.eldi-collection", func(e *colly.HTMLElement) {
		e.ForEach("li.eldi-list-item", func(i int, e *colly.HTMLElement) {
			// TODO: add handling of sub-item (2 entries at the same date)
			if date := e.Attr("data-sc-datedone"); date != "" {
				score := e.ChildText("div.epri-score")
				if score == "" { // TODO: check "done" state (no score)
					score = "✓" // e.DOM.Find("span.eins-done")
				}
				diary = append(diary, DiaryEntry{
					FrenchTitle:      trimString(e.ChildText("[id^=product-title]")),
					TitleDate:        trimString(e.ChildText("span.elco-date")),
					OriginalTitle:    trimString(e.ChildText("p.elco-original-title")),
					TitleDescription: trimString(e.ChildText("p.elco-baseline")),
					Date:             trimString(date),
					Score:            trimString(score),
				})
			}
		})
		e.Response.Ctx.Put("lastVisitedPage", page)
	})

	c.OnScraped(func(r *colly.Response) {
		if r.Ctx.GetAny("lastVisitedPage") == page {
			page++
			nextPageURL := fmt.Sprintf("https://www.senscritique.com/%s/journal/%s/all/all/page-%d.ajax", username, category, page)
			r.Request.Visit(nextPageURL)
		} else {
			fmt.Println(len(diary))
		}
	})

	url := fmt.Sprintf("https://www.senscritique.com/%s/journal/%s/all/all/page-%d.ajax", username, category, page)
	c.Visit(url)
}
