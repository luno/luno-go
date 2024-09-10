[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=luno_luno-go&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=luno_luno-go)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=luno_luno-go&metric=coverage)](https://sonarcloud.io/summary/new_code?id=luno_luno-go)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=luno_luno-go&metric=bugs)](https://sonarcloud.io/summary/new_code?id=luno_luno-go)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=luno_luno-go&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=luno_luno-go)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=luno_luno-go&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=luno_luno-go)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=luno_luno-go&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=luno_luno-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/luno/luno-go)](https://goreportcard.com/report/github.com/luno/luno-go)
[![GoDoc](https://godoc.org/github.com/luno/luno-go?status.png)](https://godoc.org/github.com/luno/luno-go)

# Luno Go SDK

This Go package provides a wrapper for the [Luno API](https://www.luno.com/api).

## Documentation

Please visit [godoc.org](https://godoc.org/github.com/luno/luno-go) for the full
package documentation.

## Authentication

Please visit the [Settings](https://www.luno.com/wallet/settings/api_keys) page
to generate an API key.

## Installation

```shell
go get github.com/luno/luno-go
```

### Example usage

A full working example of this library in action.

```go
package main

import (
	"log"
	"context"
	"time"
	"github.com/luno/luno-go"
)

func main() {
	lunoClient := luno.NewClient()
	lunoClient.SetAuth("<id>", "<secret>")

	req := luno.GetOrderBookRequest{Pair: "XBTZAR"}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	res, err := lunoClient.GetOrderBook(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
}
```

Remember to substitute `<id>` and `<secret>` for your own Id and Secret.

We recommend using environment variables rather than including your credentials in plaintext. In Bash you do so as follows:
```shell
export LUNO_API_ID="<id>"
export LUNO_API_SECRET="<secret>"
```

And then access them in Go like so:

```go
import "os"

var API_KEY_ID string = os.Getenv("LUNO_API_ID")
var API_KEY_SECRET string = os.Getenv("LUNO_API_SECRET")
```

## License

[MIT](./LICENSE.md)
