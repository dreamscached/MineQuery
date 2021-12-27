package ping

import (
	"fmt"
	"net"
	"time"
)

// Ping performs a Server List Ping interaction with modern (1.7 and newer) Minecraft server running on the specified host and the specified port.
func Ping(host string, port uint16) (*Response, error) {
	return PingWithTimeout(host, port, 0)
}

// PingWithTimeout performs a Server List Ping interaction with modern (1.7 and newer) Minecraft server running on the specified host and the specified port with
// read and write timeout.
func PingWithTimeout(host string, port uint16, timeout time.Duration) (*Response, error) {
	var deadline time.Time
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}

	conn, err := newTCPConn(host, port, deadline)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}
	defer func() { _ = conn.Close() }()

	res, err := sendServerListPing(conn, host, port)
	if err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	return res, nil
}

func sendServerListPing(conn net.Conn, host string, port uint16) (*Response, error) {
	if err := writeHandshake(conn, handshake{Host: host, Port: unsignedShort(port)}); err != nil {
		return nil, fmt.Errorf("handshake error: %w", err)
	}
	if err := writeRequest(conn); err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	res, err := readResponse(conn)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}

	return res, nil
}

// PingLegacy performs a Server List Ping interaction with legacy (1.4 to 1.6) Minecraft server running on the specified host and the specified port.
func PingLegacy(host string, port uint16) (*LegacyResponse, error) {
	return PingLegacyWithTimeout(host, port, 0)
}

// PingLegacyWithTimeout performs a Server List Ping interaction with legacy (1.4 to 1.6) Minecraft server running on the specified host and the specified port with
// read and write timeout.
func PingLegacyWithTimeout(host string, port uint16, timeout time.Duration) (*LegacyResponse, error) {
	var deadline time.Time
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}

	conn, err := newTCPConn(host, port, deadline)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}
	defer func() { _ = conn.Close() }()

	res, err := sendLegacyServerListPing(conn, host, port)
	if err != nil {
		return nil, fmt.Errorf("test: %w", err)
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

// PingAncient performs a Server List Ping interaction with old (Beta 1.8 to 1.3) Minecraft server running on the specified host and the specified port.
func PingAncient(host string, port uint16) (*AncientResponse, error) {
	return PingAncientWithTimeout(host, port, 0)
}

// PingAncientWithTimeout performs a Server List Ping interaction with old (Beta 1.8 to 1.3) Minecraft server running on the specified host and the specified port with
// read and write timeout.
func PingAncientWithTimeout(host string, port uint16, timeout time.Duration) (*AncientResponse, error) {
	var deadline time.Time
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}

	conn, err := newTCPConn(host, port, deadline)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}
	defer func() { _ = conn.Close() }()

	res, err := sendAncientServerListPing(conn)
	if err != nil {
		return nil, fmt.Errorf("test: %w", err)
	}

	return res, nil
}

func sendAncientServerListPing(conn net.Conn) (*AncientResponse, error) {
	if err := writeAncientPing(conn); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	res, err := readAncientPong(conn)
	if err != nil {
		return nil, fmt.Errorf("pong error: %w", err)
	}

	return res, nil
}
