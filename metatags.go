package fileserve

import (
	"net/http"
	"path/filepath"
)

type TagHandlerFunc func([]string, * string, http.ResponseWriter, * http.Request, * bool)

func ignore(vals []string, p * string, w http.ResponseWriter, req * http.Request, written * bool) {
	basename := filepath.Base(*p)
	for _, val := range vals {
		if val == basename || val == "*" {
			*p = ""
			break
		}
	}
}

func pseudo(vals []string, p * string, w http.ResponseWriter, req * http.Request, written * bool) {
	basename := filepath.Base(*p)
	for i := 0; i + 1 < len(vals); i += 2 {
		if vals[i] == basename || vals[i] == "*" {
			*p = filepath.Join(filepath.Dir(*p), vals[i + 1])
			break
		}
	}
}

func redirect(vals []string, p * string, w http.ResponseWriter, req * http.Request, written * bool) {
	basename := filepath.Base(*p)
	for i := 0; i + 1 < len(vals); i += 2 {
		if vals[i] == basename || vals[i] == "*" {
			if !*written {
				http.Redirect(w, req, vals[i + 1], 301)
				*written = true
				break
			}
		}
	}
}