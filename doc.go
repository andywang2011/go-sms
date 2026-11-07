// Package sms provides multi-provider SMS sending with config injection.
//
// Design constraints:
//   - No viper, environment variables, or config file paths inside the library.
//   - No global singletons or log.Fatal; callers construct senders via New*.
//   - Twilio and Aliyun accept an optional *http.Client in Config for httptest.
//   - New* returns Sender (client_iface.go) for mocks and dependency injection.
//   - ConfigDefault: N/A — use TwilioConfig / AliyunConfig per provider.
package sms
