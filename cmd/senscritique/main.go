package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/imkh/go-senscritique"
)

func main() {
	app := cli.NewApp()
	app.Name = "senscritique"
	app.Usage = "A SensCritique web scraper"

	app.Commands = []cli.Command{
		{
			Name:      "diary",
			Usage:     "Scrape a user's diary",
			UsageText: fmt.Sprintf("%s diary [command options] [username]", app.Name),
			Flags:     diaryFlags,
			Action: func(c *cli.Context) error {
				sc := senscritique.NewScraper()
				diary, err := sc.Diary.GetDiary(c.Args().First(), &senscritique.GetDiaryOptions{
					Universe: c.String("universe"),
					Year:     c.Int("year"),
					Month:    c.String("month"),
				})
				output, _ := json.Marshal(diary)
				fmt.Println(string(output))
				return err
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
