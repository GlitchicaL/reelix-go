package utils

import (
	"regexp"
	"strings"
)

func ToTitle(s string) string {
	/*
		Since the folder name could be in the
		format of "Breaking Bad" or "breaking_bad",
		we convert to snake first
	*/
	return SnakeToTitle(TitleToSnake(s))
}

func SnakeToTitle(s string) string {
	smallWords := map[string]bool{
		"and": true, "or": true, "the": true, "of": true,
		"in": true, "on": true, "at": true, "a": true, "an": true,
	}

	acronyms := map[string]bool{
		"NASA": true,
	}

	isRomanNumeral := regexp.MustCompile(`^[IVXLCDM]+$`).MatchString

	words := strings.Split(s, "_")

	for i, w := range words {
		upperWord := strings.ToUpper(w)

		/*
			There needs to be a way to handle small words,
			acronyms, and roman numerals. For example:

			call_of_duty -> Call of Duty
			nasa -> NASA
			artemis_ii -> Artemis II
		*/

		if acronyms[upperWord] {
			words[i] = upperWord
			continue
		}

		if isRomanNumeral(upperWord) {
			words[i] = upperWord
			continue
		}

		if i != 0 && smallWords[strings.ToLower(w)] {
			words[i] = strings.ToLower(w)
			continue
		}

		if len(w) > 0 {
			words[i] = strings.ToUpper(string(w[0])) + strings.ToLower(w[1:])
		}
	}

	return strings.Join(words, " ")
}

func TitleToSnake(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "_"))
}
