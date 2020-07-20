package httpfilter

import (
	"net/http"
)

type writerWrapper struct {
	http.ResponseWriter
	ok chan byte
}

func (wrapper *writerWrapper) WriteHeader(statusCode int) {
	wrapper.ResponseWriter.WriteHeader(statusCode)
	if _, ok := <-wrapper.ok; ok {
		close(wrapper.ok)
	}
}

func (wrapper *writerWrapper) Write(b []byte) (int, error) {
	n, err := wrapper.ResponseWriter.Write(b)
	if _, ok := <-wrapper.ok; ok {
		close(wrapper.ok)
	}
	return n, err
}

func wrapWriter(w http.ResponseWriter) *writerWrapper {
	wrapper := &writerWrapper{
		ResponseWriter: w,
		ok:             make(chan byte, 1),
	}
	wrapper.ok <- 0
	return wrapper
}
