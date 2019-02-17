package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/imkh/scscraper/pkg/scscraper"
)

func main() {
	app := cli.NewApp()
	app.Name = "scscraper"
	app.Usage = "A SensCritique web scraper"

	app.Commands = []cli.Command{
		{
			Name:      "diary",
			Usage:     "Scrape a user's diary",
			UsageText: fmt.Sprintf("%s diary [command options] [username]", app.Name),
			Flags:     diaryFlags,
			Action: func(c *cli.Context) error {
				scs := scscraper.New()
				scs.ScrapeDiary(c.Args().First(), &scscraper.ScrapeDiaryOptions{
					Category: c.String("category"),
					Year:     c.String("year"),
					Month:    c.String("month"),
				})
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
