<h1 align="center">📡 MineQuery</h1>
<h4 align="center">Minecraft Server List Ping library written in Go</h4>
<p align="center">
    <a href="https://github.com/alteamc/minequery/blob/master/go.mod">
        <img alt="Go version badge" src="https://img.shields.io/github/go-mod/go-version/alteamc/minequery">
    </a>
    <a href="https://github.com/alteamc/minequery/releases/latest">
        <img alt="Latest release badge" src="https://img.shields.io/github/v/release/alteamc/minequery">
    </a>
    <a href="https://pkg.go.dev/github.com/alteamc/minequery">
        <img alt="Go reference badge" src="https://pkg.go.dev/badge/github.com/alteamc/minequery.svg">
    </a>
    <a href="https://github.com/alteamc/minequery/blob/master/LICENSE">
        <img alt="License badge" src="https://img.shields.io/github/license/alteamc/minequery">
    </a>
    <a href="https://discord.gg/9ruheUG3Wg">
        <img alt="Discord server badge" src="https://discordapp.com/api/guilds/929337829610369095/widget.png?style=shield">
    </a>
    <br/>
    <a href="https://github.com/alteamc/minequery#readme">
        <img alt="Minecraft version support badge" src="https://img.shields.io/badge/minecraft%20version-Beta%201.3%20to%201.3%20%7C%201.4%20%7C%201.5%20to%201.6%20%7C%201.7%2B-brightgreen">
    </a>
</p>

# #️⃣ Minecraft Version Support

As of version 2.0.0, MineQuery supports pinging of all versions of Minecraft.

| [Beta 1.8 to 1.3] | [1.4]       | [1.6 to 1.7] | [1.7+]      |
|-------------------|-------------|--------------|-------------|
| ✅ Supported       | ✅ Supported | ✅ Supported  | ✅ Supported |

[Beta 1.8 to 1.3]: https://wiki.vg/Server_List_Ping#Beta_1.8_to_1.3

[1.4]: https://wiki.vg/Server_List_Ping#1.4_to_1.5

[1.6 to 1.7]: https://wiki.vg/Server_List_Ping#1.6

[1.7+]: https://wiki.vg/Server_List_Ping#Current

## Query Protocol Support

As of version 2.0.0, query protocol is not yet supported.
See [issue #25] to track progress.

[issue #25]: https://github.com/alteamc/minequery/issues/25


# 📚 How to use

## Basic usage

For simple pinging with default parameters, use package-global `Ping*` functions 
(where `*` is your respective Minecraft server version.)

If you're unsure about version, it is known that Notchian servers respond to
all previous version pings (e.g. 1.7+ server will respond to 1.6 ping, and so on.)

Here's a quick example how to:

```go
package main

import (
	"fmt"

	"github.com/alteamc/minequery/minequery"
)

func main() {
	res, err := minequery.Ping17("localhost", 25565)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
```

For full info on response object structure, see [documentation].

[documentation]: https://pkg.go.dev/github.com/alteamc/minequery