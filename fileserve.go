package fileserve

import (
	"net/http"
	"path/filepath"
	"io/ioutil"
	"mime"
)

type FileServer struct {
	basePath string
	tagHandlers map[string]TagHandlerFunc
}

func NewFileServer(basePath string, tagHandlers map[string]TagHandlerFunc) * FileServer {
	defaultTagHandlers := map[string]TagHandlerFunc {
		"ignore": ignore,
		"pseudo": pseudo,
		"redirect": redirect,
	}
	if tagHandlers != nil {
		for k, v := range tagHandlers {
			defaultTagHandlers[k] = v
		}
	}
	return &FileServer{
		basePath: basePath,
		tagHandlers: defaultTagHandlers,
	}
}

const metaname = "#meta.txt"

func(srvr * FileServer) ServeHTTP(w http.ResponseWriter, req * http.Request) {
	p := filepath.Join(srvr.basePath, filepath.Clean("/" + req.URL.Path))
	written := false
	metapath := filepath.Join(filepath.Dir(p), metaname)
	b, err := ioutil.ReadFile(metapath)
	metadata := make(map[string][]string)
	if err == nil {
		metadata = parseMetadata(b)
	}

	for tag, vals := range metadata {
		handler := srvr.tagHandlers[tag]
		if handler != nil {
			handler(vals, &p, w, req, &written)
		}
	}

	if !written {
		if filepath.Base(p) == metaname {
			p = ""
		}
		b, err = ioutil.ReadFile(p)
		if err != nil {
			http.Error(w, "Not found", 404)
		}
		w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(p)))
		w.Write(b)
	}
}