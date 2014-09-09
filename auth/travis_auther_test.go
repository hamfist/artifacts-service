package auth

import (
	"net/http"
	"os"
	"testing"
)

var (
	testTravisAPI = os.Getenv("ARTIFACTS_TEST_TRAVIS_API")

	travisAutherTestCases = []*travisAutherTestCase{
		&travisAutherTestCase{
			Vars: map[string]string{},
		},
	}
)

type travisAutherTestCase struct {
	Method, Timestamp, AuthToken string
	CanRead, CanWrite            bool
	Vars                         map[string]string
}

func TestTravisAutherCheck(t *testing.T) {
	ta := &TravisAuther{TravisAPI: testTravisAPI}

	for _, tc := range travisAutherTestCases {
		r, err := http.NewRequest(tc.Method, "http://localhost:9999", nil)
		if err != nil {
			panic(err)
		}

		if tc.Timestamp != "" {
			r.Header.Set("Artifacts-Timestamp", tc.Timestamp)
		}
		if tc.AuthToken != "" {
			r.Header.Set("Authorization", "token "+tc.AuthToken)
		}

		ar := ta.Check(r, tc.Vars)
		if ar.CanRead != tc.CanRead {
			t.Fatalf("CanRead %v != %v", ar.CanRead, tc.CanRead)
		}
		if ar.CanWrite != tc.CanWrite {
			t.Fatalf("CanWrite %v != %v", ar.CanWrite, tc.CanWrite)
		}
	}
}
