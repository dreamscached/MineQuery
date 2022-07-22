package minequery

import (
	"io"
)

// readAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because readAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported. Ported from Go 1.16.
func readAll(r io.Reader) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
	}
}
