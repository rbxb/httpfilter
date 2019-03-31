package fileserve

import (
	"strings"
)

func parseTagdata(b []byte) [][]string {
	entries := make([][]string, 0)
	lines := strings.Split(string(b), string(0x0A))
	tag := ""
	for _, line := range lines {
		vals := make([]string, 0)
		tokens := strings.Split(line, string(0x20))
		for _, token := range tokens {
			token = strings.TrimSpace(token)
			if len(token) < 1 {
				continue
			}
			if token[0] == 0x23 {
				entries = flushEntry(tag, vals, entries)
				tag = token[1:]
			} else {
				vals = append(vals, token)
			}
		}
		entries = flushEntry(tag, vals, entries)
	}
	return entries
}

func flushEntry(tag string, vals []string, entries [][]string) [][]string {
	if tag != "" && len(vals) > 0 {
		entry := append([]string{tag}, vals...)
		entries = append(entries, entry)
	}
	return entries
}