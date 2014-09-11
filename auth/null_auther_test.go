package auth

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewNullAuther(t *testing.T) {
	na := NewNullAuther()
	if !reflect.DeepEqual(na, &NullAuther{}) {
		t.Fatalf("unexpected new null auther")
	}
}

func TestNullAutherCheck(t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	na := NewNullAuther()
	ar := na.Check(r, map[string]string{})

	if !ar.CanRead {
		t.Fatalf("%v != true", ar.CanRead)
	}
	if !ar.CanWrite {
		t.Fatalf("%v != true", ar.CanWrite)
	}
}
