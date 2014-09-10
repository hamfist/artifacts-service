package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"strings"
	"time"
)

var (
	missingPrivateKeyError = fmt.Errorf("missing rsa private key")
)

type travisTimestamp struct {
	Raw string
	pk  *rsa.PrivateKey
}

func newTravisTimestamp(raw string, pk *rsa.PrivateKey) *travisTimestamp {
	return &travisTimestamp{
		Raw: raw,
		pk:  pk,
	}
}

func (tts *travisTimestamp) String() string {
	return fmt.Sprintf("&auth.travisTimestamp{%q, [secrets]}", tts.Raw)
}

func (tts *travisTimestamp) Value() (time.Time, error) {
	if !strings.HasPrefix(tts.Raw, "RSA:") {
		return time.Parse(time.RFC3339, tts.Raw)
	}

	if tts.pk == nil {
		return theFuture, missingPrivateKeyError
	}

	enc := strings.Replace(tts.Raw, "RSA:", "", 1)

	decr, err := rsa.DecryptPKCS1v15(rand.Reader, tts.pk, []byte(enc))
	if err != nil {
		return theFuture, err
	}

	return time.Parse(time.RFC3339, string(decr))
}
