package main

import "github.com/urfave/cli"

var (
	diaryFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "universe, u",
			Value: "all",
			Usage: "Limit the results to a specific universe",
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
		cli.BoolFlag{
			Name:  "json",
			Usage: "Print JSON output",
		},
	}
)
