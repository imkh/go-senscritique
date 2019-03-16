package main

import (
	"os"
	"strconv"

	"github.com/imkh/go-senscritique"
	"github.com/olekukonko/tablewriter"
)

func printDiaryBreakdown(diary []*senscritique.DiaryEntry) {
	universes := []senscritique.Universe{
		senscritique.Movies,
		senscritique.Shows,
		senscritique.Episodes,
		senscritique.Games,
		senscritique.Books,
		senscritique.Comics,
		senscritique.Albums,
		senscritique.Tracks,
	}

	header := []string{}
	content := []string{}
	for _, u := range universes {
		header = append(header, u.English())
		content = append(content, strconv.Itoa(countDiaryProductByUniverse(diary, u)))
	}
	header = append(header, "Total")
	content = append(content, strconv.Itoa(len(diary)))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgMagentaColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgMagentaColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgWhiteColor, tablewriter.FgBlackColor},
	)

	table.Append(content)
	table.Render()
}

func countDiaryProductByUniverse(diary []*senscritique.DiaryEntry, u senscritique.Universe) int {
	var i int
	for _, d := range diary {
		if d.Product.Universe == u {
			i++
		}
	}
	return i
}
