package server

import (
	"fmt"
	"net/http"
)

func (srv *Server) listHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) int {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "nope\n")
	return http.StatusNotImplemented
}
