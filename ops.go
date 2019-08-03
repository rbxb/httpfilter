package httpfilter

import (
	"errors"
	"net/http"
)

type OpFunc func(w http.ResponseWriter, req * http.Request, query string, args []string) string

func ignore(w http.ResponseWriter, _ * http.Request, _ string, _ []string) string {
	http.Error(w, "Not found.", 404)
	return ""
}

func pseudo(_ http.ResponseWriter, _ * http.Request, _ string, args []string) string {
	if len(args) < 1 {
		panic(errors.New("Not enough arguments."))
		return ""
	}
	return args[0]
}

func redirect(w http.ResponseWriter, req * http.Request, _ string, args []string) string {
	if len(args) < 1 {
		panic(errors.New("Not enough arguments."))
		return ""
	}
	http.Redirect(w, req, args[0], 301)
	return ""
}