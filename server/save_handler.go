package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/meatballhat/artifacts-service/artifact"
)

func (srv *Server) saveHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	repoSlug := r.Header.Get("Artifacts-Repo-Slug")
	src := r.Header.Get("Artifacts-Source")
	dest := r.Header.Get("Artifacts-Destination")
	jobNumber := r.Header.Get("Artifacts-Job-Number")
	size, _ := strconv.ParseUint(r.Header.Get("Artifacts-Size"), 10, 64)

	// TODO: validation!

	a := artifact.New(repoSlug, src, dest, jobNumber, r.Body, size)

	err := srv.store.Store(a)
	if err != nil {
		serveError(err, w, r)
		return
	}

	resp := newSaveResponse()

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		serveError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/vnd.api+json")
	fmt.Fprintf(w, string(jsonBytes)+"\n")
}

type saveResponse struct {
	Artifacts []*artifact.Artifact `json:"artifacts"`
	Meta      *saveMeta            `json:"meta"`
}

func newSaveResponse() *saveResponse {
	return &saveResponse{
		Artifacts: []*artifact.Artifact{},
		Meta: &saveMeta{
			Artifacts: &artifactsMeta{},
		},
	}
}

type saveMeta struct {
	Artifacts *artifactsMeta `json:"artifacts"`
}

type artifactsMeta struct {
}
