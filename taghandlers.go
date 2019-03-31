package fileserve

import (
	"net/http"
	"errors"
	"path/filepath"
)

var (
	ErrorNotEnoughArguments = errors.New("Not enough arguments.")
)

type TagHandler func(* Server, []string, http.ResponseWriter, * http.Request) error

func ignore(srvr * Server, vals []string, w http.ResponseWriter, req * http.Request) error {
	srvr.ServeFile("", w, req)
	return nil
}

func pseudo(srvr * Server, vals []string, w http.ResponseWriter, req * http.Request) error {
	if len(vals) < 2 {
		return ErrorNotEnoughArguments
	}
	name := filepath.Join(filepath.Dir(req.URL.Path), vals[1])
	srvr.ServeFile(name, w, req)
	return nil
}

func redirect(srvr * Server, vals []string, w http.ResponseWriter, req * http.Request) error {
	if len(vals) < 2 {
		return ErrorNotEnoughArguments
	}
	http.Redirect(w, req, vals[1], 301)
	return nil
}

func deft(srvr * Server, vals []string, w http.ResponseWriter, req * http.Request) error {
	srvr.ServeFile(req.URL.Path, w, req)
	return nil
}