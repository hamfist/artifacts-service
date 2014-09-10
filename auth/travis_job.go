package auth

import (
	"encoding/json"
	"time"
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
