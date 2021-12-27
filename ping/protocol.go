package ping

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/unicode"
)

type packetType unsignedVarInt32

type packet struct {
	id  packetType
	buf *bytes.Buffer
}

func newPacket(p packetType) *packet {
	return &packet{p, bytes.NewBuffer(nil)}
}

func (p *packet) WriteSignedVarInt(s signedVarInt32) {
	_ = writeSignedVarInt(p.buf, s)
}

func (p *packet) WriteUnsignedVarInt(u unsignedVarInt32) {
	_ = writeUnsignedVarInt(p.buf, u)
}

func (p *packet) WriteLong(l long) {
	_ = writeLong(p.buf, l)
}

func (p *packet) WriteUnsignedShort(u unsignedShort) {
	_ = writeUnsignedShort(p.buf, u)
}

func (p *packet) WriteString(s string) {
	_ = writeString(p.buf, s)
}

func (p *packet) Push(w io.Writer) error {
	buf := bytes.NewBuffer(nil)

	headerBuf := bytes.NewBuffer(nil)
	_ = writeUnsignedVarInt(headerBuf, unsignedVarInt32(p.id))
	_ = writeUnsignedVarInt(buf, unsignedVarInt32(headerBuf.Len()+p.buf.Len()))

	_, err := headerBuf.WriteTo(buf) // Writing packet header

	_, err = p.buf.WriteTo(buf) // Writing packet data

	_, err = buf.WriteTo(w)
	return err
}

// Handshake

const packetHandshake packetType = 0x0

type handshake struct {
	Host string
	Port unsignedShort
}

const (
	handshakeProtocolVersionUndefined signedVarInt32   = -1
	handshakeNextStateStatus          unsignedVarInt32 = 1
)

func writeHandshake(w io.Writer, h handshake) error {
	p := newPacket(packetHandshake)
	p.WriteSignedVarInt(handshakeProtocolVersionUndefined)
	p.WriteString(h.Host)
	p.WriteUnsignedShort(h.Port)
	p.WriteUnsignedVarInt(handshakeNextStateStatus)
	return p.Push(w)
}

// Request

const packetRequest packetType = 0x0

func writeRequest(w io.Writer) error {
	return newPacket(packetRequest).Push(w)
}

// Response

const packetResponse packetType = 0x0

// Chat represents arbitrary JSON-encoded chat components structure used in modern (1.7 and earlier)
// Minecraft server descriptions.
type Chat interface{}

// Response represents ping response from modern (1.7 and earlier) Minecraft servers.
type Response struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`

	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	}

	Description Chat `json:"description"`

	Favicon string `json:"favicon"`
}

func readResponse(r io.Reader) (*Response, error) {
	l, err := readUnsignedVarInt(r)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	if _, err = io.CopyN(buf, r, int64(l)); err != nil {
		return nil, err
	}

	p, err := readUnsignedVarInt(buf)
	if err != nil {
		return nil, err
	}
	if packetType(p) != packetResponse {
		return nil, fmt.Errorf("expected packet %#x but got %#x instead", packetResponse, p)
	}

	d, err := readString(buf)
	if err != nil {
		return nil, err
	}

	data := &Response{}
	if err = json.Unmarshal([]byte(d), data); err != nil {
		return nil, err
	}

	return data, nil
}

// Legacy (<1.6)

type legacyPing struct {
	Host string
	Port uint16
}

func writeLegacyPing(w io.Writer, l legacyPing) error {
	// https://wiki.vg/Server_List_Ping#1.6

	if _, err := w.Write([]byte{
		0xfe,       // Packet ID
		0x01,       // Ping payload
		0xfa,       // Packet identifier for plugin message
		0x00, 0x0b, // Length of MC|PingHost string (11)
		0x00, 0x4d, 0x00, 0x43, 0x00, 0x7c, 0x00, 0x50, 0x00, 0x69, 0x00, // MC|PingHost string as UTF-16BE
		0x6e, 0x00, 0x67, 0x00, 0x48, 0x00, 0x6f, 0x00, 0x73, 0x00, 0x74,
	}); err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)

	hostnameBuf := bytes.NewBuffer(nil) // Buffer for hostname encoded as UTF-16BE
	if _, err := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder().Writer(hostnameBuf).Write([]byte(l.Host)); err != nil {
		return err
	}

	_ = binary.Write(buf, binary.BigEndian, uint16(hostnameBuf.Len()+7)) // Hostname as UTF-16BE length + 7 as short

	_, _ = buf.Write([]byte{0x4a}) // Latest protocol version (74)

	_ = binary.Write(buf, binary.BigEndian, uint16(len(l.Host))) // Length of hostname in characters as short

	_, _ = hostnameBuf.WriteTo(buf) // Hostname string

	_ = binary.Write(buf, binary.BigEndian, uint32(l.Port)) // Port as int

	_, err := buf.WriteTo(w)
	return err
}

// LegacyResponse represents ping response from legacy (1.4 to 1.6) Minecraft servers.
type LegacyResponse struct {
	ProtocolVersion uint32
	Version         string
	MessageOfTheDay string
	PlayerCount     uint32
	MaxPlayers      uint32
}

func readLegacyPong(r io.Reader) (*LegacyResponse, error) {
	buf := make([]byte, 3)
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}

	buf = make([]byte, 6)
	_, err = r.Read(buf)
	if err != nil {
		return nil, err
	}

	buf, err = io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	f := bytes.Split(buf, []byte{0x00, 0x00})
	rs := &LegacyResponse{}

	rs.ProtocolVersion, err = readLegacyPongUnsignedInt(f[0])
	if err != nil {
		return nil, err
	}

	rs.Version, err = readLegacyPongString(f[1])
	if err != nil {
		return nil, err
	}

	rs.MessageOfTheDay, err = readLegacyPongString(f[2])
	if err != nil {
		return nil, err
	}

	rs.PlayerCount, err = readLegacyPongUnsignedInt(f[3])
	if err != nil {
		return nil, err
	}

	rs.MaxPlayers, err = readLegacyPongUnsignedInt(f[4])
	if err != nil {
		return nil, err
	}

	return rs, nil
}

// Ancient (Beta 1.8 to 1.3)

type ancientPacketType byte

const packetAncientPing ancientPacketType = 0xfe
const packetAncientPong ancientPacketType = 0xff

// AncientResponse represents ping response from old servers (Beta 1.8 to 1.3) Minecraft servers.
type AncientResponse struct {
	MessageOfTheDay string
	PlayerCount     uint32
	MaxPlayers      uint32
}

func writeAncientPing(w io.Writer) error {
	_, err := w.Write([]byte{byte(packetAncientPing)})
	return err
}

func readAncientPong(r io.Reader) (*AncientResponse, error) {
	p, err := (&byteReaderWrap{r}).ReadByte()
	if err != nil {
		return nil, err
	}
	if ancientPacketType(p) != packetAncientPong {
		return nil, fmt.Errorf("expected packet %#x but got %#x instead", packetAncientPong, p)
	}

	l, err := readUnsignedShort(r)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, l*2)
	if _, err = r.Read(buf); err != nil {
		return nil, err
	}

	data, err := readLegacyPongString(buf)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(data, "§")
	a := &AncientResponse{}

	a.MessageOfTheDay = parts[0]
	c, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return nil, err
	}

	a.PlayerCount = uint32(c)
	m, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return nil, err
	}

	a.MaxPlayers = uint32(m)

	return a, nil
}
