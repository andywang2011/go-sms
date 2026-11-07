# go-sms

Multi-provider SMS sending for Go (Twilio, Aliyun, log sender).

## Install

```bash
go get github.com/andywang2011/go-sms@v0.1.0
```

## Quick start

```go
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
```

## Examples

```bash
cd examples
go run .
```

## License

MIT — see [LICENSE](LICENSE).
