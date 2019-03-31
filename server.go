package fileserve

import (
	"net/http"
	"path/filepath"
	"io/ioutil"
	"mime"
)

type Server struct {
	root string
	tagfile string
	tagHandlers map[string]TagHandler
}

func NewServer(root string, tagHandlers map[string]TagHandler) * Server {
	defaultTagHandlers := map[string]TagHandler {
		"ignore": ignore,
		"pseudo": pseudo,
		"redirect": redirect,
		"default": deft,
	}
	if tagHandlers != nil {
		for k, v := range tagHandlers {
			defaultTagHandlers[k] = v
		}
	}
	return &Server{
		root: root,
		tagfile: "_tags.txt",
		tagHandlers: defaultTagHandlers,
	}
}

func(srvr * Server) ServeHTTP(w http.ResponseWriter, req * http.Request) {
	p := filepath.Join(srvr.root, filepath.Clean("/" + req.URL.Path))
	b, err := ioutil.ReadFile(filepath.Join(filepath.Dir(p), srvr.tagfile))
	if err != nil {
		http.Error(w, "Internal error.", 500)
		return
	}

	tags := parseTagdata(b)
	tags = append(tags, []string{"default", "*"})
	base := filepath.Base(p)
	for _, vals := range tags {
		if len(vals) > 1 && vals[1] == base || vals[1] == "*" {
			handler := srvr.tagHandlers[vals[0]]
			if handler != nil {
				err := handler(srvr, vals[1:], w, req)
				if err != nil {
					http.Error(w, "Internal error.", 500)
					return
				}
				break
			}
		}
	}
}

func(srvr * Server) ServeFile(name string, w http.ResponseWriter, req * http.Request) {
	name = filepath.Join(srvr.root, filepath.Clean("/" + name))
	if filepath.Base(name) == srvr.tagfile {
		name = ""
	}
	b, err := ioutil.ReadFile(name)
	if err != nil {
		http.Error(w, "Not found.", 404)
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(name)))
	w.Write(b)
}

func(srvr * Server) SetTagfileName(name string) {
	srvr.tagfile = name
}