package minequery

import (
	"fmt"
	"testing"
	"time"
)

func TestQuerier_Query(t *testing.T) {
	p := NewPinger(WithTimeout(1000 * time.Second))

	res, err := p.QueryFull("157.90.163.70", 25565)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
