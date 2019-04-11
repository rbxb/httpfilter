package request

import (
	"io/ioutil"
	"net/http"
	"github.com/rbxb/fileserve"
)

func RequestTagHandler(srvr * fileserve.Server, args []string, w http.ResponseWriter, req * http.Request) error {
	if len(args) < 1 {
		return fileserve.ErrorNotEnoughArguments
	}
	resp, err := http.Get(args[0])
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