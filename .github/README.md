<p align="center">
<h1>MineQuery</h1>
Minecraft Server List Ping library written in Go.
<img src="https://img.shields.io/github/v/release/altea-minecraft/minequery">
<img src="https://img.shields.io/github/go-mod/go-version/altea-minecraft/minequery">
<img src="https://img.shields.io/github/workflow/status/altea-minecraft/minequery/Go/master">
<img src="https://img.shields.io/codacy/grade/7a7901a7d1ee435f8cd047ed15369043">
<img src="https://img.shields.io/github/license/altea-minecraft/minequery">
</p>

[//]: # (![GitHub release &#40;latest by date&#41;]&#40;https://img.shields.io/github/v/release/altea-minecraft/minequery&#41;)

[//]: # (![GitHub go.mod Go version]&#40;https://img.shields.io/github/go-mod/go-version/altea-minecraft/minequery&#41;)

[//]: # (![GitHub Workflow Status &#40;branch&#41;]&#40;https://img.shields.io/github/workflow/status/altea-minecraft/minequery/Go/master&#41;)

[//]: # (![Codacy grade]&#40;https://img.shields.io/codacy/grade/7a7901a7d1ee435f8cd047ed15369043&#41;)

[//]: # (![GitHub]&#40;https://img.shields.io/github/license/altea-minecraft/minequery&#41;)

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