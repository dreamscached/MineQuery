package minequery

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/unicode"
)

var (
	utf16BEEncoder = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder()
	utf16BEDecoder = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
)

// readAllUntilZero reads all bytes from reader until it hits zero.
// This is a backport from newer Go stdlib for sake of minequery's compatibility with Go 1.13.
func readAllUntilZero(reader io.ByteReader) ([]byte, error) {
	buf := &bytes.Buffer{}

	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return buf.Bytes(), nil
			}
			return nil, err
		}

		if b != 0 {
			buf.WriteByte(b)
		} else {
			return buf.Bytes(), nil
		}
	}
}
