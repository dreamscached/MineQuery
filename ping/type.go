package ping

import (
	"encoding/binary"
	"io"
	"strconv"

	"golang.org/x/text/encoding/unicode"
)

type byteReaderWrap struct {
	reader io.Reader
}

func (w *byteReaderWrap) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	_, err := w.reader.Read(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], err
}

// Modern (>=1.7)

type unsignedVarInt32 uint32

func readUnsignedVarInt(r io.Reader) (unsignedVarInt32, error) {
	v, err := binary.ReadUvarint(&byteReaderWrap{r})
	if err != nil {
		return 0, err
	}
	return unsignedVarInt32(v), nil
}

func writeUnsignedVarInt(w io.Writer, u unsignedVarInt32) error {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(u))
	_, err := w.Write(buf[:n])
	return err
}

type signedVarInt32 int32

func writeSignedVarInt(w io.Writer, s signedVarInt32) error {
	buf := make([]byte, binary.MaxVarintLen32)
	n := binary.PutVarint(buf, int64(s))
	_, err := w.Write(buf[:n])
	return err
}

func writeString(w io.Writer, s string) error {
	if err := writeUnsignedVarInt(w, unsignedVarInt32(len(s))); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}

func readString(r io.Reader) (string, error) {
	l, err := readUnsignedVarInt(r)
	if err != nil {
		return "", err
	}
	buf := make([]byte, l)
	n, err := r.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

type unsignedShort uint16

func writeUnsignedShort(w io.Writer, u unsignedShort) error {
	return binary.Write(w, binary.BigEndian, uint16(u))
}

type long int64

func writeLong(w io.Writer, l long) error {
	return binary.Write(w, binary.BigEndian, l)
}

// Legacy (1.6)

func readLegacyPongString(b []byte) (string, error) {
	v, err := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder().Bytes(b)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func readLegacyPongUnsignedInt(b []byte) (uint32, error) {
	s, err := readLegacyPongString(b)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}
