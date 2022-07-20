package minequery

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/unicode"
)

var pingBeta18PingPacket = []byte{0xfe}

const (
	pingBeta18ResponsePacketID       byte = 0xff
	pingBeta18ResponseFieldSeparator      = "§"
)

// StatusBeta18 holds status response returned by Beta 1.8 to Release 1.4 (exclusively) Minecraft servers.
type StatusBeta18 struct {
	MOTD          string
	OnlinePlayers int
	MaxPlayers    int
}

// PingBeta18 pings Beta 1.8 to Release 1.4 (exclusively) Minecraft servers (Notchian servers of more late versions
// also respond to this ping packet.)
func PingBeta18(host string, port int) (StatusBeta18, error) {
	return defaultPinger.PingBeta18(host, port)
}

// PingBeta18 pings Beta 1.8 to Release 1.4 (exclusively) Minecraft servers (Notchian servers of more late versions
// also respond to this ping packet.)
func (p Pinger) PingBeta18(host string, port int) (StatusBeta18, error) {
	conn, err := p.openTCPConn(host, port)
	defer func() { _ = conn.Close() }()
	if err != nil {
		return StatusBeta18{}, err
	}

	// Write single-byte FE ping packet
	if err = writeBytes(conn, pingBeta18PingPacket); err != nil {
		return StatusBeta18{}, err
	}

	// Read packet type, return error if it isn't FF kick packet
	packetType, err := readByte(conn)
	if err != nil {
		return StatusBeta18{}, err
	} else if packetType != pingBeta18ResponsePacketID {
		return StatusBeta18{}, fmt.Errorf("expected packet ID %#x, but instead got %#x", pingBeta18ResponsePacketID, packetType)
	}

	// Read packet length, return error if it isn't readable as unsigned short
	// Worth noting that this needs to be multiplied by two further on (for encoding reasons, most probably)
	length, err := readUShort(conn)
	if err != nil {
		return StatusBeta18{}, err
	}

	// Read remainder of the status packet as raw bytes
	// This is a UTF-16BE string separated by § (paragraph sign)
	// where [0] is MOTD, [1] is online players and [2] is max players
	dataEncoded, err := readNBytes(conn, int(length*2))
	if err != nil {
		return StatusBeta18{}, err
	}

	// Decode UTF16-BE and return error if unable to decode
	utf16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
	dataString, err := utf16be.Bytes(dataEncoded)
	if err != nil {
		return StatusBeta18{}, err
	}

	// Split status string, parse and map to struct returning errors if conversions fail
	fields := strings.Split(string(dataString), pingBeta18ResponseFieldSeparator)
	if len(fields) != 3 {
		return StatusBeta18{}, fmt.Errorf("%w: expected 3 status fields, got %d", ErrInvalidStatus, len(fields))
	}
	motd, onlineString, maxString := fields[0], fields[1], fields[2]

	// Check MOTD length
	if p.UseStrict && len([]byte(motd)) > 64 {
		return StatusBeta18{}, fmt.Errorf("%w: MOTD is longer than 64 bytes", ErrInvalidStatus)
	}

	// Parse online players
	online, err := strconv.ParseInt(onlineString, 10, 32)
	if err != nil {
		return StatusBeta18{}, fmt.Errorf("%w: %s", ErrInvalidStatus, err)
	}

	// Parse max players
	max, err := strconv.ParseInt(maxString, 10, 32)
	if err != nil {
		return StatusBeta18{}, fmt.Errorf("%w: %s", ErrInvalidStatus, err)
	}

	return StatusBeta18{MOTD: motd, OnlinePlayers: int(online), MaxPlayers: int(max)}, nil
}
