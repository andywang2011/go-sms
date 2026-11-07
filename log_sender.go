package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type logSender struct {
	w io.Writer
}

// NewLogSender returns a Sender that writes JSON lines to w (for local/dev).
func NewLogSender(w io.Writer) Sender {
	if w == nil {
		w = io.Discard
	}
	return &logSender{w: w}
}

// Send writes the message as one JSON object per line.
func (s *logSender) Send(_ context.Context, msg Message) error {
	if err := validateRecipient(msg.To); err != nil {
		return err
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("log sender encode: %w", err)
	}
	_, err = fmt.Fprintf(s.w, "%s\n", b)
	return err
}
