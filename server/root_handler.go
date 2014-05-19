package server

import (
	"fmt"
	"net/http"
)

type rootHandler struct{}

func newRootHandler() *rootHandler {
	return &rootHandler{}
}

func (rh *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rh.handleGetRoot(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "whatever, meatbag\n")
	}
}

func (rh *rootHandler) handleGetRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "sure\n")
}
