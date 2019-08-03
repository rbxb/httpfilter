package reverseproxy

import (
	"net/http"
	"errors"
	"net/url"
	"net/http/httputil"
	"github.com/rbxb/httpfilter"
)

var Ops = map[string]httpfilter.FilterOpFunc{"reverseproxy":ReverseProxy}

func ReverseProxy(w http.ResponseWriter, req * http.Request, query string, args []string) string {
	if len(args) < 1 {
		panic(errors.New("Not enough arguments."))
		return ""
	}
	u, err := url.Parse(args[0])
	if err != nil {
		panic(errors.New("Couldn't parse URL."))
		return ""
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, req)
	return ""
}