# grdm

A Golang Concurrent Download Manager

[![GoDoc](https://godoc.org/github.com/lvisei/grdm?status.svg)](https://pkg.go.dev/github.com/lvisei/grdm)
[![Go Report Card](https://goreportcard.com/badge/github.com/lvisei/grdm)](https://goreportcard.com/report/github.com/lvisei/grdm)

grdm is a tool to concurrent download file manager.

## Command Line

### How to get

#### Download

You can download from GitHub [releases](https://github.com/lvisei/grdm/releases).

For example download file:

- windows: `**_windows_x86_64.zip`
- maxOS x86: `**_darwin_x86_64.tar.gz`
- maxOS M1: `**_darwin_arm64.tar.gz`

#### Build from source

```
git clone https://github.com/lvisei/grdm
cd cmd/grdm && go install
```

### Usage

```bash
grdm -u https://go.dev/dl/go1.22.1.src.tar.gz -n 2
```

Options flags:

```
grdm

Usage: grdm -u <URL> [-f] [-d] [-n]

Options:
  -d string
        Save file Directory (default ".")
  -f string
        Save file name with extension, defalut get from URL
  -n int
        Number of connections to make to the server (default 2)
  -u string
        Download URL
```

## Library

### How to get

```bash
go get github.com/lvisei/grdm
```

### Usage

```go
package main

import (
  "fmt"
  "github.com/lvisei/grdm"
)

func main() {
  d := grdm.Download{
    Url:      "https://go.dev/dl/go1.22.1.src.tar.gz",
    FileDir:  "tmp",
    Sections: 2,
  }

  if _, err := d.Do(); err != nil {
    fmt.Println(err)
  }
}
```

## LICENSE

[BSD 3-Clause](./LICENSE)
