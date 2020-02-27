package httpfilter

import (
	"errors"
	"net/http"
)

type OpFunc func(w http.ResponseWriter, req FilterRequest) string

func ignore(w http.ResponseWriter, req FilterRequest) string {
	http.Error(w, "Not found.", 404)
	return ""
}

func pseudo(w http.ResponseWriter, req FilterRequest) string {
	if len(req.Args) < 1 {
		panic(errors.New("not enough arguments"))
		return ""
	}
	return req.Args[0]
}

func redirect(w http.ResponseWriter, req FilterRequest) string {
	if len(req.Args) < 1 {
		panic(errors.New("not enough arguments"))
		return ""
	}
	http.Redirect(w, req.Request, req.Args[0], 301)
	return ""
}
