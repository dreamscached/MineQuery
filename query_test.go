package minequery

import (
	"fmt"
	"testing"
)

func TestQuerier_Query(t *testing.T) {
	var q Querier
	res, err := q.Query("157.90.163.70", 25565)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
