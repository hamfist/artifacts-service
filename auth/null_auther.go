package auth

import (
	"net/http"
)

// NullAuther implements Auther but doesn't even care!
type NullAuther struct{}

// NewNullAuther creates a new *NullAuther yey
func NewNullAuther() *NullAuther {
	return &NullAuther{}
}

// Check always allows reads and writes
func (na *NullAuther) Check(r *http.Request, vars map[string]string) *Result {
	ar := NewResult(r, vars)
	ar.CanRead = true
	ar.CanWrite = true
	return ar
}
