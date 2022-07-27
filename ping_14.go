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
	ping14ResponsePacketID              byte = 0xff
	ping14ResponsePayloadFieldSeparator      = "§"
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
func Ping14(host string, port int) (*Status14, error) {
	return defaultPinger.Ping14(host, port)
}

// Ping14 pings 1.4 to 1.6 (exclusively) Minecraft servers (Notchian servers of more late versions also respond to
// this ping packet.)
func (p *Pinger) Ping14(host string, port int) (*Status14, error) {
	conn, err := p.openTCPConn(host, port)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()

	// Send ping packet
	if err = p.ping14WritePingPacket(conn); err != nil {
		return nil, fmt.Errorf("could not write ping packet: %w", err)
	}

	// Read status response (note: uses the same packet reading approach as 1.4)
	payload, err := p.ping14ReadResponsePayload(conn)
	if err != nil {
		return nil, fmt.Errorf("could not read response packet: %w", err)
	}

	// Parse response data from status packet
	res, err := p.ping14ParseResponsePayload(payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse status from response packet: %w", err)
	}

	return res, nil
}

// Communication

func (p *Pinger) ping14WritePingPacket(writer io.Writer) error {
	// Write 2-byte FE 01 ping packet
	err := writeBytes(writer, ping14PingPacket)
	return err
}

func (p *Pinger) ping14ReadResponsePayload(reader io.Reader) ([]byte, error) {
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
	payload := bytes.NewBuffer(make([]byte, 0, length*2))
	if _, err = io.CopyN(payload, reader, int64(length*2)); err != nil {
		return nil, err
	}

	decoded, err := utf16BEDecoder.Bytes(payload.Bytes())
	if err != nil {
		return nil, err
	}

	// Return UTF16-BE decoder with data as input
	return decoded, nil
}

// Response processing

func (p *Pinger) ping14ParseResponsePayload(payload []byte) (*Status14, error) {
	// NOTE: Spigot 1.4 servers reply with 1.6 response format.
	// See https://github.com/alteamc/minequery/issues/31 for details.
	// Check if data string begins with '§1\x00' (00 a7 00 31 00 00) and pass processing to 1.6 logic in this case.
	if bytes.HasPrefix(payload, ping16ResponsePrefix) {
		if p.UseStrict {
			return nil, fmt.Errorf("%w: server unexpectedly replied with 1.6 response", ErrInvalidStatus)
		}

		res, err := p.ping16ParseResponsePayload(payload)
		if err != nil {
			return nil, fmt.Errorf("could not parse status from response packet: %w", err)
		}
		return &Status14{
			MOTD:          res.MOTD,
			OnlinePlayers: res.OnlinePlayers,
			MaxPlayers:    res.MaxPlayers,
		}, nil
	}

	// Split status string, parse and map to struct returning errors if conversions fail
	fields := strings.Split(string(payload), ping14ResponsePayloadFieldSeparator)
	if len(fields) != 3 {
		return nil, fmt.Errorf("%w: expected 3 status fields, got %d", ErrInvalidStatus, len(fields))
	}
	motd, onlineString, maxString := fields[0], fields[1], fields[2]

	// Parse online players
	online, err := strconv.ParseInt(onlineString, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: could not parse online players count: %s", ErrInvalidStatus, err)
	}

	// Parse max players
	max, err := strconv.ParseInt(maxString, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: could not parse max players count: %s", ErrInvalidStatus, err)
	}

	return &Status14{
		MOTD:          motd,
		OnlinePlayers: int(online),
		MaxPlayers:    int(max),
	}, nil
}
