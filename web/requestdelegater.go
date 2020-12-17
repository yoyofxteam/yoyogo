package web

import "net/http"

type IRequestDelegate interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
