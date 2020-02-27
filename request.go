package httpfilter

import "net/http"

type FilterRequest struct {
	*http.Request
	Query      string
	Args       []string
	setSession func(v interface{})
	getSession func() interface{}
}

func (req *FilterRequest) SetSession(v interface{}) {
	req.setSession(v)
}

func (req *FilterRequest) GetSession() interface{} {
	return req.getSession()
}
