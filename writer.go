package httpfilter

import (
	"net/http"
)

type writerWrapper struct {
	http.ResponseWriter
	http.Hijacker
	ok chan byte
}

func (wr *writerWrapper) WriteHeader(statusCode int) {
	wr.ResponseWriter.WriteHeader(statusCode)
	if _, ok := <-wr.ok; ok {
		close(wr.ok)
	}
}

func (wr *writerWrapper) Write(b []byte) (int, error) {
	n, err := wr.ResponseWriter.Write(b)
	if _, ok := <-wr.ok; ok {
		close(wr.ok)
	}
	return n, err
}

func wrapWriter(w http.ResponseWriter) *writerWrapper {
	wr := &writerWrapper{
		ResponseWriter: w,
		Hijacker:       w.(http.Hijacker),
		ok:             make(chan byte, 1),
	}
	wr.ok <- 0
	return wr
}
