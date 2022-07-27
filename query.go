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
	queryFullStatPadding          = []byte{0xff, 0xff, 0xff, 0x01}
	queryKVSectionPadding         = []byte{0x73, 0x70, 0x6c, 0x69, 0x74, 0x6e, 0x75, 0x6d, 0x00, 0x80, 0x00}
	queryPlayerSectionPadding     = []byte{0x01, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x00, 0x00}
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
	queryGameID   = "MINECRAFT"
)

// BasicQueryStatus holds query status response returned Minecraft servers via Query protocol.
type BasicQueryStatus struct {
	MOTD          string
	GameType      string
	Map           string
	OnlinePlayers int
	MaxPlayers    int
	Port          int
	IP            string
}

type FullQueryPluginEntry struct {
	Name    string
	Version string
}

type FullQueryStatus struct {
	MOTD          string
	GameType      string
	GameID        string
	Version       string
	ServerVersion string
	Plugins       []FullQueryPluginEntry
	Map           string
	OnlinePlayers int
	MaxPlayers    int
	SamplePlayers []string
	Port          int
	Host          string
}

// QueryBasic queries Minecraft servers.
//goland:noinspection GoUnusedExportedFunction
func QueryBasic(host string, port int) (*BasicQueryStatus, error) {
	return defaultPinger.QueryBasic(host, port)
}

// QueryBasic queries Minecraft servers.
//goland:noinspection GoUnusedExportedFunction
func (p *Pinger) QueryBasic(host string, port int) (*BasicQueryStatus, error) {
	conn, err := p.openUDPConn(host, port)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()

	sessionID, token, err := p.createSession(conn)
	if err != nil {
		return nil, err
	}

	res, err := p.requestBasicStat(conn, sessionID, token)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func QueryFull(host string, port int) (*FullQueryStatus, error) {
	return defaultPinger.QueryFull(host, port)
}

func (p *Pinger) QueryFull(host string, port int) (*FullQueryStatus, error) {
	conn, err := p.openUDPConn(host, port)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()

	sessionID, token, err := p.createSession(conn)
	if err != nil {
		return nil, err
	}

	res, err := p.requestFullStat(conn, sessionID, token)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Pinger) createSession(conn *net.UDPConn) (int32, int32, error) {
	sessionID := generateSessionID()
	if err := p.writeQueryHandshakePacket(conn, sessionID); err != nil {
		return 0, 0, err
	}

	content, err := p.readQueryHandshakeResponsePacket(conn, sessionID)
	if err != nil {
		return 0, 0, err
	}

	token, err := p.parseQueryHandshakeResponse(content)
	if err != nil {
		return 0, 0, err
	}

	return sessionID, token, nil
}

func (p *Pinger) requestBasicStat(conn *net.UDPConn, sessionID int32, token int32) (*BasicQueryStatus, error) {
	if err := p.writeQueryBasicStatPacket(conn, sessionID, token); err != nil {
		return nil, err
	}

	content, err := p.readQueryStatResponsePacket(conn, sessionID)
	if err != nil {
		return nil, err
	}

	return p.parseQueryBasicStatResponse(content, p.UseStrict)
}

func (p *Pinger) requestFullStat(conn *net.UDPConn, sessionID int32, token int32) (*FullQueryStatus, error) {
	if err := p.writeQueryFullStatPacket(conn, sessionID, token); err != nil {
		return nil, err
	}

	content, err := p.readQueryStatResponsePacket(conn, sessionID)
	if err != nil {
		return nil, err
	}

	return p.parseQueryFullStatResponse(content, p.UseStrict)
}

// Communication

func (p *Pinger) writeQueryHandshakePacket(conn *net.UDPConn, sessionID int32) error {
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

func (p *Pinger) readQueryHandshakeResponsePacket(conn *net.UDPConn, sessionID int32) (io.Reader, error) {
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

func (p *Pinger) writeQueryBasicStatPacket(conn *net.UDPConn, sessionID int32, token int32) error {
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

func (p *Pinger) readQueryStatResponsePacket(conn *net.UDPConn, sessionID int32) (io.Reader, error) {
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

func (p *Pinger) writeQueryFullStatPacket(conn *net.UDPConn, sessionID int32, token int32) error {
	var packet bytes.Buffer

	// Write request packet header
	_, _ = packet.Write(queryRequestHeader)

	// Write packet type
	_ = packet.WriteByte(queryPacketTypeStat)

	// Write session ID
	_ = binary.Write(&packet, binary.BigEndian, sessionID)

	// Write token
	_ = binary.Write(&packet, binary.BigEndian, token)

	// Write padding
	_, _ = packet.Write(queryFullStatPadding)

	_, err := packet.WriteTo(conn)
	return err
}

// Response processing

func (p *Pinger) parseQueryHandshakeResponse(reader io.Reader) (int32, error) {
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

func (p *Pinger) parseQueryBasicStatResponse(reader io.Reader, useStrict bool) (*BasicQueryStatus, error) {
	data, _ := io.ReadAll(reader)
	if len(data) == 0 {
		return nil, fmt.Errorf("%w: empty response body", ErrInvalidStatus)
	}
	if bytes.HasSuffix(data, queryResponseStringTerminator) {
		data = data[:len(data)-len(queryResponseStringTerminator)]
	} else if useStrict {
		return nil, fmt.Errorf("%w: response body is not NUL-termianted", ErrInvalidStatus)
	}

	fields := strings.SplitN(string(data), string(queryResponseStringTerminator), 6)
	if len(fields) != 6 {
		return nil, fmt.Errorf("%w: expected 5 first string fields in response body, got %#v", ErrInvalidStatus, len(fields)-1)
	}
	motd, gameType, mapName, onlinePlayersStr, maxPlayerStr := fields[0], fields[1], fields[2], fields[3], fields[4]

	if gameType != queryGameType && useStrict {
		return nil, fmt.Errorf("%w: expected gametype field to be %#v, got %#v", ErrInvalidStatus, queryGameType, gameType)
	}

	onlinePlayers, err := strconv.ParseInt(onlinePlayersStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: could not parse online players count: %s", ErrInvalidStatus, err)
	}

	maxPlayers, err := strconv.ParseInt(maxPlayerStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: could not parse max players count: %s", ErrInvalidStatus, err)
	}

	remReader := bytes.NewReader([]byte(fields[5]))

	var port int16
	if err = binary.Read(remReader, binary.LittleEndian, &port); err != nil {
		return nil, err
	}

	hostBytes, err := io.ReadAll(remReader)
	if err != nil {
		return nil, err
	}
	return &BasicQueryStatus{
		MOTD:          motd,
		GameType:      gameType,
		Map:           mapName,
		OnlinePlayers: int(onlinePlayers),
		MaxPlayers:    int(maxPlayers),
		Port:          int(port),
		IP:            string(hostBytes),
	}, nil
}

func (p *Pinger) parseQueryFullStatResponse(reader io.Reader, useStrict bool) (*FullQueryStatus, error) {
	data, _ := io.ReadAll(reader)
	if len(data) == 0 {
		return nil, fmt.Errorf("%w: empty response body", ErrInvalidStatus)
	}
	if bytes.HasSuffix(data, queryResponseStringTerminator) {
		data = data[:len(data)-len(queryResponseStringTerminator)]
	} else if useStrict {
		return nil, fmt.Errorf("%w: response body is not NUL-termianted", ErrInvalidStatus)
	}
	dataReader := bytes.NewReader(data)

	pb := make([]byte, len(queryKVSectionPadding))
	if _, err := dataReader.Read(pb); err != nil {
		return nil, err
	} else if !bytes.Equal(pb, queryKVSectionPadding) && p.UseStrict {
		return nil, fmt.Errorf("%w: key-value section padding is invalid", ErrInvalidStatus)
	}

	fields, err := queryReadFullStatFieldMap(dataReader)
	if err != nil {
		return nil, err
	}

	pb = make([]byte, len(queryPlayerSectionPadding))
	if _, err = dataReader.Read(pb); err != nil {
		return nil, err
	} else if !bytes.Equal(pb, queryPlayerSectionPadding) && p.UseStrict {
		return nil, fmt.Errorf("%w: player section padding is invalid", ErrInvalidStatus)
	}

	players, err := queryReadFullStatPlayerList(dataReader)
	if err != nil {
		return nil, err
	}

	motd, err := queryGetFullStatField(fields, "hostname")
	if err != nil {
		return nil, err
	}

	gameType, err := queryGetFullStatField(fields, "gametype")
	if err != nil {
		return nil, err
	} else if gameType != queryGameType && useStrict {
		return nil, fmt.Errorf("%w: expected gametype field to be %#v, got %#v", ErrInvalidStatus, queryGameType, gameType)
	}

	gameID, err := queryGetFullStatField(fields, "game_id")
	if err != nil {
		return nil, err
	} else if gameID != queryGameID && useStrict {
		return nil, fmt.Errorf("%w: expected game_id field to be %#v, got %#v", ErrInvalidStatus, queryGameID, gameID)
	}

	version, err := queryGetFullStatField(fields, "version")
	if err != nil {
		return nil, err
	}

	serverVersionStr, err := queryGetFullStatField(fields, "plugins")
	if err != nil {
		return nil, err
	}
	serverVersion, plugins, err := queryParseFullStatPluginsList(serverVersionStr)
	if err != nil {
		return nil, fmt.Errorf("%w: could not parse plugins field: %s", ErrInvalidStatus, err)
	}

	mapName, err := queryGetFullStatField(fields, "map")
	if err != nil {
		return nil, err
	}

	onlinePlayersStr, err := queryGetFullStatField(fields, "numplayers")
	if err != nil {
		return nil, err
	}
	onlinePlayers, err := strconv.ParseInt(onlinePlayersStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: could not parse numplayers field: %s", ErrInvalidStatus, err)
	}

	maxPlayersStr, err := queryGetFullStatField(fields, "maxplayers")
	if err != nil {
		return nil, err
	}
	maxPlayers, err := strconv.ParseInt(maxPlayersStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: could not parse maxplayers field: %s", ErrInvalidStatus, err)
	}

	portStr, err := queryGetFullStatField(fields, "hostport")
	if err != nil {
		return nil, err
	}
	port, err := strconv.ParseInt(portStr, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("%w: could not parse hostport field: %s", ErrInvalidStatus, err)
	}

	hostname, err := queryGetFullStatField(fields, "hostname")
	if err != nil {
		return nil, err
	}

	return &FullQueryStatus{
		MOTD:          motd,
		GameType:      gameType,
		GameID:        gameID,
		Version:       version,
		ServerVersion: serverVersion,
		Plugins:       plugins,
		Map:           mapName,
		OnlinePlayers: int(onlinePlayers),
		MaxPlayers:    int(maxPlayers),
		SamplePlayers: players,
		Port:          int(port),
		Host:          hostname,
	}, nil
}

func queryReadFullStatFieldMap(reader io.Reader) (map[string]string, error) {
	fields := make(map[string]string)
	for {
		key, err := readAllUntilZero(reader)
		if err != nil {
			return nil, err
		} else if len(key) == 0 {
			break
		}
		value, err := readAllUntilZero(reader)
		if err != nil {
			return nil, err
		}
		fields[string(key)] = string(value)
	}
	return fields, nil
}

func queryReadFullStatPlayerList(reader io.Reader) ([]string, error) {
	players := make([]string, 0, 10)
	for {
		nickname, err := readAllUntilZero(reader)
		if err != nil {
			return nil, err
		} else if len(nickname) == 0 {
			break
		}
		players = append(players, string(nickname))
	}
	return players, nil
}

func queryParseFullStatPluginsList(str string) (string, []FullQueryPluginEntry, error) {
	parts := strings.SplitN(str, ":", 2)
	if len(parts) < 2 {
		return parts[0], nil, nil
	}
	ver, rem := parts[0], parts[1]

	pluginNames := strings.Split(rem, ";")
	plugins := make([]FullQueryPluginEntry, len(pluginNames))
	for i, name := range pluginNames {
		nameParts := strings.SplitN(strings.TrimSpace(name), " ", 2)
		if len(nameParts) < 2 {
			return "", nil, fmt.Errorf("invalid plugin field syntax")
		}
		plugins[i] = FullQueryPluginEntry{nameParts[0], nameParts[1]}
	}

	return ver, plugins, nil
}

func queryGetFullStatField(fields map[string]string, key string) (string, error) {
	value, ok := fields[key]
	if !ok {
		return "", fmt.Errorf("%w: response body does not contain %s field", ErrInvalidStatus, key)
	}
	return value, nil
}

// Util

func generateSessionID() int32 { return int32(time.Now().Unix()) & querySessionIDMask }
