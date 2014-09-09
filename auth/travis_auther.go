package auth

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"
)

var (
	missingTokenError = fmt.Errorf("missing token")
)

type TravisAuther struct {
	TravisAPI  string
	PrivateKey *rsa.PrivateKey
}

func (ta *TravisAuther) Check(r *http.Request, vars map[string]string) *AuthResult {
	ar := NewAuthResult(r, vars)

	token := r.FormValue("token")
	if token == "" {
		token = strings.Replace(r.Header.Get("Authorization"), "token ", "", 1)
	}

	if token == "" {
		ar.Errors = append(ar.Errors, missingTokenError)
		return ar
	}

	jobID, ok := vars["job_id"]
	if !ok {
		jobID = ""
	}

	if ta.canReadBuilds(token, jobID) {
		ar.CanRead = true
	}

	if ta.inValidTimeWindow(token, jobID) {
		ar.CanWrite = true
	}

	return ar
}

func (ta *TravisAuther) canReadBuilds(token, jobID string) bool {
	// XXX should this be cached somewhere?
	client := &http.Client{}

	r, err := ta.buildsRequest(token, jobID)
	if err != nil {
		return false
	}

	resp, err := client.Do(r)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}

func (ta *TravisAuther) inValidTimeWindow(token, jobID string) bool {
	// TODO implementation lulz
	return false
}

func (ta *TravisAuther) buildsRequest(token, jobID string) (*http.Request, error) {
	req, err := http.NewRequest("GET", ta.TravisAPI, nil)
	if err != nil {
		return nil, err
	}

	req.URL.Path = fmt.Sprintf("/jobs/%s", jobID)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))

	return req, nil
}
