package senscritique

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"
)

func TestGetDiary(t *testing.T) {
	mux, server, scraper := setup()
	defer teardown(server)

	testDataDiary, err := ioutil.ReadFile("test/testdata/diary.html")
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/username/journal/jeuxvideo/2018/decembre/page-1.ajax", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, string(testDataDiary))
	})

	diary, err := scraper.Diary.GetDiary("username", &GetDiaryOptions{
		Universe: "jeuxvideo",
		Year:     2018,
		Month:    "decembre",
	})
	if err != nil {
		t.Errorf("Diary.GetDiary returned error: %v", err)
	}

	want := []*DiaryEntry{
		{
			Product: &DiaryProduct{
				ID:            "31252340",
				FrenchTitle:   "Super Smash Bros. Ultimate",
				ReleaseYear:   "2018",
				OriginalTitle: "Dairantō Sumasshu Burazāzu Supesharu",
				Description:   "Jeu vidéo de BANDAI NAMCO Games et Nintendo",
			},
			Date:  "2018-12-25",
			Score: "9",
		},
	}
	if !reflect.DeepEqual(want, diary) {
		t.Errorf("Diary.GetDiary returned %+v, want %+v", diary, want)
	}
}
