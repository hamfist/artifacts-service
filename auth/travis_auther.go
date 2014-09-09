package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	missingTokenError = fmt.Errorf("missing token")

	theFuture = time.Now().UTC().Add((100 * 365 * 24) * time.Hour)
)

type travisJob struct {
	StartedAt     time.Time `json:"started_at"`
	HasStartedAt  bool      `json:"-"`
	HasFinishedAt bool      `json:"-"`
}

func (tj *travisJob) UnmarshalJSON(rawBytes []byte) error {
	raw := map[string]string{}
	err := json.Unmarshal(rawBytes, &raw)
	if err != nil {
		return err
	}

	var startedAt string
	startedAt, tj.HasStartedAt = raw["started_at"]

	_, tj.HasFinishedAt = raw["finished_at"]

	if !tj.HasStartedAt {
		return nil
	}

	tj.StartedAt, err = time.Parse(time.RFC3339, startedAt)
	return err
}

type TravisAuther struct {
	TravisAPI string
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

	if r.Method == "GET" || r.Method == "HEAD" {
		return ar
	}

	if ta.inValidTimeWindow(r.Header.Get("Artifacts-Timestamp"), token, jobID) {
		ar.CanWrite = true
	}

	return ar
}

func (ta *TravisAuther) canReadBuilds(token, jobID string) bool {
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

func (ta *TravisAuther) inValidTimeWindow(reqDate, token, jobID string) bool {
	parsedDate, err := time.Parse(time.RFC3339, reqDate)
	if err != nil {
		return false
	}

	job, err := ta.getJob(token, jobID)
	if err != nil {
		return false
	}

	return (job.HasStartedAt && parsedDate.After(job.StartedAt) && !job.HasFinishedAt)
}

func (ta *TravisAuther) getJob(token, jobID string) (*travisJob, error) {
	r, err := ta.buildsRequest(token, jobID)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	job := &travisJob{
		StartedAt: theFuture,
	}
	err = json.NewDecoder(resp.Body).Decode(job)
	if err != nil {
		return nil, err
	}

	return job, nil
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
