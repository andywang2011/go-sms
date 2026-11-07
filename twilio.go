package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type twilioSender struct {
	accountSID string
	authToken  string
	from       string
	baseURL    string
	client     HTTPDoer
}

// NewTwilio creates a Twilio SMS sender from cfg.
func NewTwilio(cfg TwilioConfig) (Sender, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	client := cfg.HTTPClient
	if client == nil {
		client = defaultHTTPClient()
	}
	return &twilioSender{
		accountSID: cfg.AccountSID,
		authToken:  cfg.AuthToken,
		from:       cfg.From,
		baseURL:    cfg.resolvedBaseURL(),
		client:     client,
	}, nil
}

// Send posts a message to the Twilio Messages resource.
func (s *twilioSender) Send(ctx context.Context, msg Message) error {
	if err := validateRecipient(msg.To); err != nil {
		return err
	}
	if strings.TrimSpace(msg.Body) == "" {
		return errors.New("twilio message body is not set")
	}

	form := url.Values{}
	form.Set("To", msg.To)
	form.Set("From", s.from)
	form.Set("Body", msg.Body)

	endpoint := fmt.Sprintf("%s/2010-04-01/Accounts/%s/Messages.json", s.baseURL, url.PathEscape(s.accountSID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("twilio create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.accountSID, s.authToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("twilio request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("twilio api error: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed struct {
		SID    string `json:"sid"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("twilio decode response: %w", err)
	}
	if parsed.SID == "" {
		return errors.New("twilio response missing message sid")
	}
	return nil
}
