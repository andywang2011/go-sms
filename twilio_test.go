package sms

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewTwilioInvalidConfig(t *testing.T) {
	_, err := NewTwilio(TwilioConfig{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTwilioSendSuccess(t *testing.T) {
	const wantSID = "SM123"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "ACtest" || pass != "secret" {
			t.Errorf("basic auth = %q %q ok=%v", user, pass, ok)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/Accounts/ACtest/Messages.json") {
			t.Errorf("path = %s", r.URL.Path)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("To") != "+15559876543" || r.Form.Get("From") != "+15551111111" || r.Form.Get("Body") != "hello" {
			t.Errorf("form = %v", r.Form)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"sid": wantSID, "status": "queued"})
	}))
	t.Cleanup(srv.Close)

	sender, err := NewTwilio(TwilioConfig{
		AccountSID: "ACtest",
		AuthToken:  "secret",
		From:       "+15551111111",
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = sender.Send(context.Background(), Message{To: "+15559876543", Body: "hello"})
	if err != nil {
		t.Fatalf("Send() = %v", err)
	}
}

func TestTwilioSendAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"message":"bad"}`)
	}))
	t.Cleanup(srv.Close)

	sender, err := NewTwilio(TwilioConfig{
		AccountSID: "ACtest",
		AuthToken:  "secret",
		From:       "+1",
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = sender.Send(context.Background(), Message{To: "+1", Body: "hi"})
	if err == nil || !strings.Contains(err.Error(), "400") {
		t.Fatalf("Send() = %v", err)
	}
}

func TestTwilioSendValidation(t *testing.T) {
	sender, err := NewTwilio(TwilioConfig{AccountSID: "a", AuthToken: "b", From: "c"})
	if err != nil {
		t.Fatal(err)
	}
	if err := sender.Send(context.Background(), Message{}); err == nil {
		t.Fatal("expected recipient error")
	}
	if err := sender.Send(context.Background(), Message{To: "+1"}); err == nil {
		t.Fatal("expected body error")
	}
}
