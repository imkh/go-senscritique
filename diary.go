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
	ID            string   `json:"id"`
	Universe      Universe `json:"universe"`
	FrenchTitle   string   `json:"french_title"`
	ReleaseYear   string   `json:"release_year"`
	OriginalTitle string   `json:"original_title"`
	Details       string   `json:"details"`
	WebURL        string   `json:"web_url"`
	Poster        string   `json:"poster"`
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
		e.ForEach("li.eldi-list-item[data-sc-datedone]", func(i int, f *colly.HTMLElement) {
			// Parse date
			date, err := s.parseDate(f)
			if err != nil {
				log.Println(err)
				return
			}
			// Parse each sub-item
			f.ForEach("div[data-rel='diary-sub-item']", func(j int, g *colly.HTMLElement) {
				// Parse rating
				rating, err := s.parseRating(g)
				if err != nil {
					log.Println(err)
					return
				}
				// Parse diary entry
				diary = append(diary, &DiaryEntry{
					Product: s.parseProduct(g),
					Date:    date,
					Rating:  rating,
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

func (s *DiaryService) parseUniverse(url string) Universe {
	switch universe := strings.Split(strings.TrimPrefix(url, "/"), "/")[0]; universe {
	case "film":
		return Movies
	case "serie":
		return Shows
	case "jeuvideo":
		return Games
	case "livre":
		return Books
	case "bd":
		return Comics
	case "album":
		return Albums
	case "morceau":
		return Tracks
	default:
		return ""
	}
}

func (s *DiaryService) isEpisode(e *colly.HTMLElement) bool {
	return e.ChildAttr("span.eldi-collection-poster", "data-sc-episode-id") != ""
}

func (s *DiaryService) parseProduct(e *colly.HTMLElement) *DiaryProduct {
	// Check if the diary entry product is an episode of
	// a TV Show, in which case the parsing differs for some fields.
	var id string
	var universe Universe
	var webURL string
	if s.isEpisode(e) {
		id = trimString(e.ChildAttr(".eldi-collection-poster", "data-sc-episode-id"))
		universe = Episodes
		webURL = ""
	} else {
		href := trimString(e.ChildAttr(".eldi-collection-poster", "href"))
		id = trimString(e.ChildAttr(".eldi-collection-poster", "data-sc-product-id"))
		universe = s.parseUniverse(href)
		webURL = fmt.Sprintf("%s%s", s.scraper.baseURL, href)
	}

	return &DiaryProduct{
		ID:            id,
		Universe:      universe,
		FrenchTitle:   trimString(e.ChildText("[id^=product-title]")),
		ReleaseYear:   strings.Trim(trimString(e.ChildText("span.elco-date")), "()"),
		OriginalTitle: trimString(e.ChildText("p.elco-original-title")),
		Details:       trimString(e.ChildText("p.elco-baseline")),
		WebURL:        webURL,
		Poster:        trimString(e.ChildAttr(".eldi-collection-poster > img", "src")),
	}
}

// GetRating returns the Rating field if it's non-nil, zero value otherwise.
func (d *DiaryEntry) GetRating() int {
	if d == nil || d.Rating == nil {
		return 0
	}
	return *d.Rating
}
