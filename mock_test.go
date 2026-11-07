package sms

import (
	"context"
	"errors"
	"testing"
)

// fakeSender is a minimal Sender for copy-paste tests in downstream apps.
type fakeSender struct {
	sendFn func(context.Context, Message) error
}

func (f fakeSender) Send(ctx context.Context, msg Message) error {
	if f.sendFn == nil {
		return nil
	}
	return f.sendFn(ctx, msg)
}

func TestFakeSenderPattern(t *testing.T) {
	var called bool
	s := fakeSender{sendFn: func(_ context.Context, msg Message) error {
		called = true
		if msg.To == "" {
			return errors.New("empty")
		}
		return nil
	}}
	if err := s.Send(context.Background(), Message{To: "+1"}); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("send not called")
	}
}
