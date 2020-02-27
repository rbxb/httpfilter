package reverseproxy

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rbxb/httpfilter"
)

var Ops = map[string]httpfilter.OpFunc{"reverseproxy": ReverseProxy}

func ReverseProxy(w http.ResponseWriter, req httpfilter.FilterRequest) string {
	if len(req.Args) < 1 {
		panic(errors.New("not enough arguments"))
		return ""
	}
	u, err := url.Parse(req.Args[0])
	if err != nil {
		panic(errors.New("couldn't parse URL"))
		return ""
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, req.Request)
	return ""
}
