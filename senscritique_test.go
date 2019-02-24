package senscritique

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// setup sets up a test HTTP server along with a senscritique.Scraper that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the web page being tested.
func setup() (*http.ServeMux, *httptest.Server, *Scraper) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock web page responses.
	server := httptest.NewServer(mux)

	// scraper is the SensCritique scraper being tested.
	scraper := NewScraper()
	_ = scraper.SetBaseURL(server.URL)

	return mux, server, scraper
}

// teardown closes the test HTTP server.
func teardown(server *httptest.Server) {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %s, want %s", got, want)
	}
}
