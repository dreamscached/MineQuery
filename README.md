# minequery

Minecraft Server List Ping library written in Go.

## Features

* Modern Minecraft support (1.7 and newer)
* Legacy protocol support (1.6 and older)

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