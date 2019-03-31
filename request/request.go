package request

import (
	"io/ioutil"
	"net/http"
	"github.com/rbxb/fileserve"
)

func RequestTagHandler(srvr * fileserve.Server, vals []string, w http.ResponseWriter, req * http.Request) error {
	if len(vals) < 2 {
		return fileserve.ErrorNotEnoughArguments
	}
	resp, err := http.Get(vals[1])
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Write(b)
	return nil
}