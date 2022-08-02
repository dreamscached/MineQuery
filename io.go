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

func writeUShort(writer io.Writer, value uint16) error {
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
