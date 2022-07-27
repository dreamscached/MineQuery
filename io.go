package minequery

import (
	"bytes"
	"encoding/binary"
	"io"

	"golang.org/x/text/encoding/unicode"
)

var (
	utf16BEEncoder = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder()
	utf16BEDecoder = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
)

type byteReader struct {
	reader io.Reader
}

func newByteReader(reader io.Reader) byteReader {
	return byteReader{reader}
}

func (r byteReader) ReadByte() (byte, error) {
	byteArr := make([]byte, 1)
	if _, err := r.reader.Read(byteArr); err != nil {
		return 0, err
	}
	return byteArr[0], nil
}

func readAllUntilZero(reader io.Reader) ([]byte, error) {
	br := newByteReader(reader)
	buf := &bytes.Buffer{}

	for {
		b, err := br.ReadByte()
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

func writeBytes(writer io.Writer, bytes []byte) error {
	_, err := writer.Write(bytes)
	return err
}

func writeBuffer(writer io.Writer, buffer *bytes.Buffer) error {
	_, err := buffer.WriteTo(writer)
	return err
}

func writeByte(writer io.Writer, c byte) error {
	_, err := writer.Write([]byte{c})
	return err
}

func readByte(reader io.Reader) (byte, error) {
	byteArr := make([]byte, 1)
	if _, err := reader.Read(byteArr); err != nil {
		return 0, err
	}
	return byteArr[0], nil
}

func writeUShort(writer io.Writer, value uint16) error {
	return binary.Write(writer, binary.BigEndian, value)
}

func readUShort(reader io.Reader) (uint16, error) {
	var value uint16
	if err := binary.Read(reader, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

func writeUInt(writer io.Writer, value uint32) error {
	return binary.Write(writer, binary.BigEndian, value)
}

func writeVarInt(writer io.Writer, value int32) error {
	intBytes := make([]byte, binary.MaxVarintLen32)
	written := binary.PutVarint(intBytes, int64(value))
	_, err := writer.Write(intBytes[:written])
	return err
}

func writeUVarInt(writer io.Writer, value uint32) error {
	intBytes := make([]byte, binary.MaxVarintLen32)
	written := binary.PutUvarint(intBytes, uint64(value))
	_, err := writer.Write(intBytes[:written])
	return err
}

func readUVarInt(reader io.Reader) (uint32, error) {
	value, err := binary.ReadUvarint(newByteReader(reader))
	if err != nil {
		return 0, err
	}
	return uint32(value), nil
}

// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// readAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because readAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.
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
