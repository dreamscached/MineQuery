package minequery

import (
	"fmt"
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
	ProtocolVersion int
	ServerVersion   string
	MOTD            string
	OnlinePlayers   int
	MaxPlayers      int
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

	// Write 2-byte FE 01 ping packet
	if err = writeBytes(conn, ping14PingPacket); err != nil {
		return Status14{}, err
	}

	// Read packet type, return error if it isn't FF kick packet
	packetType, err := readByte(conn)
	if err != nil {
		return Status14{}, err
	} else if packetType != ping14ResponsePacketID {
		return Status14{}, fmt.Errorf("expected packet ID %#x, but instead got %#x", ping14ResponsePacketID, packetType)
	}

	// Read packet length, return error if it isn't readable as unsigned short
	// Worth noting that this needs to be multiplied by two further on (for encoding reasons, most probably)
	length, err := readUShort(conn)
	if err != nil {
		return Status14{}, err
	}

	// Read remainder of the status packet as raw bytes
	// This is a UTF-16BE string separated by § (paragraph sign)
	// where [0] is protocol version, [1] is Minecraft version, [3] is MOTD, [4] is online players
	// and [5] is max players
	dataEncoded, err := readNBytes(conn, int(length*2))
	if err != nil {
		return Status14{}, err
	}

	// Decode UTF16-BE and return error if unable to decode
	dataString, err := utf16BEDecoder.String(string(dataEncoded))
	if err != nil {
		return Status14{}, err
	}

	// Split status string, parse and map to struct returning errors if conversions fail
	fields := strings.Split(dataString, ping14ResponseFieldSeparator)
	if len(fields) != 5 {
		return Status14{}, fmt.Errorf("%w: expected 5 status fields, got %d", ErrInvalidStatus, len(fields))
	}
	protocolVersionString, serverVersion, motd, onlineString, maxString := fields[0], fields[1], fields[2], fields[3], fields[4]

	// Parse protocol version
	protocolVersion, err := strconv.ParseInt(protocolVersionString, 10, 32)
	if err != nil {
		return Status14{}, fmt.Errorf("%w: %s", ErrInvalidStatus, err)
	}

	// Parse online players
	online, err := strconv.ParseInt(onlineString, 10, 32)
	if err != nil {
		return Status14{}, fmt.Errorf("%w: %s", ErrInvalidStatus, err)
	}

	// Parse max players
	max, err := strconv.ParseInt(maxString, 10, 32)
	if err != nil {
		return Status14{}, fmt.Errorf("%w: %s", ErrInvalidStatus, err)
	}

	return Status14{
		ProtocolVersion: int(protocolVersion),
		ServerVersion:   serverVersion,
		MOTD:            motd,
		OnlinePlayers:   int(online),
		MaxPlayers:      int(max),
	}, nil
}
