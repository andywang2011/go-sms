package sms

import (
	"context"
	"net/http"
)

// Sender sends SMS messages (mock-friendly client surface).
type Sender interface {
	// Send delivers msg using the provider API or sink.
	Send(ctx context.Context, msg Message) error
}

// HTTPDoer performs HTTP requests; *http.Client implements this interface.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
