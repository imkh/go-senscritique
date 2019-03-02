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

	mux.HandleFunc("/go-senscritique/journal/all/2018/all/page-1.ajax", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, string(testDataDiary))
	})

	diary, err := scraper.Diary.GetDiary("go-senscritique", &GetDiaryOptions{
		Year: 2018,
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
		{
			Product: &DiaryProduct{
				ID:            "21727461",
				FrenchTitle:   "Your Name.",
				ReleaseYear:   "2016",
				OriginalTitle: "Kimi no Na wa.",
				Description:   "Long-métrage d'animation de Makoto Shinkai",
			},
			Date:  "2018-04-30",
			Score: "8",
		},
		{
			Product: &DiaryProduct{
				ID:            "10416244",
				FrenchTitle:   "The Legend of Zelda : Breath of the Wild",
				ReleaseYear:   "2017",
				OriginalTitle: "Zeruda no densetsu: Buresu obu za wairudo",
				Description:   "Jeu vidéo de Nintendo EPD, Monolith Software et Nintendo",
			},
			Date:  "2018-04-30",
			Score: "10",
		},
	}
	if !reflect.DeepEqual(want, diary) {
		t.Errorf("Diary.GetDiary returned %+v, want %+v", diary, want)
	}
}
