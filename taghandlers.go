package fileserve

import (
	"net/http"
	"path/filepath"
)

type TagHandler func(* Server, []string, http.ResponseWriter, * http.Request) error

func ignore(srvr * Server, args []string, w http.ResponseWriter, req * http.Request) error {
	srvr.ServeFile("", w, req)
	return nil
}

func pseudo(srvr * Server, args []string, w http.ResponseWriter, req * http.Request) error {
	if len(args) < 1 {
		return ErrorNotEnoughArguments
	}
	name := filepath.Join(filepath.Dir(req.URL.Path), args[0])
	srvr.ServeFile(name, w, req)
	return nil
}

func redirect(srvr * Server, args []string, w http.ResponseWriter, req * http.Request) error {
	if len(args) < 1 {
		return ErrorNotEnoughArguments
	}
	http.Redirect(w, req, args[0], 301)
	return nil
}

func deft(srvr * Server, args []string, w http.ResponseWriter, req * http.Request) error {
	srvr.ServeFile(req.URL.Path, w, req)
	return nil
}