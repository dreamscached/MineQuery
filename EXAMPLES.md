## Pinging modern Minecraft servers (1.7 and later)

NOTE: Modern servers *also* respond to older ping types.

```go
package main

import (
	"fmt"

	"github.com/dreamscached/minequery/ping"
)

func main() {
	res, err := ping.Ping("altea.land", 25565)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Description)
}
```

## Pinging legacy Minecraft servers (1.4 to 1.6)

```go
package main

import (
	"fmt"

	"github.com/dreamscached/minequery/ping"
)

func main() {
	res, err := ping.PingLegacy("altea.land", 25565)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.MessageOfTheDay)
}
```

## Pinging old Minecraft servers (Beta 1.7 to 1.3)

```go
package main

import (
	"fmt"

	"github.com/dreamscached/minequery/ping"
)

func main() {
	res, err := ping.PingAncient("altea.land", 25565)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.MessageOfTheDay)
}
```

## Pinging with timeout

All `Ping` methods have `WithTimeout` variants that let you pass a `time.Duration` value used for socket read/write
timeout.

```go
package main

import (
	"fmt"
	"time"

	"github.com/dreamscached/minequery/ping"
)

func main() {
	res, err := ping.PingWithTimeout("altea.land", 25565, 1*time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Description)
}
```
