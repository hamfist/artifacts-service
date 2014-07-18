package server

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
)

func (srv *Server) serveError(err error, w http.ResponseWriter, r *http.Request) int {
	srv.log.WithFields(logrus.Fields{
		"err": err,
	}).Error("serving 500 response")

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/vnd.api+json")
	fmt.Fprintf(w, `{"error":%q}`, err.Error())
	return http.StatusInternalServerError
}

func (srv *Server) serveUnauthorized(w http.ResponseWriter, r *http.Request) int {
	srv.log.Warn("serving 401 response")

	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Header().Set("WWW-Authenticate", "token")
	fmt.Fprintf(w, `{"error":"unauthorized"}`)
	return http.StatusUnauthorized
}
