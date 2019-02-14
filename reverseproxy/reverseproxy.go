package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ReverseProxyTagHandler(vals []string, p * string, w http.ResponseWriter, req * http.Request, written * bool) {
	if len(vals) < 0 {
		http.Error(w, "No url specified.", 500)
		*written = true
		return
	}
	u, err := url.Parse(vals[0])
	if err != nil {
		http.Error(w, "Error parsing url.", 500)
		*written = true
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, req)
	*written = true
}