package httpfilter

import (
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
)

func serveFile(w http.ResponseWriter, path string) {
	name := filepath.Base(path)
	if name == filterFileName {
		http.Error(w, "Not found.", 404)
		return
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, "Not found.", 404)
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(name))
	w.Write(b)
}
