package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/meatballhat/artifacts-service/artifact"
)

func (srv *Server) saveHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) int {
	repoSlug, repoSlugOK := vars["slug"]
	filepath, filepathOK := vars["filepath"]
	jobID, jobIDOK := vars["job_id"]

	if !repoSlugOK || !filepathOK || !jobIDOK {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"this will never work.  stop it"}`)
		return http.StatusBadRequest
	}

	src := r.Header.Get("Artifacts-Source")
	size, _ := strconv.ParseUint(r.Header.Get("Artifacts-Size"), 10, 64)

	// TODO: validation!

	a := artifact.New(repoSlug, src, filepath, jobID, r.Body, size)

	err := srv.store.Store(a)
	if err != nil {
		return serveError(err, w, r)
	}

	resp := newSaveResponse()

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		return serveError(err, w, r)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/vnd.api+json")
	fmt.Fprintf(w, string(jsonBytes)+"\n")
	return http.StatusOK
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
