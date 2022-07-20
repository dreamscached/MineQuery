package minequery

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var ping14PingPacket = []byte{0xfe}

const (
	ping14ResponsePacketID       byte = 0xff
	ping14ResponseFieldSeparator      = "§"
)

// Status14 holds status response returned by 1.4 to 1.6 (exclusively) Minecraft servers.
type Status14 struct {
	MOTD          string
	OnlinePlayers int
	MaxPlayers    int
}

// Ping14 pings 1.4 to 1.6 (exclusively) Minecraft servers (Notchian servers of more late versions also respond to
// this ping packet.)
//goland:noinspection GoUnusedExportedFunction
func Ping14(host string, port int) (Status14, error) {
	return defaultPinger.Ping14(host, port)
}

// Ping14 pings 1.4 to 1.6 (exclusively) Minecraft servers (Notchian servers of more late versions also respond to
// this ping packet.)
func (p Pinger) Ping14(host string, port int) (Status14, error) {
	conn, err := p.openTCPConn(host, port)
	defer func() { _ = conn.Close() }()
	if err != nil {
		return Status14{}, err
	}

	// Send ping packet
	if err = writePingPacket14(conn); err != nil {
		return Status14{}, fmt.Errorf("could not write ping packet: %w", err)
	}

	// Read status response (note: uses the same packet reading approach as 1.4)
	content, err := readResponsePacket14(conn)
	if err != nil {
		return Status14{}, fmt.Errorf("could not read response packet: %w", err)
	}

	// Parse response data from status packet
	res, err := parseResponseData14(content)
	if err != nil {
		return Status14{}, fmt.Errorf("could not parse status from response packet: %w", err)
	}

	return res, nil
}

// Communication

func writePingPacket14(writer io.Writer) error {
	// Write 2-byte FE 01 ping packet
	err := writeBytes(writer, ping14PingPacket)
	return err
}

func readResponsePacket14(reader io.Reader) (io.Reader, error) {
	// Read packet type, return error if it isn't FF kick packet
	id, err := readByte(reader)
	if err != nil {
		return nil, err
	} else if id != ping14ResponsePacketID {
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

func parseResponseData14(reader io.Reader) (Status14, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return Status14{}, err
	}

	// Split status string, parse and map to struct returning errors if conversions fail
	fields := strings.Split(string(data), ping14ResponseFieldSeparator)
	if len(fields) != 3 {
		return Status14{}, fmt.Errorf("%w: expected 3 status fields, got %d", ErrInvalidStatus, len(fields))
	}
	motd, onlineString, maxString := fields[0], fields[1], fields[2]

	// Parse online players
	online, err := strconv.ParseInt(onlineString, 10, 32)
	if err != nil {
		return Status14{}, fmt.Errorf("%w: could not parse online players count: %s", ErrInvalidStatus, err)
	}

	// Parse max players
	max, err := strconv.ParseInt(maxString, 10, 32)
	if err != nil {
		return Status14{}, fmt.Errorf("%w: could not parse max players count: %s", ErrInvalidStatus, err)
	}

	return Status14{
		MOTD:          motd,
		OnlinePlayers: int(online),
		MaxPlayers:    int(max),
	}, nil
}
