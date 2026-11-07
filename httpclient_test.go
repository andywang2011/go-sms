package sms

import (
	"net/http"
	"testing"
	"time"
)

func TestDefaultHTTPClient_HasTimeout(t *testing.T) {
	c := defaultHTTPClient()
	if c == http.DefaultClient {
		t.Fatal("must not return http.DefaultClient")
	}
	if c.Timeout != 10*time.Second {
		t.Fatalf("Timeout = %v, want 10s", c.Timeout)
	}
	if defaultHTTPClient() == c {
		t.Fatal("each call must return a fresh client")
	}
}

func TestNew_DefaultHTTPClientHasTimeout(t *testing.T) {
	tw, err := NewTwilio(TwilioConfig{AccountSID: "AC", AuthToken: "s", From: "+15551111111"})
	if err != nil {
		t.Fatal(err)
	}
	assertTimeout(t, "twilio", tw.(*twilioSender).client)

	al, err := NewAliyun(AliyunConfig{AccessKeyID: "id", AccessKeySecret: "s", SignName: "sign"})
	if err != nil {
		t.Fatal(err)
	}
	assertTimeout(t, "aliyun", al.(*aliyunSender).client)
}

func TestNew_CustomHTTPClientPreserved(t *testing.T) {
	custom := &http.Client{Timeout: time.Minute}
	tw, err := NewTwilio(TwilioConfig{AccountSID: "AC", AuthToken: "s", From: "+15551111111", HTTPClient: custom})
	if err != nil {
		t.Fatal(err)
	}
	if tw.(*twilioSender).client != HTTPDoer(custom) {
		t.Fatal("custom HTTPClient must be used as is")
	}
}

func assertTimeout(t *testing.T, name string, d HTTPDoer) {
	t.Helper()
	c, ok := d.(*http.Client)
	if !ok || c == http.DefaultClient || c.Timeout != 10*time.Second {
		t.Fatalf("%s default client = %#v, want *http.Client with 10s timeout", name, d)
	}
}
