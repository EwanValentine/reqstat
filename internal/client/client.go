package client

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	headers    map[string]string
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}
}

func (c *Client) AddHeader(header string) {
	parts := strings.SplitN(header, ":", 2)
	if len(parts) == 2 {
		c.headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
}

func (c *Client) Get(url string) (*Result, error) {
	return c.do("GET", url, nil)
}

func (c *Client) do(method, url string, body io.Reader) (*Result, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "reqstat/1.0")
	}

	var dnsStart, dnsEnd time.Time
	var tcpStart, tcpEnd time.Time
	var tlsStart, tlsEnd time.Time
	var firstByte time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			dnsEnd = time.Now()
		},
		ConnectStart: func(_, _ string) {
			tcpStart = time.Now()
		},
		ConnectDone: func(_, _ string, _ error) {
			tcpEnd = time.Now()
		},
		TLSHandshakeStart: func() {
			tlsStart = time.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			tlsEnd = time.Now()
		},
		GotFirstResponseByte: func() {
			firstByte = time.Now()
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start)

	result := &Result{
		URL:           url,
		StatusCode:    resp.StatusCode,
		Status:        resp.Status,
		Headers:       resp.Header,
		Body:          bodyBytes,
		ContentLength: int64(len(bodyBytes)),
		Duration:      duration,
		ContentType:   resp.Header.Get("Content-Type"),
	}

	if !dnsEnd.IsZero() && !dnsStart.IsZero() {
		result.DNSLookup = dnsEnd.Sub(dnsStart)
	}
	if !tcpEnd.IsZero() && !tcpStart.IsZero() {
		result.TCPConnection = tcpEnd.Sub(tcpStart)
	}
	if !tlsEnd.IsZero() && !tlsStart.IsZero() {
		result.TLSHandshake = tlsEnd.Sub(tlsStart)
	}
	if !firstByte.IsZero() {
		result.ServerResponse = firstByte.Sub(start)
	}

	return result, nil
}
