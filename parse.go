package httpfilter

import (
	"io/ioutil"
	"strings"
)

func parseFilterFile(path string) [][]string {
	var entries [][]string
	if b, err := ioutil.ReadFile(path); err == nil {
		lines := strings.Split(string(b), string(0x0A))
		tag := ""
		for _, line := range lines {
			vals := make([]string, 0)
			words := strings.Split(line, string(0x20))
			for _, word := range words {
				word = strings.TrimSpace(word)
				if len(word) < 1 {
					continue
				}
				if word[0] == 0x23 {
					entries = appendEntry(tag, vals, entries)
					tag = word[1:]
				} else {
					vals = append(vals, word)
				}
			}
			entries = appendEntry(tag, vals, entries)
		}
	}
	return appendEntry("deft", []string{"*"}, entries)
}

func appendEntry(tag string, vals []string, entries [][]string) [][]string {
	if tag != "" && len(vals) > 0 {
		entry := append([]string{tag}, vals...)
		entries = append(entries, entry)
	}
	return entries
}
