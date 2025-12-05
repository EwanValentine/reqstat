package client

import (
	"net/http"
	"time"
)

type Result struct {
	URL            string
	StatusCode     int
	Status         string
	Headers        http.Header
	Body           []byte
	ContentLength  int64
	Duration       time.Duration
	DNSLookup      time.Duration
	TCPConnection  time.Duration
	TLSHandshake   time.Duration
	ServerResponse time.Duration
	ContentType    string
}

func (r *Result) IsJSON() bool {
	ct := r.Headers.Get("Content-Type")
	return ct == "application/json" ||
		ct == "application/json; charset=utf-8" ||
		ct == "application/json;charset=utf-8" ||
		ct == "application/json; charset=UTF-8" ||
		len(r.Body) > 0 && (r.Body[0] == '{' || r.Body[0] == '[')
}

func (r *Result) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

func (r *Result) IsRedirect() bool {
	return r.StatusCode >= 300 && r.StatusCode < 400
}

func (r *Result) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

func (r *Result) IsServerError() bool {
	return r.StatusCode >= 500
}
