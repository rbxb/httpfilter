package httpfilter

import (
	"net/http"
)

type writerWrapper struct {
	http.ResponseWriter
	ok chan bool
}

func(wrapper * writerWrapper) WriteHeader(statusCode int) {
	<- wrapper.ok
	wrapper.ResponseWriter.WriteHeader(statusCode)
	wrapper.ok <- false
}

func wrapWriter(w http.ResponseWriter) * writerWrapper {
	wrapper := &writerWrapper{
		ResponseWriter: w,
		ok: make(chan bool, 1),
	}
	wrapper.ok <- true
	return wrapper
}