package server

import (
	"fmt"
	"io"
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

	rs, _ := a.ReadSeeker()
	if rs != nil {
		http.ServeContent(w, r, path.Base(vars["filepath"]), a.DateModified, rs)
		return http.StatusOK
	}

	rc, err := a.ReadCloser()
	if err != nil {
		return srv.serveError(err, w, r)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", a.Size))
	w.Header().Set("Content-Type", a.ContentType)
	_, err = io.Copy(w, rc)
	if err != nil {
		srv.log.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to copy all bytes")
	}

	return http.StatusOK
}
