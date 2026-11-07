package sms

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type aliyunSender struct {
	accessKeyID     string
	accessKeySecret string
	signName        string
	regionID        string
	endpoint        string
	client          HTTPDoer
}

// NewAliyun creates an Alibaba Cloud SMS sender from cfg.
func NewAliyun(cfg AliyunConfig) (Sender, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	client := cfg.HTTPClient
	if client == nil {
		client = defaultHTTPClient()
	}
	return &aliyunSender{
		accessKeyID:     cfg.AccessKeyID,
		accessKeySecret: cfg.AccessKeySecret,
		signName:        cfg.SignName,
		regionID:        cfg.resolvedRegionID(),
		endpoint:        cfg.resolvedEndpoint(),
		client:          client,
	}, nil
}

// Send calls the SendSms OpenAPI action.
func (s *aliyunSender) Send(ctx context.Context, msg Message) error {
	if err := validateRecipient(msg.To); err != nil {
		return err
	}
	if strings.TrimSpace(msg.TemplateCode) == "" {
		return errors.New("aliyun template code is not set")
	}

	params := map[string]string{
		"AccessKeyId":      s.accessKeyID,
		"Action":           "SendSms",
		"Format":           "JSON",
		"PhoneNumbers":     msg.To,
		"RegionId":         s.regionID,
		"SignName":         s.signName,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"SignatureVersion": "1.0",
		"TemplateCode":     msg.TemplateCode,
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Version":          aliyunSMSAPIVersion,
	}
	if len(msg.TemplateParams) > 0 {
		b, err := json.Marshal(msg.TemplateParams)
		if err != nil {
			return fmt.Errorf("aliyun template params: %w", err)
		}
		params["TemplateParam"] = string(b)
	}

	signature, err := signAliyunParams(params, s.accessKeySecret)
	if err != nil {
		return err
	}
	params["Signature"] = signature

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	reqURL := s.endpoint + "/?" + values.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("aliyun create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("aliyun request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("aliyun api error: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("aliyun decode response: %w", err)
	}
	if parsed.Code != "OK" {
		if parsed.Message == "" {
			parsed.Message = string(body)
		}
		return fmt.Errorf("aliyun send failed: %s", parsed.Message)
	}
	return nil
}

func signAliyunParams(params map[string]string, secret string) (string, error) {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "Signature" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var canonical strings.Builder
	for i, k := range keys {
		if i > 0 {
			canonical.WriteByte('&')
		}
		canonical.WriteString(percentEncode(k))
		canonical.WriteByte('=')
		canonical.WriteString(percentEncode(params[k]))
	}

	stringToSign := "GET&" + percentEncode("/") + "&" + percentEncode(canonical.String())
	mac := hmac.New(sha1.New, []byte(secret+"&"))
	if _, err := mac.Write([]byte(stringToSign)); err != nil {
		return "", fmt.Errorf("aliyun sign: %w", err)
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func percentEncode(s string) string {
	encoded := url.QueryEscape(s)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "*", "%2A")
	encoded = strings.ReplaceAll(encoded, "%7E", "~")
	return encoded
}
