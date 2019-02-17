package main

import "github.com/urfave/cli"

var (
	diaryFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "category, c",
			Value: "all",
			Usage: "Limit the results to a specific category",
		},
		cli.IntFlag{
			Name:  "year, y",
			Usage: "Limit the results to a specific year",
		},
		cli.StringFlag{
			Name:  "month, m",
			Value: "all",
			Usage: "Limit the results to a specific month",
		},
	}
)
