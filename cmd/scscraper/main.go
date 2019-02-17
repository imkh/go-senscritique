package main

import (
	"fmt"
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
				_, err := scs.ScrapeDiary(c.Args().First(), &scscraper.ScrapeDiaryOptions{
					Category: c.String("category"),
					Year:     c.Int("year"),
					Month:    c.String("month"),
				})
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
