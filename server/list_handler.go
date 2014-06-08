package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/hamfist/artifacts-service/metadata"
)

type listMetadataResponse struct {
	Metadata []*metadata.Metadata       `json:"data"`
	Meta     *listMetadataResponseMeta  `json:"meta"`
	Links    *listMetadataResponseLinks `json:"links"`
}

// zomg this is getting reeeeeally meta
type listMetadataResponseMeta struct {
	Href string `json:"href"`
}

type listMetadataResponseLinks struct {
	DataPath *listMetadataResponseLinksDataPath `json:"data.path"`
}

type listMetadataResponseLinksDataPath struct {
	Href string `json:"href"`
	Type string `json:"type"`
}

func (srv *Server) listHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) int {
	allMetadata, err := srv.md.LookupAll(vars["job_id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "sorry mate\n")
		return http.StatusInternalServerError
	}

	return srv.sendAllMetadataResponse(w, allMetadata)
}

func (srv *Server) sendAllMetadataResponse(w http.ResponseWriter, mds []*metadata.Metadata) int {

	resp := &listMetadataResponse{
		Metadata: mds,
		Meta: &listMetadataResponseMeta{
			Href: "TODO",
		},
		Links: &listMetadataResponseLinks{
			DataPath: &listMetadataResponseLinksDataPath{
				Href: "TODO",
				Type: "artifact",
			},
		},
	}

	respBytes, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		srv.log.WithFields(logrus.Fields{"err": err}).Error("failed to marshal list metadata response")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "sorry mate\n")
		return http.StatusInternalServerError
	}

	respString := string(respBytes)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(respString)))
	fmt.Fprintf(w, respString)
	return http.StatusOK
}
