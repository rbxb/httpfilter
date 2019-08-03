package request

import (
	"io/ioutil"
	"net/http"
	"errors"
	"github.com/rbxb/httpfilter"
)

var Ops = map[string]httpfilter.FilterOpFunc{"request":Request}

func Request(w http.ResponseWriter, req * http.Request, query string, args []string) string {
	if len(args) < 1 {
		panic(errors.New("Not enough arguments."))
		return ""
	}
	resp, err := http.Get(args[0])
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