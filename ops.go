package httpfilter

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type OpFunc func(w http.ResponseWriter, req *http.Request, args ...string)

func ignore(w http.ResponseWriter, req *http.Request, args ...string) {
	http.Error(w, "Not found.", 404)
}

func redirect(w http.ResponseWriter, req *http.Request, args ...string) {
	http.Redirect(w, req, args[0], 302)
}

func proxy(w http.ResponseWriter, req *http.Request, args ...string) {
	u, err := url.Parse(args[0])
	if err != nil {
		panic(errors.New("couldn't parse URL"))
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		http.Error(w, "Internal error.", 500)
		panic(err)
	}
	proxy.ServeHTTP(w, req)
}

func request(w http.ResponseWriter, req *http.Request, args ...string) {
	resp, err := http.Get(args[0])
	if err != nil {
		http.Error(w, "Internal error.", 500)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Internal error.", 500)
		return
	}
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Write(b)
}
