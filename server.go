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
	tags := make([][]string, 0)
	p := filepath.Join(srvr.root, filepath.Clean("/" + req.URL.Path))
	b, err := ioutil.ReadFile(filepath.Join(filepath.Dir(p), srvr.tagfile))
	if err == nil {
		tags = parseTagdata(b)
	}
	tags = append(tags, []string{"default", "*"})
	base := filepath.Base(p)
	ext := filepath.Ext(base)
	name := base[:len(base) - len(ext)]
	for _, vals := range tags {
		vbase := vals[1]
		vext := filepath.Ext(vbase)
		vname := vals[1][:len(vbase) - len(vext)]
		if (vbase == base) ||               // if name and extension match
		(vbase == "*") ||                   // or if name in tagfile is *
		(vext == ".*" && vname == name) ||  // or if name matches and extension in tagfile is *
		(vname == "*" && vext == ext) {     // or if extension matches and name in tagfile is *
			handler := srvr.tagHandlers[vals[0]]
			if handler != nil {
				err := handler(srvr, vals[2:], w, req)
				switch err {
				case nil:
					break
				default:
					http.Error(w, "Internal error.", 500)
					return
				}
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