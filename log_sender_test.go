package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestLogSenderSend(t *testing.T) {
	var buf bytes.Buffer
	sender := NewLogSender(&buf)
	err := sender.Send(context.Background(), Message{
		To:   "+15551234567",
		Body: "test body",
	})
	if err != nil {
		t.Fatal(err)
	}

	line := strings.TrimSpace(buf.String())
	var got Message
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatal(err)
	}
	if got.To != "+15551234567" || got.Body != "test body" {
		t.Errorf("got = %+v", got)
	}
}

func TestLogSenderValidation(t *testing.T) {
	sender := NewLogSender(nil)
	if err := sender.Send(context.Background(), Message{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestLogSenderNilWriter(t *testing.T) {
	sender := NewLogSender(nil)
	if err := sender.Send(context.Background(), Message{To: "+1", Body: "x"}); err != nil {
		t.Fatalf("Send() = %v", err)
	}
}
