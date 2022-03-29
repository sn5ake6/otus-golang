package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const maxCount int = 10

type wordCount struct {
	word  string
	count int
}

func Top10(s string) []string {
	wordsMap := make(map[string]int)

	for _, word := range strings.Fields(s) {
		wordsMap[word]++
	}

	words := make([]wordCount, 0, len(wordsMap))
	for word, count := range wordsMap {
		words = append(words, wordCount{word, count})
	}

	sort.Slice(words, func(i, j int) bool {
		if words[i].count > words[j].count {
			return true
		}
		if words[i].count == words[j].count {
			return words[i].word < words[j].word
		}
		return false
	})

	count := len(words)
	if count > maxCount {
		count = maxCount
	}

	top := make([]string, 0, count)

	for _, v := range words[:count] {
		top = append(top, v.word)
	}

	return top
}
