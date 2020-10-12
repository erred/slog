# slog

A repo for slog

[![pkg.go.dev](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/go.seankhliao.com/slog)
![Version](https://img.shields.io/github/v/tag/seankhliao/slog?sort=semver&style=flat-square)
[![License](https://img.shields.io/github/license/seankhliao/slog.svg?style=flat-square)](LICENSE)

```go
package main

import (
        "errors"
        "os"
        "go.seankhliao.com/slog"
)

func main() {
        l := slog.NewText(os.Stderr)
        l.Info("hello", "foo", "bar")
        l.Error(errors.New("an error"), "oops", "hello", "world")

        l = slog.NewJSON(os.Stderr)
        l.Info("hello", "foo", "bar")
        l.Error(errors.New("an error"), "oops", "hello", "world")

        http.Server{
                Errorlog: slog.StdLogger(l),
        }
}
```

output:

```txt
2020-10-12T21:46:13+02:00 INF msg="hello" foo="bar"
2020-10-12T21:46:13+02:00 ERR msg="oops" err="an error" hello="world"
{"time":"2020-10-12T21:46:13+02:00", "level":"INF", "msg":"hello", "foo":"bar"}
{"time":"2020-10-12T21:46:13+02:00", "level":"ERR", "msg":"oops", "err":"an error", "hello":"world"}
```
