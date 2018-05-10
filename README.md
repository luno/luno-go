<img src="https://www.luno.com/static/images/luno-email-336.png">

# Luno API [![GoDoc](https://godoc.org/github.com/luno/luno-go?status.png)](https://godoc.org/github.com/luno/luno-go)

This Go package provides a wrapper for the [Luno API](https://www.luno.com/api).

⚠️ *WARNING* This package is currently being tested, and should not be used in production.

### Documentation

Please visit [godoc.org](https://godoc.org/github.com/luno/luno-go) for the full
package documentation.

### Installation

```
go get github.com/luno/luno-go
```

### Authentication

Please visit the [Settings](https://www.luno.com/wallet/settings/api_keys) page
to generate an API key.

### Example usage

```go
import luno "github.com/luno/luno-go"

lunoClient := luno.NewClient()
lunoClient.SetAuth("api_key_id", "api_key_secret")

req := lunoClient.GetOrderBookRequest{
  Pair: "XBTZAR",
}
res, err := lunoClient.GetOrderBook(&req)
if err != nil {
  log.Fatal(err)
}
log.Println(res)
```

### License

[MIT](https://github.com/luno/luno-go/blob/master/LICENSE.md)
