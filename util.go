package minequery

import (
	"errors"
	"image"
	"io"
)

var errStackEmpty = errors.New("stack is empty")

type stack []interface{}

func (s *stack) Push(value interface{}) { *s = append(*s, value) }

func (s *stack) Pop() (interface{}, error) {
	if len(*s) == 0 {
		return nil, errStackEmpty
	}
	ret := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return ret, nil
}

// maxInt returns greatest of two of the passed integer numbers.
func maxInt(a int, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

// UnmarshalFunc is a function that conforms to json.Unmarshal function signature.
type UnmarshalFunc func([]byte, interface{}) error

// ImageDecodeFunc is a function that conforms to png.Decode function signature.
type ImageDecodeFunc func(io.Reader) (image.Image, error)
