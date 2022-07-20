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

func writeBytes(writer io.Writer, bytes []byte) error {
	_, err := writer.Write(bytes)
	return err
}

func writeBuffer(writer io.Writer, buffer *bytes.Buffer) error {
	_, err := buffer.WriteTo(writer)
	return err
}

func readNBytes(reader io.Reader, n int) ([]byte, error) {
	byteArr := make([]byte, n)
	if _, err := reader.Read(byteArr); err != nil {
		return nil, err
	}
	return byteArr, nil
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
