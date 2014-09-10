package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"strings"
	"time"
)

var (
	errMissingPrivKey       = fmt.Errorf("missing rsa private key")
	errUnencryptedTimestamp = fmt.Errorf("timestamp is not rsa encrypted")
)

type travisTimestamp struct {
	RawTimestamp string
	privKey      *rsa.PrivateKey
	rsaRequired  bool
}

func newTravisTimestamp(raw string, privKey *rsa.PrivateKey, rsaRequired bool) *travisTimestamp {
	return &travisTimestamp{
		RawTimestamp: raw,
		privKey:      privKey,

		rsaRequired: rsaRequired,
	}
}

func (tts *travisTimestamp) String() string {
	return fmt.Sprintf("&auth.travisTimestamp{%q, [secrets]}", tts.RawTimestamp)
}

func (tts *travisTimestamp) Value() (time.Time, error) {
	if !strings.HasPrefix(tts.RawTimestamp, "RSA:") {
		if tts.rsaRequired {
			return theFuture, errUnencryptedTimestamp
		}
		return time.Parse(time.RFC3339, tts.RawTimestamp)
	}

	if tts.privKey == nil {
		return theFuture, errMissingPrivKey
	}

	enc := strings.Replace(tts.RawTimestamp, "RSA:", "", 1)

	decr, err := rsa.DecryptPKCS1v15(rand.Reader, tts.privKey, []byte(enc))
	if err != nil {
		return theFuture, err
	}

	return time.Parse(time.RFC3339, string(decr))
}
