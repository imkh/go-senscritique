package senscritique

import (
	"fmt"
	"strconv"

	"github.com/gocolly/colly"

	"github.com/imkh/go-senscritique/internal/validator"
)

// DiaryService provides access to the diary scraping functions.
//
// Scraped page: https://www.senscritique.com/:username/journal/all/all/all
type DiaryService service

// DiaryEntry represent an entry in a SensCritique user's diary.
type DiaryEntry struct {
	FrenchTitle      *string `json:"french_title"`
	TitleDate        *string `json:"title_date"`
	OriginalTitle    *string `json:"original_title"`
	TitleDescription *string `json:"title_description"`
	Date             *string `json:"date"`
	Score            *string `json:"score"`
}

// GetDiaryOptions specifies the optional parameters to scrape a diary.
type GetDiaryOptions struct {
	Category string `default:"all" validate:"oneof=all films series episodes jeuxvideo livres bd albums morceaux"`
	Year     int    `validate:"min=0"`
	Month    string `default:"all" validate:"oneof=all janvier fevrier mars avril mai juin juillet aout septembre octobre novembre decembre"`
}

// GetDiary scrape a given user diary page.
func (s *DiaryService) GetDiary(username string, opts *GetDiaryOptions) ([]DiaryEntry, error) {
	if err := validator.ValidateStruct(opts); err != nil {
		return nil, err
	}
	if opts.Year == 0 && opts.Month != "all" {
		return nil, fmt.Errorf("A month cannot be specified without a year")
	}

	diary := make([]DiaryEntry, 0)

	// Set default values
	page := 1
	yearStr := strconv.Itoa(opts.Year)
	if opts.Year == 0 {
		yearStr = "all"
	}

	fmt.Println(opts.Category)

	s.scraper.collector.OnHTML("div.eldi-collection", func(e *colly.HTMLElement) {
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

	s.scraper.collector.OnScraped(func(r *colly.Response) {
		if r.Ctx.GetAny("lastVisitedPage") == page {
			page++
			nextPageURL := fmt.Sprintf("%s/%s/journal/%s/%s/%s/page-%d.ajax", s.scraper.baseURL, username, opts.Category, yearStr, opts.Month, page)
			r.Request.Visit(nextPageURL)
		} else {
			fmt.Println(len(diary))
		}
	})

	url := fmt.Sprintf("%s/%s/journal/%s/%s/%s/page-%d.ajax", s.scraper.baseURL, username, opts.Category, yearStr, opts.Month, page)
	s.scraper.collector.Visit(url)

	return diary, nil
}
