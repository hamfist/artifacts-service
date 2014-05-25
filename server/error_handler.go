package server

import (
	"fmt"
	"net/http"
)

func serveError(err error, w http.ResponseWriter, r *http.Request) int {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/vnd.api+json")
	fmt.Fprintf(w, fmt.Sprintf("{\"error\":%q}", err.Error()))
	return http.StatusInternalServerError
}
