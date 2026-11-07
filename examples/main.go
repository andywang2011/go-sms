package main

import (
	"context"
	"os"

	"github.com/andywang2011/go-sms"
)

func main() {
	dev := sms.NewLogSender(os.Stdout)
	_ = dev.Send(context.Background(), sms.Message{
		To:   "+15551234567",
		Body: "hello from go-sms examples",
	})
}
