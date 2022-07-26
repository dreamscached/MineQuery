package minequery

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	queryRequestHeader            = []byte{0xfe, 0xfd}
	queryResponseStringTerminator = []byte{0x0}
)

const (
	queryPacketTypeHandshake byte = 9
	queryPacketTypeStat      byte = 0
)

const (
	querySessionIDMask int32 = 0x0f0f0f0f
)

const (
	queryGameType = "SMP"
)

type QueryStatus struct {
	MOTD          string
	GameType      string
	Map           string
	OnlinePlayers int
	MaxPlayers    int
	Port          int
	IP            string
}

type Querier struct {
	Timeout   time.Duration
	UseStrict bool
}

func (q Querier) Query(host string, port int) (QueryStatus, error) {
	conn, err := q.openUDPConn(host, port)
	if err != nil {
		return QueryStatus{}, err
	}
	defer func() { _ = conn.Close() }()

	sessionID, token, err := q.createSession(conn)
	if err != nil {
		return QueryStatus{}, err
	}

	res, err := q.requestStat(conn, sessionID, token)
	if err != nil {
		return QueryStatus{}, err
	}
	return res, nil
}

func (q Querier) createSession(conn *net.UDPConn) (int32, int32, error) {
	sessionID := generateSessionID()
	if err := writeQueryHandshakePacket(conn, sessionID); err != nil {
		return 0, 0, err
	}

	content, err := readQueryHandshakeResponsePacket(conn, sessionID)
	if err != nil {
		return 0, 0, err
	}

	token, err := parseQueryHandshakeResponse(content)
	if err != nil {
		return 0, 0, err
	}

	return sessionID, token, nil
}

func (q Querier) requestStat(conn *net.UDPConn, sessionID int32, token int32) (QueryStatus, error) {
	if err := writeQueryStatPacket(conn, sessionID, token); err != nil {
		return QueryStatus{}, err
	}

	content, err := readQueryStatResponsePacket(conn, sessionID)
	if err != nil {
		return QueryStatus{}, err
	}

	return parseQueryStatResponse(content, q.UseStrict)
}

// Communication

func writeQueryHandshakePacket(conn *net.UDPConn, sessionID int32) error {
	var packet bytes.Buffer

	// Write request packet header
	_, _ = packet.Write(queryRequestHeader)

	// Write packet type
	_ = packet.WriteByte(queryPacketTypeHandshake)

	// Write session ID
	_ = binary.Write(&packet, binary.BigEndian, sessionID)

	_, err := packet.WriteTo(conn)
	return err
}

func readQueryHandshakeResponsePacket(conn *net.UDPConn, sessionID int32) (io.Reader, error) {
	b := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(b)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(b[:n])

	id, err := reader.ReadByte()
	if err != nil {
		return nil, err
	} else if id != queryPacketTypeHandshake {
		return nil, fmt.Errorf("expected packet ID %#x, but instead got %#x", queryPacketTypeHandshake, id)
	}

	var resSessionID int32
	if err = binary.Read(reader, binary.BigEndian, &resSessionID); err != nil {
		return nil, err
	} else if resSessionID != sessionID {
		return nil, fmt.Errorf("expected session ID %#x, but instead got %#x", sessionID, resSessionID)
	}

	return reader, nil
}

func writeQueryStatPacket(conn *net.UDPConn, sessionID int32, token int32) error {
	var packet bytes.Buffer

	// Write request packet header
	_, _ = packet.Write(queryRequestHeader)

	// Write packet type
	_ = packet.WriteByte(queryPacketTypeStat)

	// Write session ID
	_ = binary.Write(&packet, binary.BigEndian, sessionID)

	// Write token
	_ = binary.Write(&packet, binary.BigEndian, token)

	_, err := packet.WriteTo(conn)
	return err
}

func readQueryStatResponsePacket(conn *net.UDPConn, sessionID int32) (io.Reader, error) {
	b := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(b)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(b[:n])

	id, err := reader.ReadByte()
	if err != nil {
		return nil, err
	} else if id != queryPacketTypeStat {
		return nil, fmt.Errorf("expected packet ID %#x, but instead got %#x", queryPacketTypeStat, id)
	}

	var resSessionID int32
	if err = binary.Read(reader, binary.BigEndian, &resSessionID); err != nil {
		return nil, err
	} else if resSessionID != sessionID {
		return nil, fmt.Errorf("expected session ID %#x, but instead got %#x", sessionID, resSessionID)
	}

	return reader, nil
}

// Response processing

func parseQueryHandshakeResponse(reader io.Reader) (int32, error) {
	token, _ := io.ReadAll(reader)
	if len(token) == 0 {
		return 0, fmt.Errorf("challenge token is empty")
	} else if !bytes.HasSuffix(token, queryResponseStringTerminator) {
		return 0, fmt.Errorf("challenge token did not end with NUL byte")
	}

	tokenInt, err := strconv.ParseInt(string(token[:len(token)-1]), 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(tokenInt), nil
}

func parseQueryStatResponse(reader io.Reader, useStrict bool) (QueryStatus, error) {
	data, _ := io.ReadAll(reader)
	if len(data) == 0 {
		return QueryStatus{}, fmt.Errorf("%w: empty response body", ErrInvalidStatus)
	}
	if bytes.HasSuffix(data, queryResponseStringTerminator) {
		data = data[:len(data)-len(queryResponseStringTerminator)]
	} else if useStrict {
		return QueryStatus{}, fmt.Errorf("%w: response body is not NUL-termianted", ErrInvalidStatus)
	}

	fields := strings.SplitN(string(data), string(queryResponseStringTerminator), 6)
	if len(fields) != 6 {
		return QueryStatus{}, fmt.Errorf("%w: expected 5 first string fields in response body, got %#v", ErrInvalidStatus, len(fields)-1)
	}
	motd, gameType, mapName, onlinePlayersStr, maxPlayerStr := fields[0], fields[1], fields[2], fields[3], fields[4]

	if gameType != queryGameType && useStrict {
		return QueryStatus{}, fmt.Errorf("%w: expected gametype field to be %#v, got %#v", ErrInvalidStatus, queryGameType, gameType)
	}

	onlinePlayers, err := strconv.ParseInt(onlinePlayersStr, 10, 32)
	if err != nil {
		return QueryStatus{}, fmt.Errorf("%w: could not parse online players count: %s", ErrInvalidStatus, err)
	}

	maxPlayers, err := strconv.ParseInt(maxPlayerStr, 10, 32)
	if err != nil {
		return QueryStatus{}, fmt.Errorf("%w: could not parse max players count: %s", ErrInvalidStatus, err)
	}

	remReader := bytes.NewReader([]byte(fields[5]))

	var port int16
	if err = binary.Read(remReader, binary.LittleEndian, &port); err != nil {
		return QueryStatus{}, err
	}

	hostBytes, err := io.ReadAll(remReader)
	if err != nil {
		return QueryStatus{}, err
	}
	return QueryStatus{
		MOTD:          motd,
		GameType:      gameType,
		Map:           mapName,
		OnlinePlayers: int(onlinePlayers),
		MaxPlayers:    int(maxPlayers),
		Port:          int(port),
		IP:            string(hostBytes),
	}, nil
}

// Util

func generateSessionID() int32 { return int32(time.Now().Unix()) & querySessionIDMask }
