package ping

import (
	"fmt"
	"net"
	"time"
)

// Ping sends a sequence of packets necessary for performing so-called Server List Ping (see https://wiki.vg/Server_List_Ping.)
// This method does not set a read/write timeout; if a timeout is necessary, use PingWithTimeout.
func Ping(host string, port uint16) (*Response, error) {
	return PingWithTimeout(host, port, 0)
}

// PingWithTimeout sends a sequence of packets necessary for performing so-called Server List Ping (see https://wiki.vg/Server_List_Ping.)
//goland:noinspection GoNameStartsWithPackageName
func PingWithTimeout(host string, port uint16, timeout time.Duration) (*Response, error) {
	var deadline time.Time
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}

	conn, err := openConnection(host, port, deadline)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}

	res, err := sendServerListPing(conn, host, port)
	if err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	if err = conn.Close(); err != nil {
		return nil, fmt.Errorf("connection close error: %w", err)
	}

	return res, nil
}

// PingLegacy sends a legacy (<1.7) server list ping packet (see https://wiki.vg/Server_List_Ping.)
// This method does not set a read/write timeout; if a timeout is necessary, use PingLegacyWithTimeout.
//goland:noinspection GoNameStartsWithPackageName
func PingLegacy(host string, port uint16) (*LegacyResponse, error) {
	return PingLegacyWithTimeout(host, port, 0)
}

// PingLegacyWithTimeout sends a legacy (<1.7) server list ping packet (see https://wiki.vg/Server_List_Ping.)
//goland:noinspection GoNameStartsWithPackageName
func PingLegacyWithTimeout(host string, port uint16, timeout time.Duration) (*LegacyResponse, error) {
	var deadline time.Time
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}

	conn, err := openConnection(host, port, deadline)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}

	res, err := sendLegacyServerListPing(conn, host, port)
	if err != nil {
		return nil, fmt.Errorf("test: %w", err)
	}

	return res, nil
}

func openConnection(host string, port uint16, deadline time.Time) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	if err = conn.SetDeadline(deadline); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func sendServerListPing(conn net.Conn, host string, port uint16) (*Response, error) {
	if err := writeHandshake(conn, handshake{Host: host, Port: unsignedShort(port)}); err != nil {
		return nil, fmt.Errorf("handshake error: %w", err)
	}
	if err := writeRequest(conn); err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	if err := writePing(conn, ping{Payload: long(time.Now().Unix())}); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	res, err := readResponse(conn)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}

	return res, nil
}

func sendLegacyServerListPing(conn net.Conn, host string, port uint16) (*LegacyResponse, error) {
	if err := writeLegacyPing(conn, legacyPing{Host: host, Port: port}); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	res, err := readLegacyPong(conn)
	if err != nil {
		return nil, fmt.Errorf("pong error: %w", err)
	}

	return res, nil
}
