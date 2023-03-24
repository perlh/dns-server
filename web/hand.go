package web

import "net/http"

type HttpHandler func(res http.ResponseWriter, req *http.Request)

func (h HttpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	h(res, req)
}
