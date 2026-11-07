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

func TestNewAliyunInvalidConfig(t *testing.T) {
	_, err := NewAliyun(AliyunConfig{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAliyunSendSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s", r.Method)
		}
		q := r.URL.Query()
		if q.Get("Action") != "SendSms" {
			t.Errorf("Action = %q", q.Get("Action"))
		}
		if q.Get("PhoneNumbers") != "13800138000" {
			t.Errorf("PhoneNumbers = %q", q.Get("PhoneNumbers"))
		}
		if q.Get("TemplateCode") != "SMS_123" {
			t.Errorf("TemplateCode = %q", q.Get("TemplateCode"))
		}
		if q.Get("Signature") == "" {
			t.Error("missing Signature")
		}
		var params map[string]string
		if err := json.Unmarshal([]byte(q.Get("TemplateParam")), &params); err != nil {
			t.Fatalf("TemplateParam: %v", err)
		}
		if params["code"] != "42" {
			t.Errorf("params = %v", params)
		}
		_, _ = io.WriteString(w, `{"Code":"OK","Message":"OK"}`)
	}))
	t.Cleanup(srv.Close)

	sender, err := NewAliyun(AliyunConfig{
		AccessKeyID:     "key-id",
		AccessKeySecret: "key-secret",
		SignName:        "MySign",
		Endpoint:        srv.URL,
		HTTPClient:      srv.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = sender.Send(context.Background(), Message{
		To:             "13800138000",
		TemplateCode:   "SMS_123",
		TemplateParams: map[string]string{"code": "42"},
	})
	if err != nil {
		t.Fatalf("Send() = %v", err)
	}
}

func TestAliyunSendBusinessError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"Code":"isv.BUSINESS_LIMIT_CONTROL","Message":"rate limit"}`)
	}))
	t.Cleanup(srv.Close)

	sender, err := NewAliyun(AliyunConfig{
		AccessKeyID:     "id",
		AccessKeySecret: "secret",
		SignName:        "sign",
		Endpoint:        srv.URL,
		HTTPClient:      srv.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = sender.Send(context.Background(), Message{To: "1", TemplateCode: "T"})
	if err == nil || !strings.Contains(err.Error(), "rate limit") {
		t.Fatalf("Send() = %v", err)
	}
}

func TestAliyunSendValidation(t *testing.T) {
	sender, err := NewAliyun(AliyunConfig{AccessKeyID: "a", AccessKeySecret: "b", SignName: "c"})
	if err != nil {
		t.Fatal(err)
	}
	if err := sender.Send(context.Background(), Message{}); err == nil {
		t.Fatal("expected recipient error")
	}
	if err := sender.Send(context.Background(), Message{To: "1"}); err == nil {
		t.Fatal("expected template error")
	}
}

func TestSignAliyunParams(t *testing.T) {
	sig, err := signAliyunParams(map[string]string{
		"AccessKeyId": "testid",
		"Action":      "SendSms",
		"Format":      "JSON",
	}, "secret")
	if err != nil || sig == "" {
		t.Fatalf("signAliyunParams() = %q, %v", sig, err)
	}
}
