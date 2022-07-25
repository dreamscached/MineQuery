<h1 align="center">📡 MineQuery</h1>
<h4 align="center">Minecraft Server List Ping library written in Go</h4>
<p align="center">
    <a href="https://github.com/alteamc/minequery/actions/workflows/test.yml">
        <img alt="Test workflow badge" src="https://img.shields.io/github/workflow/status/alteamc/minequery/Test/v2?label=Test&logo=github">
    </a>
    <a href="https://github.com/alteamc/minequery/blob/v2/go.mod">
        <img alt="Go version badge" src="https://img.shields.io/github/go-mod/go-version/alteamc/minequery">
    </a>
    <a href="https://github.com/alteamc/minequery/releases/latest">
        <img alt="Latest release badge" src="https://img.shields.io/github/v/release/alteamc/minequery">
    </a>
    <a href="https://pkg.go.dev/github.com/alteamc/minequery/v2">
        <img alt="Go reference badge" src="https://pkg.go.dev/badge/github.com/alteamc/minequery.svg">
    </a>
    <a href="https://github.com/alteamc/minequery/blob/v2/LICENSE">
        <img alt="License badge" src="https://img.shields.io/github/license/alteamc/minequery">
    </a>
    <br/>
    <a href="https://github.com/alteamc/minequery#readme">
        <img alt="Minecraft version support badge" src="https://img.shields.io/badge/minecraft%20version-Beta%201.8%20to%201.3%20%7C%201.4%20%7C%201.5%20to%201.6%20%7C%201.7%2B-brightgreen">
    </a>
    <a href="https://discord.gg/9ruheUG3Wg">
        <img alt="Discord server badge" src="https://discordapp.com/api/guilds/929337829610369095/widget.png?style=shield">
    </a>
</p>

# 🚀 Migrating from v1

If you're new to MineQuery, you can skip this part. If you have used it before, you
might want to give it a read if you're planning to switch from v1.

See [docs/MIGRATING.md] for help with migrating from MineQuery v1 to v2.

[docs/MIGRATING.md]: docs/MIGRATING.md


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
res, err := minequery.Ping17("localhost", 25565)
if err != nil { panic(err) }
fmt.Println(res)
```

For full info on response object structure, see [documentation].

[documentation]: https://pkg.go.dev/github.com/alteamc/minequery


## Advanced usage

### Pinger

For more advanced usage, such as setting custom timeout or enabling more strict
response validation, you can use `Pinger` struct with `PingerOption` passed to it:

```go
pinger := minequery.NewPinger(
	minequery.WithTimeout(5 * time.Second), 
	minequery.WithUseStrict(true),
	minequery.WithProtocolVersion16(minequery.Ping16ProtocolVersion162), 
	minequery.WithProtocolVersion17(minequery.Ping17ProtocolVersion172),
)
```

Then, use `Ping*` functions on it the same way as described in [Basic usage] section:

```go
// Ping Beta 1.8+
pinger.PingBeta18("localhost", 25565)
// Ping 1.4+
pinger.Ping14("localhost", 25565)
// Ping 1.6+
pinger.Ping16("localhost", 25565)
// Ping 1.7+
pinger.Ping17("localhost", 25565)
```

[Basic usage]: #basic-usage


### WithTimeout

By default, `Pinger` has 15-second timeout before connection aborts. If you need
to customize this duration, you can use `WithTimeout` option.


### WithUseStrict

By default, `Pinger` does not validate response data it receives and silently
omits erroneous values it processes (incorrect favicon or bad player UUID).
If you need it to return an error in case of invalid response, you can use 
`WithUseStrict` option.


### WithProtocolVersion16

By default, `Pinger` sends protocol version 74 in 1.6 ping packets. If you need
to customize protocol version sent, use `WithProtocolVersion16`. MineQuery provides
a convenient set of constants you can use &mdash; see `Ping16ProtocolVersion*` constants.


### WithProtocolVersion17

By default, `Pinger` sends protocol version -1 in 1.7 ping packets. If you need
to customize protocol version sent, use `WithProtocolVersion17`. MineQuery provides
a convenient set of constants you can use &mdash; see `Ping17ProtocolVersion*` constants.