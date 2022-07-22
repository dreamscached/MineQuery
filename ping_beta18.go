package minequery

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
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
//goland:noinspection GoUnusedExportedFunction
func PingBeta18(host string, port int) (StatusBeta18, error) {
	return defaultPinger.PingBeta18(host, port)
}

// PingBeta18 pings Beta 1.8 to Release 1.4 (exclusively) Minecraft servers (Notchian servers of more late versions
// also respond to this ping packet.)
func (p Pinger) PingBeta18(host string, port int) (StatusBeta18, error) {
	conn, err := p.openTCPConn(host, port)
	if err != nil {
		return StatusBeta18{}, err
	}
	defer func() { _ = conn.Close() }()

	// Send ping packet
	if err = writePingPacketBeta18(conn); err != nil {
		return StatusBeta18{}, fmt.Errorf("could not write ping packet: %w", err)
	}

	// Read status response (note: uses the same packet reading approach as 1.4)
	content, err := readResponsePacketBeta18(conn)
	if err != nil {
		return StatusBeta18{}, fmt.Errorf("could not read response packet: %w", err)
	}

	// Parse response data from status packet
	res, err := parseResponseDataBeta18(content)
	if err != nil {
		return StatusBeta18{}, fmt.Errorf("could not parse status from response packet: %w", err)
	}

	return res, nil
}

// Communication

func writePingPacketBeta18(writer io.Writer) error {
	// Write single-byte FE ping packet
	err := writeBytes(writer, pingBeta18PingPacket)
	return err
}

func readResponsePacketBeta18(reader io.Reader) (io.Reader, error) {
	// Read packet type, return error if it isn't FF kick packet
	id, err := readByte(reader)
	if err != nil {
		return nil, err
	} else if id != pingBeta18ResponsePacketID {
		return nil, fmt.Errorf("expected packet ID %#x, but instead got %#x", ping16ResponsePacketID, id)
	}

	// Read packet length, return error if it isn't readable as unsigned short
	// Worth noting that this needs to be multiplied by two further on (for encoding reasons, most probably)
	length, err := readUShort(reader)
	if err != nil {
		return nil, err
	}

	// Read remainder of the status packet as raw bytes
	// This is a UTF-16BE string separated by § (paragraph sign)
	var data bytes.Buffer
	if _, err = io.CopyN(&data, reader, int64(length*2)); err != nil {
		return nil, err
	}

	// Return UTF16-BE decoder with data as input
	return utf16BEDecoder.Reader(&data), nil
}

// Response processing

func parseResponseDataBeta18(reader io.Reader) (StatusBeta18, error) {
	data, err := ReadAll(reader)
	if err != nil {
		return StatusBeta18{}, err
	}

	// Split status string, parse and map to struct returning errors if conversions fail
	fields := strings.Split(string(data), pingBeta18ResponseFieldSeparator)
	if len(fields) != 3 {
		return StatusBeta18{}, fmt.Errorf("%w: expected 3 status fields, got %d", ErrInvalidStatus, len(fields))
	}
	motd, onlineString, maxString := fields[0], fields[1], fields[2]

	// Parse online players
	online, err := strconv.ParseInt(onlineString, 10, 32)
	if err != nil {
		return StatusBeta18{}, fmt.Errorf("%w: could not parse online players count: %s", ErrInvalidStatus, err)
	}

	// Parse max players
	max, err := strconv.ParseInt(maxString, 10, 32)
	if err != nil {
		return StatusBeta18{}, fmt.Errorf("%w: could not parse max players count: %s", ErrInvalidStatus, err)
	}

	return StatusBeta18{
		MOTD:          motd,
		OnlinePlayers: int(online),
		MaxPlayers:    int(max),
	}, nil
}