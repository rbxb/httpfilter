package httpfilter

import (
	"bytes"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

var filterFileName = "_filters.txt"

type Server struct {
	root   string
	filter string
	ops    map[string]OpFunc
}

func NewServer(root string, filter string, ops ...map[string]OpFunc) *Server {
	sv := &Server{
		root:   root,
		filter: filter,
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
	var filter string
	if sv.filter == "" {
		filter = filepath.Join(dir, filterFileName)
	} else {
		filter = sv.filter
	}
	if p, err := ioutil.ReadFile(filter); err == nil {
		sv.parseFilter(p, wr, req)
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

func match(req *http.Request, s string) bool {
	q := ""
	if s[0] == byte('@') {
		split := strings.Split(req.Host, ".")
		if len(split) > 2 {
			q = split[0]
		}
		s = s[1:]
	} else {
		q = filepath.Base(req.URL.Path)
	}
	se := filepath.Ext(s)    //selector ext
	sn := s[:len(s)-len(se)] //selector name
	qe := filepath.Ext(q)    //query ext
	qn := q[:len(q)-len(qe)] //query name
	return (q == s) ||       //name and extension match
		(s == "*") || //selector is *
		(se == ".*" && qn == sn) || //selector ext is * and name matches
		(sn == "*" && qe == se) //selector name is * and ext matches
}

func (sv *Server) parseFilter(p []byte, wr *writerWrapper, req *http.Request) {
	lines := bytes.Split(p, []byte{'\n'})
	var op []byte
	for _, line := range lines {
		vals := make([]string, 0)
		for _, word := range bytes.Fields(line) {
			if word[0] == '#' {
				op = word[1:]
			} else {
				vals = append(vals, string(word))
			}
		}
		if len(vals) > 0 && match(req, vals[0]) {
			sv.ops[string(op)](wr, req, vals[1:]...)
			if _, ok := <-wr.ok; !ok {
				break
			}
			wr.ok <- 0
		}
	}
}
