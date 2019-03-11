package senscritique

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"

	"github.com/imkh/go-senscritique/internal/validator"
)

// DiaryService provides access to the diary scraping functions.
//
// Scraped page: https://www.senscritique.com/:username/journal/:universe/:year/:month
type DiaryService service

// DiaryProduct represents a product in a diary entry.
type DiaryProduct struct {
	ID            string `json:"id"`
	FrenchTitle   string `json:"french_title"`
	ReleaseYear   string `json:"release_year"`
	OriginalTitle string `json:"original_title"`
	Description   string `json:"title_description"`
}

// DiaryEntry represents an entry in a SensCritique user's diary.
type DiaryEntry struct {
	Product *DiaryProduct `json:"product"`
	Date    *time.Time    `json:"date"`
	Rating  *int          `json:"rating"`
}

// GetDiaryOptions specifies the optional parameters to scrape a diary.
type GetDiaryOptions struct {
	Universe string `default:"all" validate:"oneof=all films series episodes jeuxvideo livres bd albums morceaux"`
	Year     int    `validate:"min=0"`
	Month    string `default:"all" validate:"oneof=all janvier fevrier mars avril mai juin juillet aout septembre octobre novembre decembre"`
}

func (s *DiaryService) parseDate(e *colly.HTMLElement) (*time.Time, error) {
	dateStr := e.Attr("data-sc-datedone")
	if dateStr == "" {
		return nil, fmt.Errorf("no date found")
	}
	date, err := time.Parse("2006-01-02", trimString(dateStr))
	if err != nil {
		return nil, err
	}
	return &date, nil
}

func (s *DiaryService) parseRating(e *colly.HTMLElement) (*int, error) {
	ratingStr := e.ChildText("div.eldi-collection-rating")
	if ratingStr == "" {
		return nil, nil
	}
	rating, err := strconv.Atoi(trimString(ratingStr))
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

// GetDiary scrapes a given user diary page.
func (s *DiaryService) GetDiary(username string, opts *GetDiaryOptions) ([]*DiaryEntry, error) {
	if opts == nil {
		opts = new(GetDiaryOptions)
	}
	if err := validator.ValidateStruct(opts); err != nil {
		return nil, err
	}
	if opts.Year == 0 && opts.Month != "all" {
		return nil, fmt.Errorf("a month cannot be specified without a year")
	}

	var diary []*DiaryEntry

	// Set default values
	page := 1
	yearStr := strconv.Itoa(opts.Year)
	if opts.Year == 0 {
		yearStr = "all"
	}

	s.scraper.collector.OnHTML("div.eldi-collection", func(e *colly.HTMLElement) {
		e.ForEach("li.eldi-list-item[data-sc-datedone]", func(i int, e *colly.HTMLElement) {
			// Parse date
			date, err := s.parseDate(e)
			if err != nil {
				log.Println(err)
				return
			}
			// Parse each sub-item
			e.ForEach("div[data-rel='diary-sub-item']", func(i int, e *colly.HTMLElement) {
				// Parse rating
				rating, err := s.parseRating(e)
				if err != nil {
					log.Println(err)
					return
				}
				// Parse diary entry
				diary = append(diary, &DiaryEntry{
					Product: &DiaryProduct{
						ID:            trimString(e.ChildAttr("a.eldi-collection-poster", "data-sc-product-id")),
						FrenchTitle:   trimString(e.ChildText("[id^=product-title]")),
						ReleaseYear:   trimString(strings.Trim(e.ChildText("span.elco-date"), "()")),
						OriginalTitle: trimString(e.ChildText("p.elco-original-title")),
						Description:   trimString(e.ChildText("p.elco-baseline")),
					},
					Date:   date,
					Rating: rating,
				})
			})
		})
		e.Response.Ctx.Put("lastVisitedPage", page)
	})

	s.scraper.collector.OnScraped(func(r *colly.Response) {
		if r.Ctx.GetAny("lastVisitedPage") == page {
			page++
			nextPageURL := fmt.Sprintf("%s/%s/journal/%s/%s/%s/page-%d.ajax", s.scraper.baseURL, username, opts.Universe, yearStr, opts.Month, page)
			_ = r.Request.Visit(nextPageURL)
		}
	})

	url := fmt.Sprintf("%s/%s/journal/%s/%s/%s/page-%d.ajax", s.scraper.baseURL, username, opts.Universe, yearStr, opts.Month, page)
	_ = s.scraper.collector.Visit(url)

	return diary, nil
}

// GetRating returns the Rating field if it's non-nil, zero value otherwise.
func (d *DiaryEntry) GetRating() int {
	if d == nil || d.Rating == nil {
		return 0
	}
	return *d.Rating
}
