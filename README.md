pennypost-go — official Go SDK for [PennyPost](https://pennypost.io).

```go
import "github.com/pennypost-io/pennypost-go"

pp := pennypost.New("pp_live_...")
pp.SendEmail(&pennypost.SendEmailRequest{From: "Receipts <r@yourdomain.com>", To: []string{"a@b.com"}, Subject: "Hi", Text: "Hello"}, "")
```
