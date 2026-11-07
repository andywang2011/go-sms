package sms

import (
	"errors"
	"strings"
)

// Message is the outbound SMS payload shared by all providers.
type Message struct {
	// To is the destination phone number in E.164 or provider-specific format.
	To string
	// Body is plain text content (used by Twilio and LogSender).
	Body string
	// TemplateCode is the provider template identifier (required for Aliyun).
	TemplateCode string
	// TemplateParams are template variable substitutions (Aliyun TemplateParam JSON).
	TemplateParams map[string]string
}

// validateRecipient ensures To is non-empty.
func validateRecipient(to string) error {
	if strings.TrimSpace(to) == "" {
		return errors.New("sms recipient (to) is not set")
	}
	return nil
}
