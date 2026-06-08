package shared

import (
	"encoding/gob"
	"net/url"
)

func init() {
	gob.Register(&url.URL{})
}

type Request struct {
	RemoteAddr string
	Host       string
	Method     string
	URL        *url.URL
	Header     map[string][]string
}

type Requests []Request
