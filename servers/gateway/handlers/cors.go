package handlers

import (
	"log"
	"net/http"
)

/* TODO: implement a CORS middleware handler, as described
in https://drstearns.github.io/tutorials/cors/ that responds
with the following headers to all requests:

  Access-Control-Allow-Origin: *
  Access-Control-Allow-Methods: GET, PUT, POST, PATCH, DELETE
  Access-Control-Allow-Headers: Content-Type, Authorization
  Access-Control-Expose-Headers: Authorization
  Access-Control-Max-Age: 600
*/

//Begin middleware
type ResponseHeader struct {
	handler http.Handler
	k1      string
	v1      string
	k2      string
	v2      string
	k3      string
	v3      string
	k4      string
	v4      string
	k5      string
	v5      string
}

func AddFiveResponseHeaders(handlerToWrap http.Handler, k1 string, v1 string, k2 string, v2 string, k3 string, v3 string, k4 string, v4 string, k5 string, v5 string) *ResponseHeader {
	return &ResponseHeader{handlerToWrap, k1, v1, k2, v2, k3, v3, k4, v4, k5, v5}
}

func (rh *ResponseHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(rh.k1, rh.v1)
	w.Header().Set(rh.k2, rh.v2)
	w.Header().Set(rh.k3, rh.v3)
	w.Header().Set(rh.k4, rh.v4)
	w.Header().Set(rh.k5, rh.v5)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}

	rh.handler.ServeHTTP(w, r)
	log.Printf("Returning a request for %s", r.URL.Path)
}

//end middleware
