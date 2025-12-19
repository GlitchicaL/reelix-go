package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SnakeToTitle(s string) string {
	words := strings.Split(s, "_")
	caser := cases.Title(language.English)

	for i, w := range words {
		words[i] = caser.String(w)
	}

	return strings.Join(words, " ")
}

func TitleToSnake(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "_"))
}
