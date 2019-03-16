package senscritique

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/go-test/deep"
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
				Universe:      Games,
				FrenchTitle:   "Super Smash Bros. Ultimate",
				ReleaseYear:   "2018",
				OriginalTitle: "Dairantō Sumasshu Burazāzu Supesharu",
				Details:       "Jeu vidéo de BANDAI NAMCO Games et Nintendo",
				WebURL:        fmt.Sprintf("%s%s", scraper.baseURL, "/jeuvideo/Super_Smash_Bros_Ultimate/31252340"),
				Poster:        "https://media.senscritique.com/media/000017854034/90/Super_Smash_Bros_Ultimate.jpg",
			},
			Date:   Time(time.Date(2018, 12, 25, 0, 0, 0, 0, time.UTC)),
			Rating: Int(9),
		},
		{
			Product: &DiaryProduct{
				ID:            "21727461",
				Universe:      Movies,
				FrenchTitle:   "Your Name.",
				ReleaseYear:   "2016",
				OriginalTitle: "Kimi no Na wa.",
				Details:       "Long-métrage d'animation de Makoto Shinkai",
				WebURL:        fmt.Sprintf("%s%s", scraper.baseURL, "/film/Your_Name/21727461"),
				Poster:        "https://media.senscritique.com/media/000016585625/90/Your_Name.jpg",
			},
			Date:   Time(time.Date(2018, 04, 30, 0, 0, 0, 0, time.UTC)),
			Rating: Int(8),
		},
		{
			Product: &DiaryProduct{
				ID:            "10416244",
				Universe:      Games,
				FrenchTitle:   "The Legend of Zelda : Breath of the Wild",
				ReleaseYear:   "2017",
				OriginalTitle: "Zeruda no densetsu: Buresu obu za wairudo",
				Details:       "Jeu vidéo de Nintendo EPD, Monolith Software et Nintendo",
				WebURL:        fmt.Sprintf("%s%s", scraper.baseURL, "/jeuvideo/The_Legend_of_Zelda_Breath_of_the_Wild/10416244"),
				Poster:        "https://media.senscritique.com/media/000016771881/90/The_Legend_of_Zelda_Breath_of_the_Wild.jpg",
			},
			Date:   Time(time.Date(2018, 04, 30, 0, 0, 0, 0, time.UTC)),
			Rating: Int(10),
		},
	}

	if diff := deep.Equal(want, diary); diff != nil {
		t.Error(diff)
	}
}
