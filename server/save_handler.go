package server

import (
	"fmt"
	"net/http"
)

func (srv *Server) saveHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "also nope\n")
}
