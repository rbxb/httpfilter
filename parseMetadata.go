package fileserve

import (
	"strings"
)

func parseMetadata(b []byte) map[string][]string {
	metadata := make(map[string][]string)

	for i := range b {
		switch b[i] {
		case 0x09, 0x0A, 0x0D:
			b[i] = 0x20
		}
	}

	split := strings.Split(string(b), " ")

	tag := ""
	vals := make([]string, 0)
	for _, s := range split {
		if len(s) < 1 {
			continue
		} else if s[0] == 0x23 {
			if tag != "" && tag != "#" && len(vals) > 0 {
				metadata[tag] = append(metadata[tag], vals...)
			}
			if len(s[1:]) > 0 {
				tag = s[1:]
			} else {
				tag = "#"
			}
			vals = make([]string, 0)
		} else {
			switch tag {
			case "":
				//do nothing
			case "#":
				tag = s
			default:
				vals = append(vals, s)
			}
		}
	}
	if tag != "" && tag != "#" && len(vals) > 0 {
		metadata[tag] = append(metadata[tag], vals...)
	}

	return metadata
}