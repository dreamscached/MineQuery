<h1 align="center">MineQuery</h1>
<h3 align="center">Minecraft Server List Ping library written in Go.</h3>
<p align="center">
  <a href="https://github.com/alteamc/minequery/actions/workflows/go.yml"><img src="https://img.shields.io/github/workflow/status/alteamc/minequery/Go/master"></a>
  <a href="https://github.com/alteamc/minequery"><img src="https://img.shields.io/codacy/grade/7a7901a7d1ee435f8cd047ed15369043"></a>
  <a href="https://github.com/alteamc/minequery/blob/master/go.mod"><img src="https://img.shields.io/github/go-mod/go-version/alteamc/minequery"></a>
  <a href="https://github.com/alteamc/minequery/releases/latest"><img src="https://img.shields.io/github/v/release/alteamc/minequery"></a>
  <a href="https://pkg.go.dev/github.com/alteamc/minequery"><img src="https://pkg.go.dev/badge/github.com/alteamc/minequery.svg" alt="Go Reference"></a>
  <a href="https://github.com/alteamc/minequery/blob/master/LICENSE"><img src="https://img.shields.io/github/license/alteamc/minequery"></a>
</p>

## Features

* Modern Minecraft support (1.7 and newer)
* Legacy protocol support (1.6 and older)
* Older versions support (Beta 1.8 to Release 1.3)

## Usage

```go
package main

import (
	"fmt"

	"github.com/alteamc/minequery/ping"
)

func main() {
	res, err := ping.Ping("altea.land", 25565)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Description)
}
```
