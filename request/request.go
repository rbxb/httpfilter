package request

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/rbxb/httpfilter"
)

var Ops = map[string]httpfilter.OpFunc{"request": Request}

func Request(w http.ResponseWriter, req httpfilter.FilterRequest) string {
	if len(req.Args) < 1 {
		panic(errors.New("not enough arguments"))
		return ""
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
