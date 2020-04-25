package httpfilter

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type OpFunc func(w http.ResponseWriter, req FilterRequest) string

func ignore(w http.ResponseWriter, req FilterRequest) string {
	http.Error(w, "Not found.", 404)
	return ""
}

func pseudo(w http.ResponseWriter, req FilterRequest) string {
	if len(req.Args) < 1 {
		panic(errors.New("not enough arguments"))
	}
	return req.Args[0]
}

func redirect(w http.ResponseWriter, req FilterRequest) string {
	if len(req.Args) < 1 {
		panic(errors.New("not enough arguments"))
	}
	http.Redirect(w, req.Request, req.Args[0], 301)
	return ""
}

func proxy(w http.ResponseWriter, req FilterRequest) string {
	if len(req.Args) < 1 {
		panic(errors.New("not enough arguments"))
	}
	u, err := url.Parse(req.Args[0])
	if err != nil {
		panic(errors.New("couldn't parse URL"))
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		http.Error(w, "Internal error.", 500)
	}
	proxy.ServeHTTP(w, req.Request)
	return ""
}

func request(w http.ResponseWriter, req FilterRequest) string {
	if len(req.Args) < 1 {
		panic(errors.New("not enough arguments"))
	}
	resp, err := http.Get(req.Args[0])
	if err != nil {
		http.Error(w, "Internal error.", 500)
		return ""
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Internal error.", 500)
		return ""
	}
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Write(b)
	return ""
}
