package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/rbxb/fileserve"
	"errors"
)

func ReverseProxyTagHandler(srvr * fileserve.Server, vals []string, w http.ResponseWriter, req * http.Request) error {
	if len(vals) < 2 {
		return fileserve.ErrorNotEnoughArguments
	}
	u, err := url.Parse(vals[1])
	if err != nil {
		return errors.New("Couldn't parse URL.")
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, req)
	return nil
}