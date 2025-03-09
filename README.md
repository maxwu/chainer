# Chainer

[![codecov](https://codecov.io/gh/maxwu/chainer/graph/badge.svg?token=JG5TC3BJIR)](https://codecov.io/gh/maxwu/chainer)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxwu/chainer)](https://goreportcard.com/report/github.com/maxwu/chainer)

Assemble HTTP middleware functions to a http.Handler which processes the request in chained middleware order.

## Usage

```go
package main
import (
	"net/http"
	"github.com/maxwu/chainer"
)

func main() {
	http.Handle("/", chainer.Chain(
		chainer.MiddlewareFunc(logger),
		chainer.MiddlewareFunc(authenticator),
		chainer.MiddlewareFunc(authorizer),
		chainer.MiddlewareFunc(handler),
	))
	http.ListenAndServe(":8080", nil)
}
```

## References

This project also demonstrates how to use [gotest-labels](github.com/maxwu/gotest-labels) package to select tests by the tags in test case comment block.
