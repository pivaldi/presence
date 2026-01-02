package main

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func snakeToCamelCase(in string) string {
	in = strings.TrimSpace(cases.Lower(language.Und).String(in))
	if in == "" {
		return in
	}

	tokens := strings.Split(in, "_")
	caser := cases.Title(language.Und, cases.NoLower)

	var out string
	for i, token := range tokens {
		if i == 0 {
			out += token
			continue
		}
		out += caser.String(token)
	}

	return out
}

func printBar() {
	fmt.Println("\n" + strings.Repeat("=", 60))
}
