package httpfilter

import (
	"errors"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
)

var filterFileName = "_filters.txt"

type Server struct {
	root string
	ops  map[string]OpFunc
}

func NewServer(root string, ops ...map[string]OpFunc) *Server {
	sv := &Server{
		root: root,
	}
	sv.ops = map[string]OpFunc{
		"serve":    sv.serveFile,
		"ignore":   ignore,
		"redirect": redirect,
		"proxy":    proxy,
		"request":  request,
	}
	for _, m := range ops {
		for k, v := range m {
			sv.ops[k] = v
		}
	}
	return sv
}

func (sv *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	wr := wrapWriter(w)
	name := filepath.Join(sv.root, req.URL.Path)
	dir, name := filepath.Split(name)
	filters := parseFilterFile(filepath.Join(dir, filterFileName))
	for _, v := range filters {
		if _, ok := <-wr.ok; !ok {
			break
		}
		wr.ok <- 0
		if match(name, v[1]) {
			if op := sv.ops[v[0]]; op != nil {
				op(wr, req, v[2:]...)
			} else {
				panic(errors.New("Undefined operator " + v[0]))
			}
		}
	}
	if _, ok := <-wr.ok; ok {
		sv.serveFile(w, req, name)
	}
}

func (sv *Server) serveFile(w http.ResponseWriter, req *http.Request, args ...string) {
	path := filepath.Join(sv.root, filepath.Dir(req.URL.Path), args[0])
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

func match(q, s string) bool {
	se := filepath.Ext(s)    //selector ext
	sn := s[:len(s)-len(se)] //selector name
	qe := filepath.Ext(q)    //query ext
	qn := q[:len(q)-len(qe)] //query name
	return (q == s) ||       //name and extension match
		(s == "*") || //selector is *
		(se == ".*" && qn == sn) || //selector ext is * and name matches
		(sn == "*" && qe == se) //selector name is * and ext matches
}
