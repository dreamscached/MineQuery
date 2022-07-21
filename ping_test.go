package minequery

import (
	"errors"
	"os"
	"strconv"
)

func Hostname() string {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		panic(errors.New("expected to get hostname from HOST variable"))
	}
	return host
}

func Port() int {
	portStr, ok := os.LookupEnv("PORT")
	if !ok {
		panic(errors.New("expected to get port from PORT variable"))
	}
	port, err := strconv.ParseInt(portStr, 10, 16)
	if err != nil {
		panic(err)
	}
	return int(port)
}