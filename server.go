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
		"deft":     sv.serveFile,
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
	query := filepath.Join(sv.root, req.URL.Path)
	dir := filepath.Dir(query)
	query = filepath.Base(query)
	filters := parseFilterFile(filepath.Join(dir, filterFileName))
	var session interface{}
	for _, v := range filters {
		if !<-wr.ok {
			break
		}
		wr.ok <- true
		if match(query, v[1]) {
			if op := sv.ops[v[0]]; op != nil {
				query = op(wr, FilterRequest{
					Request: req,
					Query:   query,
					Args:    v[2:],
					setSession: func(v interface{}) {
						session = v
					},
					getSession: func() interface{} {
						return session
					},
				})
			} else {
				panic(errors.New("Undefined operator " + v[0]))
			}
		}
	}
}

func (sv *Server) serveFile(w http.ResponseWriter, req FilterRequest) string {
	if req.Query == filterFileName {
		http.Error(w, "Not found.", 404)
		return ""
	}
	name := filepath.Dir(req.URL.Path)
	name = filepath.Join(sv.root, name)
	name = filepath.Join(name, req.Query)
	b, err := ioutil.ReadFile(name)
	if err != nil {
		http.Error(w, "Not found.", 404)
		return ""
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(req.Query)))
	w.Write(b)
	return ""
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
