<h1 align="center">MineQuery</h1>
<h3 align="center">Minecraft Server List Ping library written in Go.</h3>
<p align="center">
  <a href="https://github.com/alteamc/minequery/actions/workflows/go.yml"><img alt="Workflow status" src="https://img.shields.io/github/workflow/status/alteamc/minequery/Go/master"></a>
  <a href="https://app.codacy.com/gh/alteamc/minequery"><img alt="Codacy grade" src="https://img.shields.io/codacy/grade/7a7901a7d1ee435f8cd047ed15369043"></a>
  <a href="https://github.com/alteamc/minequery/blob/master/go.mod"><img alt="Go version" src="https://img.shields.io/github/go-mod/go-version/alteamc/minequery"></a>
  <a href="https://github.com/alteamc/minequery/releases/latest"><img alt="Latest release" src="https://img.shields.io/github/v/release/alteamc/minequery"></a>
  <a href="https://pkg.go.dev/github.com/alteamc/minequery"><img alt="Go Reference" src="https://pkg.go.dev/badge/github.com/alteamc/minequery.svg"></a>
  <a href="https://github.com/alteamc/minequery/blob/master/LICENSE"><img alt="License" src="https://img.shields.io/github/license/alteamc/minequery"></a>
</p>

## Features

*   Modern Minecraft support (1.7 and newer)
*   Legacy protocol support (1.6 and older)
*   Older versions support (Beta 1.8 to Release 1.3)

## Usage

### Pinging modern Minecraft servers (1.7 and later)

NOTE: Modern servers *also* respond to older ping types.

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

### Pinging legacy Minecraft servers (1.4 to 1.6)

```go
package main

import (
	"fmt"

	"github.com/alteamc/minequery/ping"
)

func main() {
	res, err := ping.PingLegacy("altea.land", 25565)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.MessageOfTheDay)
}
```

### Pinging old Minecraft servers (Beta 1.7 to 1.3)

```go
package main

import (
	"fmt"

	"github.com/alteamc/minequery/ping"
)

func main() {
	res, err := ping.PingAncient("altea.land", 25565)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.MessageOfTheDay)
}
```

### Pinging with timeout

All `Ping` methods have `WithTimeout` variants that let you pass a `time.Duration` value used for socket read/write 
timeout.

```go
package main

import (
	"fmt"
	"time"

	"github.com/alteamc/minequery/ping"
)

func main() {
	res, err := ping.PingWithTimeout("altea.land", 25565, 1*time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Description)
}
```