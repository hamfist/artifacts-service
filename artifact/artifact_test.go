package artifact

import (
	"bytes"
	"io"
	"testing"
)

var (
	testArtifactInStr = "how about this"
	inSize            = uint64(len(testArtifactInStr))
	in                = bytes.NewReader([]byte(testArtifactInStr))
)

type testReadSeeker struct {
	io.ReadSeeker
}

type testReadCloser struct {
	io.ReadCloser
}

func getNewTestArtifact() *Artifact {
	return New("src", "dest", "1234", in, inSize)
}

func TestNew(t *testing.T) {
	a := getNewTestArtifact()

	for actual, expected := range map[string]string{
		a.Source:      "src",
		a.Destination: "dest",
		a.JobID:       "1234",
	} {
		if actual != expected {
			t.Fatalf("%v != %v", actual, expected)
		}
	}

	if a.Size != inSize {
		t.Fatalf("%v != %v", a.Size, inSize)
	}
}

func TestReadCloser(t *testing.T) {
	a := getNewTestArtifact()
	a.OutReadCloser = &testReadCloser{}
	rc, err := a.ReadCloser()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if rc == nil {
		t.Fatalf("ReadCloser is nil")
	}

	a.OutReadCloser = nil
	rc, err = a.ReadCloser()
	if err != errNoReader {
		t.Fatalf("ReadCloser return non-error")
	}

	if rc != nil {
		t.Fatalf("ReadCloser is not nil")
	}
}

func TestReadSeeker(t *testing.T) {
	a := getNewTestArtifact()
	a.OutReadSeeker = &testReadSeeker{}
	rc, err := a.ReadSeeker()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if rc == nil {
		t.Fatalf("ReadSeeker is nil")
	}

	a.OutReadSeeker = nil
	rc, err = a.ReadSeeker()
	if err != errNoReader {
		t.Fatalf("ReadSeeker return non-error")
	}

	if rc != nil {
		t.Fatalf("ReadSeeker is not nil")
	}
}

func TestFullDestination(t *testing.T) {
	a := getNewTestArtifact()
	if a.FullDestination() != "job/1234/dest" {
		t.Fatalf("destination != job/1234/dest")
	}
}
