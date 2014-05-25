package server

import (
	"fmt"
	"net/http"
	"path"

	"github.com/Sirupsen/logrus"
)

func (srv *Server) getPathHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) int {
	srv.log.WithFields(logrus.Fields{
		"slug":     vars["slug"],
		"filepath": vars["filepath"],
		"job_id":   vars["job_id"],
	}).Debug("serving path")

	a, err := srv.store.Fetch(vars["slug"], vars["filepath"], vars["job_id"])
	if err != nil || a == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such thing\n")
		return http.StatusNotFound
	}

	reader, err := a.Reader()
	if err != nil {
		return serveError(err, w, r)
	}

	http.ServeContent(w, r, path.Base(vars["filepath"]), a.DateModified, reader)
	return http.StatusOK
}
