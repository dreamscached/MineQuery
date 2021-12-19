<p align="center">
<h1>MineQuery</h1>
<h3>Minecraft Server List Ping library written in Go.</h3>
<br>
<a href="https://github.com/altea-minecraft/minequery/releases/latest">
<img src="https://img.shields.io/github/v/release/altea-minecraft/minequery">
</a>
<a href="https://github.com/altea-minecraft/minequery/blob/master/go.mod">
<img src="https://img.shields.io/github/go-mod/go-version/altea-minecraft/minequery">
</a>
<a href="https://github.com/altea-minecraft/minequery/actions/workflows/go.yml">
<img src="https://img.shields.io/github/workflow/status/altea-minecraft/minequery/Go/master">
</a>
<a href="https://github.com/altea-minecraft/minequery">
<img src="https://img.shields.io/codacy/grade/7a7901a7d1ee435f8cd047ed15369043">
</a>
<a href="https://github.com/altea-minecraft/minequery/blob/master/LICENSE">
<img src="https://img.shields.io/github/license/altea-minecraft/minequery">
</a>
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

	"github.com/altea-minecraft/minequery/ping"
)

func main() {
	res, err := ping.Ping("altea.land", 25565)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Description.Text)
}
```