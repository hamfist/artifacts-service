package server

import (
	"fmt"
	"net/http"
)

type uploadHandler struct{}

func newUploadHandler() *uploadHandler {
	return &uploadHandler{}
}

func (uh *uploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		uh.handlePostUpload(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "whatever, meatbag\n")
	}
}

func (uh *uploadHandler) handlePostUpload(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "why not!?\n")
}
