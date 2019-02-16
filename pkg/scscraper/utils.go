package scscraper

import "regexp"

func trimString(value string) *string {
	if value == "" {
		return nil
	}
	space := regexp.MustCompile(`\s+`)
	s := space.ReplaceAllString(value, " ")
	return &s
}
