package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/meatballhat/artifacts-service/artifact"
)

func (srv *Server) saveHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	err := srv.store.Store(r.Body)
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
