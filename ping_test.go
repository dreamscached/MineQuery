package minequery

import (
	"os"
	"strconv"
)

const (
	serverTypeVanilla     = "vanilla"
	serverTypeCraftBukkit = "craftbukkit"
	serverTypeSpigot      = "spigot"
)

const (
	defaultHostname = "localhost"
	defaultPort     = 25565
	defaultType     = serverTypeVanilla
)

func Hostname() string {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		return defaultHostname
	}
	return host
}

func Port() int {
	portStr, ok := os.LookupEnv("PORT")
	if !ok {
		return defaultPort
	}
	port, err := strconv.ParseInt(portStr, 10, 16)
	if err != nil {
		panic(err)
	}
	return int(port)
}

func Type() string {
	_type, ok := os.LookupEnv("TYPE")
	if !ok {
		return defaultType
	}
	return _type
}