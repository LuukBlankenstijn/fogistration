package client

import (
	"encoding/base64"
	"net/http"
)

type basicAuthRT struct {
	user string
	pass string
	next http.RoundTripper
}

func NewBasicAuthRT(user, pass string, next http.RoundTripper) http.RoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &basicAuthRT{user: user, pass: pass, next: next}
}

func (t *basicAuthRT) RoundTrip(req *http.Request) (*http.Response, error) {
	// clone to avoid mutating caller's req
	r := req.Clone(req.Context())
	if t.user != "" || t.pass != "" {
		cred := t.user + ":" + t.pass
		enc := base64.StdEncoding.EncodeToString([]byte(cred))
		r.Header.Set("Authorization", "Basic "+enc)
	}
	return t.next.RoundTrip(r)
}
