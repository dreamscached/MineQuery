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

type packet struct {
	id  unsignedVarInt32
	buf *bytes.Buffer
}

func newPacket(p unsignedVarInt32) *packet {
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
	_ = writeUnsignedVarInt(buf, p.id)
	if err := writeUnsignedVarInt(w, unsignedVarInt32(buf.Len()+p.buf.Len())); err != nil {
		return err
	}
	if _, err := buf.WriteTo(w); err != nil {
		return err
	} // Pushing packet ID
	_, err := p.buf.WriteTo(w) // Pushing packet data
	return err
}

// Handshake

const packetHandshake unsignedVarInt32 = 0x0

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

const packetRequest unsignedVarInt32 = 0x0

func writeRequest(w io.Writer) error {
	return newPacket(packetRequest).Push(w)
}

// Ping

const packetPing unsignedVarInt32 = 0x1

type ping struct {
	Payload long
}

func writePing(w io.Writer, p ping) error {
	pk := newPacket(packetPing)
	pk.WriteLong(p.Payload)
	return pk.Push(w)
}

// Response

const packetResponse unsignedVarInt32 = 0x0

type Description struct {
	Text string `json:"text"`
}

type descriptionObj struct {
	Text string `json:"text"`
}

func (r *Description) UnmarshalJSON(data []byte) error {
	f := &descriptionObj{}
	if err := json.Unmarshal(data, f); err != nil {
		r.Text = string(data)
		return nil
	}

	r.Text = f.Text
	return nil
}

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

	Description Description `json:"description"`

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
	if p != packetResponse {
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

	buf := bytes.NewBuffer(nil) // Buffer for hostname encoded as UTF-16BE
	if _, err := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder().Writer(buf).Write([]byte(l.Host)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, uint16(buf.Len()+7)); err != nil {
		return err
	} // Hostname as UTF-16BE length + 7 as short
	if _, err := w.Write([]byte{0x4a}); err != nil {
		return err
	} // Latest protocol version (74)
	if err := binary.Write(w, binary.BigEndian, uint16(len(l.Host))); err != nil {
		return err
	} // Length of hostname in characters as short
	if _, err := buf.WriteTo(w); err != nil {
		return err
	} // Hostname string
	if err := binary.Write(w, binary.BigEndian, uint32(l.Port)); err != nil {
		return err
	} // Port as int

	return nil
}

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

const packetAncientPing byte = 0xfe
const packetAncientPong byte = 0xff

type AncientResponse struct {
	MessageOfTheDay string
	PlayerCount     uint32
	MaxPlayers      uint32
}

func writeAncientPing(w io.Writer) error {
	_, err := w.Write([]byte{packetAncientPing})
	return err
}

func readAncientPong(r io.Reader) (*AncientResponse, error) {
	p, err := (&byteReaderWrap{r}).ReadByte()
	if err != nil {
		return nil, err
	}
	if p != packetAncientPong {
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
	m, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return nil, err
	}
	a.MaxPlayers = uint32(m)
	c, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return nil, err
	}
	a.PlayerCount = uint32(c)

	return a, nil
}
