package server

import (
	"fmt"
	"net/http"
)

func (srv *Server) rootHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "sure\n")
}
