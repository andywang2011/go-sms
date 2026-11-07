package sms_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/andywang2011/go-sms"
)

func ExampleNewLogSender() {
	buf := &bytes.Buffer{}
	dev := sms.NewLogSender(buf)
	err := dev.Send(context.Background(), sms.Message{
		To:   "+15551234567",
		Body: "hello",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.Len() > 0)
	// Output: true
}
