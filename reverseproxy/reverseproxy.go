package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/rbxb/fileserve"
	"errors"
)

func ReverseProxyTagHandler(srvr * fileserve.Server, args []string, w http.ResponseWriter, req * http.Request) error {
	if len(args) < 1 {
		return fileserve.ErrorNotEnoughArguments
	}
	u, err := url.Parse(args[0])
	if err != nil {
		return errors.New("Couldn't parse URL.")
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, req)
	return nil
}