package sms

import (
	"errors"
	"strings"
)

const (
	defaultTwilioBaseURL  = "https://api.twilio.com"
	defaultAliyunEndpoint = "https://dysmsapi.aliyuncs.com"
	defaultAliyunRegionID = "cn-hangzhou"
	aliyunSMSAPIVersion   = "2017-05-25"
)

// TwilioConfig holds Twilio REST API credentials and transport options.
type TwilioConfig struct {
	// AccountSID is the Twilio account SID.
	AccountSID string
	// AuthToken is the Twilio auth token.
	AuthToken string
	// From is the sender phone number or messaging service SID.
	From string
	// BaseURL overrides the Twilio API base URL (for httptest).
	// Doc: https://www.twilio.com/docs/messaging/api/message-resource
	BaseURL string
	// HTTPClient is the HTTP client used for requests; nil uses a client with a 10s timeout.
	HTTPClient HTTPDoer
}

// Validate checks required Twilio fields.
func (c TwilioConfig) Validate() error {
	if strings.TrimSpace(c.AccountSID) == "" {
		return errors.New("twilio account sid is not set")
	}
	if strings.TrimSpace(c.AuthToken) == "" {
		return errors.New("twilio auth token is not set")
	}
	if strings.TrimSpace(c.From) == "" {
		return errors.New("twilio from is not set")
	}
	return nil
}

// resolvedBaseURL returns BaseURL or the production default.
func (c TwilioConfig) resolvedBaseURL() string {
	if strings.TrimSpace(c.BaseURL) != "" {
		return strings.TrimRight(c.BaseURL, "/")
	}
	return defaultTwilioBaseURL
}

// AliyunConfig holds Alibaba Cloud SMS (SendSms) credentials and transport options.
type AliyunConfig struct {
	// AccessKeyID is the Alibaba Cloud access key ID.
	AccessKeyID string
	// AccessKeySecret is the Alibaba Cloud access key secret.
	AccessKeySecret string
	// SignName is the approved SMS signature name.
	SignName string
	// RegionID is the SMS region; defaults to cn-hangzhou when empty.
	RegionID string
	// Endpoint overrides the SMS API endpoint (for httptest).
	// Doc: https://help.aliyun.com/document_detail/101414.html
	Endpoint string
	// HTTPClient is the HTTP client used for requests; nil uses a client with a 10s timeout.
	HTTPClient HTTPDoer
}

// Validate checks required Aliyun fields.
func (c AliyunConfig) Validate() error {
	if strings.TrimSpace(c.AccessKeyID) == "" {
		return errors.New("aliyun access key id is not set")
	}
	if strings.TrimSpace(c.AccessKeySecret) == "" {
		return errors.New("aliyun access key secret is not set")
	}
	if strings.TrimSpace(c.SignName) == "" {
		return errors.New("aliyun sign name is not set")
	}
	return nil
}

func (c AliyunConfig) resolvedRegionID() string {
	if strings.TrimSpace(c.RegionID) != "" {
		return c.RegionID
	}
	return defaultAliyunRegionID
}

func (c AliyunConfig) resolvedEndpoint() string {
	if strings.TrimSpace(c.Endpoint) != "" {
		return strings.TrimRight(c.Endpoint, "/")
	}
	return defaultAliyunEndpoint
}
