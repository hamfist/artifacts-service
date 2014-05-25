package server

import (
	"fmt"
	"net/http"
	"path"

	"github.com/Sirupsen/logrus"
)

func (srv *Server) getPathHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	srv.log.WithFields(logrus.Fields{
		"vars": vars,
	}).Debug("not really serving")

	a, err := srv.store.Fetch(vars["slug"], vars["filepath"])
	if err != nil || a == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such thing\n")
		return
	}

	reader, err := a.Reader()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "kaboom\n")
		return
	}

	http.ServeContent(w, r, path.Base(vars["filepath"]), a.DateModified, reader)
}
