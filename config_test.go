package sms

import "testing"

func TestTwilioConfigValidate(t *testing.T) {
	valid := TwilioConfig{AccountSID: "ACxxx", AuthToken: "token", From: "+15551234567"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("Validate() = %v", err)
	}

	cases := []struct {
		name string
		cfg  TwilioConfig
	}{
		{"missing sid", TwilioConfig{AuthToken: "t", From: "f"}},
		{"missing token", TwilioConfig{AccountSID: "s", From: "f"}},
		{"missing from", TwilioConfig{AccountSID: "s", AuthToken: "t"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.cfg.Validate(); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestAliyunConfigValidate(t *testing.T) {
	valid := AliyunConfig{AccessKeyID: "id", AccessKeySecret: "secret", SignName: "sign"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("Validate() = %v", err)
	}

	cases := []struct {
		name string
		cfg  AliyunConfig
	}{
		{"missing key id", AliyunConfig{AccessKeySecret: "s", SignName: "n"}},
		{"missing secret", AliyunConfig{AccessKeyID: "id", SignName: "n"}},
		{"missing sign", AliyunConfig{AccessKeyID: "id", AccessKeySecret: "s"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.cfg.Validate(); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestConfigDefaults(t *testing.T) {
	tw := TwilioConfig{AccountSID: "a", AuthToken: "b", From: "c"}
	if tw.resolvedBaseURL() != defaultTwilioBaseURL {
		t.Fatalf("base URL = %q", tw.resolvedBaseURL())
	}
	al := AliyunConfig{AccessKeyID: "a", AccessKeySecret: "b", SignName: "s"}
	if al.resolvedEndpoint() != defaultAliyunEndpoint {
		t.Fatalf("endpoint = %q", al.resolvedEndpoint())
	}
	if al.resolvedRegionID() != defaultAliyunRegionID {
		t.Fatalf("region = %q", al.resolvedRegionID())
	}
}
