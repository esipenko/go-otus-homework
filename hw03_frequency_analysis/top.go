package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

func trimPunctuation(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsPunct(r) && r != '-'
	})
}

func Top10(text string) []string {
	words := strings.Fields(text)

	counter := make(map[string]int)
	uniqWords := make([]string, 0, len(words))

	for _, w := range words {
		w = strings.ToLower(w)
		w = trimPunctuation(w)

		if w == "" || w == "-" {
			continue
		}

		if _, ok := counter[w]; !ok {
			uniqWords = append(uniqWords, w)
		}

		counter[w]++
	}

	sort.Slice(uniqWords, func(i, j int) bool {
		a, b := uniqWords[i], uniqWords[j]

		if counter[a] == counter[b] {
			return a < b
		}

		return counter[a] > counter[b]
	})

	if len(uniqWords) < 10 {
		return uniqWords
	}

	return uniqWords[:10]
}
