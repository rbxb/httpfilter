package request

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func RequestTagHandler(vals []string, p * string, w http.ResponseWriter, req * http.Request, written * bool) {
	basename := filepath.Base(*p)
	for i := 0; i + 1 < len(vals); i += 2 {
		if vals[i] == basename || vals[i] == "*" {
			if !*written {
				resp, err := http.Get(vals[i + 1])
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()
				b, err := ioutil.ReadAll(resp.Body)
				w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
				w.Write(b)
				*written = true
				break
			}
		}
	}
}